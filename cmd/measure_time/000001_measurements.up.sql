-- Active: 1744740356603@@127.0.0.1@5432@artworks

CREATE TABLE Admins (
    id UUID PRIMARY KEY,
    username VARCHAR(50) NOT NULL,
    login VARCHAR(50) NOT NULL UNIQUE,
    hashedPassword VARCHAR(255) NOT NULL,
    createdAt TIMESTAMP NOT NULL,
    valid BOOLEAN NOT NULL DEFAULT TRUE
);
ALTER TABLE Admins ADD CONSTRAINT emptyCheck 
    CHECK(username != '' AND login != '' AND hashedPassword != ''); 

CREATE TABLE Employees (
    id UUID PRIMARY KEY,
    username VARCHAR(50) NOT NULL,
    login VARCHAR(50) NOT NULL UNIQUE,
    hashedPassword VARCHAR(255) NOT NULL,
    createdAt TIMESTAMP NOT NULL,
    valid BOOLEAN NOT NULL DEFAULT TRUE,
    adminID UUID NOT NULL,
    FOREIGN KEY (adminID) REFERENCES Admins(id)
);
ALTER TABLE Employees ADD CONSTRAINT emptyCheck 
    CHECK(username != '' AND login != '' AND hashedPassword != ''); 

CREATE TABLE Users (
    id UUID PRIMARY KEY,
    username VARCHAR(50) NOT NULL,
    login VARCHAR(50) NOT NULL UNIQUE,
    hashedPassword VARCHAR(255) NOT NULL,
    createdAt TIMESTAMP NOT NULL,
    email VARCHAR(100) UNIQUE,
    subscribeMail BOOLEAN DEFAULT FALSE
);
ALTER TABLE Users ADD CONSTRAINT emptyCheck 
    CHECK(username != '' AND login != '' AND hashedPassword != ''); 

CREATE TABLE Author (
    id UUID PRIMARY KEY,
    name VARCHAR(100) NOT NULL,
    birthYear INT,
    deathYear INT
);
ALTER TABLE Author ADD CONSTRAINT emptyCheck CHECK(name != ''); 
ALTER TABLE Author ADD CONSTRAINT birthDeathYear 
    CHECK(birthYear < deathYear AND 
            birthYear > 0);

CREATE TABLE Collection (
    id UUID PRIMARY KEY,
    title VARCHAR(255) NOT NULL
);
ALTER TABLE Collection ADD CONSTRAINT titleEmpty CHECK(title != '');

CREATE TABLE Artworks (
    id UUID PRIMARY KEY,
    title VARCHAR(255) NOT NULL,
    technic VARCHAR(100),
    material VARCHAR(100),
    size VARCHAR(50),
    creationYear INT,
    authorID UUID NOT NULL,
    collectionID UUID NOT NULL,
    FOREIGN KEY (authorID) REFERENCES Author(id),
    FOREIGN KEY (collectionID) REFERENCES Collection(id)
);
ALTER TABLE Artworks ADD CONSTRAINT emptyCheck 
    CHECK(title != '' AND creationYear > 0); 

CREATE TABLE Events (
    id UUID PRIMARY KEY,
    title VARCHAR(255) NOT NULL,
    dateBegin TIMESTAMP NOT NULL,
    dateEnd TIMESTAMP NOT NULL,
    canVisit BOOLEAN,
    adress VARCHAR(255),
    cntTickets INT,
    creatorID UUID NOT NULL,
    FOREIGN KEY (creatorID) REFERENCES Employees(id)
);
ALTER TABLE Events ADD CONSTRAINT emptyCheck 
    CHECK(title != '' AND adress != ''); 
ALTER TABLE Events ADD CONSTRAINT dateBeginEndCheck 
    CHECK(dateBegin < dateEnd);


CREATE TABLE Artwork_event (
    artworkID UUID NOT NULL,
    eventID UUID NOT NULL,
    PRIMARY KEY (artworkID, eventID),
    CONSTRAINT artworkID_notnull CHECK (artworkID IS NOT NULL),
    FOREIGN KEY (artworkID) REFERENCES Artworks(id) ON DELETE CASCADE,
    CONSTRAINT eventID_notnull CHECK (eventID IS NOT NULL),
    FOREIGN KEY (eventID) REFERENCES Events(id) ON DELETE CASCADE
);

CREATE TABLE TicketPurchases (
    id UUID PRIMARY KEY,
    customerName VARCHAR(100),
    customerEmail VARCHAR(100),
    purchaseDate TIMESTAMP NOT NULL,
    eventID UUID NOT NULL,
    FOREIGN KEY (eventID) REFERENCES Events(id)
);

ALTER TABLE TicketPurchases ADD CONSTRAINT emptyCheck 
    CHECK(customerName != '' AND customerEmail != ''); 

CREATE TABLE tickets_user (
    ticketID UUID UNIQUE,
    userID UUID,
    PRIMARY KEY (ticketID, userID),
    CONSTRAINT ticketID_notnull CHECK (ticketID IS NOT NULL),
    FOREIGN KEY (ticketID) REFERENCES TicketPurchases(id) ON DELETE CASCADE,
    CONSTRAINT userID_notnull CHECK (userID IS NOT NULL),
    FOREIGN KEY (userID) REFERENCES Users(id) ON DELETE CASCADE
);


CREATE OR REPLACE FUNCTION check_ticket_limit()
RETURNS TRIGGER AS $$
DECLARE
    max_tickets INT;
    sold_tickets INT;
