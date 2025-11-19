package main

import (
	"database/sql"
	"fmt"
	"log"
	"time"

	"github.com/CakeForKit/artworksDB.git/internal/cnfg"
	_ "github.com/ClickHouse/clickhouse-go"
	_ "github.com/lib/pq"
)

type Config struct {
	PostgresDSN   string
	ClickHouseDSN string
}

func main() {
	// ----- Config ------
	pgCreds, err := cnfg.LoadPgCredentials("../../configs/", "db", "env")
	if err != nil {
		panic(fmt.Errorf("cannot load PgCredentials: %v", err))
	}
	clhCreds, err := cnfg.LoadClickHouseCredentials("./configs/", "clickhouse", "env")
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

	migrateTables(pgDB, chDB)
	log.Println("Data migration completed!")
}

func migrateTables(pgDB, chDB *sql.DB) {
	migrations := map[string]func(*sql.DB, *sql.DB) error{
		"Admins":          migrateAdmins,
		"Employees":       migrateEmployees,
		"Users":           migrateUsers,
		"Author":          migrateAuthor,
		"Collection":      migrateCollection,
		"Artworks":        migrateArtworks,
		"Events":          migrateEvents,
		"Artwork_event":   migrateArtworkEvent,
		"TicketPurchases": migrateTicketPurchases,
		"tickets_user":    migrateTicketsUser,
	}

	for table, migrateFn := range migrations {
		start := time.Now()
		log.Printf("Starting migration of table %s...", table)

		if err := migrateFn(pgDB, chDB); err != nil {
			log.Printf("Migration of table %s failed: %v", table, err)
		} else {
			log.Printf("Migration of table %s completed in %v", table, time.Since(start))
		}
	}
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

	count, err := processUserRows(rows, chDB)
	if err != nil {
		return err
	}

	log.Printf("Migrated %d Users records", count)
	return nil
}

func processUserRows(rows *sql.Rows, chDB *sql.DB) (int, error) {
	tx, err := chDB.Begin()
	if err != nil {
		return 0, fmt.Errorf("clickhouse transaction begin error: %v", err)
	}
	defer func() {
		if err != nil {
			_ = tx.Rollback()
		}
	}()

	stmt, err := prepareUserStatement(tx)
	if err != nil {
		return 0, err
	}
	defer stmt.Close()

	count, err := insertUserRows(rows, stmt)
	if err != nil {
		return 0, err
	}

	if err := tx.Commit(); err != nil {
		return 0, fmt.Errorf("clickhouse commit error: %v", err)
	}

	return count, nil
}

func prepareUserStatement(tx *sql.Tx) (*sql.Stmt, error) {
	stmt, err := tx.Prepare(`
        INSERT INTO Users (
            id, username, login, hashedPassword, createdAt, email, subscribeMail
        ) VALUES (?, ?, ?, ?, ?, ?, ?)
    `)
	if err != nil {
		return nil, fmt.Errorf("clickhouse prepare error: %v", err)
	}
	return stmt, nil
}

func insertUserRows(rows *sql.Rows, stmt *sql.Stmt) (int, error) {
	var count int

	for rows.Next() {
		userData, err := scanUserRow(rows)
		if err != nil {
			return 0, err
		}

		if err := executeUserInsert(stmt, userData); err != nil {
			return 0, err
		}

		count++
	}

	if err := rows.Err(); err != nil {
		return 0, fmt.Errorf("postgres rows error: %v", err)
	}

	return count, nil
}

func scanUserRow(rows *sql.Rows) (*userRowData, error) {
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
		return nil, fmt.Errorf("postgres row scan error: %v", err)
	}

	return &userRowData{
		id:             id,
		username:       username,
		login:          login,
		hashedPassword: hashedPassword,
		createdAt:      createdAt,
		email:          email,
		subscribeMail:  subscribeMail,
	}, nil
}

func executeUserInsert(stmt *sql.Stmt, data *userRowData) error {
	subscribeMailUint := uint8(0)
	if data.subscribeMail {
		subscribeMailUint = 1
	}

	var emailValue interface{} = nil
	if data.email.Valid {
		emailValue = data.email.String
	}

	if _, err := stmt.Exec(
		data.id,
		data.username,
		data.login,
		data.hashedPassword,
		data.createdAt,
		emailValue,
		subscribeMailUint,
	); err != nil {
		return fmt.Errorf("clickhouse exec error: %v", err)
	}

	return nil
}

type userRowData struct {
	id             string
	username       string
	login          string
	hashedPassword string
	createdAt      time.Time
	email          sql.NullString
	subscribeMail  bool
}

// Миграция таблицы Author
func migrateAuthor(pgDB, chDB *sql.DB) error {
	rows, err := pgDB.Query("SELECT id, name, birthYear, deathYear FROM Author")
	if err != nil {
		return fmt.Errorf("postgres query error: %v", err)
	}
	defer rows.Close()

	count, err := processAuthorRows(rows, chDB)
	if err != nil {
		return err
	}

	log.Printf("Migrated %d Author records", count)
	return nil
}

