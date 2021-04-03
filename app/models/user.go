package models

type Role uint

const (
	Player Role = iota
	Owner
)

var StringToRole = map[string]Role{
	"player": Player,
	"owner":  Owner,
}

// swagger:model User
type User struct {

	// example: 10
	ID uint `json:"id" gorm:"primary_key"`

	// example: Тимофей
	Name string `json:"name" faker:"name"`

	// example: Разумов
	Surname string `json:"surname" faker:"last_name"`

	// required: true
	// example: Спамер
	Nickname string `json:"nickname" gorm:"unique;index"`

	// example: tima.razumov@gmail.com
	Email string `json:"email" gorm:"index" faker:"email"`

	// example: 8 888 888 888 888
	PhoneNumber string `json:"phoneNumber" gorm:"index"`

	// example: Moscow
	Location string `json:"location" gorm:"index"`

	// example: 1234
	Created int64 `json:"created" gorm:"autoCreateTime"`

	// example: 06.07.1997
	Birthday string `json:"birthday"`

	// example: 198
	Height uint `json:"height"`

	// example: 80
	Weight uint `json:"weight"`

	// example: кмс по баскетболу. учусь в Бауманке
	About string `json:"about"`

	// example: test123
	Password string `json:"password,omitempty" gorm:"-"`

	HashPassword []byte `json:"-"`

	TournamentOwner []Tournament `json:"-" gorm:"foreignKey:ownerId" faker:"-"`

	TeamOwner []Team `json:"-" gorm:"foreignKey:ownerId" faker:"-"`

	TeamPlayer []Team `json:"-" gorm:"many2many:team_players;" faker:"-"`

	Stats []Stats `json:"stats" gorm:"foreignkey:playerId" faker:"-"`

	Avatar Attach `json:"avatar" gorm:"foreignKey:userId"`

	Skills []Skill `json:"skills,omitempty" gorm:"many2many:user_skills;"`

	SkillsApprovals []SkillApprove `json:"-" gorm:"many2many:user_skill_approvals;"`
}

// swagger:model Users
type Users []User

func (usr *User) TableName() string {
	return "users"
}
