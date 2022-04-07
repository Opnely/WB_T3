CREATE TABLE employees.queue AS TABLE employees.employees;
ALTER TABLE employees.queue ALTER COLUMN id SET NOT NULL;
ALTER TABLE employees.queue ALTER COLUMN name SET NOT NULL;
ALTER TABLE employees.queue ALTER COLUMN last_name SET NOT NULL;
ALTER TABLE employees.queue ALTER COLUMN position SET NOT NULL;
ALTER TABLE employees.queue ALTER COLUMN good_job_count SET NOT NULL;
ALTER TABLE employees.queue ALTER COLUMN id SET DEFAULT nextval('employees.employees_id_seq'::regclass);


-- OK
CREATE OR REPLACE FUNCTION employees.get_ten_lowest_rows ()
RETURNS TABLE (id integer, last_name varchar, name varchar, patronymic varchar, phone varchar, posit varchar, good_job_count integer)  AS $$
BEGIN
    RETURN QUERY
    SELECT * FROM employees.queue ORDER BY id ASC LIMIT 10;
END;
$$ LANGUAGE 'plpgsql';


-- Ignores num, adds 10 lines
CREATE OR REPLACE FUNCTION employees.add_test_num(num integer DEFAULT 1) 
RETURNS void AS 
$$
DECLARE
    i integer;
BEGIN
    for i in 1..num LOOP
        INSERT INTO employees.queue (last_name, name, patronymic, phone, position, good_job_count)
        VALUES ('Testov', 'Test', 'Testovich', '84951112233', 'Tester', 1);
    END LOOP;
END;
$$ LANGUAGE 'plpgsql';


-- Delete 10 first rows
CREATE OR REPLACE FUNCTION employees.confirm (num integer)
RETURNS void AS
$$
BEGIN
    DELETE FROM employees.queue WHERE id IN 
    (SELECT id FROM employees.queue ORDER BY id ASC LIMIT num);
END;
$$ 
LANGUAGE 'plpgsql';
