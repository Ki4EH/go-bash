# Go-Bash

Go-Bash is a REST API built with Go that allows you to run bash scripts. It provides a simple and efficient way to manage and execute bash scripts through a RESTful interface.

Topics:
- [Features](#features)
- [Technologies Used](#technologies-used)
- [Getting Started](#getting-started)
  - [Prerequisites](#prerequisites)
  - [Installation](#installation)
- [Docker Build and Run Guide](#docker-build-and-run-guide)
- [API Endpoints](#api-endpoints)

## Features

- Get information about all commands.
- Insert a new command.
- Remove an existing command.
- Get information about specific commands by ID.

## Technologies Used

- Go: The backend is built entirely in Go.
- Gorilla Mux: A powerful HTTP router and URL matcher for building Go web servers.
- PostgreSQL: Used as the primary database for storing command information.
- Docker: Used to containerize the application.
- DATA-DOG/go-sqlmock: Used for testing the database layer.
- godotenv: Used for loading environment variables from a `.env` file.

## Getting Started

To get a local copy up and running, follow these simple steps.

### Prerequisites

- Go (latest version)
- PostgreSQL

### Installation

1. Clone the repo 
    ```sh
    git clone https://github.com/Ki4EH/go-bash.git
    ```
2. Install Go packages
   ```sh
   go mod download
   ```
3. Create a PostgreSQL database. With the following schema:
   ```sql
   CREATE TABLE IF NOT EXISTS "commands" (
    "id" SERIAL PRIMARY KEY,
    "name" TEXT NOT NULL,
    "script" TEXT NOT NULL,
    "description" TEXT NOT NULL,
    "output" TEXT NOT NULL
   );
   ```
   Or you can use init.sql file in the root directory to create the table.
   ```sh
    psql -U your_db_user -d your_db_name -a -f init.sql
    ```
4. Create a `.env` file in the root directory and add the following environment variables
   ```env
   POSTGRES_USER=your_db_user
   POSTGRES_PASSWORD=your_db_password
   POSTGRES_DB=your_db_name
   POSTGRES_HOST=your_db_host
   POSTGRES_PORT=your_db_port
   ```
5. Run the application
   ```sh
    go run main.go
    ```
6. The application should now be running on default `http://localhost:8080`
7. You can now test the API using Postman or any other API testing tool

## Docker Build and Run Guide

Before you start, make sure you have Docker and Docker Compose installed on your machine.

#### Steps

1. Navigate to the project directory where project cloned and `docker-compose.yml` file is located.
```bash
cd /path/to/your/project
```
2. Change in `.env` `POSTGRES_HOST=` to `POSTGRES_HOST=db` for saving data in the database container.

3. Change ports in `docker-compose.yml` on your own if you need to use another port for the application or database


4. Build the Docker image.

```bash
docker-compose build
```
5. After the build is complete, run the Docker container.

```bash 
docker-compose up
```


## API Endpoints

- `GET /info`: Get information about all commands.
- returns a JSON array of all commands. Example:
  ```json
  [
    {
      "id": 1,
      "name": "test"
    }
  ]
  ```
- `POST /new`: Insert a new command. The request body should be a JSON object representing the command.
- body example:
  ```json
  {
    "name": "test",
    "script": "echo 'Hello, World!'",
    "description": "test description"
  }
  ```
- `GET /remove?id={id}`: Remove an existing command by ID.
- `GET /info-by-id?id={id}`: Get information about a specific command by ID.
- returns a JSON object representing the command. Example:
  ```json
  {
    "id": 1,
    "name": "test",
    "script": "echo 'Hello, World!'",
    "description": "test description",
    "output": "Hello, World!\n"
  }
  ```
## 🚀 Roadmap
1. 🔒 Add authentication and authorization.
2. 📝 Add support for running commands with arguments.
3. 🌍 Add support for running commands with environment variables.
4. 🔄 Add support for running commands with input/output redirection.
5. ⏱️ Add support for running commands with a timeout.
6. 👤 Add support for running commands with a specific user.
7. 📂 Add support for running commands with a specific working directory.
8. 🛡️ Add protection from SQL Injection Attacks and other security vulnerabilities.
# Contact

#### Kirill Aksenov - [Poka132@yandex.ru](mailto:Poka132@yandex.ru) 
#### Telegram - [@ki4eh](https://t.me/ki4eh)