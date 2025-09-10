-- Drop unique constraint and column
ALTER TABLE vote DROP CONSTRAINT IF EXISTS ux_vote_event;
ALTER TABLE vote DROP COLUMN IF EXISTS event_id;

