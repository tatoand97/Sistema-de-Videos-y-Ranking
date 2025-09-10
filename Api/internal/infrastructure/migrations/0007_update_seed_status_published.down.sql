-- Revert promotion back to PROCESSED for the same seeded videos
UPDATE video
SET status = 'PROCESSED'
WHERE status = 'PUBLISHED'
  AND original_file IN (
    'juan_video_001.mp4',
    'maria_video_001.mp4',
    'ana_video_001.mp4'
  );

