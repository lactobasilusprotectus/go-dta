package domain

type SocialMedia struct {
	ID     int64  `json:"id" gorm:"primaryKey"`
	Name   string `json:"name" gorm:"not null"`
	URL    string `json:"url" gorm:"not null"`
	UserID int64  `json:"user_id"`
}

type SocialMediaUseCase interface {
	InsertSocialMedia(socialMedia SocialMedia) (err error)
	FindSocialMediaByUserID(UserID int64) (SocialMedia, error)
	FindAllSocialMedia() ([]SocialMedia, error)
	FindSocialMediaByID(id int64) (SocialMedia, error)
	UpdateSocialMedia(socialMedia SocialMedia, id int64) (err error)
	DeleteSocialMedia(id int64) (err error)
}

type SocialMediaRepository interface {
	InsertSocialMedia(socialMedia SocialMedia) (err error)
	FindSocialMediaByUserID(UserID int64) (SocialMedia, error)
	FindAllSocialMedia() ([]SocialMedia, error)
	FindSocialMediaByID(id int64) (SocialMedia, error)
	UpdateSocialMedia(socialMedia SocialMedia, id int64) (err error)
	DeleteSocialMedia(id int64) (err error)
}
