package repositories

import (
	"sport4all/app/models"
)

type SkillRepository interface {
	Create(skill *models.Skill) error
	GetByNamePart(namePart string, limit uint) (*[]models.Skill, error)
	//Update(sid string) (uint, error)
	//Delete(sid string) error
}
