-- Создаем тип ENUM для статусов заявок
CREATE TYPE order_status AS ENUM (
    'черновик', 
    'удалён', 
    'сформирован', 
    'завершён', 
    'отклонён'
);

-- Таблица пользователей
CREATE TABLE users (
    id SERIAL PRIMARY KEY,
    login VARCHAR(25) UNIQUE NOT NULL,
    password VARCHAR(100) NOT NULL,
    is_moderator BOOLEAN DEFAULT FALSE
);

-- Таблица услуг (формул ДЖЕЛ)
CREATE TABLE services (
    id SERIAL PRIMARY KEY,
    title VARCHAR(100) NOT NULL,
    description TEXT,
    formula TEXT NOT NULL,
    image VARCHAR(255),
    category VARCHAR(50),
    gender VARCHAR(20),
    min_age INTEGER,
    max_age INTEGER,
    is_active BOOLEAN DEFAULT TRUE
);

-- Таблица заявок (расчетов)
CREATE TABLE orders (
    id SERIAL PRIMARY KEY,
    status order_status NOT NULL,
    user_id INTEGER NOT NULL REFERENCES users(id),
    created_at TIMESTAMP NOT NULL DEFAULT CURRENT_TIMESTAMP,
    finalized_at TIMESTAMP,
    completed_at TIMESTAMP,
    moderator_id INTEGER REFERENCES users(id),
    doctor_name VARCHAR(100)
);

-- Таблица связи заявки-услуги
CREATE TABLE order_services (
    id SERIAL PRIMARY KEY,
    order_id INTEGER NOT NULL REFERENCES orders(id),
    service_id INTEGER NOT NULL REFERENCES services(id),
    quantity INTEGER DEFAULT 1,
    UNIQUE(order_id, service_id)
);

-- Создаем частичный уникальный индекс для ограничения "не более одного черновика на пользователя"
CREATE UNIQUE INDEX idx_orders_user_draft 
ON orders(user_id) 
WHERE status = 'черновик';

-- Вставляем тестовых пользователей
INSERT INTO users (login, password, is_moderator) VALUES
('user', 'user', FALSE),
('admin', 'admin', TRUE);

-- Вставляем услуги (формулы ДЖЕЛ)
INSERT INTO services (title, description, formula, image, category, gender, min_age, max_age) VALUES
('Мальчики 4-7 лет', 'Расчет ДЖЕЛ для мальчиков дошкольного возраста', 'ДЖЕЛ (л) = (0.043 × Рост) - (0.015 × Возраст) - 2.89', 'boys_4_7.png', 'дети', 'мужской', 4, 7),
('Девочки 4-7 лет', 'Расчет ДЖЕЛ для девочек дошкольного возраста', 'ДЖЕЛ (л) = (0.037 × Рост) - (0.012 × Возраст) - 2.54', 'girls_4_7.png', 'дети', 'женский', 4, 7),
('Мальчики 8-12 лет', 'Расчет ДЖЕЛ для мальчиков младшего школьного возраста', 'ДЖЕЛ (л) = (0.052 × Рост) - (0.022 × Возраст) - 4.60', 'boys_8_12.png', 'дети', 'мужской', 8, 12),
('Девочки 8-12 лет', 'Расчет ДЖЕЛ для девочек младшего школьного возраста', 'ДЖЕЛ (л) = (0.041 × Рост) - (0.018 × Возраст) - 3.70', 'girls_8_12.png', 'дети', 'женский', 8, 12),
('Юноши 13-17 лет', 'Расчет ДЖЕЛ для юношей подросткового возраста', 'ДЖЕЛ (л) = (0.052 × Рост) - (0.022 × Возраст) - 4.20', 'boys_13_17.png', 'подростки', 'мужской', 13, 17),
('Девушки 13-17 лет', 'Расчет ДЖЕЛ для девушек подросткового возраста', 'ДЖЕЛ (л) = (0.041 × Рост) - (0.018 × Возраст) - 3.20', 'girls_13_17.png', 'подростки', 'женский', 13, 17),
('Мужчины 18-60 лет', 'Расчет ДЖЕЛ для взрослых мужчин', 'ДЖЕЛ (л) = (0.052 × Рост) - (0.022 × Возраст) - 3.60', 'men_18_60.png', 'взрослые', 'мужской', 18, 60),
('Женщины 18-60 лет', 'Расчет ДЖЕЛ для взрослых женщин', 'ДЖЕЛ (л) = (0.041 × Рост) - (0.018 × Возраст) - 2.69', 'women_18_60.png', 'взрослые', 'женский', 18, 60),
('Пожилые 60+ лет', 'Расчет ДЖЕЛ для пожилых людей', 'ДЖЕЛ (л) = (0.044 × Рост) - (0.024 × Возраст) - 2.86', 'elderly_60plus.png', 'пожилые', 'унисекс', 60, 100);