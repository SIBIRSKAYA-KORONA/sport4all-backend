package models

type EventStatus uint

const (
	UnknownEvent EventStatus = iota
	NotStartedEvent
	RegistrationEvent
	InProgressEvent
	FinishedEvent
)

// swagger:model Game
type Meeting struct {
	// example: 101
	ID uint `json:"id" gorm:"primary_key"`

	Status EventStatus `json:"status"`

	Round uint `json:"round"`

	Group uint `json:"group"`

	TournamentId uint `json:"tournamentId"`

	Stats string `json:"stats"` // TODO: move to table (Anton)

	PrevMeetings []Meeting `json:"prevMeetings,omitempty" faker:"-"`

	NextMeeting *Meeting `json:"nextMeeting,omitempty" faker:"-"`

	Teams []Team `json:"teams,omitempty" gorm:"many2many:team_meetings;" faker:"-"`
}

// swagger:model Teams
type Meetings []Meeting

func (meeting *Meeting) TableName() string {
	return "meeting"
}
