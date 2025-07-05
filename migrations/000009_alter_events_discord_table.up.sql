ALTER TABLE events
ALTER COLUMN channel_id TYPE text USING channel_id::text,
ALTER COLUMN guild_id TYPE text USING guild_id::text;
