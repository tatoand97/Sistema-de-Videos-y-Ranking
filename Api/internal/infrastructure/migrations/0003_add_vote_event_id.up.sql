-- Add event_id to vote table and unique constraint for idempotency
ALTER TABLE vote ADD COLUMN IF NOT EXISTS event_id TEXT;
-- Name constraint to detect conflicts in application code
DO $$
BEGIN
    IF NOT EXISTS (
        SELECT 1 FROM pg_constraint WHERE conname = 'ux_vote_event'
    ) THEN
        ALTER TABLE vote ADD CONSTRAINT ux_vote_event UNIQUE (event_id);
    END IF;
END $$;

