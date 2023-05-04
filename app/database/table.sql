CREATE TABLE IF NOT EXISTS users (id SERIAL PRIMARY KEY, first_name VARCHAR(255), last_name VARCHAR(255), email VARCHAR(255) UNIQUE, age INT, key VARCHAR(255) UNIQUE, role INT, active BOOLEAN, created_at TIMESTAMP, updated_at TIMESTAMP);

CREATE TABLE IF NOT EXISTS accounts (id SERIAL PRIMARY KEY, owner INT, type INT, number VARCHAR(255) UNIQUE, locked_balance FLOAT, ledger_balance FLOAT, balance FLOAT, active BOOLEAN, created_at TIMESTAMP, updated_at TIMESTAMP);

CREATE TABLE IF NOT EXISTS factory (id SERIAL PRIMARY KEY, key VARCHAR(255) UNIQUE, value INT, created_at TIMESTAMP, updated_at TIMESTAMP);

CREATE TABLE IF NOT EXISTS transactions (id SERIAL PRIMARY KEY, number VARCHAR(255), amount FLOAT, session_id VARCHAR(255) UNIQUE, type INT, status INT, created_at TIMESTAMP, updated_at TIMESTAMP);