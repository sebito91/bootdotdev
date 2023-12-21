-- +goose Up
CREATE TABLE users (
    id uuid DEFAULT uuid_generate_v4(),
    created_at timestamp not null,
    updated_at timestamp not null,
    name varchar(75) unique not null,
    PRIMARY KEY (id)

);

-- +goose Down
DROP TABLE users;
