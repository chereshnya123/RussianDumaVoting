CREATE TABLE IF NOT EXISTS questions (
  id INTEGER PRIMARY KEY,
  name TEXT NOT NULL,
  tags TEXT,
  votings_id TEXT, -- json list of votings ids
  profile_committee_id INTEGER,
  responsible_committee_id INTEGER,
  other_committees TEXT, -- json list of committees ids
  authors TEXT -- json list of authors ids
);
