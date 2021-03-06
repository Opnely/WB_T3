toc.dat                                                                                             0000600 0004000 0002000 00000034321 13636322476 0014456 0                                                                                                    ustar 00postgres                        postgres                        0000000 0000000                                                                                                                                                                        PGDMP               	            x         
   postgresdb    10.4 (Debian 10.4-2.pgdg90+1)    12.0     4           0    0    ENCODING    ENCODING        SET client_encoding = 'UTF8';
                      false         5           0    0 
   STDSTRINGS 
   STDSTRINGS     (   SET standard_conforming_strings = 'on';
                      false         6           0    0 
   SEARCHPATH 
   SEARCHPATH     8   SELECT pg_catalog.set_config('search_path', '', false);
                      false         7           1262    16384 
   postgresdb    DATABASE     z   CREATE DATABASE postgresdb WITH TEMPLATE = template0 ENCODING = 'UTF8' LC_COLLATE = 'en_US.utf8' LC_CTYPE = 'en_US.utf8';
    DROP DATABASE postgresdb;
                postgres    false                     2615    16386 	   employees    SCHEMA        CREATE SCHEMA employees;
    DROP SCHEMA employees;
                postgresadmin    false                     2615    16387    test    SCHEMA        CREATE SCHEMA test;
    DROP SCHEMA test;
                postgresadmin    false         ?            1255    16411 t   employee_add(character varying, character varying, character varying, character varying, character varying, integer)    FUNCTION     ?  CREATE FUNCTION employees.employee_add(_name character varying, _last_name character varying, _patronymic character varying, _phone character varying, _position character varying, _good_job_count integer) RETURNS void
    LANGUAGE plpgsql SECURITY DEFINER
    AS $$
BEGIN

   INSERT INTO employees.employees 
   (
	
     name,
	last_name,
	patronymic,
	phone,
	position,
	good_job_count
   )
   VALUES
   (
	  
      _name,
	_last_name,
	_patronymic,
	_phone,
	_position,
	_good_job_count 
   );
 
END;
$$;
 ?   DROP FUNCTION employees.employee_add(_name character varying, _last_name character varying, _patronymic character varying, _phone character varying, _position character varying, _good_job_count integer);
    	   employees          postgresadmin    false    4         ?            1255    16389    employee_get(integer)    FUNCTION     ?  CREATE FUNCTION employees.employee_get(_id integer) RETURNS TABLE(name character varying, last_name character varying, id integer, patronymic character varying, phone character varying, _position character varying, _good_job_count integer)
    LANGUAGE plpgsql SECURITY DEFINER
    AS $$
BEGIN
     IF NOT EXISTS (SELECT 1 FROM employees.employees m WHERE m.id = _id)
    THEN 
        RAISE EXCEPTION 'Сотрудника с таким id не существует' USING ERRCODE = 50003;
    END IF;
    RETURN QUERY
    SELECT e.name,
           e.last_name,
           e.id,
           e.patronymic,
           e.phone,
		   e.position,
		   e.good_job_count
	  
         
    FROM employees.employees e
	where e.id = _id	;

END;
$$;
 3   DROP FUNCTION employees.employee_get(_id integer);
    	   employees          postgresadmin    false    4         ?            1255    16390    employee_get_all()    FUNCTION     ?  CREATE FUNCTION employees.employee_get_all() RETURNS TABLE(name character varying, last_name character varying, id integer, patronymic character varying, phone character varying, _position character varying, _good_job_count integer)
    LANGUAGE plpgsql SECURITY DEFINER
    AS $$
BEGIN

    RETURN QUERY
    SELECT e.name,
           e.last_name,
           e.id,
           e.patronymic,
           e.phone,
		   e.position,
		   e.good_job_count
          
    FROM employees.employees e;

END;
$$;
 ,   DROP FUNCTION employees.employee_get_all();
    	   employees          postgresadmin    false    4         ?            1255    16391    employee_remove(integer)    FUNCTION     ?  CREATE FUNCTION employees.employee_remove(_id integer) RETURNS void
    LANGUAGE plpgsql SECURITY DEFINER
    AS $$
BEGIN
  IF NOT EXISTS (SELECT 1 FROM employees.employees m WHERE m.id = _id)
    THEN 
        RAISE EXCEPTION 'Сотрудника с таким id не существует' USING ERRCODE = 50003;
    END IF;

   DELETE from employees.employees as e
   where e.id = _id;

END;
$$;
 6   DROP FUNCTION employees.employee_remove(_id integer);
    	   employees          postgresadmin    false    4         ?            1255    16392 }   employee_upd(integer, character varying, character varying, character varying, character varying, character varying, integer)    FUNCTION     ?  CREATE FUNCTION employees.employee_upd(_id integer, _name character varying, _last_name character varying, _patronymic character varying, _phone character varying, _position character varying, _good_job_count integer) RETURNS void
    LANGUAGE plpgsql SECURITY DEFINER
    AS $$
