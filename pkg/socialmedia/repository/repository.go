package repository

import (
	"fmt"
	commonTime "go-dts/pkg/common/time"
	"go-dts/pkg/domain"
	"go-dts/pkg/util/db"
)

type SocialMediaRepository struct {
	dbClient *db.DatabaseConnection
	time     commonTime.TimeInterface
}

func NewSocialMediaRepository(dbClient *db.DatabaseConnection, time commonTime.TimeInterface) *SocialMediaRepository {
	return &SocialMediaRepository{
		dbClient: dbClient,
		time:     time,
	}
}

func (s *SocialMediaRepository) InsertSocialMedia(socialMedia domain.SocialMedia) (err error) {
	result := s.dbClient.Master.Create(&socialMedia)

	if result.Error != nil {
		return result.Error
	}

	if result.RowsAffected == 0 {
		return fmt.Errorf("no row affected")
	}

	return nil
}

func (s *SocialMediaRepository) FindSocialMediaByUserID(UserID int64) (domain.SocialMedia, error) {
	var socialMedia domain.SocialMedia

	result := s.dbClient.Slave.Where("user_id = ?", UserID).First(&socialMedia)

	if result.Error != nil {
		return domain.SocialMedia{}, result.Error
	}

	return socialMedia, nil
}

func (s *SocialMediaRepository) FindAllSocialMedia() ([]domain.SocialMedia, error) {
	var socialMedias []domain.SocialMedia

	result := s.dbClient.Slave.Find(&socialMedias)

	if result.Error != nil {
		return []domain.SocialMedia{}, result.Error
	}

	return socialMedias, nil
}

func (s *SocialMediaRepository) FindSocialMediaByID(id int64) (domain.SocialMedia, error) {
	var socialMedia domain.SocialMedia

	result := s.dbClient.Slave.Where("id = ?", id).First(&socialMedia)

	if result.Error != nil {
		return domain.SocialMedia{}, result.Error
	}

	return socialMedia, nil
}

func (s *SocialMediaRepository) UpdateSocialMedia(socialMedia domain.SocialMedia, id int64) (err error) {
	result := s.dbClient.Master.Model(&socialMedia).Where("id = ?", id).Updates(socialMedia)

	if result.Error != nil {
		return result.Error
	}

	if result.RowsAffected == 0 {
		return fmt.Errorf("no row affected")
	}

	return nil
}

func (s *SocialMediaRepository) DeleteSocialMedia(id int64) (err error) {
	result := s.dbClient.Master.Where("id = ?", id).Delete(&domain.SocialMedia{})

	if result.Error != nil {
		return result.Error
	}

	if result.RowsAffected == 0 {
		return fmt.Errorf("no row affected")
	}

	return nil
}
