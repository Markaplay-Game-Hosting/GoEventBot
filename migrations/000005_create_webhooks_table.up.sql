CREATE TABLE IF NOT EXISTS webhooks (
    id uuid PRIMARY KEY,
    name text NOT NULL,
    url text NOT NULL
);

ALTER TABLE events 
    ADD CONSTRAINT fk_events_webhook_id FOREIGN KEY(webhook_id) REFERENCES webhooks(id) ON DELETE CASCADE;