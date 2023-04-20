package usecase

import "go-dts/pkg/domain"

type PhotoUseCase struct {
	photoRepo domain.PhotoRepository
}

func NewPhotoUseCase(photoRepo domain.PhotoRepository) *PhotoUseCase {
	return &PhotoUseCase{photoRepo: photoRepo}
}

func (p *PhotoUseCase) InsertPhoto(photo domain.Photo) (result domain.Photo, err error) {
	return p.photoRepo.InsertPhoto(photo)
}

func (p *PhotoUseCase) FindPhotoByUserID(UserID int64) (domain.Photo, error) {
	return p.photoRepo.FindPhotoByUserID(UserID)
}

func (p *PhotoUseCase) FindAllPhoto() ([]domain.Photo, error) {
	return p.photoRepo.FindAllPhoto()
}

func (p *PhotoUseCase) FindPhotoByID(id int64) (domain.Photo, error) {
	return p.photoRepo.FindPhotoByID(id)
}

func (p *PhotoUseCase) UpdatePhoto(photo domain.Photo, id int64) (err error) {
	return p.photoRepo.UpdatePhoto(photo, id)
}

func (p *PhotoUseCase) DeletePhoto(id int64) (err error) {
	return p.photoRepo.DeletePhoto(id)
}
