-- =====================================================
-- Инициализация базы данных auth
-- =====================================================

-- 1. Создание базы данных auth (если не существует)
SELECT 'CREATE DATABASE auth'
WHERE NOT EXISTS (SELECT FROM pg_database WHERE datname = 'auth')\gexec

-- 2. Подключение к базе auth
\c auth

-- 3. Создание расширения pgcrypto (для gen_random_uuid())
CREATE EXTENSION IF NOT EXISTS "pgcrypto";

-- 4. Создание таблицы auth (если не существует)
CREATE TABLE IF NOT EXISTS auth (
    id UUID PRIMARY KEY DEFAULT gen_random_uuid(),
    user_id UUID NOT NULL,
    refresh_token VARCHAR(500) NOT NULL,
    expires_at TIMESTAMP WITH TIME ZONE NOT NULL,
    revoked BOOLEAN NOT NULL DEFAULT FALSE,
    created_at TIMESTAMP WITH TIME ZONE DEFAULT CURRENT_TIMESTAMP,
    updated_at TIMESTAMP WITH TIME ZONE DEFAULT CURRENT_TIMESTAMP,
    CONSTRAINT unique_refresh_token UNIQUE (refresh_token)
);

-- 5. Создание индексов для оптимизации
CREATE INDEX IF NOT EXISTS idx_auth_user_id ON auth(user_id);
CREATE INDEX IF NOT EXISTS idx_auth_expires_at ON auth(expires_at);
CREATE INDEX IF NOT EXISTS idx_auth_revoked ON auth(revoked);

-- 6. Комментарии к таблице и колонкам
COMMENT ON TABLE auth IS 'Таблица для хранения refresh-токенов аутентификации';
COMMENT ON COLUMN auth.id IS 'Уникальный идентификатор записи';
COMMENT ON COLUMN auth.user_id IS 'ID пользователя (связь с таблицей users)';
COMMENT ON COLUMN auth.refresh_token IS 'Refresh-токен (уникальный)';
COMMENT ON COLUMN auth.expires_at IS 'Дата и время истечения токена';
COMMENT ON COLUMN auth.revoked IS 'Признак отозванного токена (TRUE - отозван)';
COMMENT ON COLUMN auth.created_at IS 'Дата и время создания записи';
COMMENT ON COLUMN auth.updated_at IS 'Дата и время последнего обновления';

-- 7. Предоставление привилегий для пользователя admin (или любого другого)
-- (Используем пользователя из POSTGRES_USER)
GRANT ALL PRIVILEGES ON DATABASE auth TO admin;
GRANT ALL PRIVILEGES ON SCHEMA public TO admin;
GRANT ALL PRIVILEGES ON ALL TABLES IN SCHEMA public TO admin;
GRANT ALL PRIVILEGES ON ALL SEQUENCES IN SCHEMA public TO admin;

-- Права для будущих таблиц
ALTER DEFAULT PRIVILEGES IN SCHEMA public GRANT ALL PRIVILEGES ON TABLES TO admin;
ALTER DEFAULT PRIVILEGES IN SCHEMA public GRANT ALL PRIVILEGES ON SEQUENCES TO admin;

-- 8. Создание триггера для автоматического обновления updated_at
CREATE OR REPLACE FUNCTION update_auth_updated_at()
RETURNS TRIGGER AS $$
BEGIN
    NEW.updated_at = CURRENT_TIMESTAMP;
    RETURN NEW;
END;
$$ LANGUAGE plpgsql;

DROP TRIGGER IF EXISTS trigger_auth_updated_at ON auth;
CREATE TRIGGER trigger_auth_updated_at
BEFORE UPDATE ON auth
FOR EACH ROW
EXECUTE FUNCTION update_auth_updated_at();

-- 9. Подтверждение успешного выполнения
DO $$
BEGIN
    RAISE NOTICE '✅ База данных auth и таблица auth успешно созданы и настроены';
END $$;