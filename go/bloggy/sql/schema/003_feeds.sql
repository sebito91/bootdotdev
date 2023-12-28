-- +goose Up
CREATE TABLE feeds (
    id uuid DEFAULT uuid_generate_v4(),
    created_at timestamp not null,
    updated_at timestamp not null,
    name varchar(125) unique not null,
    url varchar(125) unique not null,
    user_id uuid not null references users(id) on delete cascade, 
    PRIMARY KEY (id)

);

-- +goose Down
DROP TABLE feeds;
