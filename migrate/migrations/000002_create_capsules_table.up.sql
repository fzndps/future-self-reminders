CREATE TABLE IF NOT EXISTS capsules (
    id INT AUTO_INCREMENT PRIMARY KEY,
    user_id INT NOT NULL,
    title VARCHAR(200) NOT NULL,
    message TEXT NOT NULL,
    due_date DATE NOT NULL,
    delivery_method ENUM('email', 'in_app') DEFAULT 'email',
    status ENUM('pending', 'sent', 'cancelled') DEFAULT 'pending',
    category VARCHAR(50),
    mood VARCHAR(50),
    image_url VARCHAR(255),
    sent_at TIMESTAMP NULL,
    created_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP,
    updated_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP ON UPDATE CURRENT_TIMESTAMP,
    FOREIGN KEY (user_id) REFERENCES users(id) ON DELETE CASCADE
);
