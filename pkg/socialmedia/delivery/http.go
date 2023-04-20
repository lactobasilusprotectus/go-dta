package delivery

import (
	"fmt"
	"github.com/gin-gonic/gin"
	"github.com/go-playground/validator/v10"
	"go-dts/pkg/domain"
	"go-dts/pkg/util/http"
	"strconv"
)

type SocialMediaHttpHandler struct {
	socialMediaUseCase domain.SocialMediaUseCase
	authMiddleware     domain.GinAuthentication
}

func NewSocialMediaHttpHandler(socialMediaUseCase domain.SocialMediaUseCase, authMiddleware domain.GinAuthentication) *SocialMediaHttpHandler {
	return &SocialMediaHttpHandler{
		socialMediaUseCase: socialMediaUseCase,
		authMiddleware:     authMiddleware,
	}
}

func (s *SocialMediaHttpHandler) Register(g *gin.Engine) {
	g.POST("socialmedia", s.authMiddleware.MustLogin(), s.InsertSocialMedia)
	g.GET("socialmedia", s.authMiddleware.MustLogin(), s.FindAllSocialMedia)
	g.GET("socialmedia/:id", s.authMiddleware.MustLogin(), s.FindSocialMediaByID)
	g.PUT("socialmedia/:id", s.authMiddleware.MustLogin(), s.UpdateSocialMedia)
	g.DELETE("socialmedia/:id", s.authMiddleware.MustLogin(), s.DeleteSocialMedia)
}

type reqInsert struct {
	Name string `json:"name" binding:"required" validate:"required"`
	URL  string `json:"url" binding:"required" validate:"required"`
}

type reqUpdate struct {
	Name string `json:"name"`
	URL  string `json:"url"`
}

// InsertSocialMedia				godoc
//
//	@Summary					InsertSocialMedia insert socialMedia to database.
//	@Description				InsertSocialMedia insert socialMedia to database.
//	@Produce					json
//	@Tags						SocialMedia
//	@Param						body	body		reqInsert	true	"SocialMedia Request"
//	@Success					200		{object}	http.BaseResponse
//	@Failure					400		{object}	http.BaseResponse
//	@Failure					401		{object}	http.BaseResponse
//	@Failure					500		{object}	http.BaseResponse
//
//	@securityDefinitions.apikey	JWT
//	@in							header
//	@name						Authorization
//	@Security					JWT
//
//	@Router						/socialmedia [post]
func (s *SocialMediaHttpHandler) InsertSocialMedia(c *gin.Context) {
	var request reqInsert
	ctx := c.Request.Context()

	if err := c.ShouldBindJSON(&request); err != nil {
		http.WriteBadRequestResponseWithErrMsg(c, http.ResponseBadRequestError, err)
		return
	}

	if err := validator.New().Struct(&request); err != nil {
		http.WriteBadRequestResponseWithErrMsg(c, http.ResponseBadRequestError, err)
		return
	}

	userID, okUser := s.authMiddleware.GetUserIDFromCtx(ctx)

	if !okUser {
		http.WriteUnauthenticatedResponse(c)
		return
	}

	err := s.socialMediaUseCase.InsertSocialMedia(domain.SocialMedia{
		Name:   request.Name,
		URL:    request.URL,
		UserID: userID,
	})

	if err != nil {
		http.WriteServerErrorResponse(c, http.ResponseServerError, err)
		return
	}

	http.WriteOkResponse(c, "SocialMedia created")

	return
}

// FindAllSocialMedia				godoc
//
//	@Summary					FindAllSocialMedia get all SocialMedia from database.
//	@Description				FindAllSocialMedia get all SocialMedia from database.
//	@Tags						SocialMedia
//	@Success					200	{object}	http.BaseResponse
//	@Failure					400	{object}	http.BaseResponse
//	@Failure					401	{object}	http.BaseResponse
//	@Failure					500	{object}	http.BaseResponse
//
//	@securityDefinitions.apikey	JWT
//	@in							header
//	@name						Authorization
//	@Security					JWT
//
//	@Router						/socialmedia [get]
func (s *SocialMediaHttpHandler) FindAllSocialMedia(c *gin.Context) {
	socialMedia, err := s.socialMediaUseCase.FindAllSocialMedia()

	if err != nil {
		http.WriteServerErrorResponse(c, http.ResponseServerError, err)
		return
	}

	http.WriteOkResponse(c, socialMedia)

	return
}

