package modelutil

import "context"

func UserViewedMatching(ctx context.Context, userID string, matchingID string) {
}

type MatchingCandidateOpt struct {
	TopicID string
	Gender  string
}

func GetMatchingCandidate(ctx context.Context, userID string, opt MatchingCandidateOpt) {
}
