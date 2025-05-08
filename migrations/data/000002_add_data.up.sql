-- Active: 1744740356603@@127.0.0.1@5432@artworks


-- Вставляем администраторов
INSERT INTO Admins (id, username, login, hashedPassword, createdAt) VALUES
(gen_random_uuid(), 'admin1', 'admin1_login', '$2a$10$xJwL5v5Jz5U5Z5U5Z5U5Ze', NOW()),
(gen_random_uuid(), 'admin2', 'admin2_login', '$2a$10$xJwL5v5Jz5U5Z5U5Z5U5Ze', NOW());

-- Вставляем сотрудников (сначала получаем ID администраторов)
WITH admin_ids AS (
    SELECT id FROM Admins ORDER BY createdAt LIMIT 2
)
INSERT INTO Employees (id, username, login, hashedPassword, createdAt, adminID) VALUES
(gen_random_uuid(), 'employee1', 'emp1_login', '$2a$10$xJwL5v5Jz5U5Z5U5Z5U5Ze', NOW(), (SELECT id FROM admin_ids OFFSET 0 LIMIT 1)),
(gen_random_uuid(), 'employee2', 'emp2_login', '$2a$10$xJwL5v5Jz5U5Z5U5Z5U5Zd', NOW(), (SELECT id FROM admin_ids OFFSET 1 LIMIT 1));

-- Вставляем пользователей
-- INSERT INTO Users (id, username, login, hashedPassword, createdAt, email, subscribeMail) VALUES
-- (gen_random_uuid(), 'user1', 'user1_login', '$2a$10$xJwL5v5Jz5U5Z5U5Z5U5Ze', NOW(), 'user1@example.com', TRUE),
-- (gen_random_uuid(), 'user2', 'user2_login', '$2a$10$xJwL5v5Jz5U5Z5U5Z5U5Ze', NOW(), 'user2@example.com', FALSE);

-- Вставляем авторов
INSERT INTO Author (id, name, birthYear, deathYear) VALUES
(gen_random_uuid(), 'Leonardo da Vinci', 1452, 1519),
(gen_random_uuid(), 'Vincent van Gogh', 1853, 1890);

-- Вставляем коллекции
INSERT INTO Collection (id, title) VALUES
(gen_random_uuid(), 'Renaissance Masterpieces'),
(gen_random_uuid(), 'Post-Impressionism Collection');

-- Вставляем произведения искусства (сначала получаем ID авторов и коллекций)
WITH 
author_data AS (
    SELECT id FROM Author WHERE name = 'Leonardo da Vinci' LIMIT 1
),
collection_data AS (
    SELECT id FROM Collection WHERE title = 'Renaissance Masterpieces' LIMIT 1
),
van_gogh AS (
    SELECT id FROM Author WHERE name = 'Vincent van Gogh' LIMIT 1
),
impression_collection AS (
    SELECT id FROM Collection WHERE title = 'Post-Impressionism Collection' LIMIT 1
)
INSERT INTO Artworks (id, title, technic, material, size, creationYear, authorID, collectionID) VALUES
(gen_random_uuid(), 'Mona Lisa', 'Oil painting', 'Poplar wood', '77 × 53 cm', 1503, (SELECT id FROM author_data), (SELECT id FROM collection_data)),
(gen_random_uuid(), 'Starry Night', 'Oil painting', 'Canvas', '73.7 × 92.1 cm', 1889, (SELECT id FROM van_gogh), (SELECT id FROM impression_collection));

-- Вставляем события (сначала получаем ID сотрудников)
WITH employee_ids AS (
    SELECT id FROM Employees ORDER BY createdAt LIMIT 2
)
INSERT INTO Events (id, title, dateBegin, dateEnd, canVisit, adress, cntTickets, creatorID) VALUES
(gen_random_uuid(), 'Renaissance Exhibition', NOW() + INTERVAL '10 days', NOW() + INTERVAL '20 days', TRUE, '123 Art Gallery St, Museum District', 100, (SELECT id FROM employee_ids OFFSET 0 LIMIT 1)),
(gen_random_uuid(), 'Van Gogh Special', NOW() + INTERVAL '15 days', NOW() + INTERVAL '25 days', TRUE, '456 Modern Art Ave, Downtown', 150, (SELECT id FROM employee_ids OFFSET 1 LIMIT 1));

