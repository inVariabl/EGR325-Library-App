-- Start a transaction
START TRANSACTION;

-- Insert books
INSERT INTO book (isbn, title, author, publisher, publication_year, category, language, pages, shelf_location, status)
VALUES ('9780345391803', 'The Hitchhiker''s Guide to the Galaxy', 'Douglas Adams', 'Del Rey', 1995, 'Science Fiction', 'English', 224, 'SF-A1', 'available'),
       ('9780547928227', 'The Hobbit', 'J.R.R. Tolkien', 'Houghton Mifflin', 2012, 'Fantasy', 'English', 300, 'FA-T1', 'available'),
       ('9780062315007', 'The Alchemist', 'Paulo Coelho', 'HarperOne', 2014, 'Fiction', 'English', 208, 'FIC-C1', 'available'),
       ('9780132350884', 'Clean Code', 'Robert C. Martin', 'Prentice Hall', 2008, 'Technology', 'English', 464, 'TECH-M1', 'available'),
       ('9781449331818', 'Learning Python', 'Mark Lutz', 'O''Reilly Media', 2013, 'Technology', 'English', 1648, 'TECH-L1', 'available'),
       ('9780141439518', 'Pride and Prejudice', 'Jane Austen', 'Penguin Classics', 2002, 'Classic', 'English', 480, 'CL-A1', 'available'),
       ('9780061120084', 'To Kill a Mockingbird', 'Harper Lee', 'Harper Perennial', 2006, 'Fiction', 'English', 336, 'FIC-L1', 'available'),
       ('9780307474278', 'The Da Vinci Code', 'Dan Brown', 'Anchor', 2009, 'Mystery', 'English', 597, 'MYS-B1', 'available'),
       ('9780743273565', 'The Great Gatsby', 'F. Scott Fitzgerald', 'Scribner', 2004, 'Classic', 'English', 180, 'CL-F1', 'available'),
       ('9780553382563', 'A Game of Thrones', 'George R.R. Martin', 'Bantam', 2002, 'Fantasy', 'English', 864, 'FA-M1', 'available');

-- Insert corresponding inventory records
INSERT INTO book_inventory (book_id, available_copies, total_copies)
SELECT book_id,
       3,
       3
FROM book;

-- Commit the transaction
COMMIT;

-- Verify the insertions
SELECT b.*,
       bi.available_copies,
       bi.total_copies
FROM book b
LEFT JOIN book_inventory bi ON b.book_id = bi.book_id;
