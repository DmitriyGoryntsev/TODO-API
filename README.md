# TODO API

A simple TODO API built with Go using the Gin framework, PostgreSQL for data storage, and Swagger for API documentation. This project demonstrates how to create a RESTful API for managing TODO items with basic CRUD operations.

## Features

- **Create, Read, Update, Delete (CRUD)** operations for TODO items.
- **PostgreSQL** as the database for persistent storage.
- **Swagger** integrated for API documentation and testing.
- **Gin** framework for fast and efficient HTTP routing.
- Simple and clean code structure for easy understanding and maintenance.

## Technologies Used

- Go (Golang)
- Gin web framework
- PostgreSQL
- Swagger for API documentation
- Docker (optional, for containerization)

## Getting Started

### Prerequisites

- Go 1.16 or later
- PostgreSQL installed or a PostgreSQL database service
- Docker (optional, for running the database in a container)

### Installation

1. Clone the repository:

   ```bash
   git clone https://github.com/DmitriyGoryntsev/TODO-API.git
   cd TODO-API```
2. Create a PostgreSQL database and user, or use an existing one. Update the database connection settings in the config.go file.
3. Install the required Go dependencies:

  go mod tidy

4. Run the application:

  go run main.go

5. Access the API documentation via Swagger at http://localhost:8080/swagger/index.html.


