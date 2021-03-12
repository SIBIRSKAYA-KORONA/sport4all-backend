package psql

import (
	"time"

	"github.com/SIBIRSKAYA-KORONA/sport4all-backend/app/usecases"

	"github.com/jinzhu/gorm"
	_ "github.com/jinzhu/gorm/dialects/postgres"

	"github.com/SIBIRSKAYA-KORONA/sport4all-backend/app/models"
	"github.com/SIBIRSKAYA-KORONA/sport4all-backend/app/repositories"
	"github.com/SIBIRSKAYA-KORONA/sport4all-backend/pkg/errors"
	"github.com/SIBIRSKAYA-KORONA/sport4all-backend/pkg/logger"
)

type TeamStore struct {
	DB *gorm.DB
}

func CreateTeamRepository(db *gorm.DB) repositories.TeamRepository {
	return &TeamStore{DB: db}
}

func (teamStore *TeamStore) Create(team *models.Team) error {
	team.Created = time.Now().Unix()

	if err := teamStore.DB.Create(team).Error; err != nil {
		logger.Error(err)
		return errors.ErrConflict
	}

	return nil
}

func (teamStore *TeamStore) GetByID(tid uint) (*models.Team, error) {
	team := new(models.Team)
	if err := teamStore.DB.Where("id = ?", tid).First(&team).Error; err != nil {
		logger.Error(err)
		return nil, errors.ErrTeamNotFound
	}
	return team, nil
}

func (teamStore *TeamStore) GetTeamsByNamePart(namePart string, limit uint) (models.Teams, error) {
	var teams []models.Team
	err := teamStore.DB.Limit(limit).Where("name LIKE ?", namePart+"%").Find(&teams).Error
	if err != nil {
		logger.Error(err)
		return nil, errors.ErrTeamNotFound
	}
	return teams, nil
}

func (teamStore *TeamStore) InviteMember(tid uint, user *models.User, role usecases.Role) error {
	team := new(models.Team)
	err := teamStore.DB.First(team, tid).Error
	if err != nil {
		logger.Error(err)
		return errors.ErrTeamNotFound
	}

	//
	// обработать значение role
	//

	err = teamStore.DB.Model(&team).Association("Players").Append(user).Error
	if err != nil {
		logger.Error(err)
		return errors.ErrInternal
	}
	return nil
}
