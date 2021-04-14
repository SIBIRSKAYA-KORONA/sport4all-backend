package models

type InviteState uint

const (
	Opened InviteState = iota
	Rejected
	Accepted
)

type Invite struct {
	ID           uint        `json:"id" gorm:"primary_key"`
	CreatorId    uint        `json:"creator_id" gorm:"not null, index"`
	InvitedId    *uint       `json:"invited_id,omitempty" gorm:"index"`
	User         *User       `json:"user,omitempty" gorm:"-"`
	AssignedId   uint        `json:"assigned_id" gorm:"index"`
	TeamId       uint        `json:"team_id" gorm:"index"`
	Team         *Team       `json:"team,omitempty" gorm:"-"`
	TournamentId *uint       `json:"tournament_id,omitempty" gorm:"index"`
	Tournament   *Tournament `json:"tournament,omitempty" gorm:"-"`
	Type         string      `json:"type"`
	State        InviteState `json:"state"`
}
