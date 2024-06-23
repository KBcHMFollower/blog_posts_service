CREATE TABLE IF NOT EXISTS posts
(
    id        UUID PRIMARY KEY,
    user_id     UUID NOT NULL,
    title TEXT NOT NULL,
    text_content TEXT,
    images_content TEXT,
    likes INTEGER DEFAULT 0,
    created_at DATE DEFAULT CURRENT_DATE
);
CREATE INDEX IF NOT EXISTS idx_user_id ON posts (user_id);

CREATE TABLE IF NOT EXISTS comments
(
    id UUID PRIMARY KEY,
    post_id UUID NOT NULL,
    user_id UUID NOT NULL,
    content TEXT NOT NULL,
    likes INTEGER DEFAULT 0,
    created_at DATE DEFAULT  CURRENT_DATE,
    FOREIGN KEY (post_id) REFERENCES  posts(id)
);
CREATE INDEX IF NOT EXISTS  idx_post_id ON comments (post_id);