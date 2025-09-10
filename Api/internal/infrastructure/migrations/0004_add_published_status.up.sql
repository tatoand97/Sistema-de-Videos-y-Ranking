-- Add new enum value 'PUBLISHED' to video_status and relax processed_at check
DO $$
BEGIN
    -- Add enum value if not exists
    IF NOT EXISTS (
        SELECT 1
        FROM pg_type t
        JOIN pg_enum e ON t.oid = e.enumtypid
        WHERE t.typname = 'video_status' AND e.enumlabel = 'PUBLISHED'
    ) THEN
        ALTER TYPE video_status ADD VALUE 'PUBLISHED';
    END IF;
END $$;

-- Update check constraint to allow processed_at when status is PROCESSED or PUBLISHED
ALTER TABLE video DROP CONSTRAINT IF EXISTS ck_processed_at;
ALTER TABLE video
    ADD CONSTRAINT ck_processed_at
    CHECK (processed_at IS NULL OR status IN ('PROCESSED','PUBLISHED'));