select * from events

select * from artworks

SELECT 
        a.id as artwork_id, 
        e.id as event_id
    FROM 
        Artworks a
    JOIN 
        Events e ON 
        (a.title = 'Mona Lisa' AND e.title = 'Renaissance Exhibition') OR
        (a.title = 'Starry Night' AND e.title = 'Van Gogh Special')

-- Связываем произведения искусства с событиями (сначала получаем ID произведений и событий)
WITH 
artwork_event_data AS (
    SELECT 
        a.id as artwork_id, 
        e.id as event_id
    FROM 
        Artworks a
    JOIN 
        Events e ON 
        (a.title = 'Mona Lisa' AND e.title = 'Renaissance Exhibition') OR
        (a.title = 'Starry Night' AND e.title = 'Van Gogh Special')
)
INSERT INTO Artwork_event (artworkID, eventID)
SELECT artwork_id, event_id FROM artwork_event_data;

-- Вставляем покупки билетов (сначала получаем ID событий)
WITH event_data AS (
    SELECT id FROM Events ORDER BY dateBegin LIMIT 2
)
INSERT INTO TicketPurchases (id, customerName, customerEmail, eventID, purchaseDate) VALUES
(gen_random_uuid(), 'John Doe', 'john.doe@example.com', (SELECT id FROM event_data OFFSET 0 LIMIT 1), NOW()),
(gen_random_uuid(), 'Jane Smith', 'jane.smith@example.com', (SELECT id FROM event_data OFFSET 1 LIMIT 1), NOW());



----------------------------------

INSERT INTO Author (id, name, birthYear, deathYear) VALUES
(gen_random_uuid(), 'Salvador Dali', 1904, 1989),
(gen_random_uuid(), 'Johannes Vermeer', 1632, 1675),
(gen_random_uuid(), 'Edvard Munch', 1863, 1944),
(gen_random_uuid(), 'Pablo Picasso', 1881, 1973);

-- Затем добавляем произведения с правильным указанием авторов
WITH 
authors AS (
    SELECT id, name FROM Author
),
collections AS (
    SELECT id, title FROM Collection
)
INSERT INTO Artworks (id, title, technic, material, size, creationYear, authorID, collectionID) VALUES
(gen_random_uuid(), 'The Persistence of Memory', 'Oil painting', 'Canvas', '24 × 33 cm', 1931, 
    (SELECT id FROM authors WHERE name = 'Salvador Dali'), 
    (SELECT id FROM collections WHERE title = 'Post-Impressionism Collection')),

(gen_random_uuid(), 'Girl with a Pearl Earring', 'Oil painting', 'Canvas', '44 × 39 cm', 1665, 
    (SELECT id FROM authors WHERE name = 'Johannes Vermeer'), 
    (SELECT id FROM collections WHERE title = 'Renaissance Masterpieces')),

(gen_random_uuid(), 'The Scream', 'Tempera', 'Cardboard', '91 × 73 cm', 1893, 
    (SELECT id FROM authors WHERE name = 'Edvard Munch'), 
    (SELECT id FROM collections WHERE title = 'Post-Impressionism Collection')),

(gen_random_uuid(), 'Guernica', 'Oil painting', 'Canvas', '349 × 776 cm', 1937, 
    (SELECT id FROM authors WHERE name = 'Pablo Picasso'), 
    (SELECT id FROM collections WHERE title = 'Post-Impressionism Collection'));


