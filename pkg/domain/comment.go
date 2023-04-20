package domain

type Comment struct {
	ID        int64  `json:"id" gorm:"primaryKey"`
	PhotoID   int64  `json:"photo_id"`
	UserID    int64  `json:"user_id"`
	Message   string `json:"message" gorm:"not null"`
	CreatedAt string `json:"created_at"`
	UpdatedAt string `json:"updated_at"`
}

type CommentUseCase interface {
	InsertComment(comment Comment) (err error)
	FindCommentByUserID(UserID int64) (Comment, error)
	FindAllComment() ([]Comment, error)
	FindCommentByID(id int64) (Comment, error)
	UpdateComment(comment Comment, id int64) (err error)
	DeleteComment(id int64) (err error)
}

type CommentRepository interface {
	InsertComment(comment Comment) (err error)
	FindCommentByUserID(UserID int64) (Comment, error)
	FindAllComment() ([]Comment, error)
	FindCommentByID(id int64) (Comment, error)
	UpdateComment(comment Comment, id int64) (err error)
	DeleteComment(id int64) (err error)
}
