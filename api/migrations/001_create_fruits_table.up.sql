CREATE TABLE fruits (
  id UUID PRIMARY KEY,
  name TEXT NOT NULL,
  quantity INT NOT NULL,
  price NUMERIC(10,2) NOT NULL,
  created_at TIMESTAMPTZ NOT NULL,
  updated_at TIMESTAMPTZ NOT NULL
);
