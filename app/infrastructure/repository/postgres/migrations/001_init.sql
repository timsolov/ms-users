-- +migrate Up
CREATE TABLE users (
    user_id     UUID   PRIMARY KEY NOT NULL,
    email       VARCHAR(256)   NOT NULL,
    password    VARCHAR(1024)  NOT NULL,
    first_name  VARCHAR(50)    NOT NULL,
    last_name   VARCHAR(50)    NOT NULL,
    created_at  TIMESTAMPTZ    NOT NULL DEFAULT NOW(),
    updated_at  TIMESTAMPTZ    NOT NULL DEFAULT NOW()
);


-- +migrate Down
DROP TABLE users;