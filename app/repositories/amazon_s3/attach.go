package amazon_s3

import (
	"sport4all/app/models"
	"sport4all/app/repositories"
	"sport4all/pkg/errors"
	"sport4all/pkg/logger"

	"github.com/aws/aws-sdk-go/aws"
	"github.com/aws/aws-sdk-go/aws/session"
	"github.com/aws/aws-sdk-go/service/s3"
	"github.com/aws/aws-sdk-go/service/s3/s3manager"
	"github.com/jinzhu/gorm"
	"github.com/labstack/gommon/random"
)

type AttachStore struct {
	db        *gorm.DB
	sessionS3 *session.Session
	bucket    string
}

func CreateAttachRepository(db *gorm.DB, sessS3 *session.Session, bucket string) repositories.AttachRepository {
	return &AttachStore{db: db, sessionS3: sessS3, bucket: bucket}
}

func (s3Store *AttachStore) Create(attach *models.Attach) error {
	attach.Key = random.String(32, random.Alphabetic, random.Numeric)
	manager, err := s3manager.NewUploader(s3Store.sessionS3).Upload(&s3manager.UploadInput{
		Bucket:             aws.String(s3Store.bucket),
		Key:                aws.String(attach.Key),
		Body:               attach.Data,
		ContentDisposition: aws.String("attachment; filename=\"" + attach.Name + "\""),
	})
	if err != nil {
		logger.Error(err)
		return errors.ErrBadFileUploadS3
	}

	attach.URL = manager.Location

	if err = s3Store.db.Create(attach).Error; err != nil {
		logger.Error(err)
		return errors.ErrConflict
	}

	return nil
}

func (s3Store *AttachStore) Delete(key string) error {
	if err := s3Store.db.Where("key = ?", key).Delete(&models.Attach{}).Error; err != nil {
		logger.Error(err)
		return errors.ErrFileNotFound
	}

	if _, err := s3.New(s3Store.sessionS3).DeleteObject(&s3.DeleteObjectInput{
		Bucket: aws.String(s3Store.bucket),
		Key:    aws.String(key),
	}); err != nil {
		logger.Error(err)
		return errors.ErrBadFileDeleteS3
	}

	return nil
}
