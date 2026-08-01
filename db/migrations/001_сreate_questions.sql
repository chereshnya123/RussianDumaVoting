CREATE TABLE IF NOT EXISTS questions (
  id INT PRIMARY KEY, -- (example: 1208563-8)
  name TEXT NOT NULL, -- Header
  tags JSONB -- JSON array of tags
  vote_id INT -- Id of voting
  departments JSONB -- json array of departments ids
  authors JSONB -- json array of deputies ID
);  