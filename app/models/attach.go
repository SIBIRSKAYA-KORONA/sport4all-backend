package models

import (
	"mime/multipart"
)

type Attach struct {
	ID           uint           `json:"id" gorm:"primary_key"`
	URL          string         `json:"url" gorm:"not null"`
	Name         string         `json:"filename" gorm:"not null"`
	Key          string         `json:"key" gorm:"not null;index"`
	Data         multipart.File `json:"-" gorm:"-"`
	UserId       *uint          `json:"userId,omitempty" gorm:"index"`
	TeamId       *uint          `json:"teamId,omitempty" gorm:"index"`
	TournamentId *uint          `json:"tournamentId,omitempty" gorm:"index"`
	MeetingId    *uint          `json:"meetingId,omitempty" gorm:"index"`
}

type Attachments []Attach

func (attach *Attach) TableName() string {
	return "attachments"
}
