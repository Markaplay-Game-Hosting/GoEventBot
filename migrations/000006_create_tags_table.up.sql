CREATE TABLE IF NOT EXISTS tags (
    id uuid PRIMARY KEY DEFAULT gen_random_uuid(),
    name text NOT NULL,
    description text NULL,
    created_date timestamp(0) with time zone NOT NULL DEFAULT NOW(),
    updated_date timestamp(0) with time zone NOT NULL DEFAULT NOW()
);

CREATE TABLE IF NOT EXISTS event_tags (
    event_id uuid NOT NULL REFERENCES events ON DELETE CASCADE,
    tag_id uuid NOT NULL REFERENCES tags ON DELETE CASCADE,
    PRIMARY KEY (event_id, tag_id)
);