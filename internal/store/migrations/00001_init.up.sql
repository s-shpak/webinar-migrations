BEGIN TRANSACTION;

CREATE TABLE positions(
    id INT PRIMARY KEY GENERATED ALWAYS AS IDENTITY,
    title VARCHAR(200) UNIQUE NOT NULL
);

CREATE TABLE employees(
    id INT PRIMARY KEY GENERATED ALWAYS AS IDENTITY,
    first_name VARCHAR(200) NOT NULL,
    last_name VARCHAR(200) NOT NULL,
    salary NUMERIC NOT NULL,
    position INT NOT NULL REFERENCES positions(id),
    CONSTRAINT employees_salary_positive_check CHECK (salary::numeric > 0)
);

COMMIT;