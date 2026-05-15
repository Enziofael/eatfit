CREATE TABLE IF NOT EXISTS consumption_records (
    id UUID PRIMARY KEY DEFAULT uuid_generate_v4(),
    user_id UUID NOT NULL REFERENCES users(id) ON DELETE CASCADE,
    meal_id UUID REFERENCES meals(id) ON DELETE SET NULL,
    meal_name VARCHAR(200) NOT NULL,
    amount DECIMAL(7,2) NOT NULL DEFAULT 100,
    calories DECIMAL(7,2) NOT NULL DEFAULT 0,
    proteins DECIMAL(6,2) NOT NULL DEFAULT 0,
    fats DECIMAL(6,2) NOT NULL DEFAULT 0,
    carbs DECIMAL(6,2) NOT NULL DEFAULT 0,
    water DECIMAL(7,2) NOT NULL DEFAULT 0,
    consumed_at TIMESTAMPTZ NOT NULL DEFAULT NOW()
);

CREATE INDEX IF NOT EXISTS idx_consumption_user_id ON consumption_records(user_id);
CREATE INDEX IF NOT EXISTS idx_consumption_date ON consumption_records(consumed_at DESC);
CREATE INDEX IF NOT EXISTS idx_consumption_meal_id ON consumption_records(meal_id);