-- Добавляем 10 новых пользователей
INSERT INTO Users (id, username, login, hashedPassword, createdAt, email, subscribeMail) VALUES
(gen_random_uuid(), 'Alex Johnson', 'alexj', '$2a$10$xJwL5v5Jz5U5Z5U5Z5U5Ze', NOW() - INTERVAL '30 days', 'alex.johnson@example.com', TRUE),
(gen_random_uuid(), 'Maria Garcia', 'mariag', '$2a$10$xJwL5v5Jz5U5Z5U5Z5U5Ze', NOW() - INTERVAL '25 days', 'maria.garcia@example.com', FALSE),
(gen_random_uuid(), 'James Smith', 'jamess', '$2a$10$xJwL5v5Jz5U5Z5U5Z5U5Ze', NOW() - INTERVAL '20 days', 'james.smith@example.com', TRUE),
(gen_random_uuid(), 'Emma Wilson', 'emmaw', '$2a$10$xJwL5v5Jz5U5Z5U5Z5U5Ze', NOW() - INTERVAL '15 days', 'emma.wilson@example.com', FALSE),
(gen_random_uuid(), 'Michael Brown', 'michaelb', '$2a$10$xJwL5v5Jz5U5Z5U5Z5U5Ze', NOW() - INTERVAL '10 days', 'michael.brown@example.com', TRUE),
(gen_random_uuid(), 'Sophia Davis', 'sophiad', '$2a$10$xJwL5v5Jz5U5Z5U5Z5U5Ze', NOW() - INTERVAL '8 days', 'sophia.davis@example.com', FALSE),
(gen_random_uuid(), 'William Miller', 'williamm', '$2a$10$xJwL5v5Jz5U5Z5U5Z5U5Ze', NOW() - INTERVAL '5 days', 'william.miller@example.com', TRUE),
(gen_random_uuid(), 'Olivia Wilson', 'oliviaw', '$2a$10$xJwL5v5Jz5U5Z5U5Z5U5Ze', NOW() - INTERVAL '3 days', 'olivia.wilson@example.com', FALSE),
(gen_random_uuid(), 'Benjamin Taylor', 'benjamin', '$2a$10$xJwL5v5Jz5U5Z5U5Z5U5Ze', NOW() - INTERVAL '2 days', 'benjamin.taylor@example.com', TRUE),
(gen_random_uuid(), 'Ava Anderson', 'avaa', '$2a$10$xJwL5v5Jz5U5Z5U5Z5U5Ze', NOW() - INTERVAL '1 day', 'ava.anderson@example.com', FALSE);


-- Добавляем еще 3 события
WITH employee_ids AS (
    SELECT id FROM Employees ORDER BY createdAt LIMIT 2
)
INSERT INTO Events (id, title, dateBegin, dateEnd, canVisit, adress, cntTickets, creatorID) VALUES
(gen_random_uuid(), 'Surrealism Exhibition', NOW() + INTERVAL '5 days', NOW() + INTERVAL '15 days', TRUE, '789 Modern Art Blvd, Arts District', 80, (SELECT id FROM employee_ids OFFSET 0 LIMIT 1)),
(gen_random_uuid(), 'Cubism Special', NOW() + INTERVAL '12 days', NOW() + INTERVAL '22 days', TRUE, '101 Art Center Rd, Downtown', 120, (SELECT id FROM employee_ids OFFSET 1 LIMIT 1)),
(gen_random_uuid(), 'Expressionism Showcase', NOW() + INTERVAL '8 days', NOW() + INTERVAL '18 days', TRUE, '202 Gallery Lane, Museum Quarter', 90, (SELECT id FROM employee_ids OFFSET 0 LIMIT 1));


