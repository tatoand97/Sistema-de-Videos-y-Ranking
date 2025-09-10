-- Revert constraint to only allow processed_at when status is PROCESSED
ALTER TABLE video DROP CONSTRAINT IF EXISTS ck_processed_at;
ALTER TABLE video
    ADD CONSTRAINT ck_processed_at
    CHECK (processed_at IS NULL OR status = 'PROCESSED');

