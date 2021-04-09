package models

// swagger:model Sport
type Sport struct {
	// example: Football
	Name string `json:"name" gorm:"primary_key"`

	About string `json:"about"`

	// example: https://someurl
	Avatar Attach `json:"avatar" gorm:"foreignKey:sportName"`

	Tournament []Tournament `json:"tournament,omitempty" gorm:"foreignKey:sportName"`
}

func (sport *Sport) TableName() string {
	return "sport"
}
