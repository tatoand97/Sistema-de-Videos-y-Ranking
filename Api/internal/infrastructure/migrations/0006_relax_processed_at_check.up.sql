-- Update check constraint to allow processed_at when status is PROCESSED or PUBLISHED
ALTER TABLE video DROP CONSTRAINT IF EXISTS ck_processed_at;
ALTER TABLE video
    ADD CONSTRAINT ck_processed_at
    CHECK (processed_at IS NULL OR status IN ('PROCESSED','PUBLISHED'));

