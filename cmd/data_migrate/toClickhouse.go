package main

import (
	"database/sql"
	"fmt"
	"log"
	"time"

	"git.iu7.bmstu.ru/ped22u691/PPO.git/internal/cnfg"
	_ "github.com/ClickHouse/clickhouse-go"
	_ "github.com/lib/pq"
)

// Config содержит настройки подключения к базам данных
type Config struct {
	PostgresDSN   string
	ClickHouseDSN string
}

func main() {
	// ----- Config ------
	pgCreds, err := cnfg.LoadPgCredentials("../../configs/")
	if err != nil {
		panic(fmt.Errorf("cannot load PgCredentials: %v", err))
	}
	clhCreds, err := cnfg.LoadClickHouseCredentials()
	if err != nil {
		panic(fmt.Errorf("cannot load ClickHouseCredentials: %v", err))
	}
	fmt.Printf("pgCreds: %+v\n", pgCreds)
	fmt.Printf("clhCreds: %+v\n", clhCreds)
	// ------------------

	// Конфигурация подключений
	PostgresConnStr := fmt.Sprintf("postgres://%s:%s@%s:%d/%s?sslmode=disable",
		pgCreds.Username, pgCreds.Password, "localhost", pgCreds.Port, pgCreds.DbName)
	ClickHouseConnStr := fmt.Sprintf(
		"tcp://%s:%d?database=%s&username=%s&password=%s&debug=true",
		"localhost", clhCreds.Port, clhCreds.DbName, clhCreds.Username, clhCreds.Password)

	// Подключение к PostgreSQL
	pgDB, err := sql.Open("postgres", PostgresConnStr)
	if err != nil {
		log.Fatalf("Failed to connect to PostgreSQL: %v", err)
	}
	defer pgDB.Close()

	// Подключение к ClickHouse
	chDB, err := sql.Open("clickhouse", ClickHouseConnStr)
	if err != nil {
		log.Fatalf("Failed to connect to ClickHouse: %v", err)
	}
	defer chDB.Close()

	// Проверка соединений
	if err := pgDB.Ping(); err != nil {
		log.Fatalf("PostgreSQL ping failed: %v", err)
	}
	if err := chDB.Ping(); err != nil {
		log.Fatalf("ClickHouse ping failed: %v", err)
	}

	// Миграция таблиц
	tables := []string{
		"Admins",
		"Employees",
		"Users",
		"Author",
		"Collection",
		"Artworks",
		"Events",
		"Artwork_event",
		"TicketPurchases",
		"tickets_user",
	}

	for _, table := range tables {
		start := time.Now()
		log.Printf("Starting migration of table %s...", table)

		var err error
		switch table {
		case "Admins":
			err = migrateAdmins(pgDB, chDB)
		case "Employees":
			err = migrateEmployees(pgDB, chDB)
		case "Users":
			err = migrateUsers(pgDB, chDB)
		case "Author":
			err = migrateAuthor(pgDB, chDB)
		case "Collection":
			err = migrateCollection(pgDB, chDB)
		case "Artworks":
			err = migrateArtworks(pgDB, chDB)
		case "Events":
			err = migrateEvents(pgDB, chDB)
		case "Artwork_event":
			err = migrateArtworkEvent(pgDB, chDB)
		case "TicketPurchases":
			err = migrateTicketPurchases(pgDB, chDB)
		case "tickets_user":
			err = migrateTicketsUser(pgDB, chDB)
		}

		if err != nil {
			log.Printf("Migration of table %s failed: %v", table, err)
		} else {
			log.Printf("Migration of table %s completed in %v", table, time.Since(start))
		}
	}

	log.Println("Data migration completed!")
}

