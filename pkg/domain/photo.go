package domain

type Photo struct {
	ID        int64  `json:"id" gorm:"primaryKey"`
	Title     string `json:"title" gorm:"not null"`
	Caption   string `json:"caption"`
	PhotoURL  string `json:"photo_url" gorm:"not null"`
	UserID    int64  `json:"user_id"`
	CreatedAt string `json:"created_at"`
	UpdatedAt string `json:"updated_at"`
}

type PhotoUseCase interface {
	InsertPhoto(photo Photo) (result Photo, err error)
	FindPhotoByUserID(UserID int64) (Photo, error)
	FindAllPhoto() ([]Photo, error)
	FindPhotoByID(id int64) (Photo, error)
	UpdatePhoto(photo Photo, id int64) (err error)
	DeletePhoto(id int64) (err error)
}

type PhotoRepository interface {
	InsertPhoto(photo Photo) (result Photo, err error)
	FindPhotoByUserID(UserID int64) (Photo, error)
	FindAllPhoto() ([]Photo, error)
	FindPhotoByID(id int64) (Photo, error)
	UpdatePhoto(photo Photo, id int64) (err error)
	DeletePhoto(id int64) (err error)
}
