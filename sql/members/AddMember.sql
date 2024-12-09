DELIMITER //

DROP PROCEDURE IF EXISTS AddMember //

CREATE PROCEDURE AddMember(
    p_name VARCHAR(100),
    p_email VARCHAR(100),
    p_phone VARCHAR(20),
    p_address VARCHAR(255)
)
BEGIN
    DECLARE v_member_id INT;
    DECLARE v_email_exists BOOLEAN;

    -- Start transaction
    START TRANSACTION;

    -- Check if email already exists
    SELECT COUNT(*) > 0 INTO v_email_exists
    FROM member
    WHERE email = p_email;

    IF v_email_exists = 0 THEN
        -- Insert new member
        INSERT INTO member (
            name,
            email,
            phone_number,
            address,
            membership_date
        ) VALUES (
            p_name,
            p_email,
            p_phone,
            p_address,
            CURRENT_DATE
        );

        SET v_member_id = LAST_INSERT_ID();

        COMMIT;

        -- Return success with member ID
        SELECT
            'Member added successfully' as message,
            v_member_id as member_id;
    ELSE
        ROLLBACK;
        SELECT 'Email address already exists' as message;
    END IF;
END //

DELIMITER ;
