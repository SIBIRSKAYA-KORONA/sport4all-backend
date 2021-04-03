package models

type TournamentSystem uint
type TournamentRole uint

const (
	TournamentOrganizer TournamentRole = iota
	TournamentPlayer
	TournamentMember // обобщающая роль (и организатор, и игроки, и т.д.
)

//var StringToRole = map[string]Role{
//	"player": Player,
//	"owner":  Owner,
//}

// swagger:model Tournament
type Tournament struct {
	// example: 10
	ID uint `json:"id" gorm:"primary_key"`

	// example: 4
	OwnerId uint `json:"ownerId" gorm:"not null;index"`

	// example: Чемпионат мира
	Name string `json:"name" gorm:"index" faker:"name"`

	// example: Moscow
	Location string `json:"location" gorm:"index"`

	// example: olympic
	System string `json:"system"`

	// example: 1
	Status EventStatus `json:"status"`

	// example: турнир по игре с котиками
	About string `json:"about"`

	// example: 1234
	Created int64 `json:"created" gorm:"autoCreateTime"`

	Sport string `json:"sport" gorm:"index"`

	Teams []Teams `json:"teams,omitempty" gorm:"many2many:team_tournament;" faker:"-"`

	Meetings []Meeting `json:"meetings,omitempty" gorm:"foreignKey:tournamentId" faker:"-"`

	Avatar Attach `json:"avatar" gorm:"foreignKey:tournamentId"`
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
