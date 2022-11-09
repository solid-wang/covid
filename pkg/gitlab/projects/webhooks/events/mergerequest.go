package events

const (
	MergeRequestStateMerged = "merged"
)

type MergeRequestProject struct {
	ID int `json:"id"`
}

type MergeRequest struct {
	ObjectKind                   string `json:"object_kind"`
	MergeRequestProject          `json:"project"`
	MergeRequestUser             `json:"user"`
	MergeRequestObjectAttributes `json:"object_attributes"`
}

type LastCommit struct {
	ID string `json:"id"`
}

type MergeRequestObjectAttributes struct {
	TargetBranch string `json:"target_branch"`
	State        string `json:"state"`
	Description  string `json:"description"`
	LastCommit   `json:"last_commit"`
}

type MergeRequestUser struct {
	Username string `json:"username"`
}
