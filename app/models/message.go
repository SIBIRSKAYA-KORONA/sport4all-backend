package models

type MessageType string

const (
	MeetingStatusChanged MessageType = "meeting_status_changed"
	MeetingStarted       MessageType = "meeting_started"
	MeetingFinished      MessageType = "meeting_finished"
	AddedToTeam          MessageType = "added_to_team"
)

type Message struct {
	// example: 101
	ID          uint        `json:"-" gorm:"primary_key"`
	MessageType MessageType `json:"type"`

	SourceUid uint `json:"source_uid"`
	TargetUid uint `json:"target_uid"`
	TeamId    uint `json:"team_id"`
	MeetingId uint `json:"meeting_id"`

	CreateAt int64 `json:"createAt,omitempty" gorm:"not null"`
	IsRead   bool  `json:"isRead,omitempty"`
}
