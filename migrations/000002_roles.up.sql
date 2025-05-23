CREATE ROLE user_role;
GRANT SELECT
ON TABLE Author, Collection, Artworks, Events, Artwork_event 
TO user_role;
GRANT INSERT ON TABLE TicketPurchases TO user_role;

CREATE ROLE employee_role;
GRANT user_role TO employee_role;
GRANT SELECT, INSERT, UPDATE, DELETE 
ON TABLE Author, Collection, Artworks, Events, Artwork_event
TO employee_role;
GRANT SELECT ON TABLE TicketPurchases TO employee_role;

CREATE ROLE admin_role;
GRANT employee_role TO admin_role;
GRANT SELECT, INSERT, UPDATE ON TABLE Employees TO admin_role;
GRANT SELECT ON TABLE Admins TO admin_role;