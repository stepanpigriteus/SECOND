CREATE TABLE users (
    id SERIAL PRIMARY KEY,
    avatar_url TEXT NOT NULL,
    character_name TEXT NOT NULL,
    custom_name TEXT,
    created_at TIMESTAMP WITH TIME ZONE DEFAULT CURRENT_TIMESTAMP
);

CREATE TABLE sessions (
    id TEXT PRIMARY KEY,
    user_id INTEGER REFERENCES users(id) ON DELETE CASCADE,
    expires_at TIMESTAMP WITH TIME ZONE NOT NULL,
    created_at TIMESTAMP WITH TIME ZONE DEFAULT CURRENT_TIMESTAMP
);

CREATE TABLE posts (
    id SERIAL PRIMARY KEY,
    title TEXT NOT NULL,
    content TEXT NOT NULL,
    image_url TEXT,
    user_id INTEGER REFERENCES users(id),
    created_at TIMESTAMP WITH TIME ZONE DEFAULT CURRENT_TIMESTAMP,
    updated_at TIMESTAMP WITH TIME ZONE DEFAULT CURRENT_TIMESTAMP,
    deleted_at TIMESTAMP WITH TIME ZONE DEFAULT NULL,
    last_comment_at TIMESTAMP WITH TIME ZONE DEFAULT CURRENT_TIMESTAMP
);

CREATE TABLE comments (
    id SERIAL PRIMARY KEY,
    post_id INTEGER REFERENCES posts(id) ON DELETE CASCADE,
    parent_id INTEGER REFERENCES comments(id) DEFAULT NULL,
    content TEXT NOT NULL,
    user_id INTEGER REFERENCES users(id),
    created_at TIMESTAMP WITH TIME ZONE DEFAULT CURRENT_TIMESTAMP,
    file_url TEXT DEFAULT NULL
);

-- Индексы для оптимизации запросов
CREATE INDEX idx_posts_deleted_at ON posts(deleted_at);
CREATE INDEX idx_posts_last_comment_at ON posts(last_comment_at);
CREATE INDEX idx_comments_post_id ON comments(post_id);
CREATE INDEX idx_comments_parent_id ON comments(parent_id);
CREATE INDEX idx_sessions_expires_at ON sessions(expires_at);



CREATE OR REPLACE FUNCTION update_post_last_comment_time()
RETURNS TRIGGER AS $$
BEGIN
    UPDATE posts
    SET last_comment_at = NEW.created_at,
        updated_at = CURRENT_TIMESTAMP
    WHERE id = NEW.post_id;
    RETURN NEW;
END;
$$ LANGUAGE plpgsql;

-- Создаем триггер, который будет обновлять last_comment_at при добавлении нового комментария
DROP TRIGGER IF EXISTS update_post_comment_time ON comments;
CREATE TRIGGER update_post_comment_time
AFTER INSERT ON comments
FOR EACH ROW
EXECUTE FUNCTION update_post_last_comment_time();

-- Создаем функцию для периодического удаления постов
CREATE OR REPLACE FUNCTION mark_inactive_posts_as_deleted()
RETURNS void AS $$
BEGIN
    -- Отмечаем как удаленные посты без комментариев старше 10 минут
    UPDATE posts
    SET deleted_at = CURRENT_TIMESTAMP
    WHERE id NOT IN (SELECT DISTINCT post_id FROM comments)
      AND deleted_at IS NULL
      AND created_at < (CURRENT_TIMESTAMP - INTERVAL '10 minutes');

    -- Отмечаем как удаленные посты с комментариями, но без активности более 15 минут
    UPDATE posts
    SET deleted_at = CURRENT_TIMESTAMP
    WHERE id IN (SELECT DISTINCT post_id FROM comments)
      AND deleted_at IS NULL
      AND last_comment_at < (CURRENT_TIMESTAMP - INTERVAL '15 minutes');
END;
$$ LANGUAGE plpgsql;