func processAuthorRows(rows *sql.Rows, chDB *sql.DB) (int, error) {
	tx, err := chDB.Begin()
	if err != nil {
		return 0, fmt.Errorf("clickhouse transaction begin error: %v", err)
	}
	defer func() {
		if err != nil {
			_ = tx.Rollback()
		}
	}()

	stmt, err := prepareAuthorStatement(tx)
	if err != nil {
		return 0, err
	}
	defer stmt.Close()

	count, err := insertAuthorRows(rows, stmt)
	if err != nil {
		return 0, err
	}

	if err := tx.Commit(); err != nil {
		return 0, fmt.Errorf("clickhouse commit error: %v", err)
	}

	return count, nil
}

func prepareAuthorStatement(tx *sql.Tx) (*sql.Stmt, error) {
	stmt, err := tx.Prepare(`
        INSERT INTO Author (
            id, name, birthYear, deathYear
        ) VALUES (?, ?, ?, ?)
    `)
	if err != nil {
		return nil, fmt.Errorf("clickhouse prepare error: %v", err)
	}
	return stmt, nil
}

func insertAuthorRows(rows *sql.Rows, stmt *sql.Stmt) (int, error) {
	var count int

	for rows.Next() {
		authorData, err := scanAuthorRow(rows)
		if err != nil {
			return 0, err
		}

		if err := executeAuthorInsert(stmt, authorData); err != nil {
			return 0, err
		}

		count++
	}

	if err := rows.Err(); err != nil {
		return 0, fmt.Errorf("postgres rows error: %v", err)
	}

	return count, nil
}

func scanAuthorRow(rows *sql.Rows) (*authorRowData, error) {
	var (
		id        string
		name      string
		birthYear sql.NullInt64
		deathYear sql.NullInt64
	)

	if err := rows.Scan(&id, &name, &birthYear, &deathYear); err != nil {
		return nil, fmt.Errorf("postgres row scan error: %v", err)
	}

	return &authorRowData{
		id:        id,
		name:      name,
		birthYear: birthYear,
		deathYear: deathYear,
	}, nil
}

func executeAuthorInsert(stmt *sql.Stmt, data *authorRowData) error {
	var birthYearValue, deathYearValue interface{} = nil, nil
	if data.birthYear.Valid {
		birthYearValue = int32(data.birthYear.Int64)
	}
	if data.deathYear.Valid {
		deathYearValue = int32(data.deathYear.Int64)
	}

	if _, err := stmt.Exec(
		data.id,
		data.name,
		birthYearValue,
		deathYearValue,
	); err != nil {
		return fmt.Errorf("clickhouse exec error: %v", err)
	}

	return nil
}

type authorRowData struct {
	id        string
	name      string
	birthYear sql.NullInt64
	deathYear sql.NullInt64
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

	count, err := processArtworkRows(rows, chDB)
	if err != nil {
		return err
	}

	log.Printf("Migrated %d Artworks records", count)
	return nil
}

func processArtworkRows(rows *sql.Rows, chDB *sql.DB) (int, error) {
	tx, err := chDB.Begin()
	if err != nil {
		return 0, fmt.Errorf("clickhouse transaction begin error: %v", err)
	}
	defer func() {
		if err != nil {
			_ = tx.Rollback()
		}
	}()

	stmt, err := prepareArtworkStatement(tx)
	if err != nil {
		return 0, err
	}
	defer stmt.Close()

	count, err := insertArtworkRows(rows, stmt)
	if err != nil {
		return 0, err
	}

	if err := tx.Commit(); err != nil {
		return 0, fmt.Errorf("clickhouse commit error: %v", err)
	}

	return count, nil
}

func prepareArtworkStatement(tx *sql.Tx) (*sql.Stmt, error) {
	stmt, err := tx.Prepare(`
        INSERT INTO Artworks (
            id, title, technic, material, size, creationYear, authorID, collectionID
        ) VALUES (?, ?, ?, ?, ?, ?, ?, ?)
    `)
	if err != nil {
		return nil, fmt.Errorf("clickhouse prepare error: %v", err)
	}
	return stmt, nil
}

func insertArtworkRows(rows *sql.Rows, stmt *sql.Stmt) (int, error) {
	var count int

	for rows.Next() {
		artworkData, err := scanArtworkRow(rows)
		if err != nil {
			return 0, err
		}

		if err := executeArtworkInsert(stmt, artworkData); err != nil {
			return 0, err
		}

		count++
	}

	if err := rows.Err(); err != nil {
		return 0, fmt.Errorf("postgres rows error: %v", err)
	}

	return count, nil
}

func scanArtworkRow(rows *sql.Rows) (*artworkRowData, error) {
	var data artworkRowData

	if err := rows.Scan(
		&data.id,
		&data.title,
		&data.technic,
		&data.material,
		&data.size,
		&data.creationYear,
		&data.authorID,
		&data.collectionID,
	); err != nil {
		return nil, fmt.Errorf("postgres row scan error: %v", err)
	}

	return &data, nil
}