-- Сначала создадим временные таблицы для хранения ID пользователей и событий
WITH 
user_ids AS (
    SELECT id, username, email FROM Users ORDER BY createdAt LIMIT 10
),
event_ids AS (
    SELECT id, title FROM Events ORDER BY dateBegin LIMIT 5
),
-- Создаем покупки билетов
ticket_purchases AS (
    INSERT INTO TicketPurchases (id, customerName, customerEmail, purchaseDate, eventID)
    VALUES
    -- Покупки для Renaissance Exhibition
    (gen_random_uuid(), (SELECT username FROM user_ids OFFSET 0 LIMIT 1), (SELECT email FROM user_ids OFFSET 0 LIMIT 1), NOW() - INTERVAL '2 days', (SELECT id FROM event_ids WHERE title = 'Renaissance Exhibition')),
    (gen_random_uuid(), (SELECT username FROM user_ids OFFSET 1 LIMIT 1), (SELECT email FROM user_ids OFFSET 1 LIMIT 1), NOW() - INTERVAL '1 day', (SELECT id FROM event_ids WHERE title = 'Renaissance Exhibition')),
    (gen_random_uuid(), (SELECT username FROM user_ids OFFSET 2 LIMIT 1), (SELECT email FROM user_ids OFFSET 2 LIMIT 1), NOW(), (SELECT id FROM event_ids WHERE title = 'Renaissance Exhibition')),
    -- Покупки для Van Gogh Special
    (gen_random_uuid(), (SELECT username FROM user_ids OFFSET 3 LIMIT 1), (SELECT email FROM user_ids OFFSET 3 LIMIT 1), NOW() - INTERVAL '3 days', (SELECT id FROM event_ids WHERE title = 'Van Gogh Special')),
    (gen_random_uuid(), (SELECT username FROM user_ids OFFSET 4 LIMIT 1), (SELECT email FROM user_ids OFFSET 4 LIMIT 1), NOW() - INTERVAL '2 days', (SELECT id FROM event_ids WHERE title = 'Van Gogh Special')),
    -- Покупки для Surrealism Exhibition
    (gen_random_uuid(), (SELECT username FROM user_ids OFFSET 5 LIMIT 1), (SELECT email FROM user_ids OFFSET 5 LIMIT 1), NOW() - INTERVAL '1 day', (SELECT id FROM event_ids WHERE title = 'Surrealism Exhibition')),
    (gen_random_uuid(), (SELECT username FROM user_ids OFFSET 6 LIMIT 1), (SELECT email FROM user_ids OFFSET 6 LIMIT 1), NOW(), (SELECT id FROM event_ids WHERE title = 'Surrealism Exhibition')),
    -- Покупки для Cubism Special
    (gen_random_uuid(), (SELECT username FROM user_ids OFFSET 7 LIMIT 1), (SELECT email FROM user_ids OFFSET 7 LIMIT 1), NOW() - INTERVAL '4 days', (SELECT id FROM event_ids WHERE title = 'Cubism Special')),
    -- Покупки для Expressionism Showcase
    (gen_random_uuid(), (SELECT username FROM user_ids OFFSET 8 LIMIT 1), (SELECT email FROM user_ids OFFSET 8 LIMIT 1), NOW() - INTERVAL '3 days', (SELECT id FROM event_ids WHERE title = 'Expressionism Showcase')),
    (gen_random_uuid(), (SELECT username FROM user_ids OFFSET 9 LIMIT 1), (SELECT email FROM user_ids OFFSET 9 LIMIT 1), NOW() - INTERVAL '2 days', (SELECT id FROM event_ids WHERE title = 'Expressionism Showcase')),
    (gen_random_uuid(), 'Guest Visitor', 'guest1@example.com', NOW() - INTERVAL '1 day', (SELECT id FROM event_ids WHERE title = 'Expressionism Showcase'))
    RETURNING id, customerEmail
)
-- Теперь связываем билеты с пользователями (кроме гостевого билета)
INSERT INTO tickets_user (ticketID, userID)
SELECT 
    tp.id, 
    u.id
FROM 
    ticket_purchases tp
JOIN 
    Users u ON tp.customerEmail = u.email
WHERE 
    tp.customerEmail != 'guest1@example.com';












-- ЗАПОЛНЕНИЕ ДАННЫЙ ДЛЯ ЗАМЕРОВ
-- Генерация имени автора
CREATE OR REPLACE FUNCTION random_author_name() RETURNS VARCHAR(100) AS $$
DECLARE
    first_names VARCHAR[] := ARRAY['Иван', 'Петр', 'Алексей', 'Михаил', 'Сергей', 'Андрей', 'Дмитрий', 'Николай', 'Александр', 'Владимир', 'Екатерина', 'Анна', 'Мария', 'Ольга', 'Наталья', 'Елена', 'Татьяна', 'Ирина', 'Светлана', 'Юлия'];
    last_names VARCHAR[] := ARRAY['Иванов', 'Петров', 'Сидоров', 'Смирнов', 'Кузнецов', 'Попов', 'Васильев', 'Федоров', 'Морозов', 'Волков', 'Алексеев', 'Лебедев', 'Семенов', 'Егоров', 'Павлов', 'Козлов', 'Степанов', 'Николаев', 'Орлов', 'Андреев'];
