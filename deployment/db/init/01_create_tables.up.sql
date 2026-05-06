CREATE TABLE IF NOT EXISTS cinemas (
  id INTEGER PRIMARY KEY,
  name VARCHAR(255),
  address TEXT
);

CREATE TABLE IF NOT EXISTS halls (
  id INTEGER PRIMARY KEY,
  cinema_id INTEGER NOT NULL REFERENCES cinemas(id),
  name VARCHAR(255),
  hall_type VARCHAR(50)
);

CREATE TABLE IF NOT EXISTS seats (
  id INTEGER PRIMARY KEY,
  hall_id INTEGER NOT NULL REFERENCES halls(id),
  row_number INTEGER NOT NULL,
  seat_number INTEGER NOT NULL,
  seat_type VARCHAR(50)
);

CREATE TABLE IF NOT EXISTS movies (
  id INTEGER PRIMARY KEY,
  title VARCHAR(255),
  duration INTEGER,
  rating VARCHAR(20),
  description TEXT
);

CREATE TABLE IF NOT EXISTS sessions (
  id INTEGER PRIMARY KEY,
  movie_id INTEGER NOT NULL REFERENCES movies(id),
  hall_id INTEGER NOT NULL REFERENCES halls(id),
  start_time TIMESTAMP NOT NULL,
  price_base NUMERIC(10,2) NOT NULL
);

CREATE TABLE IF NOT EXISTS users (
  id INTEGER PRIMARY KEY,
  email VARCHAR(255) UNIQUE,
  password_hash VARCHAR(255),
  phone VARCHAR(50),
  is_active BOOLEAN,
  created_at TIMESTAMP,
  updated_at TIMESTAMP
);

CREATE TABLE IF NOT EXISTS bookings (
  id UUID PRIMARY KEY,
  user_id INTEGER NOT NULL REFERENCES users(id),
  total_price NUMERIC(10,2) NOT NULL,
  status VARCHAR(30),
  created_at TIMESTAMP
);

CREATE TABLE IF NOT EXISTS tickets (
  id UUID PRIMARY KEY,
  session_id INTEGER NOT NULL REFERENCES sessions(id),
  seat_id INTEGER NOT NULL REFERENCES seats(id),
  booking_id UUID REFERENCES bookings(id),
  status VARCHAR(30)
);

CREATE TABLE IF NOT EXISTS reservations (
  seat_id INTEGER NOT NULL REFERENCES seats(id),
  session_id INTEGER NOT NULL REFERENCES sessions(id),
  user_id INTEGER NOT NULL REFERENCES users(id),
  locked_until TIMESTAMP,
  PRIMARY KEY (seat_id, session_id)
);