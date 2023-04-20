package usecase

import "go-dts/pkg/domain"

type SocialMediaUseCase struct {
	socialMediaRepository domain.SocialMediaRepository
}

func NewSocialMediaUseCase(socialMediaRepository domain.SocialMediaRepository) *SocialMediaUseCase {
	return &SocialMediaUseCase{
		socialMediaRepository: socialMediaRepository,
	}
}

func (s *SocialMediaUseCase) InsertSocialMedia(socialMedia domain.SocialMedia) (err error) {
	return s.socialMediaRepository.InsertSocialMedia(socialMedia)
}

func (s *SocialMediaUseCase) FindSocialMediaByUserID(UserID int64) (domain.SocialMedia, error) {
	return s.socialMediaRepository.FindSocialMediaByUserID(UserID)
}

func (s *SocialMediaUseCase) FindAllSocialMedia() ([]domain.SocialMedia, error) {
	return s.socialMediaRepository.FindAllSocialMedia()
}

func (s *SocialMediaUseCase) FindSocialMediaByID(id int64) (domain.SocialMedia, error) {
	return s.socialMediaRepository.FindSocialMediaByID(id)
}

func (s *SocialMediaUseCase) UpdateSocialMedia(socialMedia domain.SocialMedia, id int64) (err error) {
	return s.socialMediaRepository.UpdateSocialMedia(socialMedia, id)
}

func (s *SocialMediaUseCase) DeleteSocialMedia(id int64) (err error) {
	return s.socialMediaRepository.DeleteSocialMedia(id)
}
