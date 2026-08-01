CREATE TABLE IF NOT EXISTS factions (
  id INT PRIMARY KEY,
  name TEXT -- name of faction
  head INT -- Id of deputy
  department INT -- id of departments
  target_questions JSONB -- json array of target questions
  target_tags JSONB -- json map {tag -> count of target laws}
);  