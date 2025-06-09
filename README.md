# weather-api
Software Engineering School 5.0 // Case Task

This project is a weather subscription service API. It allows users to subscribe to weather updates for a specific city at a given frequency (hourly or daily). Subscriptions require email confirmation. Users can also confirm their subscription and unsubscribe. The API also provides an endpoint to get the current weather for a city.

## Technologies Used
*   Go
*   PostgreSQL
*   Docker & Docker Compose
*   MailHog (for email testing)
*   golang-migrate (for database migrations)

## Project Structure
```
weather-api
├── .env                  # Environment variables
├── go.mod                # Go dependencies
├── go.sum
├── README.md
├── Dockerfile        # Dockerfile for building the application image
├── docker-compose.yml    # Docker Compose configuration
├── cmd/
│   └── weather-api/           # Application entry point
│       └── main.go
└── internal/             # Internal application logic
    ├── config/           # Application configuration
    ├── handler/          # HTTP handlers
    ├── model/            # Data models
    ├── repository/       # Database interaction logic
    ├── service/          # Business logic
    ├── mail/             # Email handling logic
    └── database/           # Database related files
        └── postgres/         # PostgreSQL utilities
            ├── migrations/   # Database migration files
            └── postgres.go   # PostgreSQL connection and migration logic

```

## Getting Started

### Prerequisites
*   Docker

### Launching the Application (with Docker)
1.  **Clone the repository:**
    ```bash
    git clone https://github.com/voloshmiak/weather-api.git
    cd weather-api
    ```

2.  **Environment Variables:**
    The application uses environment variables for configuration. The `docker-compose.yml` file sets these for the services. You will need to provide your own `WEATHER_API_KEY`.
    Key variables include:
    *   `DB_HOST`, `DB_PORT`, `DB_USER`, `DB_PASSWORD`, `DB_NAME` for database connection.
    *   `SMTP_HOST`, `SMTP_PORT` for the MailHog email server.
    *   `API_PORT` for the application server.
    *   `WEATHER_API_KEY` for accessing the external weather API (e.g., WeatherAPI.com). **You need to obtain an API key from a weather service provider and replace `YOUR_API_KEY` in the `docker-compose.yml` file or set it as an environment variable.**

3.  **Update `docker-compose.yml`:**
    Open the `docker-compose.yml` file and replace `YOUR_API_KEY` with your actual WeatherAPI key:
    ```yaml
    # ... other services
    services:
      app:
        # ... other app config
        environment:
          # ... other env vars
          - WEATHER_API_KEY=YOUR_ACTUAL_API_KEY_HERE # <-- Replace this
    # ...
    ```

4.  **Build and run the application using Docker Compose:**
    From the root directory of the project, run:
    ```bash
    docker-compose up --build
    ```

5.  **Accessing the Services:**
    *   **API:** `http://localhost:8080` (or the `API_PORT` you configured)
        *   `POST /subscribe` - with form data: `email`, `city`, `frequency`
        *   `GET /confirm/{token}`
        *   `GET /unsubscribe/{token}`
        *   `GET /weather?city={city_name}`
    *   **MailHog (Email UI):** `http://localhost:8025`
    *   **Database (PostgreSQL):** Accessible on `localhost:5432`

### Database Migrations
Database migrations are handled by `golang-migrate` and are applied automatically when the application starts. Migration files are located in the `/migrations` directory.

## Logic Overview

### Subscription Process (`POST /subscribe`)
1.  User submits `email`, `city`, and `frequency` (hourly/daily).
2.  The service checks if a subscription already exists for the given email.
    *   If an active, confirmed subscription exists, it returns an "Already Subscribed" error.
    *   If a subscription exists but is not confirmed, it updates the confirmation token and resends the confirmation email.
    *   If no subscription exists, it creates a new unconfirmed subscription record in the database with a unique confirmation token.
3.  A confirmation email containing a link with the unique token is sent to the user's email address (via MailHog).

### Confirmation Process (`GET /confirm/{token}`)
1.  User clicks the confirmation link in the email.
2.  The service validates the `token`.
3.  If the token is valid and found:
    *   The corresponding subscription is marked as `confirmed` in the database.
    *   A unique `unsubscribe_token` is generated and stored.
    *   The `unsubscribe_token` is returned to the user.
4.  If the token is invalid or not found, an appropriate error is returned.

### Unsubscribe Process (`GET /unsubscribe/{token}`)
1.  User uses their `unsubscribe_token` to request unsubscription.
2.  The service validates the `token`.
3.  If the token is valid and found, the corresponding subscription is deleted from the database.
4.  If the token is invalid or not found, an appropriate error is returned.

### Get Weather (`GET /weather?city={city_name}`)
1.  User requests weather for a specific `city`.
2.  The service fetches current weather data from an external weather API (WeatherAPI.com) using the configured `WEATHER_API_KEY`.
3.  The weather data (temperature, humidity, description) is returned as JSON.

### Stopping the Application
To stop the application and remove the containers, networks, and volumes created by `docker-compose up`:
```bash
docker-compose down
```