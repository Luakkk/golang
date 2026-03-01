CREATE TABLE IF NOT EXISTS movies (
  id BIGSERIAL PRIMARY KEY,
  title TEXT NOT NULL,
  genre TEXT NOT NULL,
  budget BIGINT NOT NULL DEFAULT 0,
  hero TEXT NOT NULL DEFAULT '',
  heroine TEXT NOT NULL DEFAULT '',
  created_at TIMESTAMPTZ NOT NULL DEFAULT NOW()
);

INSERT INTO movies (title, genre, budget, hero, heroine)
VALUES
('SAW', 'Horror', 500000, 'JONNY DEPP', 'Scarlett'),
('TEST', 'Romance', 1000000, 'BALE', 'ARMAS');
