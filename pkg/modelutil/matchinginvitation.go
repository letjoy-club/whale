package modelutil

import (
	"context"
	"time"
	"whale/pkg/dbquery"
	"whale/pkg/loader"
	"whale/pkg/models"
	"whale/pkg/whalecode"

	"github.com/letjoy-club/mida-tool/midacontext"
	"github.com/letjoy-club/mida-tool/shortid"
)

func AcceptMatchingInvitation(ctx context.Context, invitation *models.MatchingInvitation) error {
	quotaThunk := midacontext.GetLoader[loader.Loader](ctx).MatchingQuota.Load(ctx, invitation.InviteeID)
	quota, err := quotaThunk()
	if err != nil {
		return err
	}
	if quota.Remain <= 0 {
		return whalecode.ErrMatchingQuotaNotEnough
	}

	db := midacontext.GetDB(ctx)
	// 生成匹配和匹配结果
	Matching := dbquery.Use(db).Matching
	MatchingResult := dbquery.Use(db).MatchingResult

	m1 := models.Matching{
		ID:             shortid.New("m_", 8),
		TopicID:        invitation.TopicID,
		UserID:         invitation.UserID,
		CityID:         invitation.CityID,
		AreaIDs:        invitation.AreaIDs,
		Gender:         models.GenderN.String(),
		ChatGroupState: models.ChatGroupStateUncreated.String(),
		Remark:         invitation.Remark,
		Deadline:       time.Now(),
		State:          string(models.MatchingStateMatched),
	}
	m2 := models.Matching{
		ID:             shortid.New("m_", 8),
		TopicID:        invitation.TopicID,
		UserID:         invitation.InviteeID,
		CityID:         invitation.CityID,
		AreaIDs:        invitation.AreaIDs,
		Gender:         models.GenderN.String(),
		ChatGroupState: models.ChatGroupStateUncreated.String(),
		Deadline:       time.Now(),
		State:          string(models.MatchingStateMatched),
	}
	result := models.MatchingResult{
		MatchingIDs:    []string{m1.ID, m2.ID},
		TopicID:        invitation.TopicID,
		UserIDs:        []string{invitation.UserID, invitation.InviteeID},
		ConfirmStates:  []string{models.InvitationConfirmStateConfirmed.String(), models.InvitationConfirmStateConfirmed.String()},
		ChatGroupState: models.ChatGroupStateUncreated.String(),
	}
	err = MatchingResult.WithContext(ctx).Create(&result)
	if err != nil {
		return err
	}
	m1.ResultID = result.ID
	m2.ResultID = result.ID
	if err := Matching.WithContext(ctx).Create(&m1, &m2); err != nil {
		return err
	}
	MatchingQuota := dbquery.Use(db).MatchingQuota
	if _, err := MatchingQuota.WithContext(ctx).Where(MatchingQuota.UserID.Eq(invitation.InviteeID)).UpdateSimple(MatchingQuota.Remain.Sub(1)); err != nil {
		return err
	}
	MatchingInvitation := dbquery.Use(db).MatchingInvitation

	_, err = MatchingInvitation.
		WithContext(ctx).
		Where(MatchingInvitation.ID.Eq(invitation.ID)).
		UpdateSimple(
			MatchingInvitation.ConfirmedAt.Value(time.Now()),
			MatchingInvitation.Closed.Value(true),
			MatchingInvitation.ConfirmState.Value(models.InvitationConfirmStateConfirmed.String()),
		)
	if err != nil {
		return err
	}

	err = CheckMatchingResultAndCreateChatGroup(ctx, &result)
	if err != nil {
		return err
	}
	midacontext.GetLoader[loader.Loader](ctx).Matching.Clear(ctx, m1.ID)
	midacontext.GetLoader[loader.Loader](ctx).Matching.Clear(ctx, m2.ID)
	midacontext.GetLoader[loader.Loader](ctx).MatchingInvitation.Clear(ctx, invitation.ID)
	return nil
}

func RejectMatchingInvitation(ctx context.Context, invitation *models.MatchingInvitation) error {
	db := midacontext.GetDB(ctx)

	MatchingInvitation := dbquery.Use(db).MatchingInvitation
	_, err := MatchingInvitation.WithContext(ctx).Where(MatchingInvitation.ID.Eq(invitation.ID)).UpdateSimple(
		MatchingInvitation.ConfirmState.Value(models.InvitationConfirmStateRejected.String()),
		MatchingInvitation.ConfirmedAt.Value(time.Now()),
		MatchingInvitation.Closed.Value(true),
	)
	if err != nil {
		return err
	}

	MatchingQuota := dbquery.Use(db).MatchingQuota
	if _, err := MatchingQuota.WithContext(ctx).Where(MatchingQuota.UserID.Eq(invitation.UserID)).UpdateSimple(MatchingQuota.Remain.Add(1)); err != nil {
		return err
	}

	midacontext.GetLoader[loader.Loader](ctx).MatchingInvitation.Clear(ctx, invitation.ID)
	return nil
}
