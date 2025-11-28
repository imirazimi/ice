
CREATE TABLE IF NOT EXISTS todos (
  id VARCHAR(64) PRIMARY KEY,
  description TEXT,
  due_date DATETIME
);
  