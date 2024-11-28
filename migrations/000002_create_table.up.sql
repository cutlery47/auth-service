CREATE TABLE IF NOT EXISTS auth_schema.refresh (
    id          auth_schema.uuid_key        PRIMARY KEY,
    user_id     UUID                        NOT NULL,
    salt        UUID                        NOT NULL,
    hash        bytea                       NOT NULL,
    cost        INTEGER                     NOT NULL
);