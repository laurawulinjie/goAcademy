# Go Todo Project

A multi-app Go project that includes:

- **🖥️ `todo`** – A web-based Todo app with PostgreSQL, REST API, HTML pages, and user authentication.
- **🧾 `todo_cli`** – A command-line Todo app for managing tasks in the terminal

## 🚀 1. Run the Web App (`todo`)

### Requirements

- Go
- PostgreSQL

### Setup

#### 1. Install PostgreSQL (if needed)

```bash
brew install postgresql
brew services start postgresql
```

#### 2. Create database and table

```bash
createdb todoapp
psql todoapp
```

In the `psql` shell, run:

```bash
CREATE TABLE users (
  id SERIAL PRIMARY KEY,
  username TEXT NOT NULL UNIQUE,
  password TEXT NOT NULL
);

CREATE TABLE todos (
  id SERIAL PRIMARY KEY,
  task TEXT NOT NULL,
  status TEXT NOT NULL,
  user_id INTEGER REFERENCES users(id)
);
```

#### 3. Set up your `.env` file

Create `.env` in the root by copy & paste and rename `.env.example`.
Replace `replace_with_your_username` with your actual user.

#### 4. Run the app

```bash
cd cmd/todo
go run ./
```

open your browser at: http://localhost:8080/
