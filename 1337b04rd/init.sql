
INSERT INTO users (avatar_url, character_name, custom_name, created_at) 
VALUES 
    ('https://example.com/avatar1.jpg', 'John Doe', 'Johnny', '2025-04-24 10:00:00+00'),
    ('https://example.com/avatar2.jpg', 'Jane Smith', 'Janey', '2025-04-23 11:15:00+00'),
    ('https://example.com/avatar3.jpg', 'Sam Brown', 'Sammy', '2025-04-22 12:30:00+00'),
    ('https://example.com/avatar4.jpg', 'Alice Green', 'Ally', '2025-04-21 14:45:00+00'),
    ('https://example.com/avatar5.jpg', 'Bob White', 'Bobby', '2025-04-20 16:00:00+00');

INSERT INTO sessions (id, user_id, expires_at, created_at) 
VALUES 
    ('session_1', 1, '2025-05-01 10:00:00+00', '2025-04-24 10:00:00+00'),
    ('session_2', 2, '2025-05-02 11:15:00+00', '2025-04-23 11:15:00+00'),
    ('session_3', 3, '2025-05-03 12:30:00+00', '2025-04-22 12:30:00+00'),
    ('session_4', 4, '2025-05-04 14:45:00+00', '2025-04-21 14:45:00+00'),
    ('session_5', 5, '2025-05-05 16:00:00+00', '2025-04-20 16:00:00+00');


INSERT INTO posts (title, content, image_url, user_id, created_at, updated_at, deleted_at, last_comment_at) 
VALUES 
    ('First Post', 'This is the content of the first post.', 'https://example.com/image1.jpg', 1, '2025-04-24 10:00:00+00', '2025-04-24 10:00:00+00', NULL, '2025-04-24 10:10:00+00'),
    ('Second Post', 'Content of the second post.', 'https://example.com/image2.jpg', 2, '2025-04-23 11:15:00+00', '2025-04-23 11:15:00+00', NULL, '2025-04-23 11:30:00+00'),
    ('Third Post', 'This is the content for the third post.', 'https://example.com/image3.jpg', 3, '2025-04-22 12:30:00+00', '2025-04-22 12:30:00+00', NULL, '2025-04-22 12:40:00+00'),
    ('Fourth Post', 'Some content for the fourth post.', 'https://example.com/image4.jpg', 4, '2025-04-21 14:45:00+00', '2025-04-21 14:45:00+00', NULL, '2025-04-21 15:00:00+00'),
    ('Fifth Post', 'Content of the fifth post.', 'https://example.com/image5.jpg', 5, '2025-04-20 16:00:00+00', '2025-04-20 16:00:00+00', NULL, '2025-04-20 16:10:00+00');


INSERT INTO comments (post_id, parent_id, content, user_id, created_at) 
VALUES 
    (1, NULL, 'Great post!', 2, '2025-04-24 10:05:00+00'),
    (1, 1, 'Thanks! Glad you liked it.', 1, '2025-04-24 10:06:00+00'),
    (2, NULL, 'Very informative, thanks for sharing!', 3, '2025-04-23 11:20:00+00'),
    (3, NULL, 'Interesting perspective!', 4, '2025-04-22 12:35:00+00'),
    (4, NULL, 'I agree with your thoughts on this!', 5, '2025-04-21 14:50:00+00');
