CREATE OR REPLACE FUNCTION update_last_comment_at() 
RETURNS TRIGGER AS $$
BEGIN
    UPDATE posts
    SET last_comment_at = CURRENT_TIMESTAMP
    WHERE id = NEW.post_id;
    RETURN NEW;
END;
$$ LANGUAGE plpgsql;

-- Создаем триггер на добавление комментария
CREATE TRIGGER update_last_comment_at_trigger
AFTER INSERT ON comments
FOR EACH ROW
EXECUTE FUNCTION update_last_comment_at();
