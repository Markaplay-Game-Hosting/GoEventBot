CREATE TABLE IF NOT EXISTS settings(
    id uuid PRIMARY KEY DEFAULT gen_random_uuid(),
    name text NOT NULL,
    description text NOT NULL,
    value text NOT NULL
);

INSERT INTO settings (name, description, value)
VALUES
    ('load_interval_time', 'time to load the event before it starts', '1m')