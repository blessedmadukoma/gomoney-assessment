# Mock Premier League API - GoMoney Assessment

This project implements an API for managing teams, and providing fixtures for the Mock Premier League. It allows admins to perform CRUD operations on teams and fixtures, including generating unique links for fixtures, and basic users (or fans) to view teams, team information and fixtures.

The project is hosted on [AWS](http://ec2-13-51-79-92.eu-north-1.compute.amazonaws.com) and a [health check](http://ec2-13-51-79-92.eu-north-1.compute.amazonaws.com/api/health) to verify it works.

## Features

- **Admin Features:**
  - Signup and login.
  - Manage teams: add, remove, edit, view.
  - Create fixtures: add, remove, edit, view.
  - Generate unique links for fixtures.

- **User/Fan Features:**
  - Signup and login.
  - View teams.
  - View completed fixtures.
  - View pending fixtures.

- **Public Access:**
  - The search API for teams

## Authentication and Session Management

- Authentication and authorization for admin and user accounts are implemented using Bearer token and JWT.
- Redis is used as the session store.

## Tools/Stack
- Golang
- MongoDB
- Redis
- Docker
- Postman (for testing)

## System Design Diagrams

### Database Structure
![Database Structure](/docs/imgs/db-structure.png)

### User flow
![System Design](/docs/imgs/user-flow.png)

## Installation and Usage

1. Clone the repository:

   ```bash
   git clone https://github.com/blessedmadukoma/gomoney-assessment.git
   ```

2. Navigate to the project directory:

   ```bash
   cd gomoney-assessment
   ```

3. Install dependencies:

   ```bash
   go mod tidy
   ```

4. Set up environment variables:

   Create a `.env` file in the root directory and add the following variables:

   ```
   PORT=8089

   MONGO_INITDB_ROOT_USERNAME=admin
   MONGO_INITDB_ROOT_PASSWORD=gomoney
   MONGO_INITDB_DATABASE=gomoney_assessment_db
   MONGODB_HOST=localhost:27020
   MONGO_DB_SOURCE="mongodb://${MONGO_INITDB_ROOT_USERNAME}:${MONGO_INITDB_ROOT_PASSWORD}@${MONGODB_HOST}"

   REDIS_HOST=localhost
   REDIS_PORT=6381
   REDIS_DB_SOURCE="${REDIS_HOST}:${REDIS_PORT}"
   
   TOKEN_SYMMETRIC_KEY="12345678901234567890123456789012"
   ACCESS_TOKEN_DURATION=15m
   REFRESH_TOKEN_DURATION=24h
   
   GIN_MODE="release"
   
   LIMITER_RPS="10"
   LIMITER_BURST="5"
   LIMITER_ENABLED="false"
   ```

1. Start the server:

   ```bash
    go run main.go
   ```

   Or, having Makefile installed:

   ```bash
   make server
   ```


2. To seed the database:

   ```bash
   go run main.go -seed
   ```

   Or, having Makefile installed:

   ```bash
   make seed
   ```

3. Use Postman or any API client to interact with the endpoints. Check the file `docs/postman/Gomoney-assessement.postman_collectionn` for the postman-collection


## API Endpoints

- **Signup:**

  - `POST /api/auth/signup`

- **Login:**

  - `POST /api/auth/login`

- **Teams:**

  - `GET /api/teams`
  - `GET /api/teams/:id`
  - `POST /api/teams`
  - `PUT /api/teams/:id`
  - `DELETE /api/teams/:id`

- **Fixtures:**

  - `GET /api/fixtures`
  - `GET /api/fixtures/:id`
  - `POST /api/fixtures`
  - `PUT /api/fixtures/:id`
  - `DELETE /api/fixtures/:id`

- **Search:**
  - `GET /api/teams/search?q=query_string`
