BEGIN TRANSACTION;

DO $$
BEGIN
    IF EXISTS(
        SELECT *
        FROM information_schema.columns
        WHERE table_name='employees' and column_name='__email'
    )
    THEN
        ALTER TABLE employees
        RENAME COLUMN "__email" TO "email";
    ELSE
        ALTER TABLE employees
        ADD COLUMN email VARCHAR(200);
    END IF;
END $$;

COMMIT;