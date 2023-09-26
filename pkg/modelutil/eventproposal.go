package modelutil

import (
	"context"
	"time"
	"whale/pkg/dbquery"
	"whale/pkg/gqlient/hoopoe"
	"whale/pkg/gqlient/smew"
	"whale/pkg/loader"
	"whale/pkg/models"
	"whale/pkg/whalecode"

	"github.com/letjoy-club/mida-tool/dbutil"
	"github.com/letjoy-club/mida-tool/midacode"
	"github.com/letjoy-club/mida-tool/midacontext"
	"github.com/letjoy-club/mida-tool/shortid"
)

func CreateEventProposal(ctx context.Context, uid string, param models.CreateEventProposalParam) (*models.EventProposal, error) {
	db := dbutil.GetDB(ctx)

	now := time.Now()
	if param.Deadline == nil {
		deadline := now.Add(5 * 24 * time.Hour)
		param.Deadline = &deadline
	} else {
		if param.Deadline.Before(now) {
			return nil, whalecode.ErrDeadlineBeforeNow
		}
		if param.Deadline.Sub(now) > 30*24*time.Hour {
			return nil, whalecode.ErrDeadlineTooLong
		}
	}

	EventProposal := dbquery.Use(db).EventProposal
	proposal := &models.EventProposal{
		ID:        shortid.NewWithTime("ep_", 3),
		Active:    true,
		UserID:    uid,
		MaxNum:    param.MaxNum,
		Desc:      param.Desc,
		Title:     param.Title,
		CityID:    param.CityID,
		TopicID:   param.TopicID,
		Address:   param.Address,
		Latitude:  param.Latitude,
		Longitude: param.Longitude,
		StartAt:   param.StartAt,
		Deadline:  *param.Deadline,
		Images:    param.Images,
	}

	if len(param.Images) < 1 {
		return nil, whalecode.ErrImageEmpty
	}
	if len(param.Images) > 5 {
		return nil, whalecode.ErrImageNumExceedLimit
	}

	{
		resp, err := hoopoe.ImageProcess(ctx, midacontext.GetServices(ctx).Hoopoe, param.Images, hoopoe.BizTypeEventproposal, proposal.ID)
		if err != nil {
			return nil, err
		}
		param.Images = resp.ImagesProcess
	}

	if proposal.MaxNum <= 0 {
		proposal.MaxNum = 100
	}
	if proposal.MaxNum > 50 {
		proposal.MaxNum = 50
	}
	if proposal.Deadline.Before(time.Now()) {
		proposal.Deadline = time.Now().Add(24 * time.Hour)
	}
	if proposal.StartAt.Before(time.Now()) {
		proposal.StartAt = time.Now().Add(24 * time.Hour)
	}
	if proposal.Longitude < 70 || proposal.Longitude > 135 {
		return nil, whalecode.ErrLatLngOutOfRange
	}
	if proposal.Latitude < 0 || proposal.Latitude > 54 {
		return nil, whalecode.ErrLatLngOutOfRange
	}

	resp, err := smew.CreateEventProposalGroup(ctx, midacontext.GetServices(ctx).Smew, smew.CreateEventProposalGroupParam{
		EventProposalId: proposal.ID,
		UserId:          uid,
		TopicId:         proposal.TopicID,
	})
	if err != nil {
		return nil, err
	}

	proposal.ChatGroupID = resp.CreateEventProposalGroup

	err = EventProposal.WithContext(ctx).Create(proposal)
	if err != nil {
		return nil, err
	}
	midacontext.GetLoader[loader.Loader](ctx).AllEventProposalLoader.AppendNewData(proposal)
	return proposal, nil
}

func JoinEventProposal(ctx context.Context, uid string, eventID string) error {
	loader := midacontext.GetLoader[loader.Loader](ctx)
	proposal, err := loader.EventProposal.Load(ctx, eventID)()
	if err != nil {
		return err
	}
	if proposal.UserID == uid {
		return whalecode.ErrYouAreAlreadyInEvent
	}
	if !proposal.Active {
		return whalecode.ErrEventProposalClosed
	}
	participants, err := loader.EventProposalParticipants.Load(ctx, eventID)()

	// 已经加入了、或者已经退出了
	if participants.HasUser(uid) {
		return nil
	}

	_, err = smew.AddEventProposalGroupMember(ctx, midacontext.GetServices(ctx).Smew, proposal.ChatGroupID, uid)
	if err != nil {
		return err
	}

	db := dbutil.GetDB(ctx)
	err = dbquery.Use(db).Transaction(func(tx *dbquery.Query) error {
		EventProposal := dbquery.Use(db).EventProposal
		UserJoinEventProposal := dbquery.Use(db).UserJoinEventProposal
		err := UserJoinEventProposal.WithContext(ctx).Create(&models.UserJoinEventProposal{
			EventID: eventID,
			UserID:  uid,
			State:   string(models.JoinEventStateJoined),
		})
		if err != nil {
			return err
		}
		rx, err := EventProposal.WithContext(ctx).
			Where(EventProposal.ID.Eq(eventID)).
			UpdateSimple(
				EventProposal.JoinNum.Add(1),
			)
		if rx.RowsAffected != 1 {
			return whalecode.ErrEventProposalNotExisted
		}
		if err != nil {
			return err
		}
		return nil
	})
	if err == nil {
		loader.EventProposal.Clear(ctx, eventID)
		loader.EventProposalParticipants.Clear(ctx, eventID)
	}
	return err
}

