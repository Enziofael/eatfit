-- backend/migrations/001_create_users.sql

-- ============================================================================
-- Eatfit - Схема аутентификации и управления пользователями
-- Версия: 1.0.0
-- Описание: Базовые таблицы для регистрации, аутентификации и сессий
-- ============================================================================

BEGIN;

-- ============================================================================
-- Расширения
-- ============================================================================
CREATE EXTENSION IF NOT EXISTS "uuid-ossp";     -- Генерация UUID
CREATE EXTENSION IF NOT EXISTS "pgcrypto";       -- Криптографические функции

-- ============================================================================
-- Таблица пользователей
-- Содержит основную информацию об учётных записях
-- ============================================================================
CREATE TABLE IF NOT EXISTS users (
    -- Уникальный идентификатор пользователя
    -- Формат: UUID v4, генерируется автоматически
    user_id         UUID PRIMARY KEY DEFAULT uuid_generate_v4(),
    
    -- Email пользователя (уникальный в пределах системы)
    -- Используется для входа и отправки уведомлений
    email           VARCHAR(255) UNIQUE NOT NULL,
    
    -- Логин пользователя (уникальный в пределах системы)
    -- Используется для входа и отображения в интерфейсе
    login           VARCHAR(30) UNIQUE NOT NULL,
    
    -- Хеш пароля (bcrypt)
    -- Никогда не храним пароли в открытом виде
    password_hash   VARCHAR(255) NOT NULL,
    
    -- Статус подтверждения email адреса
    email_verified  BOOLEAN NOT NULL DEFAULT FALSE,
    
    -- Код для подтверждения email (опционально)
    verification_code VARCHAR(6),
    
    -- Срок действия кода подтверждения
    verification_code_expires_at TIMESTAMP WITH TIME ZONE,
    
    -- Временные метки
    created_at      TIMESTAMP WITH TIME ZONE NOT NULL DEFAULT NOW(),
    updated_at      TIMESTAMP WITH TIME ZONE NOT NULL DEFAULT NOW(),
    
    -- Ограничения целостности
    CONSTRAINT chk_email_format 
        CHECK (email ~* '^[A-Za-z0-9._%+-]+@[A-Za-z0-9.-]+\.[A-Za-z]{2,}$'),
    CONSTRAINT chk_login_format 
        CHECK (login ~* '^[a-zA-Z0-9_]{3,30}$'),
    CONSTRAINT chk_password_hash_not_empty 
        CHECK (length(password_hash) > 0)
);

-- ============================================================================
-- Таблица сессий
-- Хранит refresh токены для управления аутентифицированными сессиями
-- ============================================================================
CREATE TABLE IF NOT EXISTS sessions (
    -- Уникальный идентификатор сессии
    session_id      UUID PRIMARY KEY DEFAULT uuid_generate_v4(),
    
    -- Связь с пользователем
    user_id         UUID NOT NULL REFERENCES users(user_id) ON DELETE CASCADE,
    
    -- Refresh токен (уникальный)
    refresh_token   VARCHAR(512) UNIQUE NOT NULL,
    
    -- Информация об устройстве (опционально)
    device_info     VARCHAR(255),
    
    -- IP адрес, с которого создана сессия (опционально)
    ip_address      INET,
    
    -- Срок действия refresh токена
    expires_at      TIMESTAMP WITH TIME ZONE NOT NULL,
    
    -- Признак отзыва токена (мягкое удаление)
    is_revoked      BOOLEAN NOT NULL DEFAULT FALSE,
    
    -- Время создания сессии
    created_at      TIMESTAMP WITH TIME ZONE NOT NULL DEFAULT NOW(),
    
    -- Ограничения
    CONSTRAINT fk_session_user 
        FOREIGN KEY (user_id) REFERENCES users(user_id) ON DELETE CASCADE,
    CONSTRAINT chk_refresh_token_not_empty 
        CHECK (length(refresh_token) > 0)
);