BEGIN

    SELECT cntTickets INTO max_tickets
    FROM Events
    WHERE id = NEW.eventID;
    
    SELECT COUNT(*) INTO sold_tickets
    FROM TicketPurchases
    WHERE eventID = NEW.eventID;

    IF TG_OP = 'INSERT' THEN
        sold_tickets := sold_tickets + 1;
    END IF;
    
    IF sold_tickets > max_tickets THEN
        RAISE EXCEPTION 'Превышено максимальное количество билетов для события (доступно: %, пытается купить: %)', 
                        max_tickets, sold_tickets;
    END IF;
    
    RETURN NEW;
END;
$$ LANGUAGE plpgsql;


CREATE TRIGGER enforce_ticket_limit
BEFORE INSERT OR UPDATE ON TicketPurchases
FOR EACH ROW
EXECUTE FUNCTION check_ticket_limit();

CREATE OR REPLACE FUNCTION get_event_of_artwork(
    idArtwork UUID, 
    dateBeginSee TIMESTAMP, 
    dateEndSee TIMESTAMP)
RETURNS TABLE (
    event_id UUID,
    title VARCHAR(255),
    dateBegin TIMESTAMP,
    dateEnd TIMESTAMP,
    canVisit BOOLEAN,
    adress VARCHAR(255),
    cntTickets INT,
    creatorID UUID
) AS $$

    SELECT e.id, e.title, e.dateBegin, e.dateEnd, e.canVisit, e.adress, e.cntTickets, e.creatorID
    FROM Events e
    JOIN Artwork_event ae ON e.id = ae.eventID
    WHERE ae.artworkID = idArtwork
      AND e.dateBegin <= dateEndSee
      AND e.dateEnd >= dateBeginSee;

$$ LANGUAGE sql;


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
    RAISE NOTICE 'Добавлено % Author', (SELECT  COUNT(*) FROM Author LIMIT 1);
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
    RAISE NOTICE 'Добавлено % Collection', (SELECT  COUNT(*) FROM Collection);
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

-- Заполнение таблицы Events (500 записей)
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
FROM generate_series(1, 500);

-- Заполнение таблицы Artwork_event
CREATE OR REPLACE FUNCTION random_artwork_id() RETURNS UUID AS $$
BEGIN
    RETURN (SELECT id FROM Artworks ORDER BY random() LIMIT 1);
END;
$$ LANGUAGE plpgsql;

CREATE OR REPLACE FUNCTION random_event_id() RETURNS UUID AS $$
BEGIN
    RETURN (SELECT id FROM Events ORDER BY random() LIMIT 1);
END;
$$ LANGUAGE plpgsql;

-- INSERT INTO Artwork_event (artworkID, eventID)
-- SELECT random_artwork_id(), random_event_id()
-- FROM generate_series(1, 3000)
-- ON CONFLICT DO NOTHING;

CREATE OR REPLACE FUNCTION insert_artwork_events(
    target_count INTEGER DEFAULT 1000
) 
RETURNS VOID AS $$
DECLARE
    inserted_count INTEGER := 0;
    max_attempts INTEGER := target_count * 10; -- Защита от бесконечного цикла
    attempts INTEGER := 0;
    new_artwork_id UUID;
    new_event_id UUID;
BEGIN
    IF target_count < 1 THEN
        RAISE EXCEPTION 'Целевое количество записей должно быть положительным числом';
    END IF;
    
    WHILE inserted_count < target_count AND attempts < max_attempts LOOP
        new_artwork_id := random_artwork_id();
        new_event_id := random_event_id();
        
        BEGIN
            INSERT INTO Artwork_event (artworkID, eventID)
            VALUES (new_artwork_id, new_event_id);
            
            inserted_count := inserted_count + 1;
        EXCEPTION WHEN unique_violation THEN
            -- Дубликат, просто продолжаем
        END;
        
        attempts := attempts + 1;
    END LOOP;
    
    RAISE NOTICE 'Вставлено % записей из % после % попыток', 
        inserted_count, target_count, attempts;
END;
$$ LANGUAGE plpgsql;

-- Заполнение таблицы TicketPurchases (5000 записей)
DO $$
DECLARE
    event_rec RECORD;
    tickets_to_sell INT;
    tickets_available INT;
    tickets_sold INT;
    i INT;
BEGIN
    FOR event_rec IN SELECT id, cntTickets FROM Events WHERE canVisit = true LOOP
        -- Сколько билетов уже продано
        SELECT COUNT(*) INTO tickets_sold 
        FROM TicketPurchases 
        WHERE eventID = event_rec.id;
        
        -- Сколько еще можно продать
        tickets_available := event_rec.cntTickets - tickets_sold;
        
        -- Продаем случайное количество (но не больше доступного)
        IF tickets_available > 0 THEN
            tickets_to_sell := least(floor(random() * 20)::INT + 1, tickets_available);
            
            FOR i IN 1..tickets_to_sell LOOP
                INSERT INTO TicketPurchases (id, customerName, customerEmail, purchaseDate, eventID)
                VALUES (
                    gen_random_uuid(),
                    random_name(),
                    random_email(random_name()),
                    NOW() - (random() * 30 || ' days')::INTERVAL,
                    event_rec.id
                );
            END LOOP;
        END IF;
    END LOOP;
    RAISE NOTICE 'Добавлено % TicketPurchases', (SELECT  COUNT(*) FROM TicketPurchases);
END $$;

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
