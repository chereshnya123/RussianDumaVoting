CREATE TABLE IF NOT EXISTS vote_stages (
  id INTEGER PRIMARY KEY AUTOINCREMENT,
  law_draft_id INTEGER REFERENCES law_drafts(id),
  stage_name TEXT,
  for_count INTEGER,
  against_count INTEGER,
  abstained_count INTEGER,
  no_vote_count INTEGER
);