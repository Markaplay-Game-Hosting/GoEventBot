CREATE TABLE IF NOT EXISTS events(
    id uuid PRIMARY KEY DEFAULT gen_random_uuid(),
    title text NOT NULL,
    description text NOT NULL,
    start_date timestamp(0) with time zone NOT NULL,
    end_date timestamp(0) with time zone NULL,
    duration interval NOT NULL,
    repeat_interval interval NULL,
    is_active bool NOT NULL DEFAULT true,
    webhook_id uuid NOT NULL REFERENCES webhooks ON DELETE CASCADE,
    created_date timestamp(0) with time zone NOT NULL DEFAULT NOW(),
    updated_date timestamp(0) with time zone NOT NULL DEFAULT NOW()
)