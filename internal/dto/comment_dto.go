package dto

// CommentRequest struct
type CommentRequest struct {
	AuthorName string `json:"author_name" binding:"required"`
	Content    string `json:"content" binding:"required"`
}

// CommentResponse struct
type CommentResponse struct {
	ID         uint   `json:"id"`
	PostID     uint   `json:"post_id"`
	AuthorName string `json:"author_name"`
	Content    string `json:"content"`
	CreatedAt  string `json:"created_at"`
}
