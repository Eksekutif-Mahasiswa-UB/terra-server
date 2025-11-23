ALTER TABLE events 
ADD COLUMN quota INT NOT NULL DEFAULT 50 AFTER location,
MODIFY COLUMN image_url TEXT NOT NULL;

CREATE TABLE event_participants (
    id VARCHAR(73) PRIMARY KEY COMMENT 'Composite key: userID-eventID',
    user_id VARCHAR(36) NOT NULL,
    event_id VARCHAR(36) NOT NULL,
    joined_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP,
    UNIQUE KEY unique_user_event (user_id, event_id),
    FOREIGN KEY (user_id) REFERENCES users(id) ON DELETE CASCADE,
    FOREIGN KEY (event_id) REFERENCES events(id) ON DELETE CASCADE,
    INDEX idx_user_id (user_id),
    INDEX idx_event_id (event_id)
);