BEGIN
 IF NOT EXISTS (SELECT 1 FROM employees.employees m WHERE m.id = _id)
    THEN 
        RAISE EXCEPTION 'Сотрудника с таким id не существует' USING ERRCODE = 50003;
    END IF;

UPDATE employees.employees
	SET id=_id, last_name=_last_name, name=_name, patronymic=_patronymic, phone=_phone, position=_position, good_job_count=_good_job_count
	WHERE id = _id;
END;
$$;
 ?   DROP FUNCTION employees.employee_upd(_id integer, _name character varying, _last_name character varying, _patronymic character varying, _phone character varying, _position character varying, _good_job_count integer);
    	   employees          postgresadmin    false    4         ?            1255    16414    employees_get_all()    FUNCTION     ?  CREATE FUNCTION employees.employees_get_all() RETURNS TABLE(name character varying, last_name character varying, id integer, patronymic character varying, phone character varying, _position character varying, _good_job_count integer)
    LANGUAGE plpgsql SECURITY DEFINER
    AS $$
BEGIN

    RETURN QUERY
    SELECT e.name,
           e.last_name,
           e.id,
           e.patronymic,
           e.phone,
		   e.position,
		   e.good_job_count
          
    FROM employees.employees e;

END;
$$;
 -   DROP FUNCTION employees.employees_get_all();
    	   employees          postgresadmin    false    4         ?            1255    16412    employees_get_all_part1()    FUNCTION     ;  CREATE FUNCTION employees.employees_get_all_part1() RETURNS TABLE(name character varying, last_name character varying, id integer)
    LANGUAGE plpgsql SECURITY DEFINER
    AS $$
BEGIN

    RETURN QUERY
    SELECT e.name,
           e.last_name,
           e.id
          
    FROM employees.employees e;

END;
$$;
 3   DROP FUNCTION employees.employees_get_all_part1();
    	   employees          postgresadmin    false    4         ?            1255    16413    employees_get_all_part2()    FUNCTION     ?  CREATE FUNCTION employees.employees_get_all_part2() RETURNS TABLE(id integer, patronymic character varying, phone character varying, _position character varying, _good_job_count integer)
    LANGUAGE plpgsql SECURITY DEFINER
    AS $$
BEGIN

    RETURN QUERY
    SELECT 
           e.id,
           e.patronymic,
           e.phone,
		   e.position,
		   e.good_job_count
          
    FROM employees.employees e;

END;
$$;
 3   DROP FUNCTION employees.employees_get_all_part2();
    	   employees          postgresadmin    false    4         ?            1255    16393 	   get_all()    FUNCTION     ?  CREATE FUNCTION employees.get_all() RETURNS TABLE(name character varying, last_name character varying, id integer, patronymic character varying, phone character varying, _position character varying, _good_job_count integer)
    LANGUAGE plpgsql SECURITY DEFINER
    AS $$
BEGIN

    RETURN QUERY
    SELECT e.name,
           e.last_name,
           e.id,
           e.patronymic,
           e.phone,
		   e.position,
		   e.good_job_count
          
    FROM employees.employees e;

END;
$$;
 #   DROP FUNCTION employees.get_all();
    	   employees          postgresadmin    false    4         ?            1255    16394    db_error_test()    FUNCTION     ?  CREATE FUNCTION public.db_error_test() RETURNS void
    LANGUAGE plpgsql SECURITY DEFINER
    AS $$
