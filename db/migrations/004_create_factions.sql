CREATE TABLE IF NOT EXISTS factions (
  id INT PRIMARY KEY,
  name TEXT -- name of faction
  head INT -- Id of deputy
  department INT -- id of departments
  target_questions TEXT -- json array of target questions
  target_tags TEXT -- json map {tag -> count of target laws}
);  