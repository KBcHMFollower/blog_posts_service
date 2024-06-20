CREATE TABLE IF NOT EXISTS posts
(
    id        UUID PRIMARY KEY,
    user_id     UUID NOT NULL,
    title TEXT NOT NULL,
    text_content TEXT,
    images_content TEXT,
    created_at DATE NOT NULL
);
CREATE INDEX IF NOT EXISTS idx_user_id ON posts (user_id);