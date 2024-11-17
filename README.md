
# Golang Crawler Project

This project is a web crawling system written in Go that collects real estate advertisement data from various sites and displays them to users based on different filters. Users can interact with the program through a Telegram bot to view search results.

---

## Features

- **Data Collection**: Crawls and collects real estate ads from multiple sources.
- **Advanced Filters**: Users can customize their searches with filters like price range, city, neighborhood, number of rooms, area, and building age.
- **Telegram Bot**: The application includes a Telegram bot that allows users to receive their search results via Telegram.
- **Rate Limiting Management**: The application uses Redis to manage the rate of requests based on different user roles.
- **Data Validation and Normalization**: Input data is validated and normalized to ensure consistency and prevent invalid entries in the database.

---

## Tools and Technologies Used

- **Golang** (v1.16+): Primary language of the project.
- **PostgreSQL**: Database used for storing advertisement and user data.
- **Redis**: Manages request rate limits and prevents excessive user requests.
- **GORM**: ORM used for database interactions with PostgreSQL.
- **go-redis**: Library used for connecting to Redis.
- **Telegram Bot API**: For Telegram integration to handle user commands.



## ERD (Entity Relationship Diagram)

The ERD below illustrates the overall structure of the database, including the entities `Users`, `Filters`, `Ads`, `WatchList`, and `Users_Ads`. Each entity plays a critical role in storing and managing project information.

---

![IMAGE 1403-08-25 02:20:44](https://github.com/user-attachments/assets/1a6bdfda-84a9-4e75-84bf-ce28b4456311)

---

## Bootstrapping (Setup Guide)

### 1. Clone the Repository

First, clone the project:

```bash
git clone https://github.com/your-repo/crawler-project.git
cd crawler-project
```

### 2. Database Setup

Create a PostgreSQL database and configure the connection details in the `database.go` file:

```plaintext
host=localhost user=postgres password=yourpassword dbname=CrawlerDb port=5432 sslmode=disable TimeZone=Asia/Tehran
```

### 3. Install Dependencies

Install all dependencies with the following command:

```bash
go mod tidy
```

### 4. Run the Project

To run the project, use the following command:

```bash
go run main.go
```

### 5. Run the Telegram Bot

- First, set the `BotToken` in the Telegram bot settings (`api-bot/`).
- Then, start the bot so users can connect through Telegram.

---

## Tests

### Running Tests

To run tests, use the following command to execute all tests in the project:

```bash
go test ./...
```

### Unit and Integration Tests

The project includes unit tests for models and functions like data normalization and validation, as well as integration tests for verifying Telegram bot functionality and request management.

---

## Folder Structure

- **`cmd/`**: Main executable commands for the project.
- **`database/`**: Database configurations and setup.
- **`models/`**: Data models that map to PostgreSQL tables.
- **`repository/`**: Data access layer for CRUD and database operations.
- **`rate_limiter_p/`**: Manages rate limits for users and roles.
- **`api-bot/`**: Contains the Telegram bot code for user interaction.
- **`middlewares/`**: Middleware for request handling and limitations.
- **`metric/`**: Manages and logs metrics.
- **`log/`**: Manages and stores logs.

---

## Contribution

To contribute to this project:

1. Create a new branch.
2. Make your changes.
3. Submit a pull request for review.

---

© 2024 Your Company. All rights reserved.
