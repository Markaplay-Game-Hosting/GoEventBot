DROP TABLE IF EXISTS webhooks;

ALTER TABLE events 
    DROP CONSTRAINT IF EXISTS fk_events_webhook_id;