func executeArtworkInsert(stmt *sql.Stmt, data *artworkRowData) error {
	technicValue := getNullableString(data.technic)
	materialValue := getNullableString(data.material)
	sizeValue := getNullableString(data.size)

	if _, err := stmt.Exec(
		data.id,
		data.title,
		technicValue,
		materialValue,
		sizeValue,
		data.creationYear,
		data.authorID,
		data.collectionID,
	); err != nil {
		return fmt.Errorf("clickhouse exec error: %v", err)
	}

	return nil
}

func getNullableString(nullString sql.NullString) interface{} {
	if nullString.Valid {
		return nullString.String
	}
	return nil
}

type artworkRowData struct {
	id           string
	title        string
	technic      sql.NullString
	material     sql.NullString
	size         sql.NullString
	creationYear int32
	authorID     string
	collectionID string
}

// Миграция таблицы Events
func migrateEvents(pgDB, chDB *sql.DB) error {
	rows, err := pgDB.Query(`
        SELECT id, title, dateBegin, dateEnd, canVisit, address, cntTickets, creatorID, valid 
        FROM Events
    `)
	if err != nil {
		return fmt.Errorf("postgres query error: %v", err)
	}
	defer rows.Close()

	count, err := processEventRows(rows, chDB)
	if err != nil {
		return err
	}

	log.Printf("Migrated %d Events records", count)
	return nil
}

func processEventRows(rows *sql.Rows, chDB *sql.DB) (int, error) {
	tx, err := chDB.Begin()
	if err != nil {
		return 0, fmt.Errorf("clickhouse transaction begin error: %v", err)
	}
	defer func() {
		if err != nil {
			_ = tx.Rollback()
		}
	}()

	stmt, err := prepareEventStatement(tx)
	if err != nil {
		return 0, err
	}
	defer stmt.Close()

	count, err := insertEventRows(rows, stmt)
	if err != nil {
		return 0, err
	}

	if err := tx.Commit(); err != nil {
		return 0, fmt.Errorf("clickhouse commit error: %v", err)
	}

	return count, nil
}

func prepareEventStatement(tx *sql.Tx) (*sql.Stmt, error) {
	stmt, err := tx.Prepare(`
        INSERT INTO Events (
            id, title, dateBegin, dateEnd, canVisit, address, cntTickets, creatorID, valid
        ) VALUES (?, ?, ?, ?, ?, ?, ?, ?, ?)
    `)
	if err != nil {
		return nil, fmt.Errorf("clickhouse prepare error: %v", err)
	}
	return stmt, nil
}

func insertEventRows(rows *sql.Rows, stmt *sql.Stmt) (int, error) {
	var count int

	for rows.Next() {
		eventData, err := scanEventRow(rows)
		if err != nil {
			return 0, err
		}

		if err := executeEventInsert(stmt, eventData); err != nil {
			return 0, err
		}

		count++
	}

	if err := rows.Err(); err != nil {
		return 0, fmt.Errorf("postgres rows error: %v", err)
	}

	return count, nil
}

func scanEventRow(rows *sql.Rows) (*eventRowData, error) {
	var data eventRowData

	if err := rows.Scan(
		&data.id,
		&data.title,
		&data.dateBegin,
		&data.dateEnd,
		&data.canVisit,
		&data.address,
		&data.cntTickets,
		&data.creatorID,
		&data.valid,
	); err != nil {
		return nil, fmt.Errorf("postgres row scan error: %v", err)
	}

	return &data, nil
}

func executeEventInsert(stmt *sql.Stmt, data *eventRowData) error {
	canVisitValue := convertNullableBool(data.canVisit)
	addressValue := getNullableString(data.address)
	cntTicketsValue := convertNullableInt32(data.cntTickets)
	validUint := boolToUint8(data.valid)

	if _, err := stmt.Exec(
		data.id,
		data.title,
		data.dateBegin,
		data.dateEnd,
		canVisitValue,
		addressValue,
		cntTicketsValue,
		data.creatorID,
		validUint,
	); err != nil {
		return fmt.Errorf("clickhouse exec error: %v", err)
	}

	return nil
}

// Вспомогательные функции (можно вынести в общие утилиты)
func convertNullableBool(nullBool sql.NullBool) interface{} {
	if nullBool.Valid {
		return boolToUint8(nullBool.Bool)
	}
	return nil
}

func convertNullableInt32(nullInt sql.NullInt64) interface{} {
	if nullInt.Valid {
		return int32(nullInt.Int64)
	}
	return nil
}

func boolToUint8(b bool) uint8 {
	if b {
		return 1
	}
	return 0
}

type eventRowData struct {
	id         string
	title      string
	dateBegin  time.Time
	dateEnd    time.Time
	canVisit   sql.NullBool
	address    sql.NullString
	cntTickets sql.NullInt64
	creatorID  string
	valid      bool
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
