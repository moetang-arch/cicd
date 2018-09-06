package core

type PushEvent struct {
	DeliveryId string

	Ref          string `json:"ref"` // starts with 'refs/heads/' like 'refs/heads/master'
	BeforeCommit string `json:"before"`
	AfterCommit  string `json:"after"`
	CompareLink  string `json:"compare"`

	Commits []Commit `json:"commits"`

	Repository Repository `json:"repository"`
}

type Commit struct {
	Id        string `json:"id"`
	Message   string `json:"message"`
	Timestamp string `json:"timestamp"` // 2018-09-06T22:14:52+08:00

	Author struct {
		Name     string `json:"name"`
		Email    string `json:"email"`
		Username string `json:"username"`
	} `json:"author"`

	Committer struct {
		Name     string `json:"name"`
		Email    string `json:"email"`
		Username string `json:"username"`
	} `json:"committer"`

	Added    []string `json:"added"`
	Deleted  []string `json:"removed"`
	Modified []string `json:"modified"`
}

type Repository struct {
	FullName string `json:"full_name"`
}
