CREATE EXTENSION IF NOT EXISTS "uuid-ossp";

CREATE TABLE IF NOT EXISTS users
(
    user_id      UUID PRIMARY KEY                     DEFAULT uuid_generate_v4(),
    first_name   VARCHAR(32)                 NOT NULL CHECK ( first_name <> '' ),
    last_name    VARCHAR(32)                 NOT NULL CHECK ( last_name <> '' ),
    email        VARCHAR(64) UNIQUE          NOT NULL CHECK ( email <> '' ),
    password     VARCHAR(250)                NOT NULL CHECK ( octet_length(password) <> 0 ),
    role         VARCHAR(10)                 NOT NULL DEFAULT 'user',
    about        VARCHAR(1024)                        DEFAULT '',
    phone_number VARCHAR(20),
    address      VARCHAR(250),
    city         VARCHAR(30),
    country      VARCHAR(30),
    gender       VARCHAR(20)                 NOT NULL DEFAULT 'male',
    postcode     INTEGER,
    birthday     DATE                                 DEFAULT NULL,
    created_at   TIMESTAMP WITH TIME ZONE    NOT NULL DEFAULT NOW(),
    updated_at   TIMESTAMP WITH TIME ZONE             DEFAULT CURRENT_TIMESTAMP,
    login_date   TIMESTAMP(0) WITH TIME ZONE NOT NULL DEFAULT CURRENT_TIMESTAMP
);

CREATE TABLE IF NOT EXISTS companies (
    company_id          uuid PRIMARY KEY DEFAULT uuid_generate_v4(),
    company_name        VARCHAR(15) UNIQUE NOT NULL CHECK ( company_name <> '' ),
    company_description TEXT,
    amount_of_employees INT NOT NULL,
    registered          BOOLEAN NOT NULL,
    company_type        VARCHAR(50)
                        CHECK 
                        (company_type 
                        IN ('Corporations', 'NonProfit', 'Cooperative', 'Sole Proprietorship'))
                        NOT NULL
);