func KickOutFromEventProposal(ctx context.Context, eventID string, actionUser, targetUserID string) error {
	loader := midacontext.GetLoader[loader.Loader](ctx)
	thunk := loader.EventProposal.Load(ctx, eventID)
	proposal, err := thunk()
	if err != nil {
		return err
	}
	if actionUser != "" {
		if proposal.UserID != actionUser {
			return midacode.ErrNotPermitted
		}
	}
	participants, err := loader.EventProposalParticipants.Load(ctx, eventID)()
	if err != nil {
		return err
	}
	if !participants.HasUser(targetUserID) {
		return nil
	}

	_, err = smew.RemoveEventProposalGroupMember(ctx, midacontext.GetServices(ctx).Smew, proposal.ChatGroupID, targetUserID, proposal.UserID)

	db := dbutil.GetDB(ctx)
	err = dbquery.Use(db).Transaction(func(tx *dbquery.Query) error {
		UserJoinEventProposal := tx.UserJoinEventProposal
		EventProposal := tx.EventProposal

		_, err = EventProposal.WithContext(ctx).Where(EventProposal.ID.Eq(eventID)).UpdateSimple(
			EventProposal.JoinNum.Add(-1),
		)

		_, err = UserJoinEventProposal.WithContext(ctx).
			Where(UserJoinEventProposal.EventID.Eq(eventID), UserJoinEventProposal.UserID.Eq(targetUserID)).
			UpdateSimple(
				UserJoinEventProposal.State.Value(models.JoinEventStateKickedOut.String()),
				UserJoinEventProposal.LeftAt.Value(time.Now()),
			)
		if err != nil {
			return err
		}
		return nil
	})
	if err != nil {
		return err
	}

	loader.EventProposalParticipants.Clear(ctx, eventID)
	loader.EventProposal.Clear(ctx, eventID)
	return nil
}

func CloseEventProposal(ctx context.Context, userID, eventID string, isTimeout bool) error {
	loader := midacontext.GetLoader[loader.Loader](ctx)
	proposal, err := loader.EventProposal.Load(ctx, eventID)()
	if err != nil {
		return err
	}
	if !proposal.Active {
		return nil
	}
	if userID != "" {
		if proposal.UserID != userID {
			return midacode.ErrNotPermitted
		}
	}

	if userID == "" {
		_, err = smew.CloseEventProposalGroup(ctx, midacontext.GetServices(ctx).Smew, proposal.ChatGroupID, smew.GroupCloseReasonUserclose)
	} else {
		if isTimeout {
			_, err = smew.CloseEventProposalGroup(ctx, midacontext.GetServices(ctx).Smew, proposal.ChatGroupID, smew.GroupCloseReasonTimeout)
		} else {
			_, err = smew.CloseEventProposalGroup(ctx, midacontext.GetServices(ctx).Smew, proposal.ChatGroupID, smew.GroupCloseReasonSystemclose)
		}
	}

	if err != nil {
		return err
	}
	err = dbquery.Use(dbutil.GetDB(ctx)).Transaction(func(tx *dbquery.Query) error {
		EventProposal := tx.EventProposal
		_, err = EventProposal.WithContext(ctx).Where(EventProposal.ID.Eq(eventID)).UpdateSimple(EventProposal.Active.Value(false))
		if err != nil {
			return err
		}
		UserJoinEventProposal := tx.UserJoinEventProposal
		_, err = UserJoinEventProposal.WithContext(ctx).
			Where(UserJoinEventProposal.EventID.Eq(eventID), UserJoinEventProposal.State.Eq(models.JoinEventStateJoined.String())).
			UpdateSimple(UserJoinEventProposal.State.Value(models.JoinEventStateClosed.String()))
		return err
	})
	if err == nil {
		loader.EventProposal.Clear(ctx, eventID)
		loader.EventProposalParticipants.Clear(ctx, eventID)
	}
	return err
}

func LeaveEventProposal(ctx context.Context, uid, eventID string) error {
	loader := midacontext.GetLoader[loader.Loader](ctx)
	proposal, err := loader.EventProposal.Load(ctx, eventID)()
	if err != nil {
		return err
	}

	// 离开的是自己的活动，不处理
	if proposal.UserID == uid {
		return nil
	}

	_, err = smew.RemoveEventProposalGroupMember(ctx, midacontext.GetServices(ctx).Smew, proposal.ChatGroupID, uid, uid)
	if err != nil {
		return err
	}

	participants, err := loader.EventProposalParticipants.Load(ctx, eventID)()
	if err != nil {
		return err
	}
	if !participants.HasUser(uid) {
		return nil
	}

	db := dbutil.GetDB(ctx)

	err = dbquery.Use(db).Transaction(func(tx *dbquery.Query) error {
		UserJoinEventProposal := tx.UserJoinEventProposal
		EventProposal := tx.EventProposal
		_, err := EventProposal.WithContext(ctx).Where(EventProposal.ID.Eq(eventID)).UpdateSimple(EventProposal.JoinNum.Add(-1))
		if err != nil {
			return err
		}
		_, err = UserJoinEventProposal.WithContext(ctx).
			Where(UserJoinEventProposal.EventID.Eq(eventID), UserJoinEventProposal.UserID.Eq(uid)).
			UpdateSimple(
				UserJoinEventProposal.State.Value(models.JoinEventStateLeft.String()),
				UserJoinEventProposal.LeftAt.Value(time.Now()),
			)
		return err
	})
	if err != nil {
		return err
	}
	loader.EventProposal.Clear(ctx, eventID)
	loader.EventProposalParticipants.Clear(ctx, eventID)
	return nil

}
