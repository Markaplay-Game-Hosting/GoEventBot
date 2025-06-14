CREATE TABLE IF NOT EXISTS jobs (
    id uuid PRIMARY KEY DEFAULT gen_random_uuid(),
    event_id uuid NOT NULL REFERENCES events(id) ON DELETE CASCADE,
    execution_date timestamp(0) with time zone NOT NULL,
    status INTEGER NOT NULL DEFAULT 0
);