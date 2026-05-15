-- Таблица профилей
CREATE TABLE IF NOT EXISTS profiles (
    user_id UUID PRIMARY KEY REFERENCES users(id) ON DELETE CASCADE,
    first_name VARCHAR(50),
    last_name VARCHAR(50),
    avatar_url VARCHAR(500),
    height DECIMAL(5,1) CHECK (height > 0 AND height < 300),
    birth_date DATE CHECK (birth_date > '1900-01-01' AND birth_date < NOW()),
    gender VARCHAR(10) CHECK (gender IN ('male', 'female', 'other')),
    bio TEXT,
    created_at TIMESTAMPTZ NOT NULL DEFAULT NOW(),
    updated_at TIMESTAMPTZ NOT NULL DEFAULT NOW()
);

CREATE INDEX idx_profiles_user_id ON profiles(user_id);

-- Таблица истории веса
CREATE TABLE IF NOT EXISTS weight_history (
    id UUID PRIMARY KEY DEFAULT uuid_generate_v4(),
    user_id UUID NOT NULL REFERENCES users(id) ON DELETE CASCADE,
    weight DECIMAL(5,2) NOT NULL CHECK (weight > 0 AND weight < 500),
    note VARCHAR(255),
    recorded_at TIMESTAMPTZ NOT NULL DEFAULT NOW()
);

CREATE INDEX idx_weight_history_user_id ON weight_history(user_id);
CREATE INDEX idx_weight_history_recorded_at ON weight_history(recorded_at DESC);

-- Таблица норм КБЖУ
CREATE TABLE IF NOT EXISTS nutrition_norms (
    id UUID PRIMARY KEY DEFAULT uuid_generate_v4(),
    user_id UUID NOT NULL REFERENCES users(id) ON DELETE CASCADE,
    calories DECIMAL(7,2) NOT NULL DEFAULT 0 CHECK (calories >= 0),
    proteins DECIMAL(6,2) NOT NULL DEFAULT 0 CHECK (proteins >= 0),
    fats DECIMAL(6,2) NOT NULL DEFAULT 0 CHECK (fats >= 0),
    carbs DECIMAL(6,2) NOT NULL DEFAULT 0 CHECK (carbs >= 0),
    water DECIMAL(7,2) NOT NULL DEFAULT 0 CHECK (water >= 0),
    created_at TIMESTAMPTZ NOT NULL DEFAULT NOW()
);

CREATE INDEX idx_nutrition_norms_user_id ON nutrition_norms(user_id);
CREATE INDEX idx_nutrition_norms_created_at ON nutrition_norms(created_at DESC);

-- Функция для расчёта возраста
CREATE OR REPLACE FUNCTION calculate_age(birth_date DATE)
RETURNS INT AS $$
BEGIN
    RETURN DATE_PART('year', AGE(NOW(), birth_date));
END;
$$ LANGUAGE plpgsql;