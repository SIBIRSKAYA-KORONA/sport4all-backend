package models

// swagger:model Team
type Team struct {
	// example: 101
	ID uint `json:"id" gorm:"primary_key"`

	// example: 4
	OwnerId uint `json:"ownerId" gorm:"not null"`

	Players []User `json:"players,omitempty" gorm:"many2many:team_players;" faker:"-"`

	Tournaments []Tournament `json:"tournaments,omitempty" gorm:"many2many:team_tournament;" faker:"-"`

	Meetings []Meeting `json:"meetings,omitempty" gorm:"many2many:team_meetings;" faker:"-"`

	// example: ЦСКА
	Name string `json:"name" gorm:"index" faker:"name"`

	// example: Moscow
	Location string `json:"location" gorm:"index"`

	// example: super_ava.jpg
	LinkOnAvatar string `json:"linkOnAvatar"`

	// example: 1234
	Created int64 `json:"created" gorm:"autoCreateTime"`

	// example: Один из ведущих футбольных клубов Москвы
	About string `json:"about"`
}

// swagger:model Teams
type Teams []Team

// swagger:model OwnedTeams
type OwnedTeams struct {
	Owned Teams `json:"owned"`
}

// swagger:model PlayingTeams
type PlayingTeams struct {
	Owned Teams `json:"playing"`
}

func (team *Team) TableName() string {
	return "teams"
}
