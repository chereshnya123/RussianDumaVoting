CREATE TABLE IF NOT EXISTS vote (
  id INT PRIMARY KEY,
  factions JSONB -- json array of factions id
  result JSONB -- json: {id -> {accept: 10, decline: 20, ...}, ...}
);  