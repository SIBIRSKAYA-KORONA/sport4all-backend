package redis

import (
	"github.com/gomodule/redigo/redis"
	"github.com/google/uuid"

	"sport4all/app/models"
	"sport4all/app/repositories"
	"sport4all/pkg/common"
	"sport4all/pkg/errors"
	"sport4all/pkg/logger"
)

type SessionStore struct {
	DB            *redis.Pool
	ExpiresKeySec uint
}

func CreateSessionRepository(db *redis.Pool, expKeySec uint) repositories.SessionRepository {
	return &SessionStore{DB: db, ExpiresKeySec: expKeySec}
}

func (sessionStore *SessionStore) Create(session *models.Session) error {
	conn := sessionStore.DB.Get()
	defer common.Close(conn.Close)

	session.SID = uuid.New().String()
	session.ExpiresSec = sessionStore.ExpiresKeySec
	_, err := conn.Do("SETEX", session.SID, session.ExpiresSec, session.ID)
	if err != nil {
		logger.Error(err)
		return errors.ErrSessionNotFound
	}

	return nil
}

func (sessionStore *SessionStore) Get(sid string) (uint, error) {
	conn := sessionStore.DB.Get()
	defer common.Close(conn.Close)

	uid, err := redis.Uint64(conn.Do("GET", sid))
	if err != nil {
		logger.Error(err)
		return 0, errors.ErrSessionNotFound
	}

	return uint(uid), nil
}

func (sessionStore *SessionStore) Delete(sid string) error {
	conn := sessionStore.DB.Get()
	defer common.Close(conn.Close)

	if _, err := conn.Do("DEL", sid); err != nil {
		logger.Error(err)
		return errors.ErrSessionNotFound
	}

	return nil
}
