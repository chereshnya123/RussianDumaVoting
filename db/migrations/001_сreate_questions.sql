CREATE TABLE IF NOT EXISTS questions (
  id INTEGER PRIMARY KEY,
  name TEXT NOT NULL,
  tags TEXT,
  vote_id INTEGER,
  departments TEXT,
  authors TEXT
);
