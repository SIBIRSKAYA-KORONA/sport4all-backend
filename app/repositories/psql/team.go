package psql

import (
	"strings"
	"time"

	"github.com/jinzhu/gorm"

	"sport4all/app/models"
	"sport4all/app/repositories"
	"sport4all/pkg/errors"
	"sport4all/pkg/logger"
)

type TeamStore struct {
	db *gorm.DB
}

func CreateTeamRepository(db *gorm.DB) repositories.TeamRepository {
	return &TeamStore{db: db}
}

func (teamStore *TeamStore) Create(team *models.Team) error {
	team.Created = time.Now().Unix()
	if err := teamStore.db.Create(team).Error; err != nil {
		logger.Error(err)
		return errors.ErrConflict
	}

	return nil
}

func (teamStore *TeamStore) GetByID(tid uint) (*models.Team, error) {
	team := new(models.Team)
	if err := teamStore.db.Where("id = ?", tid).
		Preload("Avatar").
		Preload("Players").
		First(&team).Error; err != nil {
		logger.Error(err)
		return nil, errors.ErrTeamNotFound
	}

	return team, nil
}

func (teamStore *TeamStore) IsTeamOwner(teamID uint, userID uint) (bool, error) {
	team := new(models.Team)
	if err := teamStore.db.Where("id = ?", teamID).First(&team).Error; err != nil {
		logger.Error(err)
		return false, errors.ErrTeamNotFound
	}

	return team.OwnerId == userID, nil
}

func (teamStore *TeamStore) IsTeamPlayer(teamID uint, userID uint) (bool, error) {
	players := new(models.Users)

	if err := teamStore.db.Model(models.Team{ID: teamID}).Select("id").Related(&players, "players").Error; err != nil {
		return false, errors.ErrUserNotFound
	}

	for _, player := range *players {
		if player.ID == userID {
			return true, nil
		}
	}

	return false, nil
}

//func (teamStore *TeamStore)

func (teamStore *TeamStore) GetTeamsByUser(uid uint, role models.Role) (*models.Teams, error) {
	var userTeams models.Teams

	foreignKey := ""
	switch role {
	case models.Player:
		foreignKey = "teamPlayer"
	case models.Owner:
		foreignKey = "ownerId"
	default:
		return nil, errors.ErrInternal
	}

	if err := teamStore.db.Model(&models.User{ID: uid}).
		Preload("Avatar").
		Related(&userTeams, foreignKey).Error; err != nil {
		logger.Error(err)
		return nil, errors.ErrTeamNotFound
	}

	for i := range userTeams {
		userTeams[i].Players = nil
	}

	return &userTeams, nil
}

func (teamStore *TeamStore) GetAllTournaments(tid uint) (*models.Tournaments, error) {
	var tournamentTeams models.Tournaments
	if err := teamStore.db.Model(&models.Team{ID: tid}).
		Preload("Avatar").
		Related(&tournamentTeams, "tournaments").Error; err != nil {
		logger.Error(err)
		return nil, errors.ErrTeamNotFound
	}

	return &tournamentTeams, nil
}

func (teamStore *TeamStore) GetTeamsByNamePart(namePart string, limit uint) (*models.Teams, error) {
	teams := new(models.Teams)
	if err := teamStore.db.
		Order("name").Limit(limit).
		Where("LOWER(name) LIKE ?", "%"+strings.ToLower(namePart)+"%").
		Preload("Avatar").
		Find(&teams).Error; err != nil {
		logger.Error(err)
		return nil, errors.ErrTeamNotFound
	}

	return teams, nil
}

func (teamStore *TeamStore) InviteMember(tid uint, user *models.User, role models.Role) error {
	team := new(models.Team)
	if err := teamStore.db.First(team, tid).Error; err != nil {
		logger.Error(err)
		return errors.ErrTeamNotFound
	}

	// TODO: обработать значение role (Антон)
	if err := teamStore.db.Model(&team).Association("Players").Append(user).Error; err != nil {
		logger.Warn(err)
	}

	return nil
}

func (teamStore *TeamStore) DeleteMember(tid uint, uid uint) error {
	team := new(models.Team)
	if err := teamStore.db.First(team, tid).Error; err != nil {
		logger.Error(err)
		return errors.ErrTeamNotFound
	}

	check, err := teamStore.IsTeamOwner(tid, uid)
	if err != nil {
		logger.Error(err)
		return err
	}
	if check {
		return errors.ErrNoPermission
	}

	if err := teamStore.db.Model(&team).Association("Players").Delete(models.User{ID: uid}).Error; err != nil {
		logger.Error(err)
		return errors.ErrTeamNotFound
	}

	return nil
}

func (teamStore *TeamStore) GetUsersForInvite(tid uint, nicknamePart string, limit uint) (*models.Users, error) {
	users := new(models.Users)
	team, err := teamStore.GetByID(tid)
	if err != nil {
		logger.Error(err)
		return nil, errors.ErrTeamNotFound
	}

	teamOwnerAndPlayersIDs := []uint{team.OwnerId}
	for _, player := range team.Players {
		teamOwnerAndPlayersIDs = append(teamOwnerAndPlayersIDs, player.ID)
	}

	if err = teamStore.db.Select("id, name, surname, nickname").
		Limit(limit).
		Where("LOWER(nickname) LIKE ?", "%"+strings.ToLower(nicknamePart)+"%").
		Not("id", teamOwnerAndPlayersIDs).
		Preload("Avatar").
		Find(&users).Error; err != nil {
		logger.Error(err)
		return nil, errors.ErrUserNotFound
	}

	return users, nil
}

func (teamStore *TeamStore) GetTeamStats(tid uint) ([]models.Stats, error) {
	var stats []models.Stats
	if err := teamStore.db.Model(&models.Team{ID: tid}).
		Related(&stats, "teamId").Error; err != nil {
		logger.Error(err)
		return nil, errors.ErrTeamNotFound
	}

	return stats, nil
}
