-- +goose Up
CREATE TABLE posts (
    id uuid DEFAULT uuid_generate_v4(),
    created_at timestamp not null,
    updated_at timestamp not null,
    title varchar(125) not null,
    url varchar(125) unique not null,
    description varchar(255),
    published_at timestamp not null,
    feed_id uuid not null references feeds(id) on delete cascade, 
    unique (title, url, feed_id),
    PRIMARY KEY (id)

);

-- +goose Down
DROP TABLE posts;