-- ============================================================================
-- Таблица попыток входа (для безопасности и аудита)
-- ============================================================================
CREATE TABLE IF NOT EXISTS login_attempts (
    -- Уникальный идентификатор попытки
    attempt_id      UUID PRIMARY KEY DEFAULT uuid_generate_v4(),
    
    -- Идентификатор пользователя (может быть NULL для неудачных попыток)
    user_id         UUID REFERENCES users(user_id) ON DELETE SET NULL,
    
    -- Логин или email, который использовался при попытке входа
    login_field     VARCHAR(255) NOT NULL,
    
    -- IP адрес, с которого была попытка
    ip_address      INET,
    
    -- User-Agent браузера/приложения
    user_agent      VARCHAR(512),
    
    -- Успешность попытки
    is_successful   BOOLEAN NOT NULL DEFAULT FALSE,
    
    -- Причина неудачи (опционально)
    failure_reason  VARCHAR(100),
    
    -- Время попытки
    attempted_at    TIMESTAMP WITH TIME ZONE NOT NULL DEFAULT NOW()
);

-- ============================================================================
-- Индексы для оптимизации запросов
-- ============================================================================

-- Поиск пользователей
CREATE INDEX idx_users_email ON users(email);
CREATE INDEX idx_users_login ON users(login);
CREATE INDEX idx_users_email_verified ON users(email_verified) WHERE email_verified = FALSE;

-- Поиск сессий
CREATE INDEX idx_sessions_user_id ON sessions(user_id);
CREATE INDEX idx_sessions_refresh_token ON sessions(refresh_token);
CREATE INDEX idx_sessions_expires_at ON sessions(expires_at) WHERE is_revoked = FALSE;
CREATE INDEX idx_sessions_active ON sessions(user_id, is_revoked) WHERE is_revoked = FALSE;

-- Анализ попыток входа
CREATE INDEX idx_login_attempts_user_id ON login_attempts(user_id);
CREATE INDEX idx_login_attempts_attempted_at ON login_attempts(attempted_at);
CREATE INDEX idx_login_attempts_successful ON login_attempts(is_successful, attempted_at);

-- ============================================================================
-- Триггер для автоматического обновления updated_at
-- ============================================================================
CREATE OR REPLACE FUNCTION update_updated_at_column()
RETURNS TRIGGER AS $$
BEGIN
    NEW.updated_at = NOW();
    RETURN NEW;
END;
$$ language 'plpgsql';

CREATE TRIGGER update_users_updated_at 
    BEFORE UPDATE ON users 
    FOR EACH ROW 
    EXECUTE FUNCTION update_updated_at_column();

-- ============================================================================
-- Комментарии к таблицам и колонкам
-- ============================================================================
COMMENT ON TABLE users IS 'Основная таблица пользователей системы Eatfit';
COMMENT ON COLUMN users.user_id IS 'Уникальный идентификатор пользователя (UUID v4)';
COMMENT ON COLUMN users.email IS 'Email пользователя, используется для входа и уведомлений';
COMMENT ON COLUMN users.login IS 'Логин пользователя, отображается в интерфейсе';
COMMENT ON COLUMN users.password_hash IS 'Хеш пароля (bcrypt, cost=10)';
COMMENT ON COLUMN users.email_verified IS 'Статус подтверждения email адреса';

COMMENT ON TABLE sessions IS 'Активные сессии пользователей';
COMMENT ON COLUMN sessions.session_id IS 'Уникальный идентификатор сессии';
COMMENT ON COLUMN sessions.refresh_token IS 'Refresh токен для обновления access токена';
COMMENT ON COLUMN sessions.is_revoked IS 'Признак отзыва токена (true = недействителен)';

COMMENT ON TABLE login_attempts IS 'Журнал попыток входа в систему';
COMMENT ON COLUMN login_attempts.login_field IS 'Логин или email, использованный при входе';
COMMENT ON COLUMN login_attempts.is_successful IS 'Успешность попытки входа';

COMMIT;