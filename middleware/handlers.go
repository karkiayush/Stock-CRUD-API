package middleware

import (
	"database/sql"
	"encoding/json"
	"fmt"
	"github.com/gorilla/mux"
	"github.com/joho/godotenv"
	"github.com/karkiayush/Stock-CRUD-API/models"
	_ "github.com/lib/pq" // PostgresSQL driver
	"log"
	"net/http"
	"os"
	"strconv"
)

type Response struct {
	ID      int64  `json:"id"`
	Message string `json:"message"`
}

/**************Creating Database Connection****************/
func createConnection() *sql.DB {
	err := godotenv.Load(".env")
	if err != nil {
		log.Fatal("Error loading .env file")
	}
	db, err := sql.Open("postgres", os.Getenv("POSTGRES_URL"))
	if err != nil {
		panic(err)
	}
	err = db.Ping()
	if err != nil {
		panic(err)
	}
	fmt.Println("Successfully connected to postgres")
	return db
}

/*************CRUD methods****************/

func CreateStock(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json")

	var stock models.Stock
	err := json.NewDecoder(r.Body).Decode(&stock)
	if err != nil {
		log.Fatalf("Unable to decode request body: %v", err)
	}

	// we'll invoke the insertStock method that takes the stock struct & create record on database and returns the stock id
	insertId := insertStock(stock)
	res := Response{
		ID:      insertId,
		Message: "Stock created successfully",
	}
	// conversion of response for Response struct to json
	err = json.NewEncoder(w).Encode(res)
	if err != nil {
		log.Fatalf("Error during encoding of response while creating stock: %v", err)
	}
	return
}

func GetAllStocks(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json")

	stocks, err := getAllStock()
	if err != nil {
		log.Fatalf("Unable to get stocks: %v", err)
	}

	err = json.NewEncoder(w).Encode(stocks)
	if err != nil {
		log.Fatalf("Error during encoding of response while trying to getting all stocks: %v", err)
	}
	return
}

func GetStock(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json")

	params := mux.Vars(r)
	id, err := strconv.Atoi(params["id"])
	if err != nil {
		log.Fatalf("Unable to convert id to int: %v", err)
	}

	stock, err := getStock(int64(id))
	if err != nil {
		log.Fatalf("Unable to get stock: %v", err)
	}

	err = json.NewEncoder(w).Encode(stock)
	if err != nil {
		log.Fatalf("Error during encoding of response while trying to getting stock: %v", err)
	}
	return
}

func UpdateStock(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json")

	params := mux.Vars(r)
	id, err := strconv.Atoi(params["id"])
	if err != nil {
		log.Fatalf("Unable to convert id to int: %v", err)
	}

	var newStock models.Stock
	err = json.NewDecoder(r.Body).Decode(&newStock)
	if err != nil {
		log.Fatalf("Unable to decode request body: %v", err)
	}

	updatedRows, err := updateStockInDb(newStock, id)
	if err != nil {
		log.Fatalf("Unable to update stock in DB: %v", err)
	}

	msg := fmt.Sprintf("Stock updated successfully. Total rows affected: %v", updatedRows)
	res := Response{
		ID:      int64(id),
		Message: msg,
	}
	err = json.NewEncoder(w).Encode(res)
	if err != nil {
		log.Fatalf("Error during encoding of response while trying to getting updated stock: %v", err)
	}
	return
}

func DeleteStock(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json")
	params := mux.Vars(r)
	id, err := strconv.Atoi(params["id"])
	if err != nil {
		log.Fatalf("Unable to convert id to int: %v", err)
	}
	deletedRows, err := deleteStockDB(int64(id))
	if err != nil {
		log.Fatalf("Unable to delete stock in DB: %v", err)
	}
	msg := fmt.Sprintf("Stock deleted successfully. Total rows affected: %v", deletedRows)
	res := Response{
		ID:      int64(id),
		Message: msg,
	}
	err = json.NewEncoder(w).Encode(res)
	if err != nil {
		log.Fatalf("Error during encoding of response while trying to delete stock: %v", err)
	}
	return
}

