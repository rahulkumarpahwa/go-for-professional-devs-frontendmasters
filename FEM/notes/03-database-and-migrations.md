# 03 - Database & Migrations

## What is a Database?

A database is a persistent storage system that keeps your data safe and organized. The FEM project uses **PostgreSQL**, a powerful relational database that stores data in tables with rows and columns.

## Database Connection

### Opening a Database Connection

Located in `internal/store/database.go`:

```go
func Open() (*sql.DB, error) {
    // Create connection to PostgreSQL
    DB, err := sql.Open("pgx", "host=localhost user=postgres password=postgres dbname=postgres port=5433 sslmode=disable")
    if err != nil {
        return nil, fmt.Errorf("DB : open %w\n", err)
    }

    // Test the connection
    if err := DB.Ping(); err != nil {
        return nil, fmt.Errorf("db: open %w", err)
    }

    // Configure connection pool
    DB.SetMaxOpenConns(25)        // Max 25 open connections
    DB.SetMaxIdleConns(10)        // Keep 10 idle connections
    DB.SetConnMaxIdleTime(5 * time.Minute)    // Close idle after 5 min
    DB.SetConnMaxLifetime(30 * time.Minute)   // Close after 30 min

    fmt.Println("Connected to Database Successfully!")
    return DB, nil
}
```

**Simple Explanation:**

- Connects to PostgreSQL at `localhost:5433`
- Uses credentials: user=postgres, password=postgres
- Configures how many connections to keep open
- Tests the connection to ensure it works

### Connection String Breakdown

```
host=localhost          - Where the database runs
user=postgres           - Username to authenticate
password=postgres       - Password to authenticate
dbname=postgres         - Database name
port=5433              - Port number
sslmode=disable        - Don't require SSL (for local dev)
```

## Database Migrations

A **migration** is a version-controlled script that creates or modifies database tables. The FEM project uses **Goose** for migrations.

### Why Migrations?

Without migrations:

- ❌ Database tables created manually
- ❌ No version control
- ❌ Hard to sync with other developers
- ❌ Difficult to rollback changes

With migrations:

- ✅ Tables created by code
- ✅ Version controlled (in git)
- ✅ Everyone has same schema
- ✅ Can easily undo changes

### Migration Functions

```go
func MigrateFS(db *sql.DB, migrationsFS fs.FS, dir string) error {
    // Tell Goose where to read migration files
    goose.SetBaseFS(migrationsFS)

    defer func() {
        goose.SetBaseFS(nil)
    }()

    return Migrate(db, dir)
}

func Migrate(db *sql.DB, dir string) error {
    // Set database dialect to PostgreSQL
    err := goose.SetDialect("postgres")
    if err != nil {
        return fmt.Errorf("Migrate : %w", err)
    }

    // Run all pending migrations
    err = goose.Up(db, dir)
    if err != nil {
        return fmt.Errorf("Goose Up : %w", err)
    }
    return nil
}
```

**Simple Explanation:**

- `MigrateFS`: Uses embedded migration files
- `Migrate`: Runs all pending migrations
- `goose.Up`: Applies migrations in order

## The Three Migrations

### 1. Users Table (00001_users.sql)

```sql
-- +goose Up
CREATE TABLE IF NOT EXISTS users (
    id BIGSERIAL PRIMARY KEY,
    username VARCHAR(50) UNIQUE NOT NULL,
    email VARCHAR(255) UNIQUE NOT NULL,
    password_hash VARCHAR(255) NOT NULL,
    bio TEXT,
    created_at TIMESTAMP WITH TIME ZONE DEFAULT CURRENT_TIMESTAMP,
    updated_at TIMESTAMP WITH TIME ZONE DEFAULT CURRENT_TIMESTAMP
)

-- +goose Down
Drop TABLE users;
```

**Explanation:**

- `id`: Unique identifier (auto-incrementing)
- `username`: Must be unique, stores login name
- `email`: Must be unique, stores email address
- `password_hash`: Stores hashed password (NOT plain text!)
- `bio`: Optional text field for user bio
- `UNIQUE`: Ensures no duplicate usernames or emails
- `NOT NULL`: Field is required
- Timestamps: Track when created and last updated

**Table Structure:**

```
id | username | email | password_hash | bio | created_at | updated_at
---+----------+-------+---------------+-----+------------+-----------
1  | john_doe | j@... | $2a$10$...    | ... | 2024-01-01 | 2024-01-01
2  | jane_doe | j@... | $2a$10$...    | ... | 2024-01-02 | 2024-01-02
```

### 2. Workouts Table (00002_workouts.sql)

```sql
-- +goose Up
CREATE TABLE IF NOT EXISTS workouts (
    id BIGSERIAL PRIMARY KEY,
    user_id BIGINT NOT NULL REFERENCES users(id),
    title VARCHAR(255) NOT NULL,
    description TEXT,
    duration_minutes INTEGER NOT NULL,
    calories_burned INTEGER,
    created_at TIMESTAMP WITH TIME ZONE DEFAULT CURRENT_TIMESTAMP,
    updated_at TIMESTAMP WITH TIME ZONE DEFAULT CURRENT_TIMESTAMP
)

-- +goose Down
Drop TABLE workouts;
```

