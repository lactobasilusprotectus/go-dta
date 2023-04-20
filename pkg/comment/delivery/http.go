package delivery

import (
	"fmt"
	"github.com/gin-gonic/gin"
	"github.com/go-playground/validator/v10"
	"go-dts/pkg/domain"
	"go-dts/pkg/util/http"
	"strconv"
)

type CommentHttpHandler struct {
	commentUsecase domain.CommentUseCase
	authMiddleware domain.GinAuthentication
}

func NewCommentHttpHandler(commentUsecase domain.CommentUseCase, authMiddleware domain.GinAuthentication) *CommentHttpHandler {
	return &CommentHttpHandler{
		commentUsecase: commentUsecase,
		authMiddleware: authMiddleware,
	}
}

func (ch *CommentHttpHandler) Register(g *gin.Engine) {
	g.POST("comment", ch.authMiddleware.MustLogin(), ch.InsertComment)
	g.GET("comment", ch.authMiddleware.MustLogin(), ch.FindAllComment)
	g.GET("comment/:id", ch.authMiddleware.MustLogin(), ch.FindCommentByID)
	g.PUT("comment/:id", ch.authMiddleware.MustLogin(), ch.UpdateComment)
	g.DELETE("comment/:id", ch.authMiddleware.MustLogin(), ch.DeleteComment)
}

type reqInsert struct {
	PhotoID int64  `json:"photo_id" validate:"required" binding:"required"`
	Message string `json:"message" validate:"required" binding:"required"`
}

type reqUpdate struct {
	PhotoID int64  `json:"photo_id"`
	Message string `json:"message"`
}

// InsertComment				godoc
//
//	@Summary					InsertComment insert comment to database.
//	@Description				InsertComment insert comment to database.
//	@Produce					json
//	@Tags						Comment
//	@Param						body	body		reqInsert	true	"Commennt Request"
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
//	@Router						/comment [post]
func (ch *CommentHttpHandler) InsertComment(c *gin.Context) {
	var comment reqInsert
	ctx := c.Request.Context()

	if err := c.ShouldBindJSON(&comment); err != nil {
		http.WriteBadRequestResponseWithErrMsg(c, http.ResponseBadRequestError, err)
		return
	}

	if err := validator.New().Struct(&comment); err != nil {
		http.WriteBadRequestResponseWithErrMsg(c, http.ResponseBadRequestError, err)
		return
	}

	userID, okUser := ch.authMiddleware.GetUserIDFromCtx(ctx)

	if !okUser {
		http.WriteUnauthenticatedResponse(c)
		return
	}

	err := ch.commentUsecase.InsertComment(domain.Comment{
		PhotoID: comment.PhotoID,
		Message: comment.Message,
		UserID:  userID,
	})

	if err != nil {
		http.WriteServerErrorResponse(c, http.ResponseServerError, err)
		return
	}

	http.WriteOkResponse(c, "success insert comment")

	return
}

// FindAllComment				godoc
//
//	@Summary					FindAllComment  all comment  from database.
//	@Description				FindAllComment get all comment  from database.
//	@Tags						Comment
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
//	@Router						/comment [get]
func (ch *CommentHttpHandler) FindAllComment(c *gin.Context) {
	comments, err := ch.commentUsecase.FindAllComment()

	if err != nil {
		http.WriteBadRequestResponseWithErrMsg(c, http.ResponseBadRequestError, err)
		return
	}

	http.WriteOkResponse(c, comments)

	return
}

// FindCommentByID				godoc
//
//	@Summary					FindCommentByID get comment by ID from database.
//	@Description				FindCommentByID get comment by ID from database.
//	@ID							get-string-by-int
//	@Tags						Comment
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
//	@Router						/comment/{id} [get]
func (ch *CommentHttpHandler) FindCommentByID(c *gin.Context) {
	id := c.Param("id")

	ID, err := strconv.Atoi(id)

	if err != nil {
		http.WriteBadRequestResponseWithErrMsg(c, http.ResponseBadRequestError, err)
		return
	}

	comment, err := ch.commentUsecase.FindCommentByID(int64(ID))

	if err != nil {
		http.WriteBadRequestResponseWithErrMsg(c, http.ResponseBadRequestError, err)
		return
	}

	http.WriteOkResponse(c, comment)

	return
}

// UpdateComment				godoc
//
//	@Summary					UpdateComment update comment.
//	@Description				UpdateComment update comment.
//	@ID							get-string-by-int
//	@Tags						Comment
//	@Param						id		path		int			true	"ID"
//	@Param						body	body		reqUpdate	true	"Comment Request"
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
//	@Router						/comment/{id} [put]
func (ch *CommentHttpHandler) UpdateComment(c *gin.Context) {
	var comment reqUpdate

	id := c.Param("id")

	ID, err := strconv.Atoi(id)

	if err := c.ShouldBindJSON(&comment); err != nil {
		http.WriteBadRequestResponseWithErrMsg(c, http.ResponseBadRequestError, err)
		return
	}

	err = ch.commentUsecase.UpdateComment(domain.Comment{
		PhotoID: comment.PhotoID,
		Message: comment.Message,
	}, int64(ID))

	if err != nil {
		http.WriteBadRequestResponseWithErrMsg(c, http.ResponseBadRequestError, err)
		return
	}

	http.WriteOkResponse(c, fmt.Sprintf("success update comment with id %d", ID))
	return
}

// DeleteComment				godoc
//
//	@Summary					DeleteComment delete comment.
//	@Description				DeleteComment delete comment.
//	@ID							get-string-by-int
//	@Tags						Comment
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
//	@Router						/comment/{id} [delete]
func (ch *CommentHttpHandler) DeleteComment(c *gin.Context) {
	id := c.Param("id")

	ID, err := strconv.Atoi(id)

	if err != nil {
		http.WriteBadRequestResponseWithErrMsg(c, http.ResponseBadRequestError, err)
		return
	}

	err = ch.commentUsecase.DeleteComment(int64(ID))

	if err != nil {
		http.WriteBadRequestResponseWithErrMsg(c, http.ResponseBadRequestError, err)
		return
	}

	http.WriteOkResponse(c, fmt.Sprintf("success delete comment with id %d", ID))
	return
}
