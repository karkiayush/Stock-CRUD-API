# Stock-CRUD-API

Stock-CRUD-API is a simple RESTful API built with Golang, Gorilla Mux, and PostgreSQL to perform CRUD (Create, Read, Update, Delete) operations on stock data.

## Features
- Create a new stock entry
- Retrieve all stocks
- Retrieve a single stock by ID
- Update a stock entry
- Delete a stock entry

## Technologies Used
- Golang
- Gorilla Mux (for routing)
- PostgreSQL (for database)
- godotenv (for environment variable management)

## Installation

1. Clone the repository:
   ```sh
   git clone https://github.com/karkiayush/Stock-CRUD-API.git
   cd Stock-CRUD-API
   ```

2. Install dependencies:
   ```sh
   go mod tidy
   ```

3. Set up PostgreSQL database and update `.env` file with the correct database URL:
   ```sh
   POSTGRES_URL=your_postgres_connection_string
   ```

## Running the API

1. Start the server:
   ```sh
   go run main.go
   ```
2. The server runs on port `8080` by default.

## API Endpoints

| Method | Endpoint              | Description              |
|--------|----------------------|-------------------------|
| POST   | /api/newstock        | Create a new stock      |
| GET    | /api/stock           | Retrieve all stocks     |
| GET    | /api/stock/{id}      | Retrieve stock by ID    |
| PUT    | /api/stock/{id}      | Update stock by ID      |
| DELETE | /api/stock/{id}      | Delete stock by ID      |

## Database Schema
```sql
CREATE TABLE stocks (
    stockid SERIAL PRIMARY KEY,
    name TEXT NOT NULL,
    price NUMERIC(10,2) NOT NULL,
    company TEXT NOT NULL
);
```

## Project Structure
```
Stock-CRUD-API/
├── main.go          # Entry point
├── router/
│   ├── router.go    # Route definitions
├── middleware/
│   ├── handlers.go  # CRUD handlers
├── models/
│   ├── stock.go     # Stock model
├── .env             # Environment variables
├── go.mod           # Go modules
├── go.sum           # Dependency lock file
```

## Contributing
Contributions are welcome! Feel free to open issues or submit pull requests.

## License
This project is licensed under the MIT License.