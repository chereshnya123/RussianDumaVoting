CREATE TABLE IF NOT EXISTS votings (
  id INTEGER PRIMARY KEY,
  name TEXT,
  date DATETIME,
  question_id INT,
  factions TEXT,
  result TEXT
);
