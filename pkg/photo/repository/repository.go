package repository

import (
	"fmt"
	commonTime "go-dts/pkg/common/time"
	"go-dts/pkg/domain"
	"go-dts/pkg/util/db"
)

type PhotoRepository struct {
	dbClient *db.DatabaseConnection
	time     commonTime.TimeInterface
}

func NewPhotoRepository(dbClient *db.DatabaseConnection, time commonTime.TimeInterface) *PhotoRepository {
	return &PhotoRepository{
		dbClient: dbClient,
		time:     time,
	}
}

func (p *PhotoRepository) InsertPhoto(photo domain.Photo) (result domain.Photo, err error) {
	res := p.dbClient.Master.Save(&photo)

	if res.Error != nil {
		return domain.Photo{}, res.Error
	}

	if res.RowsAffected == 0 {
		return domain.Photo{}, fmt.Errorf("no row affected")
	}

	return photo, nil
}

func (p *PhotoRepository) FindPhotoByUserID(UserID int64) (domain.Photo, error) {
	var photo domain.Photo

	result := p.dbClient.Slave.Where("user_id = ?", UserID).First(&photo)

	if result.Error != nil {
		return domain.Photo{}, result.Error
	}

	return photo, nil
}

func (p *PhotoRepository) FindAllPhoto() ([]domain.Photo, error) {
	var photos []domain.Photo

	result := p.dbClient.Slave.Find(&photos)

	if result.Error != nil {
		return nil, result.Error
	}

	return photos, nil
}

func (p *PhotoRepository) FindPhotoByID(id int64) (domain.Photo, error) {
	var photo domain.Photo

	result := p.dbClient.Slave.Where("id = ?", id).First(&photo)

	if result.Error != nil {
		return domain.Photo{}, result.Error
	}

	return photo, nil
}

func (p *PhotoRepository) UpdatePhoto(photo domain.Photo, id int64) (err error) {
	result := p.dbClient.Master.Model(&photo).Where("id = ?", id).Updates(photo)

	if result.Error != nil {
		return result.Error
	}

	if result.RowsAffected == 0 {
		return fmt.Errorf("no row affected")
	}

	return nil
}

func (p *PhotoRepository) DeletePhoto(id int64) (err error) {
	result := p.dbClient.Master.Where("id = ?", id).Delete(&domain.Photo{})

	if result.Error != nil {
		return result.Error
	}

	if result.RowsAffected == 0 {
		return fmt.Errorf("no row affected")
	}

	return nil
}
