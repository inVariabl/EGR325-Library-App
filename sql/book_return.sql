CREATE TABLE book_return (
    return_id INT PRIMARY KEY AUTO_INCREMENT,
    checkout_id INT NOT NULL,
    return_date DATE NOT NULL,
    book_condition ENUM('good', 'damaged', 'lost') NOT NULL,
    notes TEXT,
    created_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP,
    FOREIGN KEY (checkout_id) REFERENCES checkout(checkout_id)
);
