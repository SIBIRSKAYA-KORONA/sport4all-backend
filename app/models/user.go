package models

// swagger:model User
type User struct {

	// required: true
	// example: 10
	ID uint `json:"id" gorm:"primary_key"`

	// example: Тимофей
	Name string `json:"name" gorm:"not null" faker:"name"`

	// example: Разумов
	Surname string `json:"surname" gorm:"not null" faker:"last_name"`

	// example: Спамер
	Nickname string `json:"nickname" gorm:"not null"`

	// example: 18
	Age uint `json:"age" query:"age"`

	// example: tima.razumov@gmail.com
	Email string `json:"email" gorm:"not null" faker:"email"`
}

// swagger:model Users
type Users []User

func (usr *User) TableName() string {
	return "users"
}
