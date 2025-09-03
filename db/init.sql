CREATE EXTENSION IF NOT EXISTS pgcrypto;
CREATE EXTENSION IF NOT EXISTS citext;

-- ユーザー（メール認証 + 現在のエナジーを保持）
CREATE TABLE IF NOT EXISTS users (
  id             UUID    PRIMARY KEY DEFAULT gen_random_uuid(),
  email          CITEXT  NOT NULL UNIQUE,
  password       TEXT    NOT NULL,
  display_name   TEXT    NOT NULL CHECK (length(display_name) <= 50),
  energy_value   SMALLINT NOT NULL DEFAULT 3 CHECK (energy_value BETWEEN 1 AND 5),
  created_at     TIMESTAMPTZ NOT NULL DEFAULT now(),
  updated_at     TIMESTAMPTZ NOT NULL DEFAULT now()
);
CREATE INDEX IF NOT EXISTS idx_users_updated_at ON users (updated_at DESC);

-- グループ（1曲だけ紐づけ）
CREATE TABLE IF NOT EXISTS groups (
  id            UUID PRIMARY KEY DEFAULT gen_random_uuid(),
  name          TEXT NOT NULL CHECK (length(name) <= 50),
  owner_id      UUID NOT NULL REFERENCES users(id) ON DELETE RESTRICT,
  track_id      TEXT, --spotify track id
  created_at    TIMESTAMPTZ NOT NULL DEFAULT now()
  refresh_interval_hours SMALLINT NOT NULL DEFAULT 24 
    CHECK (refresh_interval_hours IN (12,24,48,72,96,120,144,168))  
);
CREATE INDEX IF NOT EXISTS idx_groups_owner ON groups(owner_id);

-- グループ＆ユーザの一覧
CREATE TABLE IF NOT EXISTS group_members (
  group_id  UUID NOT NULL REFERENCES groups(id) ON DELETE CASCADE,
  user_id   UUID NOT NULL REFERENCES users(id)  ON DELETE CASCADE,
  role      BOOLEAN NOT NULL DEFAULT FALSE,
  joined_at TIMESTAMPTZ NOT NULL DEFAULT now(),
  PRIMARY KEY (group_id, user_id)
);
CREATE INDEX IF NOT EXISTS idx_group_members_user ON group_members(user_id);

-- セッション
CREATE TABLE IF NOT EXISTS sessions (
  token      text PRIMARY KEY,
  user_id    text NOT NULL,
  expires_at timestamptz NOT NULL,
  created_at timestamptz NOT NULL DEFAULT now()
);
CREATE INDEX IF NOT EXISTS idx_sessions_expires_at ON sessions (expires_at);
CREATE INDEX IF NOT EXISTS idx_sessions_user_id   ON sessions (user_id);