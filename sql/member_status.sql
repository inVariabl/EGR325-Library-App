CREATE TABLE IF NOT EXISTS member_status (
	status_id INT PRIMARY KEY AUTO_INCREMENT,
	member_id INT NOT NULL,
	status ENUM('active', 'suspended', 'expired') DEFAULT 'active',
	changed_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP,
	changed_by INT NOT NULL,
	reason TEXT,
	FOREIGN KEY (member_id) REFERENCES member(member_id),
	FOREIGN KEY (changed_by) REFERENCES admin(admin_id)
);
