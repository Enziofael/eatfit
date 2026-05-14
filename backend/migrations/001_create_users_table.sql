-- migrations/001_create_users_table.sql
-- Создание таблицы пользователей

-- Включаем расширение для UUID
CREATE EXTENSION IF NOT EXISTS "uuid-ossp";

-- Тип для статуса верификации (с проверкой существования)
DO $$ 
BEGIN
    IF NOT EXISTS (SELECT 1 FROM pg_type WHERE typname = 'email_verification_status') THEN
        CREATE TYPE email_verification_status AS ENUM (
            'pending',   -- ожидает подтверждения
            'verified',  -- подтверждён
            'expired'    -- код истёк (для возможной очистки)
        );
    END IF;
END $$;

-- Таблица пользователей
CREATE TABLE IF NOT EXISTS users (
    -- Первичный ключ
    id UUID PRIMARY KEY DEFAULT uuid_generate_v4(),
    
    -- Уникальные идентификаторы
    email VARCHAR(255) NOT NULL UNIQUE,
    login VARCHAR(30) NOT NULL UNIQUE,
    
    -- Безопасность
    password_hash VARCHAR(255) NOT NULL,
    
    -- Статус верификации email
    email_verified BOOLEAN NOT NULL DEFAULT FALSE,
    verification_status email_verification_status NOT NULL DEFAULT 'pending',
    verification_code VARCHAR(6),
    verification_code_expires_at TIMESTAMPTZ,
    verification_attempts INT NOT NULL DEFAULT 0,
    
    -- Временные метки
    created_at TIMESTAMPTZ NOT NULL DEFAULT NOW(),
    updated_at TIMESTAMPTZ NOT NULL DEFAULT NOW(),
    last_login_at TIMESTAMPTZ,
    
    -- Ограничения
    CONSTRAINT chk_login_format CHECK (login ~ '^[a-zA-Z0-9_]{3,30}$'),
    CONSTRAINT chk_email_format CHECK (email ~ '^[a-zA-Z0-9._%+-]+@[a-zA-Z0-9.-]+\.[a-zA-Z]{2,}$'),
    CONSTRAINT chk_verification_attempts CHECK (verification_attempts >= 0 AND verification_attempts <= 5)
);

-- Индексы для быстрого поиска
CREATE INDEX IF NOT EXISTS idx_users_email ON users(email);
CREATE INDEX IF NOT EXISTS idx_users_login ON users(login);
CREATE INDEX IF NOT EXISTS idx_users_verification_status ON users(verification_status);

-- Таблица refresh токенов (сессий)
CREATE TABLE IF NOT EXISTS refresh_tokens (
    -- Первичный ключ
    id UUID PRIMARY KEY DEFAULT uuid_generate_v4(),
    
    -- Связь с пользователем
    user_id UUID NOT NULL REFERENCES users(id) ON DELETE CASCADE,
    
    -- Токен (хешированный)
    token_hash VARCHAR(255) NOT NULL UNIQUE,
    
    -- Метаданные
    device_info VARCHAR(255),
    ip_address INET,
    user_agent TEXT,
    
    -- Временные метки
    created_at TIMESTAMPTZ NOT NULL DEFAULT NOW(),
    expires_at TIMESTAMPTZ NOT NULL,
    revoked_at TIMESTAMPTZ,
    
    -- Флаг отзыва
    is_revoked BOOLEAN NOT NULL DEFAULT FALSE
);

-- Индексы для управления сессиями
CREATE INDEX IF NOT EXISTS idx_refresh_tokens_user_id ON refresh_tokens(user_id);
CREATE INDEX IF NOT EXISTS idx_refresh_tokens_token_hash ON refresh_tokens(token_hash);
CREATE INDEX IF NOT EXISTS idx_refresh_tokens_expires ON refresh_tokens(expires_at) WHERE is_revoked = FALSE;

-- Таблица попыток входа (для безопасности)
CREATE TABLE IF NOT EXISTS login_attempts (
    id BIGSERIAL PRIMARY KEY,
    login_identifier VARCHAR(255) NOT NULL,  -- email или login
    ip_address INET NOT NULL,
    success BOOLEAN NOT NULL DEFAULT FALSE,
    attempted_at TIMESTAMPTZ NOT NULL DEFAULT NOW()
);

-- Индекс для отслеживания попыток входа
CREATE INDEX IF NOT EXISTS idx_login_attempts_identifier ON login_attempts(login_identifier, attempted_at);