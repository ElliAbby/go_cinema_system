-- данные для кинотеатров
INSERT INTO cinemas (name, address) VALUES
  ('Кинопарк Центр', 'ул. Ленина, 1'),
  ('Кинопарк Восток', 'пр. Октября, 42'),
  ('Кинопарк Премиум', 'ул. Мира, 18')
ON CONFLICT DO NOTHING;

-- данные для залов
-- Каждый кинотеатр получает собственный набор залов и форматов
INSERT INTO halls (cinema_id, name, hall_type) VALUES
  (1, 'Зал 1 - Стандарт', 'стандарт'),
  (1, 'Зал 2 - VIP', 'VIP'),
  (1, 'Зал 3 - Family', 'family'),
  (2, 'Зал 1 - IMAX', 'IMAX'),
  (2, 'Зал 2 - Стандарт', 'стандарт'),
  (2, 'Зал 3 - Dolby', 'Dolby'),
  (3, 'Зал 1 - Premium', 'premium'),
  (3, 'Зал 2 - Lounge', 'lounge'),
  (3, 'Зал 3 - Deluxe', 'deluxe')
ON CONFLICT DO NOTHING;

-- данные для мест (Кинопарк Центр)
INSERT INTO seats (hall_id, row_number, seat_number, seat_type)
SELECT 1, row_num, seat_num,
  CASE
    WHEN row_num >= 9 THEN 'диван'
    WHEN row_num IN (5, 6) AND seat_num IN (4, 5, 6, 7) THEN 'для инвалидов'
    ELSE 'обычное'
  END
FROM generate_series(1, 10) AS row_num
CROSS JOIN generate_series(1, 10) AS seat_num
ON CONFLICT DO NOTHING;

INSERT INTO seats (hall_id, row_number, seat_number, seat_type)
SELECT 2, row_num, seat_num,
  CASE
    WHEN row_num <= 2 THEN 'VIP'
    ELSE 'обычное'
  END
FROM generate_series(1, 10) AS row_num
CROSS JOIN generate_series(1, 10) AS seat_num
ON CONFLICT DO NOTHING;

INSERT INTO seats (hall_id, row_number, seat_number, seat_type)
SELECT 3, row_num, seat_num,
  CASE
    WHEN row_num <= 2 THEN 'family'
    ELSE 'обычное'
  END
FROM generate_series(1, 10) AS row_num
CROSS JOIN generate_series(1, 10) AS seat_num
ON CONFLICT DO NOTHING;

-- данные для мест (Кинопарк Восток)
INSERT INTO seats (hall_id, row_number, seat_number, seat_type)
SELECT 4, row_num, seat_num, 'обычное'
FROM generate_series(1, 10) AS row_num
CROSS JOIN generate_series(1, 10) AS seat_num
ON CONFLICT DO NOTHING;

INSERT INTO seats (hall_id, row_number, seat_number, seat_type)
SELECT 5, row_num, seat_num,
  CASE
    WHEN row_num >= 8 THEN 'диван'
    ELSE 'обычное'
  END
FROM generate_series(1, 10) AS row_num
CROSS JOIN generate_series(1, 10) AS seat_num
ON CONFLICT DO NOTHING;

INSERT INTO seats (hall_id, row_number, seat_number, seat_type)
SELECT 6, row_num, seat_num,
  CASE
    WHEN row_num <= 3 THEN 'premium'
    ELSE 'обычное'
  END
FROM generate_series(1, 10) AS row_num
CROSS JOIN generate_series(1, 10) AS seat_num
ON CONFLICT DO NOTHING;

-- данные для мест (Кинопарк Премиум)
INSERT INTO seats (hall_id, row_number, seat_number, seat_type)
SELECT 7, row_num, seat_num,
  CASE
    WHEN row_num <= 2 THEN 'VIP'
    ELSE 'премиум'
  END
FROM generate_series(1, 10) AS row_num
CROSS JOIN generate_series(1, 10) AS seat_num
ON CONFLICT DO NOTHING;

INSERT INTO seats (hall_id, row_number, seat_number, seat_type)
SELECT 8, row_num, seat_num,
  CASE
    WHEN row_num IN (1, 2) THEN 'lounge'
    ELSE 'обычное'
  END
FROM generate_series(1, 10) AS row_num
CROSS JOIN generate_series(1, 10) AS seat_num
ON CONFLICT DO NOTHING;

INSERT INTO seats (hall_id, row_number, seat_number, seat_type)
SELECT 9, row_num, seat_num,
  CASE
    WHEN row_num >= 8 THEN 'диван'
    ELSE 'обычное'
  END
FROM generate_series(1, 10) AS row_num
CROSS JOIN generate_series(1, 10) AS seat_num
ON CONFLICT DO NOTHING;

