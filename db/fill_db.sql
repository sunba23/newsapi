BEGIN;

TRUNCATE TABLE news_tags, news, tags RESTART IDENTITY CASCADE;

INSERT INTO tags (name) VALUES 
('golang'), 
('python'), 
('javascript'),
('rust'),
('docker'),
('kubernetes'),
('webdev'),
('datascience'),
('machinelearning'),
('algorithms')
ON CONFLICT (name) DO NOTHING;

INSERT INTO news (title, content, author, created_at) VALUES
('Go 1.21 Released', 'The latest version of Go introduces new features including...', 'Gopher', NOW() - interval '2 days'),
('Python 3.12 Performance Improvements', 'Significant speed boosts reported in the new Python release...', 'PyDev', NOW() - interval '5 days'),
('Rust for Web Development', 'How Rust is becoming a viable alternative for backend web services...', 'Ferris', NOW() - interval '1 week'),
('Docker Best Practices 2023', 'Updated guidelines for containerizing your applications...', 'Container Expert', NOW() - interval '3 days'),
('Machine Learning with Go', 'Exploring ML libraries available for the Go programming language...', 'AI Researcher', NOW() - interval '10 days')
RETURNING id;

INSERT INTO news_tags (news_id, tag_id) VALUES
(1, (SELECT id FROM tags WHERE name = 'golang')),
(1, (SELECT id FROM tags WHERE name = 'webdev'));

INSERT INTO news_tags (news_id, tag_id) VALUES
(2, (SELECT id FROM tags WHERE name = 'python')),
(2, (SELECT id FROM tags WHERE name = 'datascience'));

INSERT INTO news_tags (news_id, tag_id) VALUES
(3, (SELECT id FROM tags WHERE name = 'rust')),
(3, (SELECT id FROM tags WHERE name = 'webdev'));

INSERT INTO news_tags (news_id, tag_id) VALUES
(4, (SELECT id FROM tags WHERE name = 'docker')),
(4, (SELECT id FROM tags WHERE name = 'kubernetes'));

INSERT INTO news_tags (news_id, tag_id) VALUES
(5, (SELECT id FROM tags WHERE name = 'golang')),
(5, (SELECT id FROM tags WHERE name = 'machinelearning')),
(5, (SELECT id FROM tags WHERE name = 'algorithms'));

COMMIT;

SELECT n.id, n.title, string_agg(t.name, ', ') as tags
FROM news n
JOIN news_tags nt ON n.id = nt.news_id
JOIN tags t ON nt.tag_id = t.id
GROUP BY n.id, n.title
ORDER BY n.created_at DESC;
