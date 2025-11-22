DROP TABLE IF EXISTS event_participants;

ALTER TABLE events 
DROP COLUMN quota,
MODIFY COLUMN image_url TEXT NULL;
