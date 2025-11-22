-- Drop and recreate donations table with new structure
DROP TABLE IF EXISTS donations;

CREATE TABLE donations (
    id VARCHAR(36) PRIMARY KEY,
    user_id VARCHAR(36) NOT NULL,
    program_id VARCHAR(36) NOT NULL,
    amount DECIMAL(15, 2) NOT NULL,
    payment_method VARCHAR(50) NOT NULL COMMENT 'e.g., bank_transfer, ewallet, credit_card',
    status ENUM('pending', 'paid', 'failed') NOT NULL DEFAULT 'pending',
    proof_image TEXT COMMENT 'Optional URL for payment proof',
    created_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP,
    updated_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP ON UPDATE CURRENT_TIMESTAMP,
    FOREIGN KEY (user_id) REFERENCES users(id) ON DELETE CASCADE,
    FOREIGN KEY (program_id) REFERENCES programs(id) ON DELETE CASCADE,
    INDEX idx_user_id (user_id),
    INDEX idx_program_id (program_id),
    INDEX idx_status (status),
    INDEX idx_created_at (created_at)
);
