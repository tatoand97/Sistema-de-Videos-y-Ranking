-- Promote specific seeded videos from PROCESSED to PUBLISHED
-- This runs after 0004 (enum value added) and 0006 (constraint relaxed)
UPDATE video
SET status = 'PUBLISHED'
WHERE status = 'PROCESSED'
  AND original_file IN (
    'juan_video_001.mp4',
    'maria_video_001.mp4',
    'ana_video_001.mp4'
  );

