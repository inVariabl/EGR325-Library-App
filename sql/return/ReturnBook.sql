DELIMITER //

DROP PROCEDURE IF EXISTS ReturnBook //

CREATE PROCEDURE ReturnBook(
    p_checkout_id INT,
    p_book_condition VARCHAR(50),
    p_notes VARCHAR(255)
)
BEGIN
    DECLARE v_book_id INT;
    DECLARE v_checkout_exists BOOLEAN;
    DECLARE v_already_returned BOOLEAN;

    -- Start transaction
    START TRANSACTION;

    -- Check if checkout exists and get book_id
    SELECT
        COUNT(*) > 0,
        book_id,
        return_date IS NOT NULL
    INTO
        v_checkout_exists,
        v_book_id,
        v_already_returned
    FROM checkout
    WHERE checkout_id = p_checkout_id;

    -- Proceed if checkout exists and hasn't been returned
    IF v_checkout_exists = 1 AND v_already_returned = 0 THEN
        -- Create return record
        INSERT INTO book_return (
            checkout_id,
            return_date,
            book_condition,
            notes
        ) VALUES (
            p_checkout_id,
            CURRENT_DATE,
            p_book_condition,
            p_notes
        );

        -- Update checkout record
        UPDATE checkout
        SET return_date = CURRENT_DATE
        WHERE checkout_id = p_checkout_id;

        -- Update book inventory (increase available copies) if book isn't lost
        IF p_book_condition != 'lost' THEN
            UPDATE book_inventory
            SET available_copies = available_copies + 1
            WHERE book_id = v_book_id;
        END IF;

        -- If everything succeeded, commit the transaction
        COMMIT;

        -- Return success message
        SELECT 'Book returned successfully' as message;
    ELSE
        -- If there was a problem, rollback
        ROLLBACK;

        IF NOT v_checkout_exists THEN
            SELECT 'Checkout record not found' as message;
        ELSE
            SELECT 'Book has already been returned' as message;
        END IF;
    END IF;
END //

DELIMITER ;
