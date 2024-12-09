DROP PROCEDURE IF EXISTS AddBooksWithInventory;

DELIMITER //

CREATE PROCEDURE AddBooksWithInventory(
    IN p_isbn1 VARCHAR(13),
    IN p_title1 VARCHAR(255),
    IN p_author1 VARCHAR(100),
    IN p_publisher1 VARCHAR(100),
    IN p_year1 INT,
    IN p_category1 VARCHAR(50),
    IN p_language1 VARCHAR(30),
    IN p_pages1 INT,
    IN p_shelf1 VARCHAR(50),
    IN p_isbn2 VARCHAR(13),
    IN p_title2 VARCHAR(255),
    IN p_author2 VARCHAR(100),
    IN p_publisher2 VARCHAR(100),
    IN p_year2 INT,
    IN p_category2 VARCHAR(50),
    IN p_language2 VARCHAR(30),
    IN p_pages2 INT,
    IN p_shelf2 VARCHAR(50),
    IN p_isbn3 VARCHAR(13),
    IN p_title3 VARCHAR(255),
    IN p_author3 VARCHAR(100),
    IN p_publisher3 VARCHAR(100),
    IN p_year3 INT,
    IN p_category3 VARCHAR(50),
    IN p_language3 VARCHAR(30),
    IN p_pages3 INT,
    IN p_shelf3 VARCHAR(50),
    IN p_copies INT
)
BEGIN
    DECLARE v_error INT DEFAULT 0;
    DECLARE CONTINUE HANDLER FOR SQLEXCEPTION SET v_error = 1;

    START TRANSACTION;

    -- Insert books
    INSERT INTO book (isbn, title, author, publisher, publication_year, category, language, pages, shelf_location, status)
    VALUES
        (p_isbn1, p_title1, p_author1, p_publisher1, p_year1, p_category1, p_language1, p_pages1, p_shelf1, 'available'),
        (p_isbn2, p_title2, p_author2, p_publisher2, p_year2, p_category2, p_language2, p_pages2, p_shelf2, 'available'),
        (p_isbn3, p_title3, p_author3, p_publisher3, p_year3, p_category3, p_language3, p_pages3, p_shelf3, 'available');

    -- Insert corresponding inventory records
    INSERT INTO book_inventory (book_id, available_copies, total_copies)
    SELECT book_id, p_copies, p_copies
    FROM book
    WHERE book_id IN (LAST_INSERT_ID(), LAST_INSERT_ID() + 1, LAST_INSERT_ID() + 2);

    IF v_error = 1 THEN
        ROLLBACK;
        SIGNAL SQLSTATE '45000'
            SET MESSAGE_TEXT = 'Error occurred while adding books and inventory';
    ELSE
        COMMIT;

        -- Return the inserted records
        SELECT b.*, bi.available_copies, bi.total_copies
        FROM book b
        LEFT JOIN book_inventory bi ON b.book_id = bi.book_id
        WHERE b.book_id IN (LAST_INSERT_ID(), LAST_INSERT_ID() + 1, LAST_INSERT_ID() + 2);
    END IF;
END //

DELIMITER ;

-- Example usage:
-- CALL AddBooksWithInventory(
--     '9780553380163', 'Snow Crash', 'Neal Stephenson', 'Bantam Books', 2000, 'Science Fiction', 'English', 470, 'SF-S1',
--     '9781492056355', 'Architecture Patterns with Python', 'Harry Percival', 'O''Reilly Media', 2020, 'Technology', 'English', 290, 'TECH-P1',
--     '9781449340377', 'JavaScript: The Good Parts', 'Douglas Crockford', 'O''Reilly Media', 2008, 'Technology', 'English', 176, 'TECH-C1',
--     3
-- );
