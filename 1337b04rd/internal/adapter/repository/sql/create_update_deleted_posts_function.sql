CREATE OR REPLACE FUNCTION update_deleted_posts() 
RETURNS void AS $$
BEGIN
    -- Удаление постов без комментариев, если прошло 10 минут с их создания
    UPDATE posts
    SET deleted_at = CURRENT_TIMESTAMP
    WHERE deleted_at IS NULL
      AND id NOT IN (SELECT post_id FROM comments)
      AND created_at < CURRENT_TIMESTAMP - INTERVAL '10 minutes';

    -- Удаление постов с комментариями, если последний комментарий был более 15 минут назад
    UPDATE posts
    SET deleted_at = CURRENT_TIMESTAMP
    WHERE deleted_at IS NULL
      AND id IN (SELECT DISTINCT post_id FROM comments)
      AND last_comment_at < CURRENT_TIMESTAMP - INTERVAL '15 minutes';
END;
$$ LANGUAGE plpgsql;
