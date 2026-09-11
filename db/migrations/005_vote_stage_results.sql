CREATE TABLE IF NOT EXISTS vote_stage_results (
  id INTEGER PRIMARY KEY,
  for_count INTEGER,
  against_count INTEGER,
  abstained_count INTEGER,
  no_vote_count INTEGER,
  faction_id INTEGER REFERENCES factions(id)
);