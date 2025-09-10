-- Add new enum value 'PUBLISHED' to video_status
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
