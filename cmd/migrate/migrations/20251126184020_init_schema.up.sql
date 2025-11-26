CREATE TABLE routemap (
  id BIGSERIAL PRIMARY KEY,
  path TEXT UNIQUE NOT NULL,
  destination TEXT NOT NULL,
  created_at TIMESTAMPTZ DEFAULT now()
);