DECLARE _query text;
DECLARE _place_error text;
BEGIN
    CREATE LOCAL  TEMPORARY TABLE script( ex text ) ON COMMIT DROP;
	
	INSERT INTO script
	(
		ex
	)
    WITH cte AS
    (
    SELECT  n.nspname, 
            p.proname,
    		p.proargmodes, 
    		CASE WHEN p.proallargtypes IS NULL 
    		THEN  CASE WHEN array_length(string_to_array(proargtypes::text,' ')::oid[],1) > 1 THEN string_to_array(proargtypes::text,' ')::oid[] 
    	   ELSE NULL END
    		ELSE p.proallargtypes END args, 
    		p.proargnames
    FROM    pg_catalog.pg_namespace n
       INNER JOIN pg_catalog.pg_proc p
           ON p.pronamespace = n.oid
	WHERE nspowner != 10
    ),
    cte_2 AS
    (
    SELECT initcap(cte.nspname) nspname,
           cte.proname,
    	   unnest(cte.proargmodes) param_type, 
    	   format_type(unnest(cte.args), null) arg_type, 
    	   unnest(proargnames) arg_name
    FROM cte
    ),
    cte_3 AS
    (
    SELECT CASE WHEN param_type = 'i' or param_type IS NULL THEN 'input' ELSE 'output' END param_type,
    string_agg( CASE WHEN param_type = 'i' or param_type IS NULL THEN cte_2.arg_name||':= null,' END,' ' ) arg_name,
    cte_2.nspname,
    cte_2.proname,
    COUNT(proname) OVER(PARTITION BY proname) counts
    FROM cte_2
    GROUP BY nspname, proname, CASE WHEN param_type = 'i' or param_type IS NULL THEN 'input' ELSE 'output' END
    )
    SELECT DISTINCT 'SELECT '||nspname||'.'||proname||'('||CASE WHEN arg_name IS NULL THEN '' ELSE left(arg_name, (length(arg_name)-1)) END ||');' ex
    FROM cte_3
    WHERE param_type = 'input' or counts = 1;
    
	FOR _query IN (SELECT ex FROM script) LOOP 
	    BEGIN
	    EXECUTE _query;
		EXCEPTION
            WHEN not_null_violation or SQLSTATE '50003' or SQLSTATE '50001' or SQLSTATE '50002' or SQLSTATE '50000' THEN 
			RAISE NOTICE 'Пропущена бизнесовая ошибка или not null.
';
	        WHEN others THEN
	         GET STACKED DIAGNOSTICS _place_error = PG_EXCEPTION_CONTEXT;
	         _place_error = 'Объект: '|| _place_error||' С ошибкой: '||SQLERRM||'
			 ';
	         RAISE NOTICE '%', _place_error;
	    END;
	END LOOP;
	
	RAISE NOTICE 'Ошибок Нет!
';
	
END;
$$;
 &   DROP FUNCTION public.db_error_test();
       public          postgresadmin    false         ?            1255    16395 9   do_something(integer, jsonb, timestamp without time zone)    FUNCTION     =  CREATE FUNCTION test.do_something(_employee_id integer, _test_data jsonb, _date timestamp without time zone) RETURNS void
    LANGUAGE plpgsql SECURITY DEFINER
    AS $$
BEGIN

 IF _employee_id IS NULL 
 THEN
   RAISE EXCEPTION '_employee_id не должен быть null' USING ERRCODE = 50000;
 END IF;
  IF _test_data IS NULL 
 THEN
   RAISE EXCEPTION '_test_data не должен быть null' USING ERRCODE = 50000;
    END IF;
  IF _date IS NULL 
 THEN
   RAISE EXCEPTION '_date не должен быть null' USING ERRCODE = 50000;
    END IF;	
   


END;
$$;
 l   DROP FUNCTION test.do_something(_employee_id integer, _test_data jsonb, _date timestamp without time zone);
       test          postgresadmin    false    5         ?            1255    16396    get_db_error(integer)    FUNCTION     ?  CREATE FUNCTION test.get_db_error(_id integer) RETURNS boolean
    LANGUAGE plpgsql SECURITY DEFINER
    AS $$