// Миграция таблицы Admins
func migrateAdmins(pgDB, chDB *sql.DB) error {
	// Получение данных из PostgreSQL
	rows, err := pgDB.Query("SELECT id, username, login, hashedPassword, createdAt, valid FROM Admins")
	if err != nil {
		return fmt.Errorf("postgres query error: %v", err)
	}
	defer rows.Close()

	// Подготовка запроса для ClickHouse
	tx, err := chDB.Begin()
	if err != nil {
		return fmt.Errorf("clickhouse transaction begin error: %v", err)
	}

	stmt, err := tx.Prepare(`
		INSERT INTO Admins (
			id, username, login, hashedPassword, createdAt, valid
		) VALUES (?, ?, ?, ?, ?, ?)
	`)
	if err != nil {
		return fmt.Errorf("clickhouse prepare error: %v", err)
	}
	defer stmt.Close()

	// Чтение и вставка данных
	var count int
	for rows.Next() {
		var (
			id             string
			username       string
			login          string
			hashedPassword string
			createdAt      time.Time
			valid          bool
		)

		if err := rows.Scan(&id, &username, &login, &hashedPassword, &createdAt, &valid); err != nil {
			return fmt.Errorf("postgres row scan error: %v", err)
		}

		validUint := uint8(0)
		if valid {
			validUint = 1
		}

		if _, err := stmt.Exec(
			id,
			username,
			login,
			hashedPassword,
			createdAt,
			validUint,
		); err != nil {
			return fmt.Errorf("clickhouse exec error: %v", err)
		}

		count++
	}

	if err := rows.Err(); err != nil {
		return fmt.Errorf("postgres rows error: %v", err)
	}

	if err := tx.Commit(); err != nil {
		return fmt.Errorf("clickhouse commit error: %v", err)
	}

	log.Printf("Migrated %d Admins records", count)
	return nil
}

// Миграция таблицы Employees
func migrateEmployees(pgDB, chDB *sql.DB) error {
	rows, err := pgDB.Query("SELECT id, username, login, hashedPassword, createdAt, valid, adminID FROM Employees")
	if err != nil {
		return fmt.Errorf("postgres query error: %v", err)
	}
	defer rows.Close()

	tx, err := chDB.Begin()
	if err != nil {
		return fmt.Errorf("clickhouse transaction begin error: %v", err)
	}

	stmt, err := tx.Prepare(`
		INSERT INTO Employees (
			id, username, login, hashedPassword, createdAt, valid, adminID
		) VALUES (?, ?, ?, ?, ?, ?, ?)
	`)
	if err != nil {
		return fmt.Errorf("clickhouse prepare error: %v", err)
	}
	defer stmt.Close()

	var count int
	for rows.Next() {
		var (
			id             string
			username       string
			login          string
			hashedPassword string
			createdAt      time.Time
			valid          bool
			adminID        string
		)

		if err := rows.Scan(&id, &username, &login, &hashedPassword, &createdAt, &valid, &adminID); err != nil {
			return fmt.Errorf("postgres row scan error: %v", err)
		}

		validUint := uint8(0)
		if valid {
			validUint = 1
		}

		if _, err := stmt.Exec(
			id,
			username,
			login,
			hashedPassword,
			createdAt,
			validUint,
			adminID,
		); err != nil {
			return fmt.Errorf("clickhouse exec error: %v", err)
		}

		count++
	}

	if err := rows.Err(); err != nil {
		return fmt.Errorf("postgres rows error: %v", err)
	}

	if err := tx.Commit(); err != nil {
		return fmt.Errorf("clickhouse commit error: %v", err)
	}

	log.Printf("Migrated %d Employees records", count)
	return nil
}

// Миграция таблицы Users
func migrateUsers(pgDB, chDB *sql.DB) error {
	rows, err := pgDB.Query("SELECT id, username, login, hashedPassword, createdAt, email, subscribeMail FROM Users")
	if err != nil {
		return fmt.Errorf("postgres query error: %v", err)
	}
	defer rows.Close()

	tx, err := chDB.Begin()
	if err != nil {
		return fmt.Errorf("clickhouse transaction begin error: %v", err)
	}

	stmt, err := tx.Prepare(`
		INSERT INTO Users (
			id, username, login, hashedPassword, createdAt, email, subscribeMail
		) VALUES (?, ?, ?, ?, ?, ?, ?)
	`)
	if err != nil {
		return fmt.Errorf("clickhouse prepare error: %v", err)
	}
	defer stmt.Close()

	var count int
	for rows.Next() {
		var (
			id             string
			username       string
			login          string
			hashedPassword string
			createdAt      time.Time
			email          sql.NullString
			subscribeMail  bool
		)

		if err := rows.Scan(&id, &username, &login, &hashedPassword, &createdAt, &email, &subscribeMail); err != nil {
			return fmt.Errorf("postgres row scan error: %v", err)
		}

		subscribeMailUint := uint8(0)
		if subscribeMail {
			subscribeMailUint = 1
		}

		var emailValue interface{} = nil
		if email.Valid {
			emailValue = email.String
		}

		if _, err := stmt.Exec(
			id,
			username,
			login,
			hashedPassword,
			createdAt,
			emailValue,
			subscribeMailUint,
		); err != nil {
			return fmt.Errorf("clickhouse exec error: %v", err)
		}

		count++
	}

	if err := rows.Err(); err != nil {
		return fmt.Errorf("postgres rows error: %v", err)
	}

	if err := tx.Commit(); err != nil {
		return fmt.Errorf("clickhouse commit error: %v", err)
	}

	log.Printf("Migrated %d Users records", count)
	return nil
}

