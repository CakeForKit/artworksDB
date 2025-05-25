-- Active: 1744740356603@@127.0.0.1@5432@artworks

-- SELECT TABLE_NAME
-- FROM INFORMATION_SCHEMA.TABLES

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
    valid BOOLEAN NOT NULL DEFAULT TRUE,
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

CREATE OR REPLACE FUNCTION validate_artwork_creation_year()
RETURNS TRIGGER AS $$
DECLARE
    author_birth INT;
    author_death INT;
BEGIN
    -- Получаем годы жизни автора
    SELECT birthYear, deathYear INTO author_birth, author_death
    FROM Author
    WHERE id = NEW.authorID;
    
    -- Проверяем, что год создания artwork находится после года рождения автора
    IF NEW.creationYear < author_birth THEN
        RAISE EXCEPTION 'Год создания произведения (%) не может быть раньше года рождения автора (%)', 
                        NEW.creationYear, author_birth;
    END IF;
    
    -- Если год смерти автора указан, проверяем что год создания не позже
    IF author_death IS NOT NULL AND NEW.creationYear > author_death THEN
        RAISE EXCEPTION 'Год создания произведения (%) не может быть позже года смерти автора (%)', 
                        NEW.creationYear, author_death;
    END IF;
    
    RETURN NEW;
END;
$$ LANGUAGE plpgsql;

-- Триггер для вставки
CREATE TRIGGER check_artwork_year_on_insert
BEFORE INSERT ON Artworks
FOR EACH ROW
EXECUTE FUNCTION validate_artwork_creation_year();

-- Триггер для обновления
CREATE TRIGGER check_artwork_year_on_update
BEFORE UPDATE ON Artworks
FOR EACH ROW
WHEN (NEW.creationYear IS DISTINCT FROM OLD.creationYear OR NEW.authorID IS DISTINCT FROM OLD.authorID)
EXECUTE FUNCTION validate_artwork_creation_year();

CREATE OR REPLACE FUNCTION validate_existing_artworks_on_author_update()
RETURNS TRIGGER AS $$
DECLARE
    invalid_artwork RECORD;
BEGIN
    -- Проверяем все произведения этого автора при изменении дат жизни
    IF (NEW.birthYear IS DISTINCT FROM OLD.birthYear) OR 
       (NEW.deathYear IS DISTINCT FROM OLD.deathYear) THEN
        
        -- Ищем произведения, которые выходят за новые границы жизни автора
        FOR invalid_artwork IN 
            SELECT a.id, a.title, a.creationYear
            FROM Artworks a
            WHERE a.authorID = NEW.id
              AND (a.creationYear < NEW.birthYear OR 
                  (NEW.deathYear IS NOT NULL AND a.creationYear > NEW.deathYear))
        LOOP
            RAISE EXCEPTION 
                'Произведение "%" (год создания: %) не соответствует новым датам жизни автора (рождение: %, смерть: %)',
                invalid_artwork.title, invalid_artwork.creationYear, NEW.birthYear, NEW.deathYear;
        END LOOP;
    END IF;
    
    RETURN NEW;
END;
$$ LANGUAGE plpgsql;

-- Триггер для обновления автора
CREATE TRIGGER check_artworks_on_author_update
BEFORE UPDATE ON Author
FOR EACH ROW
EXECUTE FUNCTION validate_existing_artworks_on_author_update();