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

	TournamentId uint `json:"tournamentId" gorm:"index"`

	Stats []Stats `json:"stats,omitempty" gorm:"foreignKey:meetingId"`

	NextMeetingID *uint `json:"nextMeetingID,omitempty"`

	PrevMeetings []Meeting `json:"prevMeetings,omitempty" gorm:"foreignKey:nextMeetingID" faker:"-"`

	Teams []Team `json:"teams,omitempty" gorm:"many2many:team_meetings;" faker:"-"`

	Attachments []Attach `json:"attachments,omitempty" gorm:"foreignKey:meetingId" faker:"-"`

	Protocol Attach `json:"protocol" gorm:"foreignKey:meetingProtocolId" faker:"-"`
}

// swagger:model Teams
type Meetings []Meeting

func (meeting *Meeting) TableName() string {
	return "meetings"
}
