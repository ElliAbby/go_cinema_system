-- Seed данные для кинотеатров
INSERT INTO cinemas (name, address) VALUES
  ('Кинопарк Центр', 'ул. Ленина, 1'),
  ('Кинопарк Восток', 'пр. Октября, 42')
ON CONFLICT DO NOTHING;

-- Seed данные для залов
INSERT INTO halls (cinema_id, name, hall_type) VALUES
  (1, 'Зал 1 - Стандарт', 'стандарт'),
  (1, 'Зал 2 - VIP', 'VIP'),
  (2, 'Зал 1 - IMAX', 'IMAX'),
  (2, 'Зал 2 - Стандарт', 'стандарт')
ON CONFLICT DO NOTHING;

-- Seed данные для мест (Зал 1 кинопарка 1 - 10 рядов х 12 мест)
INSERT INTO seats (hall_id, row_number, seat_number, seat_type)
SELECT 1, row_num, seat_num, 
  CASE 
    WHEN row_num >= 7 AND row_num <= 10 THEN 'диван'
    WHEN row_num IN (5, 6) AND seat_num IN (6, 7) THEN 'для инвалидов'
    ELSE 'обычное'
  END
FROM generate_series(1, 10) AS row_num
CROSS JOIN generate_series(1, 12) AS seat_num
ON CONFLICT DO NOTHING;

-- Seed данные для мест (Зал 2 кинопарка 1 - 8 рядов х 10 мест)
INSERT INTO seats (hall_id, row_number, seat_number, seat_type)
SELECT 2, row_num, seat_num, 'обычное'
FROM generate_series(1, 8) AS row_num
CROSS JOIN generate_series(1, 10) AS seat_num
ON CONFLICT DO NOTHING;

-- Seed данные для мест (Зал 1 кинопарка 2 - 12 рядов х 16 мест)
INSERT INTO seats (hall_id, row_number, seat_number, seat_type)
SELECT 3, row_num, seat_num, 'обычное'
FROM generate_series(1, 12) AS row_num
CROSS JOIN generate_series(1, 16) AS seat_num
ON CONFLICT DO NOTHING;

-- Seed данные для мест (Зал 2 кинопарка 2 - 9 рядов х 12 мест)
INSERT INTO seats (hall_id, row_number, seat_number, seat_type)
SELECT 4, row_num, seat_num, 'обычное'
FROM generate_series(1, 9) AS row_num
CROSS JOIN generate_series(1, 12) AS seat_num
ON CONFLICT DO NOTHING;

-- Seed данные для фильмов
INSERT INTO movies (title, duration, rating, description) VALUES
  ('Интерстеллар', 169, '9.5', 'Научно-фантастический фильм о путешествии человечества в далёкий космос'),
  ('Темный рыцарь', 152, '8.0', 'Криминальная драма о борьбе с преступностью в Готэме'),
  ('Начало', 148, '8.2', 'Научно-фантастический триллер о кражах из подсознания'),
  ('Матрица', 136, '7.8', 'Классический научно-фантастический фильм про виртуальный мир'),
  ('Аватар', 162, '7.6', 'Эпический научно-фантастический фильм на планете Пандора')
ON CONFLICT DO NOTHING;

-- Seed данные для сеансов (сегодня и завтра)
-- Зал 1, Кинопарка 1 - Интерстеллар
INSERT INTO sessions (movie_id, hall_id, start_time, price_base) VALUES
  (1, 1, NOW() + INTERVAL '2 hours', 350.00),
  (1, 1, NOW() + INTERVAL '6 hours', 350.00),
  (1, 1, NOW() + INTERVAL '1 day 2 hours', 350.00)
ON CONFLICT DO NOTHING;

-- Зал 2, Кинопарка 1 - Темный рыцарь
INSERT INTO sessions (movie_id, hall_id, start_time, price_base) VALUES
  (2, 2, NOW() + INTERVAL '3 hours', 450.00),
  (2, 2, NOW() + INTERVAL '7 hours', 450.00),
  (2, 2, NOW() + INTERVAL '1 day 3 hours', 450.00)
ON CONFLICT DO NOTHING;

-- Зал 1, Кинопарка 2 - Начало
INSERT INTO sessions (movie_id, hall_id, start_time, price_base) VALUES
  (3, 3, NOW() + INTERVAL '4 hours', 500.00),
  (3, 3, NOW() + INTERVAL '1 day 4 hours', 500.00)
ON CONFLICT DO NOTHING;

-- Зал 2, Кинопарка 2 - Матрица
INSERT INTO sessions (movie_id, hall_id, start_time, price_base) VALUES
  (4, 4, NOW() + INTERVAL '5 hours', 350.00),
  (4, 4, NOW() + INTERVAL '1 day 5 hours', 350.00)
ON CONFLICT DO NOTHING;

-- Seed данные для тестовых пользователей
INSERT INTO users (email, password_hash, phone, is_active, created_at, updated_at) VALUES
  ('test@example.com', '$2a$10$N9qo8uLOickgx2ZMRZoMyeIjZAgcg7b3XeKeUxWdeS86E36P4/KLm', '+7-999-123-45-67', true, NOW(), NOW()),
  ('user@example.com', '$2a$10$abcdefghijklmnopqrstuvwxyzABCDEFGHIJKLMNOPQRSTUVWX', '+7-999-234-56-78', true, NOW(), NOW()),
  ('admin@example.com', '$2a$10$XYZABCDEFGHIJKLMNOPQRSTUVWXYZabcdefghijklmnopqrstuvwxyz', '+7-999-345-67-89', true, NOW(), NOW())
ON CONFLICT DO NOTHING;