package election

type State int

const (
	Follower State = iota
	Candidate
	Leader
)
