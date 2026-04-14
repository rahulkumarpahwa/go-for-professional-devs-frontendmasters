# 10 - Integration & Complete Flow

## How All Parts Connect

The FEM project is composed of many parts working together. This final section shows how everything connects and how a request flows through the entire system.

## Architecture Overview

```
┌─────────────────────────────────────────────────────────────┐
│                    HTTP Client                              │
│              (Browser, Postman, Mobile App)                 │
└────────────────────┬────────────────────────────────────────┘
                     │ HTTP Request
                     ▼
┌─────────────────────────────────────────────────────────────┐
│                   HTTP Server                               │
│                  (Port 8080)                                │
└────────────────────┬────────────────────────────────────────┘
                     │
                     ▼
┌─────────────────────────────────────────────────────────────┐
│                   Chi Router                                │
│        r.Get(), r.Post(), r.Put(), r.Delete()              │
└────────────────────┬────────────────────────────────────────┘
                     │
                     ▼
┌─────────────────────────────────────────────────────────────┐
│                   Middleware                                │
│         Authenticate, RequireUser                           │
└────────────────────┬────────────────────────────────────────┘
                     │
                     ▼
┌─────────────────────────────────────────────────────────────┐
│                  Handler                                    │
│         (WorkoutHandler, UserHandler, TokenHandler)        │
└────────────────────┬────────────────────────────────────────┘
                     │
                     ▼
┌─────────────────────────────────────────────────────────────┐
│                   Store                                     │
│       (Database operations, queries, transactions)          │
└────────────────────┬────────────────────────────────────────┘
                     │
                     ▼
┌─────────────────────────────────────────────────────────────┐
│                PostgreSQL Database                          │
│          (users, workouts, workout_entries tables)          │
└─────────────────────────────────────────────────────────────┘
```

## Component Dependencies

```
main.go
  │
  ├─→ app.NewApplication()
  │     ├─→ store.Open()        (PostgreSQL connection)
  │     ├─→ store.MigrateFS()   (Create tables)
  │     ├─→ Stores              (WorkoutStore, UserStore, TokenStore)
  │     ├─→ Handlers            (WorkoutHandler, UserHandler, TokenHandler)
  │     └─→ Middleware          (Authentication, RequireUser)
  │
  ├─→ routes.SetupRoutes(app)  (Register all endpoints)
  │
  └─→ Server.ListenAndServe()  (Listen for requests)
```

## Complete Request-Response Cycle

### Example: Create Workout with Full Flow

