CREATE TABLE admin (
    admin_id INT PRIMARY KEY AUTO_INCREMENT,          -- Unique identifier for each admin
    username VARCHAR(50) NOT NULL UNIQUE,             -- Username for admin login, must be unique
    password_hash VARCHAR(255) NOT NULL,              -- Hashed password for security
    email VARCHAR(100) NOT NULL UNIQUE,               -- Contact email, must be unique
    first_name VARCHAR(50) NOT NULL,                  -- Admin's first name
    last_name VARCHAR(50) NOT NULL,                   -- Admin's last name
    phone_number VARCHAR(20),                         -- Optional phone number
    role ENUM('superadmin', 'manager', 'staff')       -- Role of the admin (e.g., superadmin, manager, staff)
        DEFAULT 'staff',
    created_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP,   -- Timestamp for account creation
    updated_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP ON UPDATE CURRENT_TIMESTAMP,  -- Timestamp for last update
    last_login TIMESTAMP                              -- Timestamp for last login, initially NULL
);
CREATE TABLE book (
    book_id INT PRIMARY KEY AUTO_INCREMENT,
    isbn VARCHAR(20) UNIQUE NOT NULL,
    title VARCHAR(255) NOT NULL,
    author VARCHAR(100),
    publisher VARCHAR(100),
    publication_year YEAR,
    category VARCHAR(50),
    language VARCHAR(30) DEFAULT 'English',
    pages INT,
    shelf_location VARCHAR(50),
    status ENUM('available', 'checked_out') DEFAULT 'available',  -- Indicates if book is available
    created_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP,
    updated_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP ON UPDATE CURRENT_TIMESTAMP
);
CREATE TABLE book_inventory (
    inventory_id INT PRIMARY KEY AUTO_INCREMENT,      -- Unique identifier for each inventory record
    book_id INT NOT NULL,                             -- Foreign key referencing the book table
    available_copies INT DEFAULT 1,                   -- Number of currently available copies
    total_copies INT DEFAULT 1,                       -- Total number of copies owned
    FOREIGN KEY (book_id) REFERENCES book(book_id)    -- Ensures data integrity
        ON DELETE CASCADE                              -- Deletes inventory record if book is deleted
);
CREATE TABLE member (
    member_id INT PRIMARY KEY AUTO_INCREMENT,
    name VARCHAR(100) NOT NULL,
    email VARCHAR(100) UNIQUE NOT NULL,
    phone_number VARCHAR(20),
    membership_date DATE DEFAULT CURRENT_DATE,
    address VARCHAR(255)
);

DELIMITER $$

CREATE PROCEDURE UpdateBookDetails(
    IN p_book_id INT,
    IN p_isbn VARCHAR(20),
    IN p_title VARCHAR(255),
    IN p_author VARCHAR(100),
    IN p_publisher VARCHAR(100),
    IN p_publication_year YEAR,
    IN p_category VARCHAR(50),
    IN p_language VARCHAR(30),
    IN p_pages INT,
    IN p_available_copies INT,
    IN p_total_copies INT,
    IN p_shelf_location VARCHAR(50),
    IN p_status ENUM('available', 'checked_out')
)

BEGIN
    UPDATE book
    SET
        isbn = p_isbn,
        title = p_title,
        author = p_author,
        publisher = p_publisher,
        publication_year = p_publication_year,
        category = p_category,
        language = p_language,
        pages = p_pages,
        available_copies = p_available_copies,
        total_copies = p_total_copies,
        shelf_location = p_shelf_location,
        status = p_status,
        updated_at = CURRENT_TIMESTAMP
    WHERE
        book_id = p_book_id;
END$$

-- Stored Procedure to Remove a Book
CREATE PROCEDURE RemoveBook(
    IN p_book_id INT
)
BEGIN
    -- Optional: Check if the book exists
    IF EXISTS (SELECT 1 FROM book WHERE book_id = p_book_id) THEN
        DELETE FROM book WHERE book_id = p_book_id;
    ELSE
        SIGNAL SQLSTATE '45000' SET MESSAGE_TEXT = 'Book ID does not exist.';
    END IF;
END$$

-- Reset the delimiter to default
DELIMITER ;
DELIMITER $$

CREATE PROCEDURE UpdateInventory(
    IN p_book_id INT,
    IN p_available_copies INT,
    IN p_total_copies INT
)
BEGIN
    UPDATE book_inventory
    SET
        available_copies = p_available_copies,
        total_copies = p_total_copies
    WHERE
        book_id = p_book_id;
END$$

DELIMITER ;