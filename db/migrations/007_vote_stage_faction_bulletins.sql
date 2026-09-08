CREATE TABLE IF NOT EXISTS vote_stage_bulletins (
  id INTEGER PRIMARY KEY,
  for_count INTEGER,
  AGAINST INTEGER,
  ABSTAINED INTEGER,
  NO_VOTE INTEGER,
  faction_id INTEGER REFERENCES factions(id),
  vote_stage_id INTEGER REFERENCES vote_stages(id),
  UNIQUE (vote_stage_id, faction_id)
);