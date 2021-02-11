DROP TABLE IF EXISTS test;
CREATE TABLE IF NOT EXISTS test (
    flag1 boolean,
    flag2 bool not null default true
);

INSERT INTO test VALUES (true, false);
INSERT INTO test VALUES ('true', 'false');
INSERT INTO test VALUES ('yes', 'no');
INSERT INTO test VALUES ('t', 'f');
INSERT INTO test VALUES ('y', 'n');
INSERT INTO test VALUES ('1', '0');
INSERT INTO test VALUES (null, false);

TRUNCATE test;
SELECT * FROM test;

DROP TABLE IF EXISTS test;

CREATE TABLE IF NOT EXISTS test (
    name varchar(16),
    sex char(6),
    introduction text
);

INSERT INTO test VALUES ('Sam', 'male', 'I''am Sam');
INSERT INTO test VALUES ('Judy', 'female', 'I''am Judy');
INSERT INTO test VALUES ('Sam1', ' male', 'I''am Sam');
INSERT INTO test VALUES ('Sam2', ' male ', 'I''am Sam');

SELECT
       pg_column_size(t.name),
       pg_column_size(t.sex),
       pg_column_size(t.introduction)
   FROM test as t;

SELECT name,
       length(name)         length_of_name,
       sex,
       length(sex)          length_of_sex,
       introduction,
       length(introduction) length_of_introduction
FROM test;

SELECT pg_column_size('abc'::varchar(4));


DROP TABLE IF EXISTS test;

CREATE TABLE IF NOT EXISTS test (
    id serial primary key,
    age smallint,
    population int,
    size bigint
);

INSERT INTO test (age, population, size) VALUES (-32768, -2147483648, -9223372036854775808);
INSERT INTO test (age, population, size) VALUES (32767, 2147483647, 9223372036854775807);

SELECT * FROM test;


DROP TABLE IF EXISTS test;

CREATE TABLE IF NOT EXISTS test (
    price1 float(4),
    price2 real,
    price3 float4,
    price4 float8,
    price5 numeric,
    price6 numeric(4, 2)
);

INSERT INTO test VALUES (1.1234, 1.1234, 1.1234, 1.1234, 1.1234, 1.12);
INSERT INTO test VALUES (1.1234, 1.1234, 1.1234, 1.1234, 1.1234, 1.123);
INSERT INTO test VALUES (1.1234, 1.1234, 1.1234, 1.1234, 1.1234, 12.123);
INSERT INTO test VALUES (1.1234, 1.1234, 1.1234, 1.1234, 1.1234, 123.123); -- ±¨´í

DROP TABLE IF EXISTS test;

CREATE TABLE IF NOT EXISTS test (
    money1 money,
    money2 money
);

INSERT INTO test VALUES (-92233720368547758.08, 92233720368547758.07);

SELECT * FROM test;



DROP TABLE IF EXISTS test;

CREATE TABLE IF NOT EXISTS test (
    t1 date,
    t2 time,
    t3 timestamp,
    d interval
);

INSERT INTO test VALUES ('2020-01-01', '12:12:12', current_timestamp, interval '40' minute);

SELECT * FROM test;


DROP TABLE IF EXISTS test;

CREATE TABLE IF NOT EXISTS test (
    arr int[]
);

INSERT INTO test VALUES ('{1, 2, 3}');

SELECT arr[1] as idx_1, arr[2] as idx_2, arr[3] as idx_3 FROM test;
SELECT * FROM test WHERE arr[1]=1;


DROP TABLE IF EXISTS test;

CREATE TABLE IF NOT EXISTS test (
    json1 json,
    jsonb1 jsonb
);

INSERT INTO test VALUES ('{"a": 1}', '{"b": 2}');
SELECT * FROM test;

CREATE INDEX test_index_1 ON test ((json1 ->> 'a'));
CREATE INDEX test_index_2 ON test ((jsonb1 ->> 'b'));

EXPLAIN INSERT INTO test VALUES ('{"a": 1}', '{"b": 2}');

SELECT * FROM pg_indexes WHERE tablename='test';

DROP TABLE IF EXISTS test;

CREATE TABLE IF NOT EXISTS test (
    uuid uuid
);

CREATE EXTENSION IF NOT EXISTS "uuid-ossp";
INSERT INTO test VALUES (uuid_generate_v4());
SELECT * FROM test;


DROP TABLE IF EXISTS test;

CREATE TABLE IF NOT EXISTS test (
    point point,
    line line,
    lseg lseg,
    box box,
    path1 path,
    path2 path,
    polygon polygon,
    circle circle
);

INSERT INTO test VALUES (
    '(1, 2)', -- pointpublic
    '((1,2),(3,4))', -- line
    '((1,2),(3,4))',  -- lseq
    '((1,2),(3,4))',  -- box
    '((1,2),(3,4))',  -- path ·â±Õ
    '[(1,2),(3,4)]',  -- path
    '((1,2),(3,4))',  -- polygon
    '<(1,2), 3>');  -- circle

SELECT * FROM test;