```
STEP 1: CLIENT SENDS REQUEST
═══════════════════════════════════
POST /workouts HTTP/1.1
Authorization: Bearer eyJhbGciOiJIUzI1NiIsInR5cCI6IkpXVCJ9...
Content-Type: application/json

{
  "title": "Morning Strength Training",
  "description": "Full body workout",
  "duration_minutes": 60,
  "calories_burned": 500,
  "entries": [
    {
      "exercise_name": "Squats",
      "sets": 3,
      "reps": 12,
      "weight": 100.5,
      "order_index": 1
    }
  ]
}


STEP 2: HTTP SERVER RECEIVES REQUEST
═════════════════════════════════════
Server := &http.Server{
    Addr: ":8080",
    Handler: r,  // chi router
}
Server.ListenAndServe()  ← Request received here


STEP 3: CHI ROUTER MATCHES ROUTE
═════════════════════════════════════════
Route matching:
  Method: POST ✓
  Path: /workouts ✓

Inside r.Group() with r.Use(app.Middleware.Authenticate) ✓
Uses r.Middleware.RequireUser wrapper ✓

Route found: r.Post("/workouts",
    app.Middleware.RequireUser(
        app.WorkoutHandler.HandleCreateWorkout))


STEP 4: AUTHENTICATE MIDDLEWARE RUNS
═════════════════════════════════════════════
func (um *UserMiddleware) Authenticate(next http.Handler) http.Handler {
    return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {

        // Read Authorization header
        authHeader := r.Header.Get("Authorization")
        // authHeader = "Bearer eyJhbGc..."

        // Parse "Bearer <token>"
        headerParts := strings.Split(authHeader, " ")
        // headerParts = ["Bearer", "eyJhbGc..."]

        token := headerParts[1]
        // token = "eyJhbGc..."

        // Get user from database using token
        user, err := um.UserStore.GetUserToken(tokens.ScopeAuth, token)
        // Query: SELECT * FROM users WHERE id = (SELECT user_id FROM tokens WHERE token = ?)
        // Result: user = &User{ID: 1, Username: "john_doe", ...}

        // Store user in context
        SetUser(r, user)
        // r.Context() now contains user info

        next.ServeHTTP(w, r)  // Continue to RequireUser
    })
}


STEP 5: REQUIREUSER MIDDLEWARE RUNS
═════════════════════════════════════════════
func (um *UserMiddleware) RequireUser(next http.HandlerFunc) http.HandlerFunc {
    return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {

        // Get user from context
        user := GetUser(r)
        // user = &User{ID: 1, Username: "john_doe", ...}

        // Check if anonymous
        if user.IsAnonymous() {  // ID == 0
            // Not executed because user is authenticated
            return
        }

        next.ServeHTTP(w, r)  // Continue to handler
    })
}


STEP 6: HANDLER RUNS - HandleCreateWorkout
═════════════════════════════════════════════════════════════
func (h *WorkoutHandler) HandleCreateWorkout(w http.ResponseWriter, r *http.Request) {

    // Parse JSON request body
    var workout store.Workout
    json.NewDecoder(r.Body).Decode(&workout)
    // workout = Workout{
    //   Title: "Morning Strength Training",
    //   Description: "Full body workout",
    //   DurationMinutes: 60,
    //   CaloriesBurned: 500,
    //   Entries: [{exercise_name: "Squats", ...}]
    // }

    // Get current user from context
    currentUser := middleware.GetUser(r)
    // currentUser = &User{ID: 1, Username: "john_doe", ...}

    // Set user ID
    workout.UserID = currentUser.ID
    // workout.UserID = 1

    // Call store to save
    createdWorkout, err := h.store.CreateWorkout(&workout)
    // (continues to STEP 7)

    // Return response
    utils.WriteJson(w, http.StatusOK,
        utils.Envelope{"CreatedWorkout": createdWorkout})
}


STEP 7: STORE OPERATION - CreateWorkout
═════════════════════════════════════════════════════════════
func (pgws *PostgresWorkoutStore) CreateWorkout(workout *Workout) (*Workout, error) {

    // Begin transaction
    tx, err := pgws.db.Begin()

    // Insert workout
    query := `INSERT INTO workouts
              (title, description, duration_minutes, calories_burned, user_id)
              VALUES ($1, $2, $3, $4, $5)
              RETURNING id;`

    err = tx.QueryRow(query,
        "Morning Strength Training",
        "Full body workout",
        60,
        500,
        1,
    ).Scan(&workout.ID)
    // workout.ID = 1 (assigned by database)

    // Insert entries
    for index := range workout.Entries {
        entry := &workout.Entries[index]

        query := `INSERT INTO workout_entries
                  (workout_id, exercise_name, sets, reps, duration_seconds, weight, notes, order_index)
                  VALUES ($1, $2, $3, $4, $5, $6, $7, $8)
                  RETURNING id;`

        err := tx.QueryRow(query,
            1,              // workout_id
            "Squats",       // exercise_name
            3,              // sets
            &12,            // reps
            nil,            // duration_seconds
            &100.5,         // weight
            "",             // notes
            1,              // order_index
        ).Scan(&entry.ID)
        // entry.ID = 1 (assigned by database)
    }

    // Commit transaction
    err = tx.Commit()

    return workout, nil
}


STEP 8: DATABASE OPERATIONS
═════════════════════════════════════════════════════════════
PostgreSQL received:
  1. BEGIN TRANSACTION
  2. INSERT INTO workouts (...)  VALUES (...)
     ├─ Check NOT NULL constraints ✓
     ├─ Check unique constraints ✓
     └─ INSERT successful, return id=1
  3. INSERT INTO workout_entries (...) VALUES (...)
     ├─ Check foreign key (workout_id=1 exists) ✓
     ├─ Check CHECK constraint (reps XOR duration_seconds) ✓
     └─ INSERT successful, return id=1
  4. COMMIT TRANSACTION ✓

Database state after:
  workouts table:
    id | user_id | title | ... | created_at
    1  | 1       | Morning Strength Training | ... | 2024-01-15 10:30:00

  workout_entries table:
    id | workout_id | exercise_name | sets | reps | ...
    1  | 1          | Squats        | 3    | 12   | ...


STEP 9: RESPONSE PREPARED
═════════════════════════════════════════════════════════════
Handler receives:
  createdWorkout = Workout{
    ID: 1,
    UserID: 1,
    Title: "Morning Strength Training",
    Description: "Full body workout",
    DurationMinutes: 60,
    CaloriesBurned: 500,
    CreatedAt: "2024-01-15T10:30:00Z",
    UpdatedAt: "2024-01-15T10:30:00Z",
    Entries: [{
      ID: 1,
      WorkoutId: 1,
      ExerciseName: "Squats",
      Sets: 3,
      Reps: &12,
      DurationSeconds: nil,
      Weight: &100.5,
      Notes: "",
      OrderIndex: 1,
    }]
  }

Handler calls:
  utils.WriteJson(w, http.StatusOK, Envelope{"CreatedWorkout": createdWorkout})

  This function:
    1. Sets header: Content-Type: application/json
    2. Sets status: 200 OK
    3. Encodes workout struct to JSON
    4. Writes JSON to response body


STEP 10: HTTP RESPONSE SENT
═════════════════════════════════════════════════════════════
HTTP/1.1 200 OK
Content-Type: application/json
Content-Length: 487

{
  "CreatedWorkout": {
    "id": 1,
    "user_id": 1,
    "title": "Morning Strength Training",
    "description": "Full body workout",
    "duration_minutes": 60,
    "calories_burned": 500,
    "created_at": "2024-01-15T10:30:00Z",
    "updated_at": "2024-01-15T10:30:00Z",
    "entries": [
      {
        "id": 1,
        "workout_id": 1,
        "exercise_name": "Squats",
        "sets": 3,
        "reps": 12,
        "duration_seconds": null,
        "weight": 100.5,
        "notes": "",
        "order_index": 1
      }
    ]
  }
}


STEP 11: CLIENT RECEIVES RESPONSE
═════════════════════════════════════
Client parses JSON response
Shows workout created successfully
Displays: Workout 1 - Morning Strength Training
```