**Explanation:**

- `id`: Unique workout identifier
- `user_id`: Foreign key linking to users table
- `REFERENCES users(id)`: Ensures user_id exists in users table
- `title`: Name of the workout (e.g., "Morning Run")
- `description`: Details about the workout
- `duration_minutes`: How long the workout took
- `calories_burned`: Energy expended

**Relationship:**

```
users table          workouts table
───────────          ──────────────
id (1) ──────────┐   id (1)
username         └──→ user_id (1 to many)
email                title
password_hash        description
```

One user can have many workouts!

### 3. Workout Entries Table (00003_workout_entries.sql)

```sql
-- +goose Up
CREATE TABLE IF NOT EXISTS workout_entries (
    id BIGSERIAL PRIMARY KEY,
    workout_id BIGINT NOT NULL REFERENCES workouts(id) ON DELETE CASCADE,
    exercise_name VARCHAR(255) NOT NULL,
    sets INTEGER NOT NULL,
    reps INTEGER,
    duration_seconds INTEGER,
    weight DECIMAL(5,2),
    notes TEXT,
    order_index INTEGER NOT NULL,
    created_at TIMESTAMP WITH TIME ZONE DEFAULT CURRENT_TIMESTAMP,
    CONSTRAINT valid_workout_entry CHECK (
        (reps IS NOT NULL OR duration_seconds IS NOT NULL) AND
        (reps IS NULL OR duration_seconds IS NULL)
    )
)

-- +goose Down
Drop TABLE workout_entries;
```

**Explanation:**

- `id`: Unique exercise entry identifier
- `workout_id`: Which workout this exercise belongs to
- `ON DELETE CASCADE`: If workout deleted, delete its entries
- `exercise_name`: Exercise name (e.g., "Squats")
- `sets`: Number of sets performed
- `reps`: Repetitions (for strength training)
- `duration_seconds`: Duration (for cardio)
- `weight`: Weight lifted (e.g., 100.5 kg)
- `notes`: Additional notes about the exercise
- `order_index`: Order of exercises in workout
- **CHECK CONSTRAINT**: Either reps OR duration_seconds, but NOT both

**The Constraint Explained:**

```
✓ Valid:   reps=12, duration_seconds=NULL
✓ Valid:   reps=NULL, duration_seconds=300
✗ Invalid: reps=12, duration_seconds=300
✗ Invalid: reps=NULL, duration_seconds=NULL
```

**Relationship:**

```
workouts table           workout_entries table
──────────────          ─────────────────────
id (1) ────────┐        id (1)
title          └───────→ workout_id (1 to many)
description            exercise_name
duration_minutes       sets
                       reps
                       duration_seconds
```

One workout can have many entries!

## Complete Data Model

```
users (1) ───── (many) workouts (1) ───── (many) workout_entries
  └── username            └── title              └── exercise_name
  └── email               └── duration_minutes   └── sets
  └── password_hash       └── calories_burned    └── reps
  └── bio                 └── created_at         └── weight
```

## Migration Flow

```
┌──────────────────────────────────────┐
│   Application Starts                 │
└────────────┬─────────────────────────┘
             │
             ▼
┌──────────────────────────────────────┐
│   Connect to PostgreSQL              │
└────────────┬─────────────────────────┘
             │
             ▼
┌──────────────────────────────────────┐
│   Read Migration Files               │
│   (00001, 00002, 00003)              │
└────────────┬─────────────────────────┘
             │
             ▼
┌──────────────────────────────────────┐
│   Check Goose History Table          │
│   (Track which migrations ran)       │
└────────────┬─────────────────────────┘
             │
             ▼
┌──────────────────────────────────────┐
│   Run New Migrations                 │
│   (Skip already-run ones)            │
└────────────┬─────────────────────────┘
             │
             ├──► Migration 1: Create users
             │
             ├──► Migration 2: Create workouts
             │
             └──► Migration 3: Create entries
             │
             ▼
┌──────────────────────────────────────┐
│   Database Ready                     │
└──────────────────────────────────────┘
```

## Thinking About Data Flow

```
HTTP Request with JSON
    │
    ▼
Handler parses JSON into Go struct
    │
    ▼
Handler calls Store method
    │
    ▼
Store creates SQL INSERT/UPDATE/DELETE
    │
    ▼
SQL query runs against PostgreSQL
    │
    ▼
Database validates and stores data
    │
    ▼
Returns result to Go application
    │
    ▼
Handler converts to JSON response
    │
    ▼
HTTP Response sent to client
```

## Key Takeaways

1. **Connection Pool**: Manages multiple database connections efficiently
2. **Migrations**: Version-controlled database schema changes
3. **Foreign Keys**: Ensure data relationships are valid
4. **Constraints**: Enforce business rules at database level
5. **Cascade Delete**: Automatically clean up related data
6. **Three Tables**: Users → Workouts → Workout Entries

## Next Steps

- Part 4: Learn about **Authentication & Middleware** - securing access
- Part 5: Learn about **User Management** - handling user operations

---

**Key Concept**: Migrations are like git commits for your database schema. Every change is tracked and can be rolled back!
