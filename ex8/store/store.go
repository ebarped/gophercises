package store

import (
	"database/sql"
	"fmt"
	"log"

	"phone/phone"
)

var db *sql.DB

func GetConnection() *sql.DB {
	// Para evitar realizar una nueva conexión en cada llamada a
	// la función GetConnection.
	if db != nil {
		return db
	}

	db, err := sql.Open("sqlite3", "phones.db")
	if err != nil {
		panic(err)
	}
	return db
}

func InitDB(numbers []string) error {
	if isDBInitialized() {
		log.Println("DB already initialized!")
		return nil
	}

	db := GetConnection()

	createQuery := `CREATE TABLE IF NOT EXISTS phone_numbers (
					id INTEGER PRIMARY KEY AUTOINCREMENT,
					number VARCHAR(20)
	 );`

	_, err := db.Exec(createQuery)
	if err != nil {
		return err
	}

	insertQuery := `INSERT INTO phone_numbers (number) VALUES(?)`

	insertStmt, err := db.Prepare(insertQuery)
	if err != nil {
		return err
	}
	defer insertStmt.Close()

	for _, number := range numbers {
		fmt.Println("inserting phone number:", number)
		r, err := insertStmt.Exec(number)
		if err != nil {
			return err
		}

		rowsAffected, err := r.RowsAffected()
		if err != nil || rowsAffected != 1 {
			return fmt.Errorf("error: failed to insert phone number %s", number)
		}
	}
	log.Println("DB successfully initialized!")
	return nil
}

// isDBInitialized returns true if the DB is already initialized
// (this is a workaround, better methods must exist)
func isDBInitialized() bool {
	db := GetConnection()
	row := db.QueryRow("SELECT COUNT(*) FROM 'phone_numbers'")

	var temp string
	err := row.Scan(&temp)
	if err != nil {
		return false
	}

	if temp != "" {
		return true
	}
	return false
}

// GetPhones returns all the phones stored in the datastore
func GetPhones() []phone.Phone {
	var result []phone.Phone

	db := GetConnection()
	rows, err := db.Query("SELECT id,number FROM 'phone_numbers'")
	if err != nil {
		log.Fatalf("error getting phone numbers: %s\n", err)
	}
	defer rows.Close()

	for rows.Next() {
		var id int
		var num string
		err := rows.Scan(&id, &num)
		if err != nil {
			log.Fatalf("error getting value of row: %s\n", err)
		}
		p := phone.New(id, num)
		result = append(result, p)
	}
	return result
}

// UpdatePhoneNumber updates the phone number of the phone identified by "id", setting its new value as "num"
func UpdatePhoneNumber(id int, num string) error {
	db := GetConnection()
	updateQuery := `UPDATE phone_numbers 
					SET number = ?
					WHERE id = ?`

	updateStmt, err := db.Prepare(updateQuery)
	if err != nil {
		return err
	}
	defer updateStmt.Close()

	r, err := updateStmt.Exec(num, id)
	if err != nil {
		return err
	}

	rowsAffected, err := r.RowsAffected()
	if err != nil || rowsAffected != 1 {
		return fmt.Errorf("failed to update phone number id=%d to %s: %v", id, num, err)
	}

	return nil
}

// RemovePhone removes the phone identified by "id"
func RemovePhone(id int) error {
	db := GetConnection()
	removeQuery := `DELETE FROM phone_numbers 
					WHERE id = ?`

	removeStmt, err := db.Prepare(removeQuery)
	if err != nil {
		return err
	}
	defer removeStmt.Close()

	r, err := removeStmt.Exec(id)
	if err != nil {
		return err
	}

	rowsAffected, err := r.RowsAffected()
	if err != nil || rowsAffected != 1 {
		return fmt.Errorf("failed to remove phone number id=%d: %v", id, err)
	}

	return nil
}