## Data Flow Diagram

```
┌──────────────────┐
│  HTTP Request    │
│  POST /workouts  │
│  + Token         │
│  + JSON body     │
└────────┬─────────┘
         │
         ▼
┌──────────────────────────────────────┐
│  Router Match                        │
│  POST /workouts ✓                    │
│  Group with Authenticate ✓           │
└────────┬─────────────────────────────┘
         │
         ▼
┌──────────────────────────────────────┐
│  Middleware 1: Authenticate          │
│  - Read header                       │
│  - Validate token                    │
│  - Get user from DB                  │
│  - Set in context                    │
└────────┬─────────────────────────────┘
         │
         ▼
┌──────────────────────────────────────┐
│  Middleware 2: RequireUser           │
│  - Get user from context             │
│  - Check not anonymous               │
│  - Allow access                      │
└────────┬─────────────────────────────┘
         │
         ▼
┌──────────────────────────────────────┐
│  Handler                             │
│  - Decode JSON                       │
│  - Set user ID                       │
│  - Call store                        │
└────────┬─────────────────────────────┘
         │
         ▼
┌──────────────────────────────────────┐
│  Store (Database Layer)              │
│  - Begin transaction                 │
│  - INSERT workout                    │
│  - INSERT workout_entries            │
│  - Commit                            │
└────────┬─────────────────────────────┘
         │
         ▼
┌──────────────────────────────────────┐
│  Database                            │
│  - Validate constraints              │
│  - Insert rows                       │
│  - Return new IDs                    │
└────────┬─────────────────────────────┘
         │
         ▼
┌──────────────────────────────────────┐
│  Response Builder                    │
│  - Convert to JSON                   │
│  - Set status 200                    │
│  - Write response                    │
└────────┬─────────────────────────────┘
         │
         ▼
┌──────────────────┐
│  HTTP Response   │
│  200 OK          │
│  JSON body       │
└──────────────────┘
```

## Error Handling Integration

```
Any step fails:
  │
  ├─ Validation fails (step 6)
  │   └─→ 400 Bad Request + error message
  │
  ├─ Token invalid (step 4)
  │   └─→ 401 Unauthorized
  │
  ├─ User not authenticated (step 5)
  │   └─→ 401 Unauthorized
  │
  ├─ Database constraint fails (step 8)
  │   └─→ 500 Internal Server Error
  │
  └─ Other error (any step)
      └─→ 500 Internal Server Error
```

