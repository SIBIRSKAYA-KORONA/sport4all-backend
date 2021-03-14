package psql

import (
	"time"

	"github.com/jinzhu/gorm"

	"sport4all/app/models"
	"sport4all/app/repositories"
	"sport4all/pkg/errors"
	"sport4all/pkg/logger"
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

	err := teamStore.DB.Model(team).Related(&team.Players, "Players").Order("id").Error
	if err != nil {
		logger.Error(err)
		return nil, errors.ErrInternal
	}

	return team, nil
}

func (teamStore *TeamStore) GetTeamsByUser(uid uint, role models.Role) (models.Teams, error) {
	var userTeams []models.Team
	usr := &models.User{ID: uid}

	if role == models.Player {
		err := teamStore.DB.Model(usr).Preload("Players").Related(&userTeams, "Player").Error
		if err != nil {
			logger.Error(err)
			return nil, errors.ErrTeamNotFound
		}
	} else if role == models.Owner {
		err := teamStore.DB.Model(&models.User{ID: uid}).Related(&userTeams, "owner_id").Error
		if err != nil {
			logger.Error(err)
			return nil, errors.ErrTeamNotFound
		}
	}

	for i := range userTeams {
		userTeams[i].Players = nil
	}
	return userTeams, nil
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

func (teamStore *TeamStore) InviteMember(tid uint, user *models.User, role models.Role) error {
	team := new(models.Team)
	err := teamStore.DB.First(team, tid).Error
	if err != nil {
		logger.Error(err)
		return errors.ErrTeamNotFound
	}

	// TODO: обработать значение role
	err = teamStore.DB.Model(&team).Association("Players").Append(user).Error
	if err != nil {
		logger.Error(err)
		return errors.ErrInternal
	}
	return nil
}

func (teamStore *TeamStore) GetUsersForInvite(tid uint, nicknamePart string, limit uint) (models.Users, error) {
	var users []models.User
	team, err := teamStore.GetByID(tid)
	if err != nil {
		logger.Error(err)
		return nil, errors.ErrTeamNotFound
	}
	var teamOwnerAndPlayersIDs []uint
	teamOwnerAndPlayersIDs = append(teamOwnerAndPlayersIDs, team.OwnerId)

	for _, player := range team.Players {
		teamOwnerAndPlayersIDs = append(teamOwnerAndPlayersIDs, player.ID)
	}

	err = teamStore.DB.Select("id, name, surname, nickname, link_on_avatar").
		Limit(limit).
		Where("nickname LIKE ?", nicknamePart+"%").
		Not("id", teamOwnerAndPlayersIDs).
		Find(&users).Error
	if err != nil {
		logger.Error(err)
		return nil, errors.ErrUserNotFound
	}

	return users, nil
}
