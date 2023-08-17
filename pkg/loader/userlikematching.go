package loader

type UserThumbsUpMotions struct {
	UserID    string
	MotionIDs []string
}

func (u *UserThumbsUpMotions) Size() int {
	return len(u.MotionIDs)
}

func (u *UserThumbsUpMotions) DoThumbsUp(motionID string) {
	if u.ThumbsUp(motionID) {
		return
	}
	// 有并发问题，但是不影响稳定性
	u.MotionIDs = insert(u.MotionIDs, motionID)
}

func (u *UserThumbsUpMotions) UnThumbsUp(motionID string) {
	if !u.ThumbsUp(motionID) {
		return
	}
	// 有并发问题，但是不影响稳定性
	u.MotionIDs = remove(u.MotionIDs, motionID)
}

func (u *UserThumbsUpMotions) ThumbsUp(motionID string) bool {
	return searchString(motionID, u.MotionIDs) != -1
}