-- данные для фильмов
-- Оставляем исходный набор и добавляем новые фильмы
INSERT INTO movies (title, duration, rating, description) VALUES
  ('Интерстеллар', 169, '9.5', 'Научно-фантастический фильм о путешествии человечества в далёкий космос'),
  ('Темный рыцарь', 152, '8.0', 'Криминальная драма о борьбе с преступностью в Готэме'),
  ('Начало', 148, '8.2', 'Научно-фантастический триллер о кражах из подсознания'),
  ('Матрица', 136, '7.8', 'Классический научно-фантастический фильм про виртуальный мир'),
  ('Аватар', 162, '7.6', 'Эпический научно-фантастический фильм на планете Пандора'),
  ('Оппенгеймер', 180, '8.7', 'История создания атомной бомбы и судьбы ее создателя'),
  ('Дюна: Часть вторая', 166, '8.5', 'Продолжение эпической истории о пустынной планете Арракис'),
  ('Человек-паук: Нет пути домой', 148, '8.1', 'Мультивселенское приключение Питера Паркера'),
  ('Джон Уик 4', 169, '8.0', 'Стильный экшен о возвращении легендарного киллера'),
  ('Гладиатор 2', 148, '8.4', 'Возвращение в мир Римской империи и новых арен'),
  ('Фуриоса', 148, '8.2', 'Постапокалиптическая история о дороге к свободе'),
  ('Дэдпул и Росомаха', 127, '8.0', 'Безумный супергеройский экшен с фирменным юмором')
ON CONFLICT DO NOTHING;

-- данные для сеансов
-- У каждого кинотеатра свои залы, а в залах показываются разные фильмы в разное время
INSERT INTO sessions (movie_id, hall_id, start_time, price_base) VALUES
  -- Кинопарк Центр
  (1, 1, NOW() + INTERVAL '2 hours', 350.00),
  (3, 1, NOW() + INTERVAL '1 day 1 hour', 350.00),
  (2, 2, NOW() + INTERVAL '3 hours', 450.00),
  (6, 2, NOW() + INTERVAL '1 day 3 hours', 450.00),
  (8, 3, NOW() + INTERVAL '4 hours', 400.00),
  (10, 3, NOW() + INTERVAL '1 day 4 hours', 400.00),

  -- Кинопарк Восток
  (4, 4, NOW() + INTERVAL '5 hours', 500.00),
  (7, 4, NOW() + INTERVAL '1 day 5 hours', 500.00),
  (5, 5, NOW() + INTERVAL '6 hours', 360.00),
  (11, 5, NOW() + INTERVAL '1 day 6 hours', 360.00),
  (6, 6, NOW() + INTERVAL '7 hours', 520.00),
  (12, 6, NOW() + INTERVAL '1 day 7 hours', 520.00),

  -- Кинопарк Премиум
  (9, 7, NOW() + INTERVAL '8 hours', 700.00),
  (10, 7, NOW() + INTERVAL '1 day 8 hours', 700.00),
  (11, 8, NOW() + INTERVAL '9 hours', 540.00),
  (1, 8, NOW() + INTERVAL '1 day 9 hours', 540.00),
  (12, 9, NOW() + INTERVAL '10 hours', 620.00),
  (5, 9, NOW() + INTERVAL '1 day 10 hours', 620.00)
ON CONFLICT DO NOTHING;

-- данные для тестовых пользователей
INSERT INTO users (email, password_hash, phone, is_active, created_at, updated_at) VALUES
  ('test@example.com', '$2a$10$N9qo8uLOickgx2ZMRZoMyeIjZAgcg7b3XeKeUxWdeS86E36P4/KLm', '+7-999-123-45-67', true, NOW(), NOW()),
  ('user@example.com', '$2a$10$abcdefghijklmnopqrstuvwxyzABCDEFGHIJKLMNOPQRSTUVWX', '+7-999-234-56-78', true, NOW(), NOW()),
  ('admin@example.com', '$2a$10$XYZABCDEFGHIJKLMNOPQRSTUVWXYZabcdefghijklmnopqrstuvwxyz', '+7-999-345-67-89', true, NOW(), NOW())
ON CONFLICT DO NOTHING;

-- Демонстрационные бронирования и занятые места
INSERT INTO bookings (id, user_id, total_price, status, created_at) VALUES
  (1001, 1, 1050.00, 'pending', NOW() - INTERVAL '20 minutes'),
  (1002, 2, 900.00, 'pending', NOW() - INTERVAL '18 minutes'),
  (1003, 3, 1280.00, 'pending', NOW() - INTERVAL '15 minutes')
ON CONFLICT DO NOTHING;

INSERT INTO reservations (seat_id, session_id, user_id, booking_id, locked_until)
SELECT s.id, 1, 1, 1001, NOW() + INTERVAL '25 minutes'
FROM seats s
WHERE s.hall_id = 1 AND s.row_number = 1 AND s.seat_number IN (1, 2, 3)
ON CONFLICT DO NOTHING;

INSERT INTO reservations (seat_id, session_id, user_id, booking_id, locked_until)
SELECT s.id, 3, 2, 1002, NOW() + INTERVAL '25 minutes'
FROM seats s
WHERE s.hall_id = 2 AND s.row_number = 2 AND s.seat_number IN (4, 5, 6)
ON CONFLICT DO NOTHING;

INSERT INTO reservations (seat_id, session_id, user_id, booking_id, locked_until)
SELECT s.id, 7, 3, 1003, NOW() + INTERVAL '25 minutes'
FROM seats s
WHERE s.hall_id = 4 AND s.row_number = 3 AND s.seat_number IN (2, 3, 4, 5)
ON CONFLICT DO NOTHING;