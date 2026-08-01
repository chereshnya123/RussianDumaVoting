CREATE TABLE IF NOT EXISTS questions (
  id TEXT PRIMARY KEY, -- (example: 1208563-8)
  name TEXT NOT NULL, -- Header
  tags TEXT -- JSON array of tags
  vote_id INT -- Id of voting
  departments TEXT -- json array of departments ids
  authors TEXT -- json array of deputies ID
);  