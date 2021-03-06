package models

// swagger:model Stats
type Stats struct {
	// example: 101
	ID uint `json:"id" gorm:"primary_key"`

	// example: 81
	Score uint `json:"score"`

	// more stats
	// ...

	// example: 1234
	Created int64 `json:"created" gorm:"autoCreateTime"`

	// example: 4
	MeetingId uint `json:"meetingId" gorm:"not null;index"`

	// example: 3
	TeamId uint `json:"teamId" gorm:"not null;index"`

	// example: 24
	PlayerId *uint `json:"playerId,omitempty" gorm:"index"`
}

// swagger:model StatsSet
type StatsSet []Stats

func (stats *Stats) TableName() string {
	return "stats"
}
