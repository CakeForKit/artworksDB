-- Active: 1744740356603@@127.0.0.1@5432@artworks


-- Вставляем администраторов
INSERT INTO Admins (id, username, login, hashedPassword, createdAt) VALUES
(gen_random_uuid(), 'admin1', 'admin1_login', '$2a$10$xJwL5v5Jz5U5Z5U5Z5U5Ze', NOW()),
(gen_random_uuid(), 'admin2', 'admin2_login', '$2a$10$xJwL5v5Jz5U5Z5U5Z5U5Ze', NOW());

-- Вставляем сотрудников (сначала получаем ID администраторов)
WITH admin_ids AS (
    SELECT id FROM Admins ORDER BY createdAt LIMIT 2
)
INSERT INTO Employee (id, username, login, hashedPassword, createdAt, adminID) VALUES
(gen_random_uuid(), 'employee1', 'emp1_login', '$2a$10$xJwL5v5Jz5U5Z5U5Z5U5Ze', NOW(), (SELECT id FROM admin_ids OFFSET 0 LIMIT 1)),
(gen_random_uuid(), 'employee2', 'emp2_login', '$2a$10$xJwL5v5Jz5U5Z5U5Z5U5Ze', NOW(), (SELECT id FROM admin_ids OFFSET 1 LIMIT 1));

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
INSERT INTO Artwork (id, title, technic, material, size, creationYear, authorID, collectionID) VALUES
(gen_random_uuid(), 'Mona Lisa', 'Oil painting', 'Poplar wood', '77 × 53 cm', 1503, (SELECT id FROM author_data), (SELECT id FROM collection_data)),
(gen_random_uuid(), 'Starry Night', 'Oil painting', 'Canvas', '73.7 × 92.1 cm', 1889, (SELECT id FROM van_gogh), (SELECT id FROM impression_collection));

-- Вставляем события (сначала получаем ID сотрудников)
WITH employee_ids AS (
    SELECT id FROM Employee ORDER BY createdAt LIMIT 2
)
INSERT INTO Events (id, title, dateBegin, dateEnd, access, adress, cntTickets, creatorID) VALUES
(gen_random_uuid(), 'Renaissance Exhibition', NOW() + INTERVAL '10 days', NOW() + INTERVAL '20 days', TRUE, '123 Art Gallery St, Museum District', 100, (SELECT id FROM employee_ids OFFSET 0 LIMIT 1)),
(gen_random_uuid(), 'Van Gogh Special', NOW() + INTERVAL '15 days', NOW() + INTERVAL '25 days', TRUE, '456 Modern Art Ave, Downtown', 150, (SELECT id FROM employee_ids OFFSET 1 LIMIT 1));

-- Связываем произведения искусства с событиями (сначала получаем ID произведений и событий)
WITH 
artwork_event_data AS (
    SELECT 
        a.id as artwork_id, 
        e.id as event_id
    FROM 
        Artwork a
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
INSERT INTO Artwork (id, title, technic, material, size, creationYear, authorID, collectionID) VALUES
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