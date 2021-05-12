package models

type Entity uint
type MessageTrigger uint

const (
	MeetingEntity Entity = iota
	TournamentEntity
	TeamEntity
	UserEntity
)

const (
	EventStatusChanged MessageTrigger = iota
	AddToTeam
	InviteStatusChanged
	SkillApproved
)

var EntityToStr = map[Entity]string{
	MeetingEntity:    "meeting",
	TournamentEntity: "tournament",
	TeamEntity:       "team",
	UserEntity:       "user",
}

var StrToEntity = map[string]Entity{
	"meeting":    MeetingEntity,
	"tournament": TournamentEntity,
	"team":       TeamEntity,
	"user":       UserEntity,
}

var StatusToStr = map[EventStatus]string{
	InProgressEvent: "started",
	FinishedEvent:   "finished",
}

type Message struct {
	// example: 101
	ID         uint   `json:"-" gorm:"primary_key"`
	MessageStr string `json:"type"`

	SourceUid    uint         `json:"source_uid"`
	TargetUid    uint         `json:"target_uid"`
	TeamId       uint         `json:"team_id"`
	MeetingId    uint         `json:"meeting_id"`
	TournamentId uint         `json:"tournament_id"`
	InviteState  *InviteState `json:"invite_state,omitempty"`

	CreateAt int64 `json:"createAt,omitempty" gorm:"not null"`
	IsRead   bool  `json:"isRead,omitempty"`
}
