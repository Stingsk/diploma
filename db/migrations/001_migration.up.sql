CREATE SCHEMA IF NOT EXISTS gophermart;
SET SEARCH_PATH TO gophermart;

CREATE TABLE IF NOT EXISTS users(
  id SERIAL PRIMARY KEY,
  login VARCHAR (50) UNIQUE,
  password VARCHAR (255)
);
CREATE UNIQUE INDEX IF NOT EXISTS users_idx ON users USING btree (id);
CREATE UNIQUE INDEX IF NOT EXISTS users_login_uniq_idx ON users USING btree (login);

CREATE TABLE IF NOT EXISTS orders(
  number BIGINT PRIMARY KEY,
  login  VARCHAR (50) REFERENCES users(login),
  status VARCHAR (50) DEFAULT 'NEW',
  accrual DECIMAL DEFAULT NULL,
  withdraw DECIMAL DEFAULT NULL,
  uploaded_at TIMESTAMP DEFAULT now()
);
CREATE UNIQUE INDEX IF NOT EXISTS orders_login_idx ON orders USING btree (login);
CREATE UNIQUE INDEX IF NOT EXISTS order_number_idx ON orders USING btree (number);