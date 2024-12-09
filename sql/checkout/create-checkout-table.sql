-- Create checkout table if it doesn't exist
CREATE TABLE IF NOT EXISTS checkout (
    checkout_id INT PRIMARY KEY AUTO_INCREMENT,
    book_id INT NOT NULL,
    member_id INT NOT NULL,
    checkout_date TIMESTAMP DEFAULT CURRENT_TIMESTAMP,
    due_date DATE NOT NULL,
    return_date TIMESTAMP NULL,
    notes TEXT,
    FOREIGN KEY (book_id) REFERENCES book(book_id),
    FOREIGN KEY (member_id) REFERENCES member(member_id)
);
