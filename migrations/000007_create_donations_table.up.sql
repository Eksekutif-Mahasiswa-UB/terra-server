CREATE TABLE donations (
    id VARCHAR(36) PRIMARY KEY,
    order_id VARCHAR(255) NOT NULL UNIQUE,
    user_id VARCHAR(36) NULL,
    amount BIGINT NOT NULL,
    status ENUM('pending', 'success', 'failed') NOT NULL DEFAULT 'pending',
    payment_gateway VARCHAR(50) NOT NULL,
    created_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP,
    updated_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP ON UPDATE CURRENT_TIMESTAMP,
    FOREIGN KEY (user_id) REFERENCES users(id) ON DELETE SET NULL
);