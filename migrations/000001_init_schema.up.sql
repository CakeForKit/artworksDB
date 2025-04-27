-- -- Active: 1744740356603@@127.0.0.1@5432@artworks
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

CREATE TABLE Employee (
    id UUID PRIMARY KEY,
    username VARCHAR(50) NOT NULL,
    login VARCHAR(50) NOT NULL UNIQUE,
    hashedPassword VARCHAR(255) NOT NULL,
    createdAt TIMESTAMP NOT NULL,
    valid BOOLEAN NOT NULL DEFAULT TRUE,
    adminID UUID NOT NULL,
    FOREIGN KEY (adminID) REFERENCES Admins(id)
);
ALTER TABLE Employee ADD CONSTRAINT emptyCheck 
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
ALTER TABLE Artwork ADD CONSTRAINT emptyCheck 
    CHECK(title != '' AND creationYear > 0); 

CREATE TABLE Events (
    id UUID PRIMARY KEY,
    title VARCHAR(255) NOT NULL,
    dateBegin TIMESTAMP NOT NULL,
    dateEnd TIMESTAMP NOT NULL,
    access BOOLEAN,
    adress VARCHAR(255),
    cntTickets INT,
    creatorID UUID NOT NULL,
    FOREIGN KEY (creatorID) REFERENCES Employee(id)
);
ALTER TABLE Events ADD CONSTRAINT emptyCheck 
    CHECK(title != '' AND adress != ''); 
ALTER TABLE Events ADD CONSTRAINT dateBeginEndCheck 
    CHECK(dateBegin < dateEnd);


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
    purchaseDate TIMESTAMP NOT NULL,
    eventID UUID NOT NULL,
    userID UUID DEFAULT NULL,
    FOREIGN KEY (eventID) REFERENCES Events(id),
    FOREIGN KEY (userID) REFERENCES Users(id)
);
ALTER TABLE TicketPurchases ADD CONSTRAINT emptyCheck 
    CHECK(customerName != '' AND customerEmail != ''); 