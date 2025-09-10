-- Revert constraint to only allow processed_at when status is PROCESSED
ALTER TABLE video DROP CONSTRAINT IF EXISTS ck_processed_at;
ALTER TABLE video
    ADD CONSTRAINT ck_processed_at
    CHECK (processed_at IS NULL OR status = 'PROCESSED');

-- Optional: convert any PUBLISHED statuses back to PROCESSED for compatibility
UPDATE video SET status = 'PROCESSED' WHERE status = 'PUBLISHED';

-- Note: PostgreSQL does not support dropping an enum value easily; we keep 'PUBLISHED' value present.

