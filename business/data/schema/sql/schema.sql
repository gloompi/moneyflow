-- Version: 1.1
-- Description: Create table users
CREATE TABLE users (
	user_id       UUID,
	name          TEXT,
	email         TEXT UNIQUE,
	roles         TEXT[],
	password_hash TEXT,
	date_created  TIMESTAMP,
	date_updated  TIMESTAMP,

	PRIMARY KEY (user_id)
);

-- Version: 1.2
-- Description: Create type reoccurrence enum and duration enum
CREATE TYPE reoccurrence_type AS ENUM ('Monthly', 'Daily', 'Once');
CREATE TYPE duration_type AS ENUM ('Months', 'Days', 'None');

-- Version: 1.3
-- Description: Create table incomes
CREATE TABLE incomes (
	income_id           UUID,
  user_id             UUID,
	name                TEXT,
  category            TEXT,
  currency            TEXT,
	amount              INT,
	reoccurrence        INT,
  duration            INT,
  "reoccurrence_type" reoccurrence_type DEFAULT 'Monthly',
  "duration_type"     duration_type DEFAULT 'None',
	date_created        TIMESTAMP,
	date_updated        TIMESTAMP,

	PRIMARY KEY (income_id),
	FOREIGN KEY (user_id) REFERENCES users(user_id) ON DELETE CASCADE
);

-- Version: 1.4
-- Description: Create table expenses
CREATE TABLE expenses (
	expense_id          UUID,
  user_id             UUID,
	name                TEXT,
  category            TEXT,
  currency            TEXT,
	amount              INT,
	reoccurrence        INT,
  duration            INT,
  "reoccurrence_type" reoccurrence_type DEFAULT 'Monthly',
  "duration_type"     duration_type,
	date_created        TIMESTAMP,
	date_updated        TIMESTAMP,

	PRIMARY KEY (expense_id),
	FOREIGN KEY (user_id) REFERENCES users(user_id) ON DELETE CASCADE
);
