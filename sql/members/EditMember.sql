DELIMITER //

DROP PROCEDURE IF EXISTS EditMember //

CREATE PROCEDURE EditMember(
    p_member_id INT,
    p_name VARCHAR(100),
    p_email VARCHAR(100),
    p_phone VARCHAR(20),
    p_address VARCHAR(255)
)
BEGIN
    DECLARE v_member_exists BOOLEAN;
    DECLARE v_email_exists BOOLEAN;

    -- Start transaction
    START TRANSACTION;

    -- Check if member exists
    SELECT COUNT(*) > 0 INTO v_member_exists
    FROM member
    WHERE member_id = p_member_id;

    -- Check if email exists for another member
    SELECT COUNT(*) > 0 INTO v_email_exists
    FROM member
    WHERE email = p_email
    AND member_id != p_member_id;

    IF v_member_exists = 0 THEN
        SELECT 'Member not found' as message;
    ELSEIF v_email_exists = 1 THEN
        SELECT 'Email address already in use by another member' as message;
    ELSE
        -- Update member
        UPDATE member
        SET
            name = p_name,
            email = p_email,
            phone_number = p_phone,
            address = p_address
        WHERE member_id = p_member_id;

        COMMIT;

        -- Return success
        SELECT 'Member updated successfully' as message;
    END IF;

    IF v_member_exists = 0 OR v_email_exists = 1 THEN
        ROLLBACK;
    END IF;
END //

DELIMITER ;
