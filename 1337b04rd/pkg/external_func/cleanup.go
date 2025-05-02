package externalfunc

import (
	"database/sql"
	"log"
	"time"
)

func StartCleanupRoutine(db *sql.DB) {
	log.Println("Начало работы функции очистки")
	for {

		log.Println("Запуск цикла очистки постов")

		err := callMarkInactivePosts(db)
		if err != nil {
			log.Println("Ошибка вызова функции:", err)
		} else {
			log.Println("mark_inactive_posts_as_deleted() вызвана")
		}

		time.Sleep(1 * time.Minute)
	}
}

func callMarkInactivePosts(db *sql.DB) error {
	err := db.Ping()
	if err != nil {
		log.Printf("Ошибка соединения с БД: %v", err)
		return err
	}

	_, err = db.Exec("SELECT public.mark_inactive_posts_as_deleted()")
	if err != nil {
		log.Printf("Детальная ошибка вызова функции: %v", err)
		return err
	}

	log.Println("Функция mark_inactive_posts_as_deleted() успешно вызвана")
	return nil
}

func InitPostTriggers(db *sql.DB) error {
	_, err := db.Exec(`
		CREATE OR REPLACE FUNCTION mark_inactive_posts_as_deleted() RETURNS void AS $$
DECLARE
    deleted_without_comments INTEGER := 0;
    deleted_inactive INTEGER := 0;
BEGIN
    -- Посты без комментариев
    WITH deleted AS (
        UPDATE posts
        SET deleted_at = CURRENT_TIMESTAMP
        WHERE id NOT IN (SELECT DISTINCT post_id FROM comments)
          AND deleted_at IS NULL
          AND created_at < (CURRENT_TIMESTAMP - INTERVAL '10 minutes')
        RETURNING id
    )
    SELECT COUNT(*) INTO deleted_without_comments FROM deleted;
    
    -- Посты с неактивными комментариями
    WITH deleted AS (
        UPDATE posts
        SET deleted_at = CURRENT_TIMESTAMP
        WHERE id IN (SELECT DISTINCT post_id FROM comments)
          AND deleted_at IS NULL
          AND last_comment_at < (CURRENT_TIMESTAMP - INTERVAL '15 minutes')
        RETURNING id
    )
    SELECT COUNT(*) INTO deleted_inactive FROM deleted;

    -- Логируем результаты
    RAISE NOTICE 'Удалено постов без комментариев: %, с неактивными комментариями: %', 
                 deleted_without_comments, deleted_inactive;

    -- Дополнительный вывод для отладки
    IF deleted_without_comments = 0 AND deleted_inactive = 0 THEN
        RAISE NOTICE 'Нет постов для удаления по заданным условиям.';
    END IF;
END;
$$ LANGUAGE plpgsql;
	`)
	if err != nil {
		return err
	}

	_, err = db.Exec(`
		CREATE OR REPLACE FUNCTION update_last_comment_at() RETURNS trigger AS $$
		BEGIN
			UPDATE posts
			SET last_comment_at = CURRENT_TIMESTAMP
			WHERE id = NEW.post_id;
			RETURN NEW;
		END;
		$$ LANGUAGE plpgsql;
	`)
	if err != nil {
		return err
	}

	_, err = db.Exec(`
		DO $$
		BEGIN
			IF NOT EXISTS (
				SELECT 1 FROM pg_trigger WHERE tgname = 'trigger_update_last_comment'
			) THEN
				CREATE TRIGGER trigger_update_last_comment
				AFTER INSERT OR UPDATE ON comments
				FOR EACH ROW
				EXECUTE FUNCTION update_last_comment_at();
			END IF;
		END
		$$;
	`)
	return err
}
