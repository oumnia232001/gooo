-- +goose Up
CREATE TABLE todo_models (
    id BIGINT(20) AUTO_INCREMENT PRIMARY KEY,
    title LONGTEXT NOT NULL,
    completed TINYINT(1) DEFAULT 0,
    created_at DATETIME(3) NOT NULL,
    updated_at DATETIME(3) NOT NULL
);

-- +goose Down
DROP TABLE todo_models;
