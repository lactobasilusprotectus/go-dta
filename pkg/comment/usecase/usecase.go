package usecase

import "go-dts/pkg/domain"

type CommentUseCase struct {
	commentRepo domain.CommentRepository
}

func NewCommentUseCase(commentRepo domain.CommentRepository) *CommentUseCase {
	return &CommentUseCase{
		commentRepo: commentRepo,
	}
}

func (c *CommentUseCase) InsertComment(comment domain.Comment) (err error) {
	return c.commentRepo.InsertComment(comment)
}

func (c *CommentUseCase) FindCommentByUserID(UserID int64) (domain.Comment, error) {
	return c.commentRepo.FindCommentByUserID(UserID)
}

func (c *CommentUseCase) FindAllComment() ([]domain.Comment, error) {
	return c.commentRepo.FindAllComment()
}

func (c *CommentUseCase) FindCommentByID(id int64) (domain.Comment, error) {
	return c.commentRepo.FindCommentByID(id)
}

func (c *CommentUseCase) UpdateComment(comment domain.Comment, id int64) (err error) {
	return c.commentRepo.UpdateComment(comment, id)
}

func (c *CommentUseCase) DeleteComment(id int64) (err error) {
	return c.commentRepo.DeleteComment(id)
}