// FindSocialMediaByID				godoc
//
//	@Summary					FindSocialMediaByID get SocialMedia by ID from database.
//	@Description				FindSocialMediaByID get SocialMedia by ID from database.
//	@Tags						SocialMedia
//	@Param						id	path		int	true	"ID"
//	@Success					200	{object}	http.BaseResponse
//	@Failure					400	{object}	http.BaseResponse
//	@Failure					401	{object}	http.BaseResponse
//	@Failure					500	{object}	http.BaseResponse
//
//	@securityDefinitions.apikey	JWT
//	@in							header
//	@name						Authorization
//	@Security					JWT
//
//	@Router						/socialmedia/{id} [get]
func (s *SocialMediaHttpHandler) FindSocialMediaByID(c *gin.Context) {
	id := c.Param("id")

	ID, err := strconv.Atoi(id)

	if err != nil {
		http.WriteBadRequestResponseWithErrMsg(c, http.ResponseBadRequestError, err)
		return
	}

	socialMedia, err := s.socialMediaUseCase.FindSocialMediaByID(int64(ID))

	if err != nil {
		http.WriteServerErrorResponse(c, http.ResponseServerError, err)
		return
	}

	http.WriteOkResponse(c, socialMedia)

	return
}

// UpdateSocialMedia				godoc
//
//	@Summary					UpdateSocialMedia update SocialMedia by ID.
//	@Description				UpdateSocialMedia update SocialMedia by ID.
//	@Tags						SocialMedia
//	@Param						id		path		int			true	"ID"
//	@Param						body	body		reqUpdate	true	"socialmedia Request"
//	@Success					200		{object}	http.BaseResponse
//	@Failure					400		{object}	http.BaseResponse
//	@Failure					401		{object}	http.BaseResponse
//	@Failure					500		{object}	http.BaseResponse
//
//	@securityDefinitions.apikey	JWT
//	@in							header
//	@name						Authorization
//	@Security					JWT
//
//	@Router						/socialmedia/{id} [put]
func (s *SocialMediaHttpHandler) UpdateSocialMedia(c *gin.Context) {
	id := c.Param("id")

	ID, err := strconv.Atoi(id)

	if err != nil {
		http.WriteBadRequestResponseWithErrMsg(c, http.ResponseBadRequestError, err)
		return
	}

	var request reqUpdate

	if err = c.ShouldBindJSON(&request); err != nil {
		http.WriteBadRequestResponseWithErrMsg(c, http.ResponseBadRequestError, err)
		return
	}

	err = s.socialMediaUseCase.UpdateSocialMedia(domain.SocialMedia{
		Name: request.Name,
		URL:  request.URL,
	}, int64(ID))

	if err != nil {
		http.WriteServerErrorResponse(c, http.ResponseServerError, err)
		return
	}

	http.WriteOkResponse(c, fmt.Sprintf("SocialMedia with id %d updated", ID))

	return
}

// DeleteSocialMedia				godoc
//
//	@Summary					DeleteSocialMedia SocialMedia photo by ID.
//	@Description				DeleteSocialMedia SocialMedia photo by ID.
//	@Tags						SocialMedia
//	@Param						id	path		int	true	"ID"
//	@Success					200	{object}	http.BaseResponse
//	@Failure					400	{object}	http.BaseResponse
//	@Failure					401	{object}	http.BaseResponse
//	@Failure					500	{object}	http.BaseResponse
//
//	@securityDefinitions.apikey	JWT
//	@in							header
//	@name						Authorization
//	@Security					JWT
//
//	@Router						/socialmedia/{id} [delete]
func (s *SocialMediaHttpHandler) DeleteSocialMedia(c *gin.Context) {
	id := c.Param("id")

	ID, err := strconv.Atoi(id)

	if err != nil {
		http.WriteBadRequestResponseWithErrMsg(c, http.ResponseBadRequestError, err)
		return
	}

	err = s.socialMediaUseCase.DeleteSocialMedia(int64(ID))

	if err != nil {
		http.WriteServerErrorResponse(c, http.ResponseServerError, err)
		return
	}

	http.WriteOkResponse(c, fmt.Sprintf("SocialMedia with id %d deleted", ID))

	return
}
