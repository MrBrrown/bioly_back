CREATE DATABASE profile_db;

\connect profile_db

CREATE TABLE IF NOT EXISTS user_page (
  id          BIGSERIAL PRIMARY KEY,
  user_id     BIGINT      NOT NULL,
  page        JSONB       NOT NULL,
  created_at  TIMESTAMPTZ NOT NULL DEFAULT now(),
  updated_at  TIMESTAMPTZ NOT NULL DEFAULT now(),
  CONSTRAINT user_page_user_fk
    FOREIGN KEY (user_id) REFERENCES auth.users(id) ON DELETE CASCADE,
  CONSTRAINT user_page_size_limit CHECK (octet_length(page::text) <= 1048576)
);