CREATE TABLE IF NOT EXISTS vote (
  id TEXT PRIMARY KEY,
  factions TEXT -- json array of factions id
  result TEXT -- json: {id -> {accept: 10, decline: 20, ...}, ...}
);  