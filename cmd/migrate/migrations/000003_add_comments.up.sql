CREATE TABLE IF NOT EXISTS comments (
  id bigserial PRIMARY KEY,
  content text NOT NULL,
  author_id INT NOT NULL,
  post_id INT NOT NULL,
  version INT NOT NULL DEFAULT 0,
  updated_at timestamp(0) with time zone NOT NULL DEFAULT NOW(),
  created_at timestamp(0) with time zone NOT NULL DEFAULT NOW(),
  CONSTRAINT fk_user FOREIGN KEY (author_id) REFERENCES users (id) ON DELETE CASCADE,
  CONSTRAINT fk_post FOREIGN KEY (post_id) REFERENCES posts (id) ON DELETE CASCADE
);
