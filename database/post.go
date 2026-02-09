package database

import "time"

type Post struct {
	ID          int64      `gorm:"primaryKey" json:"id"`
	AuthorID    int64      `gorm:"not null;index" json:"author_id"`
	Title       string     `gorm:"size:200;not null" json:"title"`
	ContentMD   string     `gorm:"type:MEDIUMTEXT;not null" json:"content_md"`
	Status      string     `gorm:"type:enum('draft','published');default:'draft';index" json:"status"`
	CreatedAt   time.Time  `json:"created_at"`
	UpdatedAt   time.Time  `json:"updated_at"`
	PublishedAt *time.Time `json:"published_at,omitempty"`
}

func CreatePost(post *Post) error {
	result := DB.Create(post)
	return result.Error
}
func GetPublishedPostByID(id uint64) (*Post, error) {
	var post Post
	result := DB.Where("id = ? AND status = ?", id, "published").First(&post)
	if result.Error != nil {
		return nil, result.Error
	}
	return &post, nil
}
func ListPublishedPosts(page, limit int) ([]Post, error) {
	var posts []Post
	result := DB.Where("status = ?", "published").Order("published_at DESC").Find(&posts)
	return posts, result.Error
}
func GetPostByID(id uint64) (*Post, error) {
	var post Post
	result := DB.First(&post, id)
	if result.Error != nil {
		return nil, result.Error
	}
	return &post, nil
}
func UpdatePost(updatedPost *Post) error {
	var post Post
	result := DB.First(&post, updatedPost.ID)
	if result.Error != nil {
		return result.Error
	}

	post.Title = updatedPost.Title
	post.ContentMD = updatedPost.ContentMD
	post.Status = updatedPost.Status

	return DB.Save(&post).Error
}
func DeletePostByID(id uint64) error {
	result := DB.Delete(&Post{}, id)
	if result.RowsAffected == 0 {
		return result.Error
	}

	return nil
}