## Parts Summary

| Part | Component        | Purpose                | Files                                |
| ---- | ---------------- | ---------------------- | ------------------------------------ |
| 1    | Project Overview | Understand what FEM is | (this doc)                           |
| 2    | Entry Point      | How app starts         | main.go, app.go                      |
| 3    | Database         | Data storage           | database.go, migrations/\*.sql       |
| 4    | Auth/Middleware  | Security               | middleware.go                        |
| 5    | User Management  | User operations        | user_handler.go, token_handler.go    |
| 6    | Workout CRUD     | Workout operations     | workout_handler.go, workout_store.go |
| 7    | Data Models      | Struct definitions     | store structs                        |
| 8    | Routing          | URL matching           | routes.go                            |
| 9    | Error Handling   | Error management       | Error patterns throughout            |
| 10   | Integration      | How it all connects    | (this doc)                           |

## Key Connections

```
User Registration
├─ Handler: HandleRegisterUser
├─ Store: CreateUser
├─ Database: INSERT into users
└─ Response: User created successfully

User Authentication (Login)
├─ Handler: HandleCreateToken
├─ Store: GetUserByUsername
├─ Store: CreateNewToken
├─ Database: INSERT into tokens
└─ Response: Token returned

Workout Creation
├─ Middleware: Authenticate (get user)
├─ Middleware: RequireUser (verify authenticated)
├─ Handler: HandleCreateWorkout
├─ Store: CreateWorkout
├─ Database:
│   ├─ INSERT into workouts
│   └─ INSERT into workout_entries
└─ Response: Workout created

Workout Retrieval
├─ Middleware: Authenticate (get user)
├─ Middleware: RequireUser (verify authenticated)
├─ Handler: HandleGetWorkoutById
├─ Store: GetWorkoutById
├─ Database:
│   ├─ SELECT from workouts
│   └─ SELECT from workout_entries
└─ Response: Workout data

Workout Update
├─ Middleware: Authenticate (get user)
├─ Middleware: RequireUser (verify authenticated)
├─ Handler: HandleUpdateWorkoutByID
├─ Store: GetWorkoutOwner (verify ownership)
├─ Store: UpdateWorkout
├─ Database:
│   ├─ UPDATE workouts
│   ├─ DELETE from workout_entries
│   └─ INSERT into workout_entries
└─ Response: Update successful

Workout Deletion
├─ Middleware: Authenticate (get user)
├─ Middleware: RequireUser (verify authenticated)
├─ Handler: HandleDeleteWorkoutByID
├─ Store: GetWorkoutOwner (verify ownership)
├─ Store: DeleteWorkout
├─ Database:
│   ├─ Database triggers CASCADE
│   └─ DELETE from workouts (and entries)
└─ Response: Deletion successful
```

## Technology Stack Integration

```
HTTP Request (JSON)
        │
        ▼
json.NewDecoder() ◄─── encoding/json
        │
        ▼
Chi Router ◄─── github.com/go-chi/chi/v5
        │
        ▼
Middleware Layer
        │
        ▼
Handler (Go function)
        │
        ▼
Store Interface
        │
        ▼
SQL Queries ◄─── database/sql
        │
        ▼
PGX Driver ◄─── github.com/jackc/pgx/v4
        │
        ▼
PostgreSQL
        │
        ▼
Response Builder
        │
        ▼
json.NewEncoder() ◄─── encoding/json
        │
        ▼
HTTP Response (JSON)
```

## Conclusion

The FEM project follows a clean layered architecture:

1. **HTTP Layer** - Receives requests, sends responses
2. **Routing Layer** - Matches URLs to handlers
3. **Middleware Layer** - Handles authentication
4. **Handler Layer** - Processes requests
5. **Store Layer** - Abstracts database operations
6. **Database Layer** - PostgreSQL persistence

Each layer has clear responsibilities and interfaces. This makes the code:

- **Maintainable** - Easy to find and fix bugs
- **Testable** - Each layer can be tested independently
- **Scalable** - Easy to add new features
- **Secure** - Auth checked at middleware layer

---

## Next Steps

You've completed all 10 parts! Review this documentation when:

- Learning how the project works
- Implementing new features
- Debugging issues
- Onboarding new developers

**Key Concept**: FEM is a complete, production-ready API. Every request goes through a choreographed dance of authentication, validation, processing, and response!
