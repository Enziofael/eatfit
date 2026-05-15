-- Таблица блюд (без is_composite — определяется наличием компонентов)
CREATE TABLE IF NOT EXISTS meals (
    id UUID PRIMARY KEY DEFAULT uuid_generate_v4(),
    user_id UUID NOT NULL REFERENCES users(id) ON DELETE CASCADE,
    name VARCHAR(200) NOT NULL,
    description TEXT,
    recipe TEXT,
    image_url VARCHAR(500),
    calories DECIMAL(7,2) NOT NULL DEFAULT 0,
    proteins DECIMAL(6,2) NOT NULL DEFAULT 0,
    fats DECIMAL(6,2) NOT NULL DEFAULT 0,
    carbs DECIMAL(6,2) NOT NULL DEFAULT 0,
    water DECIMAL(7,2) NOT NULL DEFAULT 0,
    created_at TIMESTAMPTZ NOT NULL DEFAULT NOW(),
    updated_at TIMESTAMPTZ NOT NULL DEFAULT NOW()
);

CREATE INDEX IF NOT EXISTS idx_meals_user_id ON meals(user_id);
CREATE INDEX IF NOT EXISTS idx_meals_name ON meals(name);
CREATE INDEX IF NOT EXISTS idx_meals_calories ON meals(calories);
CREATE INDEX IF NOT EXISTS idx_meals_created_at ON meals(created_at DESC);

-- Компоненты блюд (связь многие-ко-многим между meals)
CREATE TABLE IF NOT EXISTS meal_components (
    id UUID PRIMARY KEY DEFAULT uuid_generate_v4(),
    meal_id UUID NOT NULL REFERENCES meals(id) ON DELETE CASCADE,
    component_meal_id UUID NOT NULL REFERENCES meals(id) ON DELETE RESTRICT, -- RESTRICT: нельзя удалить блюдо, если оно используется как компонент
    amount DECIMAL(7,2) NOT NULL DEFAULT 100,
    sort_order INT NOT NULL DEFAULT 0,
    UNIQUE(meal_id, component_meal_id)
);

CREATE INDEX IF NOT EXISTS idx_meal_components_meal_id ON meal_components(meal_id);
CREATE INDEX IF NOT EXISTS idx_meal_components_component_id ON meal_components(component_meal_id);