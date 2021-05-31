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
  expire_idle TIMESTAMP NOT NULL,
  expire_abs TIMESTAMP NOT NULL
);

-- TODO: see if we can set user privileges

-- Add some phony data
INSERT INTO Accounts VALUES (
  "1f5e53d9-fccf-4e9c-a3bc-4f6d527327f9",
  "admin@example.com",
  "36871e27dc012f640cdd158c7c207d49",   -- password is sneakyadminpassword
  "8bc78e90a114942e38ee62a89b2f22cf"
);
