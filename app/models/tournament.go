package models

type TournamentSystem uint

const (
	UnknownSystem TournamentSystem = iota
	OlympicSystem
	CircularSystem
)

// swagger:model Tournament
type Tournament struct {
	// example: 10
	ID uint `json:"id" gorm:"primary_key"`

	// example: 4
	OwnerId uint `json:"ownerId" gorm:"not null"`

	// example: Чемпионат мира
	Name string `json:"name" gorm:"index" faker:"name"`

	// example: Moscow
	Location string `json:"location" gorm:"index"`

	// example: 2
	System TournamentSystem `json:"system"`

	// example: 3
	Status EventStatus `json:"status"`

	// example: турнир по игре с котиками
	About string `json:"about"`

	// example: ЦСКА, Зенит
	Teams []Teams `json:"teams,omitempty" gorm:"many2many:team_tournament;" faker:"-"`

	// example: игра1, игра2
	Meetings []Meeting `json:"meetings,omitempty" gorm:"foreignKey:tournamentId" faker:"-"`

	// example: 1234
	Created int64 `json:"created" gorm:"autoCreateTime"`
}

// swagger:model Teams
type Tournaments []Tournament

type UserTournament struct {
	Owner      Tournaments `json:"owner"`
	TeamMember Tournaments `json:"teamMember"`
	TeamOwner  Tournaments `json:"teamOwner"`
}

func (tournament *Tournament) TableName() string {
	return "tournaments"
}