// Миграция таблицы Author
func migrateAuthor(pgDB, chDB *sql.DB) error {
	rows, err := pgDB.Query("SELECT id, name, birthYear, deathYear FROM Author")
	if err != nil {
		return fmt.Errorf("postgres query error: %v", err)
	}
	defer rows.Close()

	tx, err := chDB.Begin()
	if err != nil {
		return fmt.Errorf("clickhouse transaction begin error: %v", err)
	}

	stmt, err := tx.Prepare(`
		INSERT INTO Author (
			id, name, birthYear, deathYear
		) VALUES (?, ?, ?, ?)
	`)
	if err != nil {
		return fmt.Errorf("clickhouse prepare error: %v", err)
	}
	defer stmt.Close()

	var count int
	for rows.Next() {
		var (
			id        string
			name      string
			birthYear sql.NullInt64
			deathYear sql.NullInt64
		)

		if err := rows.Scan(&id, &name, &birthYear, &deathYear); err != nil {
			return fmt.Errorf("postgres row scan error: %v", err)
		}

		var birthYearValue, deathYearValue interface{} = nil, nil
		if birthYear.Valid {
			birthYearValue = int32(birthYear.Int64)
		}
		if deathYear.Valid {
			deathYearValue = int32(deathYear.Int64)
		}

		if _, err := stmt.Exec(
			id,
			name,
			birthYearValue,
			deathYearValue,
		); err != nil {
			return fmt.Errorf("clickhouse exec error: %v", err)
		}

		count++
	}

	if err := rows.Err(); err != nil {
		return fmt.Errorf("postgres rows error: %v", err)
	}

	if err := tx.Commit(); err != nil {
		return fmt.Errorf("clickhouse commit error: %v", err)
	}

	log.Printf("Migrated %d Author records", count)
	return nil
}

// Миграция таблицы Collection
func migrateCollection(pgDB, chDB *sql.DB) error {
	rows, err := pgDB.Query("SELECT id, title FROM Collection")
	if err != nil {
		return fmt.Errorf("postgres query error: %v", err)
	}
	defer rows.Close()

	tx, err := chDB.Begin()
	if err != nil {
		return fmt.Errorf("clickhouse transaction begin error: %v", err)
	}

	stmt, err := tx.Prepare(`
		INSERT INTO Collection (
			id, title
		) VALUES (?, ?)
	`)
	if err != nil {
		return fmt.Errorf("clickhouse prepare error: %v", err)
	}
	defer stmt.Close()

	var count int
	for rows.Next() {
		var (
			id    string
			title string
		)

		if err := rows.Scan(&id, &title); err != nil {
			return fmt.Errorf("postgres row scan error: %v", err)
		}

		if _, err := stmt.Exec(
			id,
			title,
		); err != nil {
			return fmt.Errorf("clickhouse exec error: %v", err)
		}

		count++
	}

	if err := rows.Err(); err != nil {
		return fmt.Errorf("postgres rows error: %v", err)
	}

	if err := tx.Commit(); err != nil {
		return fmt.Errorf("clickhouse commit error: %v", err)
	}

	log.Printf("Migrated %d Collection records", count)
	return nil
}

