\connect bioly

CREATE TABLE IF NOT EXISTS profiles.user_page (
  id          BIGSERIAL PRIMARY KEY,
  user_id     BIGINT      NOT NULL,
  page        JSONB       NOT NULL,
  created_at  TIMESTAMPTZ NOT NULL DEFAULT now(),
  updated_at  TIMESTAMPTZ NOT NULL DEFAULT now(),
  CONSTRAINT user_page_size_limit CHECK (octet_length(page::text) <= 1048576)
);

INSERT INTO profiles.user_page (user_id, page)
VALUES (
    1,
    '{
       "title": "Моя страница",
       "bio": "Backend developer",
       "links": [
         { "type": "telegram", "url": "https://t.me/username" },
         { "type": "github", "url": "https://github.com/username" }
       ],
       "theme": "dark"
     }'::jsonb
);
