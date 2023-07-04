package modelutil

import (
	"context"
	"fmt"
	"time"
	"whale/pkg/dbquery"
	"whale/pkg/gqlient/scream"
	"whale/pkg/loader"
	"whale/pkg/models"
	"whale/pkg/whalecode"

	"github.com/letjoy-club/mida-tool/dbutil"
	"github.com/letjoy-club/mida-tool/graphqlutil"
	"github.com/letjoy-club/mida-tool/midacode"
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

	db := dbutil.GetDB(ctx)
	// 生成匹配和匹配结果
	Matching := dbquery.Use(db).Matching
	MatchingResult := dbquery.Use(db).MatchingResult

	m1 := models.Matching{
		ID:             shortid.NewWithTime("m_", 6),
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
		ID:             shortid.NewWithTime("m_", 6),
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
		CreatedBy:      models.ResultCreatedByInvitation.String(),
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
			MatchingInvitation.MatchingIds.Value(graphqlutil.ElementList[string]([]string{m1.ID, m2.ID})),
			MatchingInvitation.MatchingResultId.Value(result.ID),
			MatchingInvitation.ConfirmState.Value(models.InvitationConfirmStateConfirmed.String()),
		)
	if err != nil {
		return err
	}

	err = CheckMatchingResultAndCreateChatGroup(ctx, &result)
	if err != nil {
		return err
	}

	midacontext.GetLoader[loader.Loader](ctx).MatchingResult.Clear(ctx, result.ID)
	midacontext.GetLoader[loader.Loader](ctx).Matching.Clear(ctx, m1.ID)
	midacontext.GetLoader[loader.Loader](ctx).Matching.Clear(ctx, m2.ID)
	midacontext.GetLoader[loader.Loader](ctx).MatchingInvitation.Clear(ctx, invitation.ID)

	topicName := GetTopicName(ctx, invitation.TopicID)
	if topicName == "" {
		return nil
	}

	thunk := midacontext.GetLoader[loader.Loader](ctx).UserAvatarNickname.Load(ctx, invitation.InviteeID)
	invitee, err := thunk()

	if err != nil {
		fmt.Println("GetUserAvatarNickname err:", err)
		return nil
	}

	_, err = scream.SendUserNotification(ctx, midacontext.GetServices(ctx).Scream,
		scream.UserNotificationKindInvitationaccepted,
		invitation.UserID,
		map[string]interface{}{
			"invitationId": invitation.ID,
			"userId":       invitee.ID,
			"userName":     invitee.Nickname,
			"topicName":    topicName,
		},
	)

	if err != nil {
		fmt.Println("SendUserNotification err:", err)
	}
	return nil
}

func RejectMatchingInvitation(ctx context.Context, invitation *models.MatchingInvitation) error {
	db := dbutil.GetDB(ctx)

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

	thunk := midacontext.GetLoader[loader.Loader](ctx).UserAvatarNickname.Load(ctx, invitation.InviteeID)
	invitee, err := thunk()
	if err == nil {
		topicName := GetTopicName(ctx, invitation.TopicID)
		if topicName != "" {
			_, err = scream.SendUserNotification(ctx,
				midacontext.GetServices(ctx).Scream,
				scream.UserNotificationKindInvitationdenied,
				invitation.UserID,
				map[string]interface{}{
					"userName":     invitee.Nickname,
					"userId":       invitee.ID,
					"invitationId": invitation.ID,
					"topicName":    topicName,
				},
			)
			if err != nil {
				fmt.Println("err", err)
			}
		}
	} else {
		fmt.Println("err", err)
	}

	midacontext.GetLoader[loader.Loader](ctx).MatchingInvitation.Clear(ctx, invitation.ID)
	return nil
}

func CancelMatchingInvitation(ctx context.Context, uid string, invitationID string) error {
	db := dbutil.GetDB(ctx)
	MatchingInvitation := dbquery.Use(db).MatchingInvitation
	matchingInvitation, err := MatchingInvitation.WithContext(ctx).Where(MatchingInvitation.ID.Eq(invitationID)).Take()
	if err != nil {
		return midacode.ItemMayNotFound(err)
	}
	if uid != "" {
		if uid != matchingInvitation.UserID {
			return midacode.ErrNotPermitted
		}
	}

	if matchingInvitation.Closed {
		return nil
	}

	if matchingInvitation.ConfirmState != models.InvitationConfirmStateUnconfirmed.String() {
		return nil
	}

	_, err = MatchingInvitation.WithContext(ctx).Where(MatchingInvitation.ID.Eq(invitationID)).UpdateSimple(
		MatchingInvitation.Closed.Value(true),
		MatchingInvitation.ConfirmState.Value(models.InvitationConfirmStateRejected.String()),
	)
	if err != nil {
		return err
	}
	// 发起者的配对配额+1
	MatchingQuota := dbquery.Use(db).MatchingQuota
	_, err = MatchingQuota.WithContext(ctx).Where(MatchingQuota.UserID.Eq(matchingInvitation.UserID)).UpdateSimple(MatchingQuota.Remain.Add(1))
	if err != nil {
		return err
	}

	midacontext.GetLoader[loader.Loader](ctx).MatchingQuota.Clear(ctx, matchingInvitation.UserID)
	midacontext.GetLoader[loader.Loader](ctx).MatchingInvitation.Clear(ctx, matchingInvitation.ID)

	topicName := GetTopicName(ctx, matchingInvitation.TopicID)
	if topicName == "" {
		return nil
	}

	thunk := midacontext.GetLoader[loader.Loader](ctx).UserAvatarNickname.Load(ctx, matchingInvitation.UserID)
	inviter, err := thunk()
	if err != nil {
		fmt.Println("failed to get inviter", err)
		return nil
	}

	_, err = scream.SendUserNotification(ctx, midacontext.GetServices(ctx).Scream,
		scream.UserNotificationKindInvitationcanceled,
		matchingInvitation.InviteeID,
		map[string]interface{}{
			"invitationId": matchingInvitation.ID,
			"userId":       inviter.ID,
			"userName":     inviter.Nickname,
			"topicName":    topicName,
		},
	)
	if err != nil {
		fmt.Println("failed to send user notification", err)
	}

	return nil
}
