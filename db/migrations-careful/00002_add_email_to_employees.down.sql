BEGIN TRANSACTION;

ALTER TABLE employees
RENAME COLUMN "email" TO "__email";

COMMIT;