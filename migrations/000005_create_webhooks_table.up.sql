CREATE TABLE IF NOT EXISTS webhooks (
    id uuid PRIMARY KEY,
    name text NOT NULL,
    url text NOT NULL,
);