// Миграция таблицы Artworks
func migrateArtworks(pgDB, chDB *sql.DB) error {
	rows, err := pgDB.Query(`
		SELECT id, title, technic, material, size, creationYear, authorID, collectionID 
		FROM Artworks
	`)
	if err != nil {
		return fmt.Errorf("postgres query error: %v", err)
	}
	defer rows.Close()

	tx, err := chDB.Begin()
	if err != nil {
		return fmt.Errorf("clickhouse transaction begin error: %v", err)
	}

	stmt, err := tx.Prepare(`
		INSERT INTO Artworks (
			id, title, technic, material, size, creationYear, authorID, collectionID
		) VALUES (?, ?, ?, ?, ?, ?, ?, ?)
	`)
	if err != nil {
		return fmt.Errorf("clickhouse prepare error: %v", err)
	}
	defer stmt.Close()

	var count int
	for rows.Next() {
		var (
			id           string
			title        string
			technic      sql.NullString
			material     sql.NullString
			size         sql.NullString
			creationYear int32
			authorID     string
			collectionID string
		)

		if err := rows.Scan(&id, &title, &technic, &material, &size, &creationYear, &authorID, &collectionID); err != nil {
			return fmt.Errorf("postgres row scan error: %v", err)
		}

		var technicValue, materialValue, sizeValue interface{} = nil, nil, nil
		if technic.Valid {
			technicValue = technic.String
		}
		if material.Valid {
			materialValue = material.String
		}
		if size.Valid {
			sizeValue = size.String
		}

		if _, err := stmt.Exec(
			id,
			title,
			technicValue,
			materialValue,
			sizeValue,
			creationYear,
			authorID,
			collectionID,
		); err != nil {
			return fmt.Errorf("clickhouse exec error: %v", err)
		}

		count++
	}

	if err := rows.Err(); err != nil {
		return fmt.Errorf("postgres rows error: %v", err)
	}

	if err := tx.Commit(); err != nil {
		return fmt.Errorf("clickhouse commit error: %v", err)
	}

	log.Printf("Migrated %d Artworks records", count)
	return nil
}

// Миграция таблицы Events
func migrateEvents(pgDB, chDB *sql.DB) error {
	rows, err := pgDB.Query(`
		SELECT id, title, dateBegin, dateEnd, canVisit, adress, cntTickets, creatorID, valid 
		FROM Events
	`)
	if err != nil {
		return fmt.Errorf("postgres query error: %v", err)
	}
	defer rows.Close()

	tx, err := chDB.Begin()
	if err != nil {
		return fmt.Errorf("clickhouse transaction begin error: %v", err)
	}

	stmt, err := tx.Prepare(`
		INSERT INTO Events (
			id, title, dateBegin, dateEnd, canVisit, adress, cntTickets, creatorID, valid
		) VALUES (?, ?, ?, ?, ?, ?, ?, ?, ?)
	`)
	if err != nil {
		return fmt.Errorf("clickhouse prepare error: %v", err)
	}
	defer stmt.Close()

	var count int
	for rows.Next() {
		var (
			id         string
			title      string
			dateBegin  time.Time
			dateEnd    time.Time
			canVisit   sql.NullBool
			adress     sql.NullString
			cntTickets sql.NullInt64
			creatorID  string
			valid      bool
		)

		if err := rows.Scan(&id, &title, &dateBegin, &dateEnd, &canVisit, &adress, &cntTickets, &creatorID, &valid); err != nil {
			return fmt.Errorf("postgres row scan error: %v", err)
		}

		var canVisitValue, adressValue, cntTicketsValue interface{} = nil, nil, nil
		if canVisit.Valid {
			canVisitUint := uint8(0)
			if canVisit.Bool {
				canVisitUint = 1
			}
			canVisitValue = canVisitUint
		}
		if adress.Valid {
			adressValue = adress.String
		}
		if cntTickets.Valid {
			cntTicketsValue = int32(cntTickets.Int64)
		}

		validUint := uint8(0)
		if valid {
			validUint = 1
		}

		if _, err := stmt.Exec(
			id,
			title,
			dateBegin,
			dateEnd,
			canVisitValue,
			adressValue,
			cntTicketsValue,
			creatorID,
			validUint,
		); err != nil {
			return fmt.Errorf("clickhouse exec error: %v", err)
		}

		count++
	}

	if err := rows.Err(); err != nil {
		return fmt.Errorf("postgres rows error: %v", err)
	}

	if err := tx.Commit(); err != nil {
		return fmt.Errorf("clickhouse commit error: %v", err)
	}

	log.Printf("Migrated %d Events records", count)
	return nil
}

