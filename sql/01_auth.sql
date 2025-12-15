\connect bioly

CREATE TABLE IF NOT EXISTS auth.users (
  id             BIGSERIAL    PRIMARY KEY,
  username       TEXT         NOT NULL,
  password_hash  TEXT         NOT NULL,
  last_login_at  TIMESTAMPTZ  NULL,
  created_at     TIMESTAMPTZ  NOT NULL DEFAULT NOW(),
  updated_at     TIMESTAMPTZ  NOT NULL DEFAULT NOW()
);

CREATE UNIQUE INDEX IF NOT EXISTS users_username_lower_uidx
  ON auth.users (LOWER(username));

CREATE TABLE IF NOT EXISTS auth.refresh_tokens (
  id          BIGSERIAL    PRIMARY KEY,
  user_id     BIGINT       NOT NULL REFERENCES auth.users(id) ON DELETE CASCADE,
  jti         UUID         NOT NULL,
  token_hash  TEXT         NOT NULL,
  user_agent  TEXT         NULL,
  ip          INET         NULL,
  expires_at  TIMESTAMPTZ  NOT NULL,
  revoked_at  TIMESTAMPTZ  NULL,
  created_at  TIMESTAMPTZ  NOT NULL DEFAULT NOW(),
  updated_at  TIMESTAMPTZ  NOT NULL DEFAULT NOW()
);

CREATE UNIQUE INDEX IF NOT EXISTS refresh_tokens_jti_uidx
  ON auth.refresh_tokens (jti);

CREATE INDEX IF NOT EXISTS refresh_tokens_user_id_idx
  ON auth.refresh_tokens (user_id);

INSERT INTO auth.users (username, password_hash)
VALUES ('test', '$argon2id$v=19$m=65536,t=1,p=10$N+0U3LXewHdjFrkjrvn6NQ$5lowDuhO6KuqRdveEFIdOWe81KtJPTkANvgD4F/aqzk'); -- plain: password123

INSERT INTO auth.users (username, password_hash)
VALUES ('admin', '$argon2id$v=19$m=65536,t=1,p=10$DgmMFWnKxCF9Lv4jz90L1w$cb3nu9Wqf0pMTHiuEW6DR3F9KBNlMZd7bct7luZi0ws'); -- plain: admin

INSERT INTO auth.users (username, password_hash)
VALUES ('rootuser', '$argon2id$v=19$m=65536,t=1,p=10$GsrhMNY5iHQAxqO9d3nZMw$B6dXaGjlBes7n5cIhw93+CPzp25ionS1nWT5pykNqu4'); -- plain: rootuser

INSERT INTO auth.users (username, password_hash)
VALUES ('123456', '$argon2id$v=19$m=65536,t=1,p=10$0r1L1QlKkC+ZnUn/JAhxcA$TpjiJ8Qi4fRRoQ1IL0jPKtVCOV0xfA6o7mUphDhIbvw'); -- plain: 123456

INSERT INTO auth.users (username, password_hash)
VALUES ('login', '$argon2id$v=19$m=65536,t=1,p=10$KrmAusLsUEckKjkGkwKGsQ$m7yTibqGtf0MpPYMIuJnT0vAu0e6YUdj1AmfyLq7Stc'); -- plain: pass