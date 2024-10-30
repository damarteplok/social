CREATE TABLE IF NOT EXISTS roles (
    id BIGSERIAL PRIMARY key,
    name VARCHAR(256) NOT NULL UNIQUE,
    level int NOT NULL DEFAULT 0,
    description TEXT
);

INSERT INTO
    roles (name, description, level)
VALUES 
    (
        'user',
        'A user can create posts and comments',
        1
    );

INSERT INTO
    roles (name, description, level)
VALUES 
    (
        'moderator',
        'A moderator can update posts other users posts',
        2
    );

INSERT INTO
    roles (name, description, level)
VALUES 
    (
        'admin',
        'A admin can update and delete other users posts',
        3
    );