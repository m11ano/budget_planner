-- +goose Up

-- Таблица account
CREATE TABLE "account" (
    id                              UUID PRIMARY KEY DEFAULT gen_random_uuid(),
    email                           TEXT NOT NULL UNIQUE,
    password_hash                   TEXT NOT NULL,
    is_confirmed                    BOOLEAN NOT NULL DEFAULT FALSE,
    is_blocked                      BOOLEAN NOT NULL DEFAULT FALSE,
    last_login_at                   TIMESTAMPTZ,
    last_request_at                 TIMESTAMPTZ,
    last_request_ip                 INET,
    profile_name                    TEXT NOT NULL,
    profile_surname                 TEXT NOT NULL,
    created_at                      TIMESTAMPTZ NOT NULL DEFAULT NOW(),
    updated_at                      TIMESTAMPTZ NOT NULL DEFAULT NOW(),
    deleted_at                      TIMESTAMPTZ
);

-- Таблица auth_session
CREATE TABLE "auth_session" (
    id                          UUID PRIMARY KEY DEFAULT gen_random_uuid(),
    account_id                  UUID NOT NULL REFERENCES "account"(id) ON DELETE CASCADE,
    refresh_token               UUID NOT NULL,
    refresh_version             INT NOT NULL DEFAULT 1,
    refresh_token_issued_at     TIMESTAMPTZ NOT NULL,
    refresh_token_expires_at    TIMESTAMPTZ NOT NULL,
    refresh_token_request_ip    INET,
    create_request_ip           INET,
    created_at                  TIMESTAMPTZ NOT NULL DEFAULT NOW(),
    updated_at                  TIMESTAMPTZ NOT NULL DEFAULT NOW(),
    deleted_at                  TIMESTAMPTZ
);
CREATE INDEX idx_auth_session_account_id ON "auth_session" (account_id) WHERE deleted_at IS NULL;

-- +goose Down
DROP INDEX IF EXISTS idx_auth_session_account_id;
DROP TABLE IF EXISTS "auth_session";

DROP TABLE IF EXISTS "account";
