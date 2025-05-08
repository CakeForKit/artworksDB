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
-- Создаем временную функцию для генерации случайных имен
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

-- Создаем временные функции для генерации данных
CREATE OR REPLACE FUNCTION random_artwork_title() RETURNS VARCHAR(255) AS $$
DECLARE
    prefixes VARCHAR[] := ARRAY['Портрет', 'Пейзаж', 'Натюрморт', 'Композиция', 'Этюд', 'Эскиз', 'Абстракция', 'Импровизация', 'Фантазия', 'Вид', 'Сон', 'Мечта', 'Воспоминание', 'Ода', 'Поэма', 'Симфония', 'Ритм', 'Гармония', 'Контраст', 'Форма'];
    suffixes VARCHAR[] := ARRAY['весны', 'осени', 'зимы', 'лета', 'утра', 'вечера', 'дня', 'ночи', 'света', 'тьмы', 'цвета', 'линии', 'формы', 'пространства', 'времени', 'движения', 'покоя', 'радости', 'грусти', 'любви'];
BEGIN
    RETURN prefixes[1 + floor(random() * array_length(prefixes, 1))] || ' ' || 
           suffixes[1 + floor(random() * array_length(suffixes, 1))];
END;
$$ LANGUAGE plpgsql;

CREATE OR REPLACE FUNCTION random_technique() RETURNS VARCHAR(100) AS $$
DECLARE
    techniques VARCHAR[] := ARRAY['Масляная живопись', 'Акварель', 'Гуашь', 'Темпера', 'Акрил', 'Гравюра', 'Литография', 'Офорт', 'Фреска', 'Мозаика', 'Витраж', 'Коллаж', 'Ассамбляж', 'Инсталляция', 'Цифровое искусство', 'Фотография', 'Скульптура', 'Резьба', 'Лепка', 'Рисунок'];
BEGIN
    RETURN techniques[1 + floor(random() * array_length(techniques, 1))];
END;
$$ LANGUAGE plpgsql;

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