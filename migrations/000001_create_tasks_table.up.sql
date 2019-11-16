CREATE TABLE IF NOT EXISTS tasks(
  id serial PRIMARY KEY,
  title text not null,
  status integer not null unique,

  created_at  timestamp not null,
  updated_at  timestamp not null,
  deleted_at  timestamp
);
