#!/bin/bash
set -e
export PGPASSWORD=$POSTGRES_PASSWORD;

psql -v ON_ERROR_STOP=1 --username "$POSTGRES_USER" --dbname "$POSTGRES_DB" <<-EOSQL
  CREATE USER $APP_DB_USER WITH PASSWORD '$APP_DB_PASS';
  CREATE DATABASE $APP_DB_NAME;
  GRANT ALL PRIVILEGES ON DATABASE $APP_DB_NAME TO $APP_DB_USER;
  \connect $APP_DB_NAME $APP_DB_USER
  BEGIN;
  CREATE TABLE IF NOT EXISTS users (
    acct varchar(30) NOT NULL PRIMARY KEY,
    pwd varchar(100) NOT NULL CHECK (CHAR_LENGTH(pwd) >= 50),
    fullname varchar(50) NOT NULL,
    created_at timestamptz NULL DEFAULT CURRENT_TIMESTAMP,
    updated_at timestamptz NULL DEFAULT CURRENT_TIMESTAMP
	);
	CREATE INDEX users_fullname_idx ON public.users (fullname);
  COMMIT;
EOSQL