DELIMITER //

DROP PROCEDURE IF EXISTS CheckoutBook //

CREATE PROCEDURE CheckoutBook(
    p_book_id INT,
    p_member_id INT,
    p_notes VARCHAR(255)
)
BEGIN
    DECLARE v_available_copies INT DEFAULT NULL;
    DECLARE v_member_exists BOOLEAN;
    DECLARE v_checkout_id INT;

    -- Start transaction
    START TRANSACTION;

    -- Check if book is available
    SELECT available_copies INTO v_available_copies
    FROM book_inventory
    WHERE book_id = p_book_id
    FOR UPDATE;

    -- Simple member verification
    SELECT COUNT(*) > 0 INTO v_member_exists
    FROM member
    WHERE member_id = p_member_id;

    -- Proceed if book is available and member exists
    IF v_available_copies > 0 AND v_member_exists = 1 THEN
        -- Create checkout record
        INSERT INTO checkout (
            book_id,
            member_id,
            checkout_date,
            due_date,
            notes
        ) VALUES (
            p_book_id,
            p_member_id,
            CURRENT_TIMESTAMP,
            DATE_ADD(CURRENT_DATE, INTERVAL 14 DAY),
            p_notes
        );

        SET v_checkout_id = LAST_INSERT_ID();

        -- Update book inventory
        UPDATE book_inventory
        SET available_copies = available_copies - 1
        WHERE book_id = p_book_id;

        -- If everything succeeded, commit the transaction
        COMMIT;

        -- Return single success message
        SELECT 'Checkout successful' as message, v_checkout_id as checkout_id;
    ELSE
        -- If there was a problem, rollback
        ROLLBACK;

        IF v_available_copies IS NULL OR v_available_copies <= 0 THEN
            SELECT 'Book is not available' as message;
        ELSE
            SELECT 'Member not found' as message;
        END IF;
    END IF;
END //

DELIMITER ;
