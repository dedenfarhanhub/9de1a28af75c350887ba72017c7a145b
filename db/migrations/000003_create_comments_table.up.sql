CREATE TABLE comments
(
    id          INT AUTO_INCREMENT PRIMARY KEY,
    post_id     INT,
    author_name VARCHAR(255) NOT NULL,
    content     TEXT         NOT NULL,
    created_at  TIMESTAMP DEFAULT CURRENT_TIMESTAMP,
    FOREIGN KEY (post_id) REFERENCES posts (id) ON DELETE CASCADE,
    INDEX       idx_comments_post_id (post_id),
    INDEX       idx_comments_created_at (created_at)
);