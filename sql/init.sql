-- Create tables
CREATE TABLE IF NOT EXISTS Accounts (
  account_id CHAR(36) PRIMARY KEY,    -- UUID
  email VARCHAR(256) NOT NULL,
  password_hash CHAR(32) NOT NULL,
  salt CHAR(32) NOT NULL
);

CREATE TABLE IF NOT EXISTS Sessions (
  session_token CHAR(32) PRIMARY KEY,
  account_id CHAR(36),
  csrf_token CHAR(32) NOT NULL,
  expire_id TIMESTAMP NOT NULL,
  expire_abs TIMESTAMP NOT NULL
);

-- Docker will initialize non-root user w/ all privileges
REVOKE ALL PRIVILEGES, GRANT OPTION FROM 'tp_int'@'%';
GRANT SELECT, INSERT, UPDATE, DELETE FROM 'tp_int'@'%';
FLUSH PRIVILEGES;

-- Add some phony data
INSERT INTO Accounts VALUES (
  "1f5e53d9-fccf-4e9c-a3bc-4f6d527327f9",
  "admin@example.com",
  "380a45771726817b2372bcd389e7e9b9",   -- password is sneakyadminpassword
  "8bc78e90a114942e38ee62a89b2f22cf"
);