BEGIN
    RETURN first_names[1 + floor(random() * array_length(first_names, 1))] || ' ' || 
           last_names[1 + floor(random() * array_length(last_names, 1))];
END;
$$ LANGUAGE plpgsql;

-- Генерация 500 авторов
DO $$
DECLARE
    i INTEGER;
    birth_year INT;
    death_year INT;
BEGIN
    FOR i IN 1..500 LOOP
        birth_year := 1200 + floor(random() * 800)::INT;
        death_year := birth_year + 20 + floor(random() * 80)::INT;
        
        INSERT INTO Author (id, name, birthYear, deathYear)
        VALUES (
            gen_random_uuid(),
            random_author_name(),
            birth_year,
            death_year
        );
    END LOOP;
END $$;

-- Генерация 100 коллекций
DO $$
DECLARE
    i INTEGER;
    styles VARCHAR[] := ARRAY['Импрессионизм', 'Экспрессионизм', 'Кубизм', 'Сюрреализм', 'Футуризм', 'Дадаизм', 'Поп-арт', 'Минимализм', 'Концептуализм', 'Барокко', 'Рококо', 'Классицизм', 'Романтизм', 'Реализм', 'Символизм', 'Модерн', 'Постимпрессионизм', 'Фовизм', 'Абстракционизм', 'Супрематизм'];
    prefixes VARCHAR[] := ARRAY['Великие произведения', 'Шедевры', 'Коллекция', 'Сокровища', 'Архив', 'Галерея', 'Музей', 'Собрание', 'Альбом', 'Выставка'];
BEGIN
    FOR i IN 1..100 LOOP
        INSERT INTO Collection (id, title)
        VALUES (
            gen_random_uuid(),
            prefixes[1 + floor(random() * array_length(prefixes, 1))] || ' ' || 
            styles[1 + floor(random() * array_length(styles, 1))]
        );
    END LOOP;
END $$;

-- Название экспоната
CREATE OR REPLACE FUNCTION random_artwork_title() RETURNS VARCHAR(255) AS $$
DECLARE
    prefixes VARCHAR[] := ARRAY['Портрет', 'Пейзаж', 'Натюрморт', 'Композиция', 'Этюд', 'Эскиз', 'Абстракция', 'Импровизация', 'Фантазия', 'Вид', 'Сон', 'Мечта', 'Воспоминание', 'Ода', 'Поэма', 'Симфония', 'Ритм', 'Гармония', 'Контраст', 'Форма'];
    suffixes VARCHAR[] := ARRAY['весны', 'осени', 'зимы', 'лета', 'утра', 'вечера', 'дня', 'ночи', 'света', 'тьмы', 'цвета', 'линии', 'формы', 'пространства', 'времени', 'движения', 'покоя', 'радости', 'грусти', 'любви'];
BEGIN
    RETURN prefixes[1 + floor(random() * array_length(prefixes, 1))] || ' ' || 
           suffixes[1 + floor(random() * array_length(suffixes, 1))];
END;
$$ LANGUAGE plpgsql;

-- Техника
CREATE OR REPLACE FUNCTION random_technique() RETURNS VARCHAR(100) AS $$
DECLARE
    techniques VARCHAR[] := ARRAY['Масляная живопись', 'Акварель', 'Гуашь', 'Темпера', 'Акрил', 'Гравюра', 'Литография', 'Офорт', 'Фреска', 'Мозаика', 'Витраж', 'Коллаж', 'Ассамбляж', 'Инсталляция', 'Цифровое искусство', 'Фотография', 'Скульптура', 'Резьба', 'Лепка', 'Рисунок'];
BEGIN
    RETURN techniques[1 + floor(random() * array_length(techniques, 1))];
END;
$$ LANGUAGE plpgsql;