/***********Helper methods*************/
func insertStock(stock models.Stock) int64 {
	db := createConnection()
	defer func(db *sql.DB) {
		err := db.Close()
		if err != nil {
			log.Fatalf("Error closing database connection: %v", err)
		}
	}(db)

	var id int64
	sqlStatement := `INSERT INTO stocks(name,price,company) VALUES($1,$2,$3) RETURNING stockid`
	err := db.QueryRow(sqlStatement, stock.Name, stock.Price, stock.Company).Scan(&id)
	if err != nil {
		log.Fatalf("Error inserting stock: %v", err)
	}

	fmt.Printf("Inserted a single record with id: %v\n", id)
	return id
}

func getStock(id int64) (models.Stock, error) {
	db := createConnection()
	defer func(db *sql.DB) {
		err := db.Close()
		if err != nil {
			log.Fatalf("Error closing database connection: %v", err)
		}
	}(db)

	var stock models.Stock
	sqlStatement := `SELECT * FROM stocks WHERE stockid = $1`
	err := db.QueryRow(sqlStatement, id).Scan(&stock.StockId, &stock.Name, &stock.Price, &stock.Company)
	if err != nil {
		log.Fatalf("Error getting stock: %v", err)
	}

	fmt.Printf("Got a single record with id: %v\n", stock.StockId)
	return stock, nil
}

func getAllStock() ([]models.Stock, error) {
	db := createConnection()
	defer func(db *sql.DB) {
		err := db.Close()
		if err != nil {
			log.Fatalf("Error closing database connection: %v", err)
		}
	}(db)

	var stocks []models.Stock
	sqlStatement := `SELECT * FROM stocks`
	rows, err := db.Query(sqlStatement)
	if err != nil {
		log.Fatalf("Unable to execute query: %v", err)
	}
	defer func(rows *sql.Rows) {
		err := rows.Close()
		if err != nil {
			log.Fatalf("Error closing rows: %v", err)
		}
	}(rows)

	// Iterating over the obtained rows
	for rows.Next() {
		var stock models.Stock
		err := rows.Scan(&stock.StockId, &stock.Name, &stock.Price, &stock.Company)
		if err != nil {
			log.Fatalf("Error scanning rows: %v", err)
		}
		stocks = append(stocks, stock)
	}
	return stocks, err
}

func updateStockInDb(stock models.Stock, id int) (int, error) {
	db := createConnection()
	defer func(db *sql.DB) {
		err := db.Close()
		if err != nil {
			log.Fatalf("Error closing database connection: %v", err)
		}
	}(db)

	sqlStatement := `update stocks set name=$1, price=$2, company=$3 where stockid=$4`
	res, err := db.Exec(sqlStatement, stock.Name, stock.Price, stock.Company, id)
	if err != nil {
		log.Fatalf("Error updating stock: %v", err)
		return 0, err
	}
	// Getting affected row count
	rowsAffected, err := res.RowsAffected()
	if err != nil {
		log.Fatalf("Error getting affected rows: %v", err)
		return 0, err
	}
	return int(rowsAffected), nil
}

func deleteStockDB(id int64) (int, error) {
	db := createConnection()
	defer func(db *sql.DB) {
		err := db.Close()
		if err != nil {
			log.Fatalf("Error closing database connection: %v", err)
		}
	}(db)

	sqlStatement := `DELETE FROM stocks WHERE stockid=$1`
	res, err := db.Exec(sqlStatement, id)
	if err != nil {
		log.Fatalf("Error deleting stock: %v", err)
		return 0, err
	}
	rowsAffected, err := res.RowsAffected()
	if err != nil {
		log.Fatalf("Error getting affected rows: %v", err)
	}
	return int(rowsAffected), nil
}
