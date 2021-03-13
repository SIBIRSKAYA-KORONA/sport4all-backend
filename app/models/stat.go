package models

// swagger:model MeetingStat
type MeetingStat struct {
	// example: 101
	ID uint `json:"id" gorm:"primary_key"`

	// example: 4
	MeetingId uint `json:"meetingId" gorm:"not null"`

	// example: 1
	LeftTeamId uint `json:"leftTeamId" gorm:"not null"`

	// example: https://someurl
	LeftTeamThumbURL string `json:"leftThumb" gorm:"not null"`

	// example: 0
	LeftTeamScores uint `json:"leftScores"`

	// example: 1
	RightTeamId uint `json:"rightTeamId" gorm:"not null"`

	// example: https://someurl
	RightTeamThumbURL string `json:"rightThumb" gorm:"not null"`

	// example: 0
	RightTeamScores uint `json:"rightScores"`

	// example: 12
	MVP uint `json:"mvpID"`

	// example: 1655344454
	Date int64 `json:"date" gorm:"autoCreateTime"`
}

func (meetingStat *MeetingStat) TableName() string {
	return "meeting_stat"
}