BEGIN

    IF _id = 1
    THEN
        RAISE EXCEPTION 'Пользовательская ошибка' USING ERRCODE = 50000;
	END IF;	
	IF _id = 2
    THEN
        RAISE EXCEPTION 'Внутренняя ошибка' USING ERRCODE = 10000;	
    ELSE
	   
       RETURN TRUE;
    END IF;

END;
$$;
 .   DROP FUNCTION test.get_db_error(_id integer);
       test          postgresadmin    false    5         ?            1255    16397    long_time_process()    FUNCTION     ?   CREATE FUNCTION test.long_time_process() RETURNS void
    LANGUAGE plpgsql SECURITY DEFINER
    AS $$
BEGIN

  select pg_sleep(500);

END;
$$;
 (   DROP FUNCTION test.long_time_process();
       test          postgresadmin    false    5         ?            1259    16405 	   employees    TABLE     2  CREATE TABLE employees.employees (
    id integer NOT NULL,
    last_name character varying(100) NOT NULL,
    name character varying(100) NOT NULL,
    patronymic character varying(15),
    phone character varying(15),
    "position" character varying(20) NOT NULL,
    good_job_count integer NOT NULL
);
     DROP TABLE employees.employees;
    	   employees            postgresadmin    false    4         ?            1259    16403    employees_id_seq    SEQUENCE     ?   CREATE SEQUENCE employees.employees_id_seq
    AS integer
    START WITH 1
    INCREMENT BY 1
    NO MINVALUE
    NO MAXVALUE
    CACHE 1;
 *   DROP SEQUENCE employees.employees_id_seq;
    	   employees          postgresadmin    false    4    199         8           0    0    employees_id_seq    SEQUENCE OWNED BY     K   ALTER SEQUENCE employees.employees_id_seq OWNED BY employees.employees.id;
       	   employees          postgresadmin    false    198         ?
           2604    16408    employees id    DEFAULT     r   ALTER TABLE ONLY employees.employees ALTER COLUMN id SET DEFAULT nextval('employees.employees_id_seq'::regclass);
 >   ALTER TABLE employees.employees ALTER COLUMN id DROP DEFAULT;
    	   employees          postgresadmin    false    199    198    199         1          0    16405 	   employees 
   TABLE DATA           j   COPY employees.employees (id, last_name, name, patronymic, phone, "position", good_job_count) FROM stdin;
 	   employees          postgresadmin    false    199       2865.dat 9           0    0    employees_id_seq    SEQUENCE SET     A   SELECT pg_catalog.setval('employees.employees_id_seq', 5, true);
       	   employees          postgresadmin    false    198         ?
           2606    16410    employees employees_pkey 
   CONSTRAINT     Y   ALTER TABLE ONLY employees.employees
    ADD CONSTRAINT employees_pkey PRIMARY KEY (id);
 E   ALTER TABLE ONLY employees.employees DROP CONSTRAINT employees_pkey;
    	   employees            postgresadmin    false    199                                                                                                                                                                                                                                                                                                                       2865.dat                                                                                            0000600 0004000 0002000 00000000465 13636322476 0014277 0                                                                                                    ustar 00postgres                        postgres                        0000000 0000000                                                                                                                                                                        1	ггг	ппп	 fff	79299999999	 стажер	1
2	ггг	ппп	 fff	79299999999	 стажер	1
3	ггг	ппп	 fff	79299999999	 стажер	1
4	Петр	Иванов	Сидорович	79299999999	 стажер	10
5	Денис	Геритрудович	Юрьевич	79299999999	дизайнер	10
\.


                                                                                                                                                                                                           restore.sql                                                                                         0000600 0004000 0002000 00000033505 13636322476 0015406 0                                                                                                    ustar 00postgres                        postgres                        0000000 0000000                                                                                                                                                                        --
-- NOTE:
--
-- File paths need to be edited. Search for $$PATH$$ and
-- replace it with the path to the directory containing
-- the extracted data files.
--
--
-- PostgreSQL database dump
--

-- Dumped from database version 10.4 (Debian 10.4-2.pgdg90+1)
-- Dumped by pg_dump version 12.0

SET statement_timeout = 0;
SET lock_timeout = 0;
SET idle_in_transaction_session_timeout = 0;
SET client_encoding = 'UTF8';
SET standard_conforming_strings = on;
SELECT pg_catalog.set_config('search_path', '', false);
SET check_function_bodies = false;
SET xmloption = content;
SET client_min_messages = warning;
SET row_security = off;

DROP DATABASE postgresdb;
--
-- Name: postgresdb; Type: DATABASE; Schema: -; Owner: postgres
--

CREATE DATABASE postgresdb WITH TEMPLATE = template0 ENCODING = 'UTF8' LC_COLLATE = 'en_US.utf8' LC_CTYPE = 'en_US.utf8';


ALTER DATABASE postgresdb OWNER TO postgres;

\connect postgresdb

SET statement_timeout = 0;
SET lock_timeout = 0;
SET idle_in_transaction_session_timeout = 0;
SET client_encoding = 'UTF8';
SET standard_conforming_strings = on;
SELECT pg_catalog.set_config('search_path', '', false);
SET check_function_bodies = false;
SET xmloption = content;
SET client_min_messages = warning;
SET row_security = off;

--
-- Name: employees; Type: SCHEMA; Schema: -; Owner: postgresadmin
--

CREATE SCHEMA employees;


ALTER SCHEMA employees OWNER TO postgresadmin;

--
-- Name: test; Type: SCHEMA; Schema: -; Owner: postgresadmin
--

CREATE SCHEMA test;


ALTER SCHEMA test OWNER TO postgresadmin;

--
-- Name: employee_add(character varying, character varying, character varying, character varying, character varying, integer); Type: FUNCTION; Schema: employees; Owner: postgresadmin
--

CREATE FUNCTION employees.employee_add(_name character varying, _last_name character varying, _patronymic character varying, _phone character varying, _position character varying, _good_job_count integer) RETURNS void
    LANGUAGE plpgsql SECURITY DEFINER
    AS $$
BEGIN

   INSERT INTO employees.employees 
   (
	
     name,
	last_name,
	patronymic,
	phone,
	position,
	good_job_count
   )
   VALUES
   (
	  
      _name,
	_last_name,
	_patronymic,
	_phone,
	_position,
	_good_job_count 
   );
 
END;
$$;


ALTER FUNCTION employees.employee_add(_name character varying, _last_name character varying, _patronymic character varying, _phone character varying, _position character varying, _good_job_count integer) OWNER TO postgresadmin;

--
-- Name: employee_get(integer); Type: FUNCTION; Schema: employees; Owner: postgresadmin
--

CREATE FUNCTION employees.employee_get(_id integer) RETURNS TABLE(name character varying, last_name character varying, id integer, patronymic character varying, phone character varying, _position character varying, _good_job_count integer)
    LANGUAGE plpgsql SECURITY DEFINER
    AS $$
BEGIN
     IF NOT EXISTS (SELECT 1 FROM employees.employees m WHERE m.id = _id)
    THEN 
        RAISE EXCEPTION 'Сотрудника с таким id не существует' USING ERRCODE = 50003;
    END IF;
    RETURN QUERY
    SELECT e.name,
           e.last_name,
           e.id,
           e.patronymic,
           e.phone,
		   e.position,
		   e.good_job_count
	  
         
    FROM employees.employees e
	where e.id = _id	;

END;
$$;


ALTER FUNCTION employees.employee_get(_id integer) OWNER TO postgresadmin;

--
-- Name: employee_get_all(); Type: FUNCTION; Schema: employees; Owner: postgresadmin
--

CREATE FUNCTION employees.employee_get_all() RETURNS TABLE(name character varying, last_name character varying, id integer, patronymic character varying, phone character varying, _position character varying, _good_job_count integer)
    LANGUAGE plpgsql SECURITY DEFINER
    AS $$
BEGIN

    RETURN QUERY
    SELECT e.name,
           e.last_name,
           e.id,
           e.patronymic,
           e.phone,
		   e.position,
		   e.good_job_count
          
    FROM employees.employees e;

END;
$$;


ALTER FUNCTION employees.employee_get_all() OWNER TO postgresadmin;

--
-- Name: employee_remove(integer); Type: FUNCTION; Schema: employees; Owner: postgresadmin
--

CREATE FUNCTION employees.employee_remove(_id integer) RETURNS void
    LANGUAGE plpgsql SECURITY DEFINER
    AS $$
BEGIN
  IF NOT EXISTS (SELECT 1 FROM employees.employees m WHERE m.id = _id)
    THEN 
        RAISE EXCEPTION 'Сотрудника с таким id не существует' USING ERRCODE = 50003;
    END IF;

   DELETE from employees.employees as e
   where e.id = _id;

END;
$$;


ALTER FUNCTION employees.employee_remove(_id integer) OWNER TO postgresadmin;

--
-- Name: employee_upd(integer, character varying, character varying, character varying, character varying, character varying, integer); Type: FUNCTION; Schema: employees; Owner: postgresadmin
--

CREATE FUNCTION employees.employee_upd(_id integer, _name character varying, _last_name character varying, _patronymic character varying, _phone character varying, _position character varying, _good_job_count integer) RETURNS void
    LANGUAGE plpgsql SECURITY DEFINER
    AS $$
BEGIN
 IF NOT EXISTS (SELECT 1 FROM employees.employees m WHERE m.id = _id)
    THEN 
        RAISE EXCEPTION 'Сотрудника с таким id не существует' USING ERRCODE = 50003;
    END IF;

UPDATE employees.employees
	SET id=_id, last_name=_last_name, name=_name, patronymic=_patronymic, phone=_phone, position=_position, good_job_count=_good_job_count
	WHERE id = _id;
END;
$$;


ALTER FUNCTION employees.employee_upd(_id integer, _name character varying, _last_name character varying, _patronymic character varying, _phone character varying, _position character varying, _good_job_count integer) OWNER TO postgresadmin;

--
-- Name: employees_get_all(); Type: FUNCTION; Schema: employees; Owner: postgresadmin
--

CREATE FUNCTION employees.employees_get_all() RETURNS TABLE(name character varying, last_name character varying, id integer, patronymic character varying, phone character varying, _position character varying, _good_job_count integer)
    LANGUAGE plpgsql SECURITY DEFINER
    AS $$
BEGIN

    RETURN QUERY
    SELECT e.name,
           e.last_name,
           e.id,
           e.patronymic,
           e.phone,
		   e.position,
		   e.good_job_count
          
    FROM employees.employees e;

END;
$$;


ALTER FUNCTION employees.employees_get_all() OWNER TO postgresadmin;

--
-- Name: employees_get_all_part1(); Type: FUNCTION; Schema: employees; Owner: postgresadmin
--

CREATE FUNCTION employees.employees_get_all_part1() RETURNS TABLE(name character varying, last_name character varying, id integer)
    LANGUAGE plpgsql SECURITY DEFINER
    AS $$
BEGIN

    RETURN QUERY
    SELECT e.name,
           e.last_name,
           e.id
          
    FROM employees.employees e;

END;
$$;


ALTER FUNCTION employees.employees_get_all_part1() OWNER TO postgresadmin;

--
-- Name: employees_get_all_part2(); Type: FUNCTION; Schema: employees; Owner: postgresadmin
--

CREATE FUNCTION employees.employees_get_all_part2() RETURNS TABLE(id integer, patronymic character varying, phone character varying, _position character varying, _good_job_count integer)
    LANGUAGE plpgsql SECURITY DEFINER
    AS $$
BEGIN

    RETURN QUERY
    SELECT 
           e.id,
           e.patronymic,
           e.phone,
		   e.position,
		   e.good_job_count
          
    FROM employees.employees e;

END;
$$;


ALTER FUNCTION employees.employees_get_all_part2() OWNER TO postgresadmin;

--
-- Name: get_all(); Type: FUNCTION; Schema: employees; Owner: postgresadmin
--

CREATE FUNCTION employees.get_all() RETURNS TABLE(name character varying, last_name character varying, id integer, patronymic character varying, phone character varying, _position character varying, _good_job_count integer)
    LANGUAGE plpgsql SECURITY DEFINER
    AS $$
BEGIN

    RETURN QUERY
    SELECT e.name,
           e.last_name,
           e.id,
           e.patronymic,
           e.phone,
		   e.position,
		   e.good_job_count
          
    FROM employees.employees e;

END;
$$;


ALTER FUNCTION employees.get_all() OWNER TO postgresadmin;

--
-- Name: db_error_test(); Type: FUNCTION; Schema: public; Owner: postgresadmin
--

CREATE FUNCTION public.db_error_test() RETURNS void
    LANGUAGE plpgsql SECURITY DEFINER
    AS $$
DECLARE _query text;
DECLARE _place_error text;
BEGIN
    CREATE LOCAL  TEMPORARY TABLE script( ex text ) ON COMMIT DROP;
	
	INSERT INTO script
	(
		ex
	)
    WITH cte AS
    (
    SELECT  n.nspname, 
            p.proname,
    		p.proargmodes, 
    		CASE WHEN p.proallargtypes IS NULL 
    		THEN  CASE WHEN array_length(string_to_array(proargtypes::text,' ')::oid[],1) > 1 THEN string_to_array(proargtypes::text,' ')::oid[] 
    	   ELSE NULL END
    		ELSE p.proallargtypes END args, 
    		p.proargnames
    FROM    pg_catalog.pg_namespace n
       INNER JOIN pg_catalog.pg_proc p
           ON p.pronamespace = n.oid
	WHERE nspowner != 10
    ),
    cte_2 AS
    (
    SELECT initcap(cte.nspname) nspname,
           cte.proname,
    	   unnest(cte.proargmodes) param_type, 
    	   format_type(unnest(cte.args), null) arg_type, 
    	   unnest(proargnames) arg_name
    FROM cte
    ),
    cte_3 AS
    (
    SELECT CASE WHEN param_type = 'i' or param_type IS NULL THEN 'input' ELSE 'output' END param_type,
    string_agg( CASE WHEN param_type = 'i' or param_type IS NULL THEN cte_2.arg_name||':= null,' END,' ' ) arg_name,
    cte_2.nspname,
    cte_2.proname,
    COUNT(proname) OVER(PARTITION BY proname) counts
    FROM cte_2
    GROUP BY nspname, proname, CASE WHEN param_type = 'i' or param_type IS NULL THEN 'input' ELSE 'output' END
    )
    SELECT DISTINCT 'SELECT '||nspname||'.'||proname||'('||CASE WHEN arg_name IS NULL THEN '' ELSE left(arg_name, (length(arg_name)-1)) END ||');' ex
    FROM cte_3
    WHERE param_type = 'input' or counts = 1;
    
	FOR _query IN (SELECT ex FROM script) LOOP 
	    BEGIN
	    EXECUTE _query;
		EXCEPTION
            WHEN not_null_violation or SQLSTATE '50003' or SQLSTATE '50001' or SQLSTATE '50002' or SQLSTATE '50000' THEN 
			RAISE NOTICE 'Пропущена бизнесовая ошибка или not null.
';
	        WHEN others THEN
	         GET STACKED DIAGNOSTICS _place_error = PG_EXCEPTION_CONTEXT;
	         _place_error = 'Объект: '|| _place_error||' С ошибкой: '||SQLERRM||'
			 ';
	         RAISE NOTICE '%', _place_error;
	    END;
	END LOOP;
	
	RAISE NOTICE 'Ошибок Нет!
';
	
END;
$$;


ALTER FUNCTION public.db_error_test() OWNER TO postgresadmin;

--
-- Name: do_something(integer, jsonb, timestamp without time zone); Type: FUNCTION; Schema: test; Owner: postgresadmin
--

CREATE FUNCTION test.do_something(_employee_id integer, _test_data jsonb, _date timestamp without time zone) RETURNS void
    LANGUAGE plpgsql SECURITY DEFINER
    AS $$
BEGIN

 IF _employee_id IS NULL 
 THEN
   RAISE EXCEPTION '_employee_id не должен быть null' USING ERRCODE = 50000;
 END IF;
  IF _test_data IS NULL 
 THEN
   RAISE EXCEPTION '_test_data не должен быть null' USING ERRCODE = 50000;
    END IF;
  IF _date IS NULL 
 THEN
   RAISE EXCEPTION '_date не должен быть null' USING ERRCODE = 50000;
    END IF;	
   


END;
$$;


ALTER FUNCTION test.do_something(_employee_id integer, _test_data jsonb, _date timestamp without time zone) OWNER TO postgresadmin;

--
-- Name: get_db_error(integer); Type: FUNCTION; Schema: test; Owner: postgresadmin
--

CREATE FUNCTION test.get_db_error(_id integer) RETURNS boolean
    LANGUAGE plpgsql SECURITY DEFINER
    AS $$
BEGIN

    IF _id = 1
    THEN
        RAISE EXCEPTION 'Пользовательская ошибка' USING ERRCODE = 50000;
	END IF;	
	IF _id = 2
    THEN
        RAISE EXCEPTION 'Внутренняя ошибка' USING ERRCODE = 10000;	
    ELSE
	   
       RETURN TRUE;
    END IF;

END;
$$;


ALTER FUNCTION test.get_db_error(_id integer) OWNER TO postgresadmin;

--
-- Name: long_time_process(); Type: FUNCTION; Schema: test; Owner: postgresadmin
--

CREATE FUNCTION test.long_time_process() RETURNS void
    LANGUAGE plpgsql SECURITY DEFINER
    AS $$
BEGIN

  select pg_sleep(500);

END;
$$;


ALTER FUNCTION test.long_time_process() OWNER TO postgresadmin;

SET default_tablespace = '';

--
-- Name: employees; Type: TABLE; Schema: employees; Owner: postgresadmin
--

CREATE TABLE employees.employees (
    id integer NOT NULL,
    last_name character varying(100) NOT NULL,
    name character varying(100) NOT NULL,
    patronymic character varying(15),
    phone character varying(15),
    "position" character varying(20) NOT NULL,
    good_job_count integer NOT NULL
);


ALTER TABLE employees.employees OWNER TO postgresadmin;

--
-- Name: employees_id_seq; Type: SEQUENCE; Schema: employees; Owner: postgresadmin
--

CREATE SEQUENCE employees.employees_id_seq
    AS integer
    START WITH 1
    INCREMENT BY 1
    NO MINVALUE
    NO MAXVALUE
    CACHE 1;


ALTER TABLE employees.employees_id_seq OWNER TO postgresadmin;

--
-- Name: employees_id_seq; Type: SEQUENCE OWNED BY; Schema: employees; Owner: postgresadmin
--

ALTER SEQUENCE employees.employees_id_seq OWNED BY employees.employees.id;


--
-- Name: employees id; Type: DEFAULT; Schema: employees; Owner: postgresadmin
--

ALTER TABLE ONLY employees.employees ALTER COLUMN id SET DEFAULT nextval('employees.employees_id_seq'::regclass);


--
-- Data for Name: employees; Type: TABLE DATA; Schema: employees; Owner: postgresadmin
--

COPY employees.employees (id, last_name, name, patronymic, phone, "position", good_job_count) FROM stdin;
\.
COPY employees.employees (id, last_name, name, patronymic, phone, "position", good_job_count) FROM '$$PATH$$/2865.dat';

--
-- Name: employees_id_seq; Type: SEQUENCE SET; Schema: employees; Owner: postgresadmin
--

SELECT pg_catalog.setval('employees.employees_id_seq', 5, true);


--
-- Name: employees employees_pkey; Type: CONSTRAINT; Schema: employees; Owner: postgresadmin
--

ALTER TABLE ONLY employees.employees
    ADD CONSTRAINT employees_pkey PRIMARY KEY (id);


--
-- PostgreSQL database dump complete
--

                                                                                                                                                                                                                                                                                                                                                                                                                                                                                                                                                                                                                                                                                                                                                                                                                                                                                                                                                                                                                                                                                                                                                                                                                                                           