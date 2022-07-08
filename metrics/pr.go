package metrics

import "time"

type PR struct {
	Repo           string
	StartTS        time.Time
	RaiseTS        time.Time
	MergeTS        time.Time
	ApproveTS      time.Time
	NCodeLines     int
	Pipelines      []string
	FirstCommentTS time.Time
}
