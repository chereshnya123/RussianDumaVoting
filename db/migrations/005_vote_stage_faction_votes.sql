CREATE TABLE IF NOT EXISTS vote_stage_votes (
  id INTEGER PRIMARY KEY,
  for_count INTEGER,
  against_count INTEGER,
  abstained_count INTEGER,
  no_vote_count INTEGER,
  faction_id INTEGER REFERENCES factions(id),
  vote_stage_id INTEGER REFERENCES vote_stages(id)
);