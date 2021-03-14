package models

type MeetingStatus uint

const (
	Unknown MeetingStatus = iota
	New
	Progress
	Finished
)

// swagger:model Game
type Meeting struct {
	// example: 101
	ID uint `json:"id" gorm:"primary_key"`

	Status MeetingStatus `json:"status"`

	Stats string `json:"stats"` // TODO: move to table (Anton)

	TournamentId uint `json:"tournamentId"`

	PrevGame *Meeting `json:"prevGame,omitempty" faker:"-"`

	NextGame *Meeting `json:"nextGame,omitempty" faker:"-"`

	Teams []Team `json:"teams,omitempty" gorm:"many2many:team_meetings;" faker:"-"`
}

// swagger:model Teams
type Meetings []Meeting

func (meeting *Meeting) TableName() string {
	return "meeting"
}
