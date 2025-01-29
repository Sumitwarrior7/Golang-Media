ALTER TABLE
  IF EXISTS users
ADD
  COLUMN role_id INT REFERENCES roles(id) DEFAULT 1;

-- Sets the role_id for every row in the users table to the id from the roles table where the role's name is 'user'
UPDATE
  users
SET
  role_id = (
    SELECT id FROM roles WHERE name = 'user'
  );

-- DEFAULT condition is removed from role_id
ALTER TABLE
  users
ALTER COLUMN
  role_id DROP DEFAULT;

ALTER TABLE
  users
ALTER COLUMN
  role_id
SET
  NOT NULL;