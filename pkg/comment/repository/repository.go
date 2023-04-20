package repository

import (
	"fmt"
	commonTime "go-dts/pkg/common/time"
	"go-dts/pkg/domain"
	"go-dts/pkg/util/db"
)

type CommentRepository struct {
	dbClient *db.DatabaseConnection
	time     commonTime.TimeInterface
}

func NewCommentRepository(dbClient *db.DatabaseConnection, time commonTime.TimeInterface) *CommentRepository {
	return &CommentRepository{
		dbClient: dbClient,
		time:     time,
	}
}

func (c *CommentRepository) InsertComment(comment domain.Comment) (err error) {
	result := c.dbClient.Master.Create(&comment)

	if result.Error != nil {
		return result.Error
	}

	if result.RowsAffected == 0 {
		return fmt.Errorf("no row affected")
	}

	return nil
}

func (c *CommentRepository) FindCommentByUserID(UserID int64) (domain.Comment, error) {
	var comment domain.Comment

	result := c.dbClient.Slave.Where("user_id = ?", UserID).First(&comment)

	if result.Error != nil {
		return domain.Comment{}, result.Error
	}

	return comment, nil
}

func (c *CommentRepository) FindAllComment() ([]domain.Comment, error) {
	var comments []domain.Comment

	result := c.dbClient.Slave.Find(&comments)

	if result.Error != nil {
		return nil, result.Error
	}

	return comments, nil
}

func (c *CommentRepository) FindCommentByID(id int64) (domain.Comment, error) {
	var comment domain.Comment

	result := c.dbClient.Slave.Where("id = ?", id).First(&comment)

	if result.Error != nil {
		return domain.Comment{}, result.Error
	}

	return comment, nil
}

func (c *CommentRepository) UpdateComment(comment domain.Comment, id int64) (err error) {
	result := c.dbClient.Master.Model(&comment).Where("id = ?", id).Updates(comment)

	if result.Error != nil {
		return result.Error
	}

	if result.RowsAffected == 0 {
		return fmt.Errorf("no row affected")
	}

	return nil
}

func (c *CommentRepository) DeleteComment(id int64) (err error) {
	result := c.dbClient.Master.Where("id = ?", id).Delete(&domain.Comment{})

	if result.Error != nil {
		return result.Error
	}

	if result.RowsAffected == 0 {
		return fmt.Errorf("no row affected")
	}

	return nil
}
