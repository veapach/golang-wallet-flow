CREATE TABLE IF NOT EXISTS users (
                                     id              SERIAL                PRIMARY KEY,
                                     email           VARCHAR(255) NOT NULL UNIQUE,
                                     password_hashed VARCHAR(255) NOT NULL,
                                     name            VARCHAR(100) NOT NULL
);

CREATE TABLE IF NOT EXISTS user_tokens (
                                           id              SERIAL                PRIMARY KEY,
                                           user_id         INT                   NOT NULL REFERENCES users(id),
                                           access_token    VARCHAR(255)          NOT NULL UNIQUE,
                                           refresh_token   VARCHAR(255)          NOT NULL UNIQUE,
                                           created_at      TIMESTAMP             DEFAULT CURRENT_TIMESTAMP
);

CREATE TABLE IF NOT EXISTS currencies (
                                          id              SERIAL                PRIMARY KEY,
                                          name            VARCHAR(10)           NOT NULL UNIQUE
);

INSERT INTO currencies (name) VALUES
                                  ('RUB'),
                                  ('USD'),
                                  ('EUR'),
                                  ('GBP'),
                                  ('CNY');

CREATE TABLE IF NOT EXISTS wallets (
                                       id              SERIAL                PRIMARY KEY,
                                       user_id         INT                   NOT NULL REFERENCES users(id),
                                       currency_id     INT                   NOT NULL DEFAULT 1 REFERENCES currencies(id),
                                       balance         DECIMAL(10, 2)        NOT NULL DEFAULT 0,
                                       created_at      TIMESTAMP             DEFAULT CURRENT_TIMESTAMP,
                                       updated_at      TIMESTAMP             DEFAULT CURRENT_TIMESTAMP
);

CREATE TABLE IF NOT EXISTS transactions_history (
                                                    id              SERIAL                PRIMARY KEY,
                                                    type            VARCHAR(20)           NOT NULL CHECK (type IN ('deposit', 'withdraw', 'transfer', 'exchange')),
                                                    amount          DECIMAL(18, 2)        NOT NULL CHECK (amount > 0),
                                                    wallet_id_from  INT                   REFERENCES wallets(id),
                                                    wallet_id_to    INT                   REFERENCES wallets(id),
                                                    created_at      TIMESTAMP             DEFAULT CURRENT_TIMESTAMP,

                                                    CONSTRAINT chk_wallets_not_both_null CHECK (wallet_id_from IS NOT NULL OR wallet_id_to IS NOT NULL),
                                                    CONSTRAINT chk_wallets_not_self CHECK (wallet_id_from IS DISTINCT FROM wallet_id_to)
);