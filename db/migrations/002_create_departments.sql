CREATE TABLE IF NOT EXISTS departments (
  id TEXT PRIMARY KEY,
  head INT NOT NULL, -- Id of deputy
  name TEXT -- name of department
  size INT -- count of deputies in department
  members TEXT -- json array of deputies ID
);  