// Миграция таблицы Artwork_event
func migrateArtworkEvent(pgDB, chDB *sql.DB) error {
	rows, err := pgDB.Query("SELECT artworkID, eventID FROM Artwork_event")
	if err != nil {
		return fmt.Errorf("postgres query error: %v", err)
	}
	defer rows.Close()

	tx, err := chDB.Begin()
	if err != nil {
		return fmt.Errorf("clickhouse transaction begin error: %v", err)
	}

	stmt, err := tx.Prepare(`
		INSERT INTO Artwork_event (
			artworkID, eventID
		) VALUES (?, ?)
	`)
	if err != nil {
		return fmt.Errorf("clickhouse prepare error: %v", err)
	}
	defer stmt.Close()

	var count int
	for rows.Next() {
		var (
			artworkID string
			eventID   string
		)

		if err := rows.Scan(&artworkID, &eventID); err != nil {
			return fmt.Errorf("postgres row scan error: %v", err)
		}

		if _, err := stmt.Exec(
			artworkID,
			eventID,
		); err != nil {
			return fmt.Errorf("clickhouse exec error: %v", err)
		}

		count++
	}

	if err := rows.Err(); err != nil {
		return fmt.Errorf("postgres rows error: %v", err)
	}

	if err := tx.Commit(); err != nil {
		return fmt.Errorf("clickhouse commit error: %v", err)
	}

	log.Printf("Migrated %d Artwork_event records", count)
	return nil
}

// Миграция таблицы TicketPurchases
func migrateTicketPurchases(pgDB, chDB *sql.DB) error {
	rows, err := pgDB.Query(`
		SELECT id, customerName, customerEmail, purchaseDate, eventID 
		FROM TicketPurchases
	`)
	if err != nil {
		return fmt.Errorf("postgres query error: %v", err)
	}
	defer rows.Close()

	tx, err := chDB.Begin()
	if err != nil {
		return fmt.Errorf("clickhouse transaction begin error: %v", err)
	}

	stmt, err := tx.Prepare(`
		INSERT INTO TicketPurchases (
			id, customerName, customerEmail, purchaseDate, eventID
		) VALUES (?, ?, ?, ?, ?)
	`)
	if err != nil {
		return fmt.Errorf("clickhouse prepare error: %v", err)
	}
	defer stmt.Close()

	var count int
	for rows.Next() {
		var (
			id            string
			customerName  string
			customerEmail string
			purchaseDate  time.Time
			eventID       string
		)

		if err := rows.Scan(&id, &customerName, &customerEmail, &purchaseDate, &eventID); err != nil {
			return fmt.Errorf("postgres row scan error: %v", err)
		}

		if _, err := stmt.Exec(
			id,
			customerName,
			customerEmail,
			purchaseDate,
			eventID,
		); err != nil {
			return fmt.Errorf("clickhouse exec error: %v", err)
		}

		count++
	}

	if err := rows.Err(); err != nil {
		return fmt.Errorf("postgres rows error: %v", err)
	}

	if err := tx.Commit(); err != nil {
		return fmt.Errorf("clickhouse commit error: %v", err)
	}

	log.Printf("Migrated %d TicketPurchases records", count)
	return nil
}

// Миграция таблицы tickets_user
func migrateTicketsUser(pgDB, chDB *sql.DB) error {
	rows, err := pgDB.Query("SELECT ticketID, userID FROM tickets_user")
	if err != nil {
		return fmt.Errorf("postgres query error: %v", err)
	}
	defer rows.Close()

	tx, err := chDB.Begin()
	if err != nil {
		return fmt.Errorf("clickhouse transaction begin error: %v", err)
	}

	stmt, err := tx.Prepare(`
		INSERT INTO tickets_user (
			ticketID, userID
		) VALUES (?, ?)
	`)
	if err != nil {
		return fmt.Errorf("clickhouse prepare error: %v", err)
	}
	defer stmt.Close()

	var count int
	for rows.Next() {
		var (
			ticketID string
			userID   string
		)

		if err := rows.Scan(&ticketID, &userID); err != nil {
			return fmt.Errorf("postgres row scan error: %v", err)
		}

		if _, err := stmt.Exec(
			ticketID,
			userID,
		); err != nil {
			return fmt.Errorf("clickhouse exec error: %v", err)
		}

		count++
	}

	if err := rows.Err(); err != nil {
		return fmt.Errorf("postgres rows error: %v", err)
	}

	if err := tx.Commit(); err != nil {
		return fmt.Errorf("clickhouse commit error: %v", err)
	}

	log.Printf("Migrated %d tickets_user records", count)
	return nil
}
