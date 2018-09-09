package model

import "encoding/json"

type PushEvent struct {
	PushSource string `json:"push_source"`

	SourceId int64 `json:"source_id"` // the id of the source of events. stored in somewhere. 0 means no source id found.

	Branch        string `json:"branch"`
	CommitVersion string `json:"commit_version"`
	CloneSource   string `json:"clone_source"`
}

func (this *PushEvent) FromBytes(data []byte) error {
	return json.Unmarshal(data, this)
}

func (this *PushEvent) ToBytes() ([]byte, error) {
	return json.Marshal(this)
}
