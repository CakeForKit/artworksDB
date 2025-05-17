



-- create table if not exists prolog.phone(
--     id SERIAL PRIMARY KEY,
--     phone VARCHAR(255) NOT NULL
-- );


-- create table if not exists prolog.address(
--  id SERIAL PRIMARY KEY,
--     city VARCHAR(255) NOT NULL,
--     street VARCHAR(255) NOT NULL,
--     house int
-- );


-- CREATE TABLE IF NOT EXISTS prolog.person_phone (
--     id SERIAL PRIMARY KEY,
--     person_id INT REFERENCES prolog.surname(id) ON DELETE CASCADE,
--     phone_id INT REFERENCES prolog.phone(id) ON DELETE CASCADE
-- );

-- Промежуточная таблица: Человек ↔ Машина
-- CREATE TABLE IF NOT EXISTS prolog.person_car (
--     id SERIAL PRIMARY KEY,
--     person_id INT REFERENCES prolog.surname(id) ON DELETE CASCADE,
--     car_id INT REFERENCES prolog.car(id) ON DELETE CASCADE
-- );

create schema if not exists prolog;

create table if not exists prolog.surname(
    id SERIAL PRIMARY KEY,
    surname VARCHAR(255) NOT NULL
);

create table if not exists prolog.car(
    id SERIAL PRIMARY KEY,
    model VARCHAR(255) NOT NULL,
    colour VARCHAR(255) NOT NULL,
    price int,
    number VARCHAR(255)
);

create table if not exists prolog.car_person(
    idCar int,
    idSurname int,
    FOREIGN KEY (idCar) REFERENCES prolog.car(id),
    FOREIGN KEY (idSurname) REFERENCES prolog.surname(id)
);

drop Table prolog.car
drop Table prolog.car_person

INSERT INTO prolog.surname (surname) VALUES
('surname1'),
('surname2');

INSERT INTO prolog.car_person (idCar, idSurname) VALUES
(1, 6);

select * from prolog.surname

select * from prolog.car

select * from prolog.car_person


SELECT s.surname, c.number
FROM 
    prolog.surname s,
    prolog.car c
WHERE 
    EXISTS (
        SELECT 1 
        FROM prolog.car_person cp 
        WHERE cp.idSurname = s.id AND cp.idCar = c.id
    )
    AND c.model = 'm1'
    AND c.colour = 'black';

SELECT c.model, c.colour, c.price
FROM 
    prolog.surname s,
    prolog.car c
WHERE 
    EXISTS (
        SELECT 1 
        FROM prolog.car_person cp 
        WHERE cp.idSurname = s.id AND cp.idCar = c.id
    )
    AND s.surname = 'surname1'
    AND c.number = 'qwe33';

SELECT 
    (SELECT surname FROM prolog.surname WHERE id = cp.idSurname) AS surname,
    (SELECT number FROM prolog.car WHERE id = cp.idCar AND model = 'm1' AND colour = 'black') AS number
FROM 
    prolog.car_person cp
WHERE 
    cp.idCar IN (SELECT id FROM prolog.car WHERE model = 'm1' AND colour = 'black');

select surname, number
from car
where number = 'qwe33q' and 
    id in (select idCar 
            from car_person
            where idSurname in (select id
                                from surname
                                where model = 'm1' ans colour = 'black'))

-- INSERT INTO prolog.phone (phone) VALUES
-- ('+79261111444'),
-- ('+79262211554'),
-- ('+79261123456'),
-- ('+79267778899'),
-- ('+79265571444');

INSERT INTO prolog.car (model, colour, price, number) VALUES
('m1', 'black', 100, 'qwe33q'),
('m3', 'black', 110, 'asd22r'),
('m3', 'red', 44, 'ert44r');


INSERT INTO prolog.address (city, street, house) VALUES
('moscow', 'rokossovsky', 10),
('moscow', 'lenina', 12),
('saint_petersburg', 'nevsky', 5),
('smolensk', 'lenina', 15),
('moscow', 'tverskaya', 25);


-- Привязка Иванова к телефонам
INSERT INTO prolog.person_phone (person_id, phone_id) VALUES
(1, 1), (1, 2);

-- Привязка Иванова к машинам
INSERT INTO prolog.person_car (person_id, car_id) VALUES
(1, 1), -- BMW
(1, 2); -- Mercedes


INSERT INTO prolog.person_phone (person_id, phone_id) VALUES
(2, 3);
INSERT INTO prolog.person_car (person_id, car_id) VALUES
(2, 7); -- Toyota white


SELECT marka, colour, cost, number
FROM prolog.car
WHERE id IN (
    SELECT car_id
    FROM prolog.person_car
    WHERE person_id = (
        SELECT id
        FROM prolog.surname
        WHERE surname = 'ivanov'
    )
);