-- Материал
CREATE OR REPLACE FUNCTION random_material() RETURNS VARCHAR(100) AS $$
DECLARE
    materials VARCHAR[] := ARRAY['Холст, масло', 'Бумага, акварель', 'Дерево, масло', 'Картон, гуашь', 'Металл', 'Камень', 'Глина', 'Стекло', 'Ткань', 'Пластик', 'Комбинированные материалы', 'Золото', 'Серебро', 'Бронза', 'Мрамор', 'Гранит', 'Керамика', 'Фарфор', 'Воск', 'Папье-маше'];
BEGIN
    RETURN materials[1 + floor(random() * array_length(materials, 1))];
END;
$$ LANGUAGE plpgsql;

-- Генерация 10,000 произведений искусства
DO $$
DECLARE
    i INTEGER;
    author_rec RECORD;
    collection_rec RECORD;
    creation_year INT;
    size_w INT;
    size_h INT;
BEGIN
    FOR i IN 1..10000 LOOP
        -- Выбираем случайного автора
        SELECT id INTO author_rec FROM Author ORDER BY random() LIMIT 1;
        
        -- Выбираем случайную коллекцию
        SELECT id INTO collection_rec FROM Collection ORDER BY random() LIMIT 1;
        
        -- Генерируем год создания (между годом рождения автора + 20 и годом смерти)
        creation_year := (SELECT birthYear + 20 + floor(random() * (deathYear - birthYear - 20))::INT 
                          FROM Author WHERE id = author_rec.id);
        
        -- Генерируем размеры
        size_w := 10 + floor(random() * 500)::INT;
        size_h := 10 + floor(random() * 500)::INT;
        
        INSERT INTO Artworks (id, title, technic, material, size, creationYear, authorID, collectionID)
        VALUES (
            gen_random_uuid(),
            random_artwork_title(),
            random_technique(),
            random_material(),
            size_w || ' × ' || size_h || ' см',
            creation_year,
            author_rec.id,
            collection_rec.id
        );
        
        -- Выводим прогресс каждые 1000 записей
        IF i % 1000 = 0 THEN
            RAISE NOTICE 'Добавлено % произведений', i;
        END IF;
    END LOOP;
END $$;

-- Случайное имя
CREATE OR REPLACE FUNCTION random_name() RETURNS VARCHAR(50) AS $$
DECLARE
    first_names VARCHAR[] := ARRAY['Иван', 'Алексей', 'Дмитрий', 'Сергей', 'Андрей', 'Михаил', 'Екатерина', 'Анна', 'Мария', 'Ольга', 'Наталья', 'Елена'];
    last_names VARCHAR[] := ARRAY['Иванов', 'Петров', 'Сидоров', 'Смирнов', 'Кузнецов', 'Попов', 'Васильев', 'Морозов', 'Волков', 'Лебедев'];
BEGIN
    RETURN first_names[1 + floor(random() * array_length(first_names, 1))] || ' ' || 
           last_names[1 + floor(random() * array_length(last_names, 1))];
END;
$$ LANGUAGE plpgsql;

-- Случайный логин
CREATE OR REPLACE FUNCTION random_login(name VARCHAR) RETURNS VARCHAR(50) AS $$
BEGIN
    RETURN lower(replace(name, ' ', '_')) || floor(random() * 1000)::INT;
END;
$$ LANGUAGE plpgsql;

-- Случайный email
CREATE OR REPLACE FUNCTION random_email(name VARCHAR) RETURNS VARCHAR(100) AS $$
BEGIN
    RETURN lower(replace(name, ' ', '.')) || floor(random() * 1000)::INT || '@' || 
           (ARRAY['gmail.com', 'yahoo.com', 'mail.ru', 'yandex.ru'])[1 + floor(random() * 4)];
END;
$$ LANGUAGE plpgsql;

-- Случайны хеш пароля (имитация)
CREATE OR REPLACE FUNCTION random_password_hash() RETURNS VARCHAR(255) AS $$
BEGIN
    RETURN 'hashed_password_' || substr(md5(random()::text), 0, 20);
END;
$$ LANGUAGE plpgsql;

