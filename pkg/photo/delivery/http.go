package delivery

import (
	"fmt"
	"github.com/gin-gonic/gin"
	"github.com/go-playground/validator/v10"
	"go-dts/pkg/domain"
	"go-dts/pkg/util/http"
	"strconv"
)

type PhotoHttpHandler struct {
	photoUseCase   domain.PhotoUseCase
	authMiddleware domain.GinAuthentication
}

func NewPhotoHttpHandler(photoUseCase domain.PhotoUseCase, authMiddleware domain.GinAuthentication) *PhotoHttpHandler {
	return &PhotoHttpHandler{
		photoUseCase:   photoUseCase,
		authMiddleware: authMiddleware,
	}
}

type reqInsert struct {
	Title    string `json:"title" validate:"required"`
	Caption  string `json:"caption"`
	PhotoURL string `json:"photo_url" validate:"required"`
}

type reqUpdate struct {
	Title    string `json:"title"`
	Caption  string `json:"caption"`
	PhotoURL string `json:"photo_url"`
}

func (p *PhotoHttpHandler) Register(g *gin.Engine) {
	g.POST("photo", p.authMiddleware.MustLogin(), p.InsertPhoto)
	g.GET("photo", p.authMiddleware.MustLogin(), p.FindAllPhoto)
	g.GET("photo/:id", p.authMiddleware.MustLogin(), p.FindPhotoByID)
	g.PUT("photo/:id", p.authMiddleware.MustLogin(), p.UpdatePhoto)
	g.DELETE("photo/:id", p.authMiddleware.MustLogin(), p.DeletePhoto)
}

// InsertPhoto				godoc
//
//	@Summary					InsertPhoto insert photo to database.
//	@Description				InsertPhoto insert photo to database.
//	@Produce					json
//	@Tags						Photo
//	@Param						body	body		reqInsert	true	"Photo Request"
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
//	@Router						/photo [post]
func (p *PhotoHttpHandler) InsertPhoto(c *gin.Context) {
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

	userID, okUser := p.authMiddleware.GetUserIDFromCtx(ctx)

	if !okUser {
		http.WriteUnauthenticatedResponse(c)
		return
	}

	res, err := p.photoUseCase.InsertPhoto(domain.Photo{
		Title:    request.Title,
		Caption:  request.Caption,
		PhotoURL: request.PhotoURL,
		UserID:   userID,
	})

	if err != nil {
		http.WriteServerErrorResponse(c, http.ResponseServerError, err)
		return
	}

	http.WriteOkResponse(c, fmt.Sprintf("Photo with id %d created", res.ID))

	return

}

// FindAllPhoto				godoc
//
//	@Summary					FindAllPhoto get all photo from database.
//	@Description				FindAllPhoto get all photo from database.
//	@Tags						Photo
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
//	@Router						/photo [get]
func (p *PhotoHttpHandler) FindAllPhoto(c *gin.Context) {
	allPhoto, err := p.photoUseCase.FindAllPhoto()

	if err != nil {
		http.WriteServerErrorResponse(c, http.ResponseServerError, err)
		return
	}

	http.WriteOkResponse(c, allPhoto)

	return
}

// FindPhotoByID				godoc
//
//	@Summary					FindPhotoByID get photo by ID from database.
//	@Description				FindPhotoByID get photo by ID from database.
//	@Tags						Photo
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
//	@Router						/photo/{id} [get]
func (p *PhotoHttpHandler) FindPhotoByID(c *gin.Context) {
	id := c.Param("id")

	photoID, err := strconv.Atoi(id)

	if err != nil {
		http.WriteBadRequestResponseWithErrMsg(c, http.ResponseBadRequestError, err)
		return
	}

	photo, err := p.photoUseCase.FindPhotoByID(int64(photoID))

	if err != nil {
		http.WriteServerErrorResponse(c, http.ResponseServerError, err)
		return
	}

	http.WriteOkResponse(c, photo)

	return
}

// UpdatePhoto				godoc
//
//	@Summary					UpdatePhoto update photo by ID.
//	@Description				UpdatePhoto update photo by ID.
//	@Tags						Photo
//	@Param						id		path		int			true	"ID"
//	@Param						body	body		reqUpdate	true	"Photo Request"
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
//	@Router						/photo/{id} [put]
func (p *PhotoHttpHandler) UpdatePhoto(c *gin.Context) {
	var request reqUpdate

	id := c.Param("id")

	photoID, err := strconv.Atoi(id)

	if err != nil {
		http.WriteBadRequestResponseWithErrMsg(c, http.ResponseBadRequestError, err)
		return
	}

	if err = c.ShouldBindJSON(&request); err != nil {
		http.WriteBadRequestResponseWithErrMsg(c, http.ResponseBadRequestError, err)
		return
	}

	if err = validator.New().Struct(&request); err != nil {
		http.WriteBadRequestResponseWithErrMsg(c, http.ResponseBadRequestError, err)
		return
	}

	err = p.photoUseCase.UpdatePhoto(domain.Photo{
		ID:       int64(photoID),
		Title:    request.Title,
		Caption:  request.Caption,
		PhotoURL: request.PhotoURL,
	}, int64(photoID))

	if err != nil {
		http.WriteServerErrorResponse(c, http.ResponseServerError, err)
		return
	}

	http.WriteOkResponse(c, fmt.Sprintf("Photo with id %d updated", photoID))

	return
}

// DeletePhoto				godoc
//
//	@Summary					DeletePhoto delete photo by ID.
//	@Description				DeletePhoto delete photo by ID.
//	@Tags						Photo
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
//	@Router						/photo/{id} [delete]
func (p *PhotoHttpHandler) DeletePhoto(c *gin.Context) {
	id := c.Param("id")

	photoID, err := strconv.Atoi(id)

	if err != nil {
		http.WriteBadRequestResponseWithErrMsg(c, http.ResponseBadRequestError, err)
		return
	}

	err = p.photoUseCase.DeletePhoto(int64(photoID))

	if err != nil {
		http.WriteServerErrorResponse(c, http.ResponseServerError, err)
		return
	}

	http.WriteOkResponse(c, fmt.Sprintf("Photo with id %d deleted", photoID))

	return
}
