# 01 - FEM Project Overview

## What is the FEM Project?

The **FEM (Frontend Masters)** project is a **Fitness Workout Tracking API** built with **Go** and **PostgreSQL**. It's a RESTful backend service that allows users to:

- Register and authenticate
- Create, read, update, and delete workouts
- Track individual exercises within each workout
- Store detailed exercise metrics (reps, sets, weight, duration, etc.)

## Project Type

This is a **Backend REST API** built with:

- **Language**: Go
- **Database**: PostgreSQL
- **Framework**: Chi (HTTP router)
- **Database Migration**: Goose
- **Authentication**: Token-based (Bearer tokens)

## Project Structure

```
FEM/
├── main.go                          # Application entry point
├── docker-compose.yml               # Docker configuration
├── go.mod                           # Go module dependencies
├── internal/
│   ├── app/
│   │   └── app.go                   # Application initialization
│   ├── api/
│   │   ├── workout_handler.go       # Workout API endpoints
│   │   ├── user_handler.go          # User registration endpoints
│   │   └── token_handler.go         # Authentication endpoints
│   ├── middleware/
│   │   └── middleware.go            # Authentication middleware
│   ├── store/
│   │   ├── database.go              # Database connection & migrations
│   │   ├── workout_store.go         # Workout database operations
│   │   ├── user_store.go            # User database operations
│   │   └── token_store.go           # Token database operations
│   ├── utils/
│   │   └── utils.go                 # Helper functions
│   └── tokens/
│       └── tokens.go                # Token generation & validation
├── migrations/
│   ├── 00001_users.sql              # Users table migration
│   ├── 00002_workouts.sql           # Workouts table migration
│   ├── 00003_workout_entries.sql    # Workout entries table migration
│   └── fs.go                        # Embedded migration files
└── notes/
    └── (documentation files)
```

## Core Concepts

### 1. **Users**

- Register with username, email, password, and bio
- Credentials are stored securely with password hashing
- Each user has a unique ID

### 2. **Authentication**

- Token-based system using Bearer tokens
- Tokens expire after specified duration (24 hours)
- Tokens are validated on protected routes

### 3. **Workouts**

- Each workout belongs to a specific user
- Contains metadata: title, description, duration, calories burned
- Can have multiple exercises (entries)

### 4. **Workout Entries**

- Individual exercises within a workout
- Track: exercise name, sets, reps, weight, duration
- Must have either reps or duration_seconds (but not both)

## API Endpoints

### Public Routes

- `POST /users` - Register a new user
- `POST /tokens/authentication` - Login and get auth token
- `GET /health` - Health check

### Protected Routes (Require Authentication)

- `GET /workouts/{id}` - Get a specific workout
- `POST /workouts` - Create a new workout
- `PUT /workouts/{id}` - Update a workout
- `DELETE /workouts/{id}` - Delete a workout

## How It Works (Simple Overview)

1. **User registers** → Credentials stored in database
2. **User logs in** → Receives token
3. **User creates workout** → Sends request with Bearer token
4. **Server validates token** → Identifies user
5. **Workout saved** → Linked to authenticated user
6. **Data persisted** → Stored in PostgreSQL

## Technologies Used

| Technology | Purpose                    |
| ---------- | -------------------------- |
| Go         | Programming language       |
| Chi        | HTTP routing               |
| PostgreSQL | Database                   |
| Goose      | Database migration tool    |
| PGX        | PostgreSQL driver          |
| Context    | Request context management |

## Next Steps

This documentation is organized into 10 parts:

1. **Project Overview** (This file)
2. Application Entry Point & Initialization
3. Database & Migrations
4. Authentication & Middleware
5. User Management
6. Workout Operations
7. Data Models & Structures
8. Routing System
9. Error Handling & Utils
10. Integration & Complete Flow

---

**Start with Part 2** to understand how the application initializes and runs.
