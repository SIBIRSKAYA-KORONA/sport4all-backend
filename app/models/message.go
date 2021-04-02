package models

type Message struct {
	MessageType string `json:"type"`
	SourceUid   uint   `json:"source_uid"`
	TargetUid   uint   `json:"target_uid"`
	Tid         uint   `json:"tid"`
	Mid         uint   `json:"mid"`
}
