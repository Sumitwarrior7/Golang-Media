CREATE TABLE IF NOT EXISTS roles (
    id bigserial PRIMARY KEY,
    name varchar(255) not null unique,
    level int not null default 0,
    description text
);


INSERT INTO 
    roles(name, description, level)
VALUES 
    (
        'user',
        'user can create posts and comments',
        1
    );

INSERT INTO 
    roles(name, description, level)
VALUES 
    (
        'moderator',
        'moderator can update other user posts',
        2
    );

INSERT INTO 
    roles(name, description, level)
VALUES 
    (
        'admin',
        'admin can update and delete other user posts',
        3
    );