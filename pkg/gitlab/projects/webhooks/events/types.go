package events

//type Commit struct {
//	ID      string `json:"id"`
//	Message string `json:"message"`
//}

//type Push struct {
//	ObjectKind  string `json:"object_kind"`
//	EventName   string `json:"event_name"`
//	Ref         string `json:"ref"`
//	CheckoutSHA string `json:"checkout_sha"`
//	UserName    string `json:"user_name"`
//	Project     `json:"project"`
//	Commits     []Commit `json:"commits"`
//}

type TagPush struct {
	ObjectKind  string `json:"object_kind"`
	Ref         string `json:"ref"`
	CheckoutSHA string `json:"checkout_sha"`
	Message     string `json:"message"`
	UserName    string `json:"user_name"`
	Project     `json:"project"`
}

const (
	PipelineSuccess = "success"
	PipelineFailed  = "failed"
)

type Pipeline struct {
	ObjectKind       string `json:"object_kind"`
	ObjectAttributes `json:"object_attributes"`
	User             `json:"user"`
	Project          `json:"project"`
	//Commit           `json:"commit"`
}

type ObjectAttributes struct {
	ID        int        `json:"id"`
	Ref       string     `json:"ref"`
	Tag       bool       `json:"tag"`
	SHA       string     `json:"sha"`
	Status    string     `json:"status"`
	Variables []Variable `json:"variables"`
}

type Variable struct {
	Key   string `json:"key"`
	Value string `json:"value"`
}

type User struct {
	Username string `json:"username"`
}

type Project struct {
	ID   int    `json:"id"`
	Name string `json:"name"`
}

//func NewPush() *Push {
//	return &Push{}
//}

func NewTagPush() *TagPush {
	return &TagPush{}
}

func NewPipeline() *Pipeline {
	return &Pipeline{}
}

//func (p *Push) GetBranch() string {
//	s := strings.Split(p.Ref, "/")
//	return s[len(s)-1]
//}
