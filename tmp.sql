-- Active: 1744740356603@@127.0.0.1@5432@artworks
SELECT *
FROM INFORMATION_SCHEMA.TABLES
WHERE schemaname = 'public';

select * from Artwork_event

select * 
from public.artworks

-- Исследуемый запрос
EXPLAIN ANALYZE
SELECT public.Artworks.title, Author.name
FROM public.Artworks
JOIN public.Author
ON Artworks.authorID = Author.id

CREATE INDEX idx_artworks_authorid ON Artworks(authorID);

DROP INDEX IF EXISTS idx_artworks_authorid;

-- Все индексы в базе
SELECT * FROM pg_indexes
WHERE schemaname = 'public';

EXPLAIN ANALYZE
SELECT a.title, e.title, e.dateBegin, e.dateEnd
FROM Events e
JOIN Artwork_event ae
ON e.id = ae.eventID
JOIN artworks a
ON ae.artworkID = a.id

select a.id, a.title, e.title, e.dateBegin, e.dateEnd
from artworks a
join Artwork_event ae
on a.id = ae.artworkid
join events e
on ae.eventid = e.id


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
    JOIN Artwork_event ae 
    ON e.id = ae.eventID
    WHERE ae.artworkID = idArtwork
      AND e.dateBegin <= dateEndSee
      AND e.dateEnd >= dateBeginSee;

$$ LANGUAGE sql;


select tp.id, tp.customername, tp.customeremail, tp.purchasedate, tp.eventid, tu.userid
from TicketPurchases tp
join tickets_user tu
on tp.id = tu.ticketID
where tu.userID = '8472a25c-464e-4730-9b2f-d4970a838310'

INSERT INTO TicketPurchases (id, customerName, customerEmail, purchaseDate, eventID)
VALUES ()


SELECT 
    e.id,
    COUNT(tp.id) AS tickets_sold
FROM Events e
LEFT JOIN TicketPurchases tp 
ON e.id = tp.eventID
GROUP BY e.id
ORDER BY 
    e.dateBegin;

SELECT COUNT(tp.id)
    -- COUNT(tp.id) AS tickets_sold
FROM Events e
LEFT JOIN TicketPurchases tp 
ON e.id = tp.eventID
WHERE e.id = 'e851464a-3c58-4a19-b269-58dbf619f01d';

-- Посмотрим количество билетов на каждое событие
SELECT 
    e.id, e.title AS event_title,
    e.cntTickets AS max_tickets,
    COUNT(tp.id) AS tickets_sold,
    e.cntTickets - COUNT(tp.id) AS tickets_available
FROM 
    Events e
LEFT JOIN 
    TicketPurchases tp ON e.id = tp.eventID
GROUP BY 
    e.id, e.title, e.cntTickets
ORDER BY 
    e.dateBegin;

-- Посмотрим распределение билетов среди пользователей
SELECT 
    u.id, u.username,
    COUNT(tu.ticketID) AS tickets_purchased,
    STRING_AGG(e.title, ', ') AS events_attending
FROM 
    Users u
JOIN 
    tickets_user tu ON u.id = tu.userID
JOIN 
    TicketPurchases tp ON tu.ticketID = tp.id
JOIN 
    Events e ON tp.eventID = e.id
GROUP BY 
    u.id, u.username
ORDER BY 
    tickets_purchased DESC;

drop Function if exists get_event_of_artwork

SELECT e.id, e.title, e.dateBegin, e.dateEnd, e.canVisit, e.adress, e.cntTickets, e.creatorID
    FROM Events e
    JOIN Artwork_event ae ON e.id = ae.eventID
    WHERE ae.artworkID = '30154661-36c5-4761-96ea-691abb9bb407'
      AND e.dateBegin <= '2025-05-22 00:00:00'
      AND e.dateEnd >= '2025-05-01 00:00:00';

select event_id, title, datebegin
from get_event_of_artwork('30154661-36c5-4761-96ea-691abb9bb407', '2025-04-01 00:00:00', '2025-06-22 00:00:00')

select id, title, dateBegin, dateEnd, canVisit, adress, cntTickets, creatorID
from events
WHERE 
    dateBegin >= '2025-05-01 00:00:00'::timestamp AND
    dateEnd <= '2023-11-30 23:59:59'::timestamp
-- ORDER BY dateBegin;

select art.id, art.title, art.technic, art.material, art.size, art.creationYear, 
    au.id, au.name, au.birthyear, au.deathyear, col.id, col.title
from artwork art
join author au
on art.authorid = au.id
join collection col
on art.collectionid = col.id
where exists (select 1
                from Artwork_event ae
                where art.id = ae.artworkID and 
                    '750f41af-0125-4807-8515-fed3828e2f0e' = ae.eventID)


-- Создание таблицы User
CREATE TABLE User (
    id UUID PRIMARY KEY,
    username VARCHAR(50) NOT NULL,
    login VARCHAR(50) NOT NULL UNIQUE,
    hashedPassword VARCHAR(255) NOT NULL,
    createdAt DATETIME NOT NULL,
    email VARCHAR(100) UNIQUE,
    subscribeMail BOOLEAN DEFAULT FALSE
);

CREATE TABLE Admin (
    id UUID PRIMARY KEY,
    username VARCHAR(50) NOT NULL,
    login VARCHAR(50) NOT NULL UNIQUE,
    hashedPassword VARCHAR(255) NOT NULL,
    createdAt DATETIME NOT NULL,
    valid BOOLEAN NOT NULL DEFAULT TRUE
);

CREATE TABLE Employee (
    id UUID PRIMARY KEY,
    username VARCHAR(50) NOT NULL,
    login VARCHAR(50) NOT NULL UNIQUE,
    hashedPassword VARCHAR(255) NOT NULL,
    createdAt DATETIME NOT NULL,
    valid BOOLEAN NOT NULL DEFAULT TRUE,
    adminID UUID NOT NULL,
    FOREIGN KEY (adminID) REFERENCES Admin(id)
);

CREATE TABLE Author (
    id UUID PRIMARY KEY,
    name VARCHAR(100) NOT NULL,
    birthYear INT,
    deathYear INT
);

CREATE TABLE Collection (
    id UUID PRIMARY KEY,
    title VARCHAR(255) NOT NULL
);

CREATE TABLE Artwork (
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

CREATE TABLE Events (
    id UUID PRIMARY KEY,
    title VARCHAR(255) NOT NULL,
    dateBegin DATETIME NOT NULL,
    dateEnd DATETIME NOT NULL,
    canVisit BOOLEAN,
    adress VARCHAR(255),
    cntTickets INT DEFAULT 0,
    creatorID INT NOT NULL,
    FOREIGN KEY (creatorID) REFERENCES Employee(id)
);

CREATE TABLE Artwork_event (
    artworkID UUID NOT NULL,
    eventID UUID NOT NULL,
    FOREIGN KEY (artworkID) REFERENCES Artwork(id),
    FOREIGN KEY (eventID) REFERENCES Events(id)
);

CREATE TABLE TicketPurchases (
    id UUID PRIMARY KEY,
    customerName VARCHAR(100),
    customerEmail VARCHAR(100),
    eventID UUID NOT NULL,
    purchaseDate DATETIME NOT NULL,
    FOREIGN KEY (eventID) REFERENCES Events(id)
);