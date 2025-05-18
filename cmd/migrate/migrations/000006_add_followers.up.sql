CREATE TABLE IF NOT EXISTS followers(
  user_id bigserial NOT NULL,
  follower_id bigserial NOT NULL, 
  PRIMARY KEY (user_id, follower_id),
  CONSTRAINT fk_user FOREIGN KEY (user_id) REFERENCES users (id) ON DELETE CASCADE,
  CONSTRAINT fk_follower FOREIGN KEY (follower_id) REFERENCES users (id) ON DELETE CASCADE
)