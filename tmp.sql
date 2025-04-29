-- Active: 1744740356603@@127.0.0.1@5432@artworks
-- SELECT TABLE_NAME
-- FROM INFORMATION_SCHEMA.TABLES

select * from events

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
                    '750f41af-0125-4807-8515-fed3828e2f0e' == ae.eventID)


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

CREATE TABLE Event (
    id UUID PRIMARY KEY,
    title VARCHAR(255) NOT NULL,
    dateBegin DATETIME NOT NULL,
    dateEnd DATETIME NOT NULL,
    access BOOLEAN,
    adress VARCHAR(255),
    cntTickets INT DEFAULT 0,
    creatorID INT NOT NULL,
    FOREIGN KEY (creatorID) REFERENCES Employee(id)
);

CREATE TABLE Artwork_event (
    artworkID UUID NOT NULL,
    eventID UUID NOT NULL,
    FOREIGN KEY (artworkID) REFERENCES Artwork(id),
    FOREIGN KEY (eventID) REFERENCES Event(id)
);

CREATE TABLE TicketPurchases (
    id UUID PRIMARY KEY,
    customerName VARCHAR(100),
    customerEmail VARCHAR(100),
    eventID UUID NOT NULL,
    purchaseDate DATETIME NOT NULL,
    FOREIGN KEY (eventID) REFERENCES Event(id)
);