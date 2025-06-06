-- db-seed.sql: Insert fake data for testing

-- Insert 200 fake users with bytea password
INSERT INTO users (username, email, password)
SELECT 'user' || g, 'user' || g || '@example.com', convert_to('password' || g, 'UTF8')
FROM generate_series(1, 200) AS g;

-- Insert 300 fake posts (randomly assigned to valid users)
INSERT INTO posts (title, user_id, content, tags)
SELECT 
  'Post #' || g, 
  u.id, 
  'This is the content of post #' || g, 
  ARRAY['tag' || ((g % 10) + 1), 'tag' || ((g % 20) + 1)]
FROM generate_series(1, 300) AS g
JOIN (SELECT id FROM users ORDER BY random() LIMIT 300) u ON TRUE;

-- Insert 500 fake comments (randomly assigned to valid users and posts)
INSERT INTO comments (content, author_id, post_id, likes)
SELECT 
  'Comment #' || g, 
  u.id, 
  p.id, 
  (random() * 10)::int
FROM generate_series(1, 500) AS g
JOIN (SELECT id FROM users ORDER BY random() LIMIT 500) u ON TRUE
JOIN (SELECT id FROM posts ORDER BY random() LIMIT 500) p ON TRUE;

-- Insert 400 fake followers (random pairs of users, no self-follow, no duplicates)
WITH user_pairs AS (
  SELECT u1.id AS user_id, u2.id AS follower_id
  FROM users u1 CROSS JOIN users u2
  WHERE u1.id <> u2.id
  ORDER BY random()
  LIMIT 400
)
INSERT INTO followers (user_id, follower_id)
SELECT user_id, follower_id FROM user_pairs;
