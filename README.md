# weather-api
Software Engineering School 5.0 // Case Task

This project is a weather subscription service API. It allows users to subscribe to weather updates for a specific city at a given frequency (hourly or daily). Subscriptions require email confirmation. Users can also confirm their subscription and unsubscribe. The API also provides an endpoint to get the current weather for a city.

## Features
*   Subscribe to weather updates (requires email confirmation).
*   Confirm subscription via a unique token sent to email.
*   Unsubscribe from weather updates using a unique token.
*   Get current weather information for a specified city.

## Technologies Used
*   Go
*   PostgreSQL
*   Docker & Docker Compose
*   MailHog (for email testing)
*   golang-migrate (for database migrations)

## Project Structure
*   `cmd/weather-api/main.go`: Main application entry point.
*   `internal/`: Contains the core application logic.
    *   `config/`: Configuration management (environment variables).
    *   `handler/`: HTTP handlers for API endpoints.
    *   `repository/`: Database interaction logic.
    *   `service/`: Business logic.
    *   `models/`: Data structures.
*   `migrations/`: SQL database migration files.
*   `Dockerfile`: Defines the Docker image for the application.
*   `docker-compose.yml`: Defines services, networks, and volumes for local development.

## Getting Started

### Prerequisites
*   Docker
*   Docker Compose

### Launching the Application
1.  **Clone the repository:**
    ```bash
    git clone <repository-url>
    cd weather-api
    ```

2.  **Environment Variables:**
    The application uses environment variables for configuration. A `.env` file can be used to set these locally. The `docker-compose.yml` file already sets the necessary environment variables for the services.
    Key variables include:
    *   `DB_HOST`, `DB_PORT`, `DB_USER`, `DB_PASSWORD`, `DB_NAME` for database connection.
    *   `SMTP_HOST`, `SMTP_PORT` for the MailHog email server.
    *   `API_PORT` for the application server.

3.  **Build and run the application using Docker Compose:**
    From the root directory of the project, run:
    ```bash
    docker-compose up --build
    ```
    This command will:
    *   Build the Go application Docker image.
    *   Start the application container.
    *   Start a PostgreSQL database container.
    *   Start a MailHog container for email testing.

4.  **Accessing the Services:**
    *   **API:** `http://localhost:8080` (or the `API_PORT` you configured)
        *   `POST /subscribe` - with form data: `email`, `city`, `frequency`
        *   `GET /confirm/{token}`
        *   `GET /unsubscribe/{token}`
        *   `GET /weather?city={city_name}`
    *   **MailHog (Email UI):** `http://localhost:8025`
    *   **Database (PostgreSQL):** Accessible on `localhost:5432` (from your host machine if needed, or `db:5432` from within the app network).

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
2.  The service fetches current weather data from an external weather API (WeatherAPI.com).
3.  The weather data (temperature, humidity, description) is returned as JSON.

### Stopping the Application
To stop the application and remove the containers, networks, and volumes created by `docker-compose up`:
```bash
docker-compose down