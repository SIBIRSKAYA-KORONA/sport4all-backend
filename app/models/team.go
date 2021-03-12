package models

// swagger:model Team
type Team struct {
	// example: 101
	ID uint `json:"id" gorm:"primary_key"`

	// example: 4
	OwnerId uint `json:"owner_id" gorm:"not null"`

	Players []User `json:"players,omitempty" gorm:"many2many:team_players;" faker:"-"`

	// example: ЦСКА
	Name string `json:"name" gorm:"unique;index" faker:"name"`

	// example: Moscow
	Location string `json:"location" gorm:"index"`

	// example: super_ava.jpg
	LinkOnAvatar string `json:"link_on_avatar"`

	// example: 1234
	Created int64 `json:"created" gorm:"autoCreateTime"`

	// example: Один из ведущих футбольных клубов Москвы
	About string `json:"about"`
}

// swagger:model Teams
type Teams []Team

func (team *Team) TableName() string {
	return "teams"
}
