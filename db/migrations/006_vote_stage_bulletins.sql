CREATE TABLE IF NOT EXISTS vote_stage_bulletins (
  id INTEGER PRIMARY KEY,
  verdict VARCHAR(10) CHECK(verdict in ('FOR', 'AGAINST', 'ABSTAINED', 'NO_VOTE')),
  deputy_id INTEGER REFERENCES deputies(id),
  vote_stage_id INTEGER REFERENCES vote_stages(id),
  UNIQUE (vote_stage_id, deputy_id)
);