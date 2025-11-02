-- Удаление тестовых данных в обратном порядке зависимостей

-- Удаляем связи многие-ко-многим
DELETE FROM tickets_user;
DELETE FROM Artwork_event;

-- Удаляем основные данные
DELETE FROM TicketPurchases;
DELETE FROM Events;
DELETE FROM Users;
DELETE FROM Employees;
DELETE FROM Admins;
DELETE FROM Artworks;
DELETE FROM Author;
DELETE FROM Collection;

-- Удаляем вспомогательные функции
DROP FUNCTION IF EXISTS random_author_name();
DROP FUNCTION IF EXISTS random_artwork_title();
DROP FUNCTION IF EXISTS random_technique();
DROP FUNCTION IF EXISTS random_material();
DROP FUNCTION IF EXISTS random_name();
DROP FUNCTION IF EXISTS random_login(VARCHAR);
DROP FUNCTION IF EXISTS random_email(VARCHAR);
DROP FUNCTION IF EXISTS random_password_hash();
DROP FUNCTION IF EXISTS random_artwork_id();
DROP FUNCTION IF EXISTS random_event_id();