package dto

// PostRequest represents the request body for a post.
type PostRequest struct {
	Title    string `json:"title" binding:"required"`
	Content  string `json:"content" binding:"required"`
	AuthorID uint   `json:"author_id"`
}

// PostResponse represents the response body for a post.
type PostResponse struct {
	ID        uint            `json:"id"`
	Title     string          `json:"title"`
	Content   string          `json:"content"`
	AuthorID  uint            `json:"author_id"`
	Author    *AuthorResponse `json:"author"`
	CreatedAt string          `json:"created_at"`
	UpdatedAt string          `json:"updated_at"`
}
