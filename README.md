# JSON Parser

A Go-based backend service with a React frontend for parsing workflow JSON files using spreadsheet-defined service mappings.

---

## Project Structure

```
.
├── backend
│   ├── cmd
│   ├── internal
│   ├── scripts
│   ├── Dockerfile
│   ├── go.mod
│   ├── go.sum
│   ├── tern.conf.sample
│   └── .env.sample
├── frontend
│   ├── src
│   ├── public
│   ├── Dockerfile
│   ├── package.json
│   └── ...
└── docker-compose.yml
```

---

# Running with Docker

## Build Images

```bash
docker compose build --no-cache
```

## Start the Application

```bash
docker compose up
```

Services:

- **Backend:** http://localhost:5555
- **Frontend:** http://localhost:5173

---

# Local Development

## Backend

Navigate to the backend directory.

```bash
cd backend
```

### 1. Configure Environment

Copy the sample environment file.

```bash
cp .env.sample .env
```

Fill in the required database configuration.

### 2. Configure Tern

Copy the sample configuration.

```bash
cp tern.conf.sample tern.conf
```

Update the database credentials in `tern.conf`.

### 3. Run Database Migrations

```bash
tern migrate \
  --migrations ./internal/db/migrations
```

### 4. Start the Backend

```bash
go run ./cmd
```

The backend will be available at:

```
http://localhost:5555
```

---

## Frontend

Navigate to the frontend directory.

```bash
cd frontend
```

Set the development environment.

```bash
export NODE_ENV=development
```

Install dependencies.

```bash
npm install
```



Start the development server.

```bash
npm run dev
```

The frontend will be available at:

```
http://localhost:5173
```

---

# Requirements

### Backend

- Go
- PostgreSQL
- Tern

### Frontend

- Node.js
- npm

---

# Notes

- Database migrations must be executed before starting the backend when running locally.
- Ensure PostgreSQL is running and matches the credentials specified in both `.env` and `tern.conf`.