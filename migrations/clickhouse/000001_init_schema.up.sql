-- Active: 1748883456071@@127.0.0.1@8123@artworks
-- Создаем базу данных

-- Таблица Admins
CREATE TABLE IF NOT EXISTS artworks.Admins
(
    id UUID,
    username String,
    login String,
    hashedPassword String,
    createdAt DateTime,
    valid UInt8 DEFAULT 1,
    CONSTRAINT emptyCheck CHECK empty(username) = 0 AND empty(login) = 0 AND empty(hashedPassword) = 0
)
ENGINE = MergeTree()
ORDER BY id
PRIMARY KEY id;

-- Таблица Employees
CREATE TABLE IF NOT EXISTS artworks.Employees
(
    id UUID,
    username String,
    login String,
    hashedPassword String,
    createdAt DateTime,
    valid UInt8 DEFAULT 1,
    adminID UUID,
    CONSTRAINT emptyCheck CHECK empty(username) = 0 AND empty(login) = 0 AND empty(hashedPassword) = 0
)
ENGINE = MergeTree()
ORDER BY id
PRIMARY KEY id;

-- Таблица Users
CREATE TABLE IF NOT EXISTS artworks.Users
(
    id UUID,
    username String,
    login String,
    hashedPassword String,
    createdAt DateTime,
    email Nullable(String),
    subscribeMail UInt8 DEFAULT 0,
    CONSTRAINT emptyCheck CHECK empty(username) = 0 AND empty(login) = 0 AND empty(hashedPassword) = 0
)
ENGINE = MergeTree()
ORDER BY id
PRIMARY KEY id;

-- Таблица Author
CREATE TABLE IF NOT EXISTS artworks.Author
(
    id UUID,
    name String,
    birthYear Nullable(Int32),
    deathYear Nullable(Int32),
    CONSTRAINT emptyCheck CHECK empty(name) = 0,
    CONSTRAINT birthDeathYear CHECK (isNull(birthYear) OR (isNull(deathYear) OR birthYear < deathYear) AND birthYear > 0)
)
ENGINE = MergeTree()
ORDER BY id
PRIMARY KEY id;

-- Таблица Collection
CREATE TABLE IF NOT EXISTS artworks.Collection
(
    id UUID,
    title String,
    CONSTRAINT titleEmpty CHECK empty(title) = 0
)
ENGINE = MergeTree()
ORDER BY id
PRIMARY KEY id;

-- Таблица Artworks
CREATE TABLE IF NOT EXISTS artworks.Artworks
(
    id UUID,
    title String,
    technic Nullable(String),
    material Nullable(String),
    size Nullable(String),
    creationYear Int32,
    authorID UUID,
    collectionID UUID,
    CONSTRAINT emptyCheck CHECK empty(title) = 0 AND creationYear > 0
)
ENGINE = MergeTree()
ORDER BY id
PRIMARY KEY id;

-- Таблица Events
CREATE TABLE IF NOT EXISTS artworks.Events
(
    id UUID,
    title String,
    dateBegin DateTime,
    dateEnd DateTime,
    canVisit Nullable(UInt8),
    adress Nullable(String),
    cntTickets Nullable(Int32),
    creatorID UUID,
    valid UInt8 DEFAULT 1,
    CONSTRAINT emptyCheck CHECK empty(title) = 0 AND empty(adress) = 0,
    CONSTRAINT dateBeginEndCheck CHECK dateBegin < dateEnd
)
ENGINE = MergeTree()
ORDER BY id
PRIMARY KEY id;

-- Таблица Artwork_event (многие-ко-многим)
CREATE TABLE IF NOT EXISTS artworks.Artwork_event
(
    artworkID UUID,
    eventID UUID,
    CONSTRAINT artworkID_notnull CHECK artworkID IS NOT NULL,
    CONSTRAINT eventID_notnull CHECK eventID IS NOT NULL
)
ENGINE = MergeTree()
ORDER BY (artworkID, eventID)
PRIMARY KEY (artworkID, eventID);

-- Таблица TicketPurchases
CREATE TABLE IF NOT EXISTS artworks.TicketPurchases
(
    id UUID,
    customerName String,
    customerEmail String,
    purchaseDate DateTime,
    eventID UUID,
    CONSTRAINT emptyCheck CHECK empty(customerName) = 0 AND empty(customerEmail) = 0
)
ENGINE = MergeTree()
ORDER BY id
PRIMARY KEY id;

-- Таблица tickets_user
CREATE TABLE IF NOT EXISTS artworks.tickets_user
(
    ticketID UUID,
    userID UUID,
    CONSTRAINT ticketID_notnull CHECK ticketID IS NOT NULL,
    CONSTRAINT userID_notnull CHECK userID IS NOT NULL
)
ENGINE = MergeTree()
ORDER BY (ticketID, userID)
PRIMARY KEY (ticketID, userID);