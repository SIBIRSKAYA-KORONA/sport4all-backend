package models

// swagger:model Tournament
type Tournament struct {
	// example: 10
	ID uint `json:"id" gorm:"primary_key"`

	// example: 4
	OwnerId uint `json:"ownerId" gorm:"not null"`

	// example: Чемпионат мира
	Name string `json:"name" gorm:"unique;index" faker:"name"`

	// example: Moscow
	Location string `json:"location" gorm:"index"`

	// example: Teams
	Teams []Teams `json:"teams,omitempty" gorm:"many2many:team_tournament;" faker:"-"`

	// example: Meetings
	Meetings []Meeting `json:"meetings,omitempty" gorm:"foreignKey:tournamentId" faker:"-"`

	// example: 1234
	Created int64 `json:"created" gorm:"autoCreateTime"`
}

// swagger:model Teams
type Tournaments []Tournament

func (tournament *Tournament) TableName() string {
	return "tournaments"
}
