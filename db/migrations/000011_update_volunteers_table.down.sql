DROP TABLE IF EXISTS volunteers;

CREATE TABLE volunteers (
    id VARCHAR(36) PRIMARY KEY,
    full_name VARCHAR(255) NOT NULL,
    email VARCHAR(255) NOT NULL UNIQUE,
    phone_number VARCHAR(20) NOT NULL,
    birth_date DATE NOT NULL,
    gender ENUM('Male', 'Female') NOT NULL,
    domicile VARCHAR(255) NOT NULL,
    status ENUM('Student', 'High School', 'Employee') NOT NULL,
    interest VARCHAR(255) NOT NULL,
    certificate_name VARCHAR(255) NOT NULL,
    created_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP,
    updated_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP ON UPDATE CURRENT_TIMESTAMP
);
