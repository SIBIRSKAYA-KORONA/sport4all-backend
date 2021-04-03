package models

type MessageEntity uint
type MessageTrigger uint

const (
	MeetingEntity MessageEntity = iota
	TournamentEntity
	TeamEntity
)

const (
	EventStatusChanged MessageTrigger = iota
	AddToTeam
)

var EntityToStr = map[MessageEntity]string{
	MeetingEntity: "meeting",
	TournamentEntity: "tournament",
	TeamEntity: "team",
}

var StatusToStr = map[EventStatus]string{
	InProgressEvent: "started",
	FinishedEvent: "finished",
}

type Message struct {
	// example: 101
	ID          uint        `json:"-" gorm:"primary_key"`
	MessageStr  string `json:"type"`

	SourceUid uint `json:"source_uid"`
	TargetUid uint `json:"target_uid"`
	TeamId    uint `json:"team_id"`
	MeetingId uint `json:"meeting_id"`
	TournamentId uint `json:"tournament_id"`

	CreateAt int64 `json:"createAt,omitempty" gorm:"not null"`
	IsRead   bool  `json:"isRead,omitempty"`
}
