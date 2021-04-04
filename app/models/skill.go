package models

type Skill struct {
	// example: 10
	ID uint `json:"id" gorm:"primary_key"`

	// example: Go
	Name string `json:"name" gorm:"unique;index"`

	Users []User `json:"users" gorm:"many2many:user_skills;"`

	Approvals []SkillApprove `json:"approvals" gorm:"foreignKey:skillId;"`
}

func (skill *Skill) TableName() string {
	return "skills"
}

type SkillApprove struct {
	ID uint `json:"id" gorm:"primary_key"`

	SkillId uint `json:"skillId" gorm:"index"`

	FromUid uint `json:"fromUid"`

	ToUid uint `json:"-"`

	ToUser *User `json:"toUser" gorm:"-"`

	CreateAt int64 `json:"createAt"`
}

func (approve *SkillApprove) TableName() string {
	return "skill_approvals"
}