-- Заполнение таблицы Admins (5 записей)
INSERT INTO Admins (id, username, login, hashedPassword, createdAt, valid)
SELECT 
    gen_random_uuid(),
    random_name(),
    random_login(random_name()),
    random_password_hash(),
    NOW() - (random() * 365 * 5 || ' days')::INTERVAL,
    random() > 0.1 -- 90% активных
FROM generate_series(1, 5);

-- Заполнение таблицы Employees (100 записей)
INSERT INTO Employees (id, username, login, hashedPassword, createdAt, valid, adminID)
SELECT 
    gen_random_uuid(),
    random_name(),
    random_login(random_name()),
    random_password_hash(),
    NOW() - (random() * 365 * 3 || ' days')::INTERVAL,
    random() > 0.2, -- 80% активных
    (SELECT id FROM Admins ORDER BY random() LIMIT 1)
FROM generate_series(1, 100);

-- Заполнение таблицы Users (5000 записей)
INSERT INTO Users (id, username, login, hashedPassword, createdAt, email, subscribeMail)
SELECT 
    gen_random_uuid(),
    random_name(),
    random_login(random_name()),
    random_password_hash(),
    NOW() - (random() * 365 * 5 || ' days')::INTERVAL,
    random_email(random_name()),
    random() > 0.7 -- 30% подписаны
FROM generate_series(1, 5000)
ON CONFLICT DO NOTHING;

-- Заполнение таблицы Events (200 записей)
INSERT INTO Events (id, title, dateBegin, dateEnd, canVisit, adress, cntTickets, creatorID)
SELECT 
    gen_random_uuid(),
    'Выставка "' || 
    (ARRAY['Импрессионистов', 'Современного искусства', 'Русских мастеров', 'Европейской живописи', 'Авангарда'])[1 + floor(random() * 5)] || 
    '" в ' || 
    (ARRAY['Москве', 'Санкт-Петербурге', 'Париже', 'Нью-Йорке', 'Лондоне', 'Берлине'])[1 + floor(random() * 6)],
    NOW() + (random() * 30 || ' days')::INTERVAL,
    NOW() + ((30 + random() * 60) || ' days')::INTERVAL,
    random() > 0.3, -- 70% можно посетить
    'ул. ' || 
    (ARRAY['Ленина', 'Пушкина', 'Гоголя', 'Толстого', 'Достоевского'])[1 + floor(random() * 5)] || 
    ', д. ' || (1 + floor(random() * 100))::TEXT,
    50 + floor(random() * 500)::INT, -- от 50 до 550 билетов
    (SELECT id FROM Employees WHERE valid = true ORDER BY random() LIMIT 1)
FROM generate_series(1, 200);

-- Заполнение таблицы Artwork_event
INSERT INTO Artwork_event (artworkID, eventID)
SELECT 
    (SELECT id FROM Artworks ORDER BY random() LIMIT 1),
    (SELECT id FROM Events ORDER BY random() LIMIT 1)
FROM generate_series(1, 3000)
ON CONFLICT DO NOTHING; -- пропускать дубликаты

-- Заполнение таблицы TicketPurchases (5000 записей)
INSERT INTO TicketPurchases (id, customerName, customerEmail, purchaseDate, eventID)
SELECT 
    gen_random_uuid(),
    random_name(),
    random_email(random_name()),
    NOW() - (random() * 30 || ' days')::INTERVAL,
    (SELECT id FROM Events WHERE canVisit = true ORDER BY random() LIMIT 1)
FROM generate_series(1, 5000);

-- Заполнение таблицы tickets_user
-- Сначала обновим TicketPurchases, добавив userID для 70% билетов
UPDATE TicketPurchases 
SET customerEmail = random_email(customerName),
    customerName = random_name()
WHERE random() > 0.3; -- обновим 70% записей

-- Затем создадим связи
INSERT INTO tickets_user (ticketID, userID)
SELECT 
    tp.id,
    (SELECT id FROM Users ORDER BY random() LIMIT 1)
FROM TicketPurchases tp
WHERE random() > 0.5 -- 50% билетов с привязкой к пользователям
ON CONFLICT DO NOTHING;







