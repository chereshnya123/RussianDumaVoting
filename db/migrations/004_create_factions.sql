CREATE TABLE IF NOT EXISTS factions (
  id INTEGER PRIMARY KEY AUTOINCREMENT,
  api_id INTEGER NOT NULL UNIQUE,
  name TEXT,
  head INTEGER,
  target_questions TEXT,
  target_tags TEXT
);
