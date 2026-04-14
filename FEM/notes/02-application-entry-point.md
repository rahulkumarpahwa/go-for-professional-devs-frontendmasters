# 02 - Application Entry Point & Initialization

## Understanding the Entry Point

The entry point of the FEM application is in `main.go`. This is where the entire application starts, initializes all components, and begins listening for HTTP requests.

## main.go - Step by Step

### 1. Parse CLI Flags

```go
var port int
flag.IntVar(&port, "port", 8080, "go backend server port")
flag.Parse()
```

**What it does:**

- Creates a command-line flag for specifying the port
- Default port is **8080**
- Allows users to run: `go run main.go -port 9000`

**Simple Explanation:**
Think of CLI flags as settings you can pass when starting the app. Instead of hardcoding the port, we let users customize it.

### 2. Initialize the Application

```go
app, err := app.NewApplication()
if err != nil {
    panic(err)
}
defer app.DB.Close()
```

**What it does:**

- Calls `NewApplication()` which initializes everything
- `defer app.DB.Close()` ensures database connection closes when main exits

**Simple Explanation:**
`NewApplication()` is like a constructor that sets up all the pieces we need (database, handlers, routes, etc.).

### 3. Setup Routes

```go
r := routes.SetupRoutes(app)
```

**What it does:**

- Configures all HTTP routes (API endpoints)
- Registers middleware for authentication
- Associates routes with handler functions

**Simple Explanation:**
This tells the application what endpoints exist and what code should handle each one.

### 4. Create HTTP Server

```go
Server := &http.Server{
    Addr:         fmt.Sprintf(":%d", port),
    Handler:      r,
    IdleTimeout:  time.Minute,
    ReadTimeout:  10 * time.Second,
    WriteTimeout: 30 * time.Second,
}
```

**What it does:**

- Creates an HTTP server with specific timeouts
- `Addr`: The address and port to listen on (e.g., `:8080`)
- `Handler`: The router (all our routes)
- `IdleTimeout`: How long before closing idle connections
- `ReadTimeout`: How long to wait for client to send request
- `WriteTimeout`: How long to wait for server to send response

**Simple Explanation:**
This is like setting up the actual web server with rules about how long it waits for connections.

### 5. Start the Server

```go
app.Logger.Printf("App Started Running on PORT %d!\n", port)

err = Server.ListenAndServe()
if err != nil {
    app.Logger.Fatal(err)
}
```

**What it does:**

- Logs that the server is starting
- Begins listening for HTTP requests
- Blocks indefinitely (the server runs continuously)

**Simple Explanation:**
This starts the server and keeps it running forever until an error occurs.

## Complete main.go Code

```go
package main

import (
    "flag"
    "fmt"
    "net/http"
    "time"

    "github.com/rahulkumarpahwa/femProject/internal/app"
    "github.com/rahulkumarpahwa/femProject/internal/routes"
)

func main() {
    // Step 1: Parse CLI flags
    var port int
    flag.IntVar(&port, "port", 8080, "go backend server port")
    flag.Parse()

    // Step 2: Initialize the application
    app, err := app.NewApplication()
    if err != nil {
        panic(err)
    }
    defer app.DB.Close()

    // Step 3: Setup routes
    r := routes.SetupRoutes(app)

    // Step 4: Create HTTP server with timeouts
    Server := &http.Server{
        Addr:         fmt.Sprintf(":%d", port),
        Handler:      r,
        IdleTimeout:  time.Minute,
        ReadTimeout:  10 * time.Second,
        WriteTimeout: 30 * time.Second,
    }

    // Step 5: Start the server
    app.Logger.Printf("App Started Running on PORT %d!\n", port)

    err = Server.ListenAndServe()
    if err != nil {
        app.Logger.Fatal(err)
    }
}
```

## Application Initialization - app.NewApplication()

### What Gets Initialized

The `NewApplication()` function in `internal/app/app.go` does the following:

```go
func NewApplication() (*Application, error) {
    // 1. Connect to PostgreSQL
    pgDB, err := store.Open()

    // 2. Run database migrations
    err = store.MigrateFS(pgDB, migrations.FS, ".")

    // 3. Create logger
    logger := log.New(os.Stdout, "", log.Ldate|log.Ltime)

    // 4. Initialize stores (database access layers)
    workoutStore := store.NewPostgresWorkoutStore(pgDB, logger)
    userStore := store.NewPostgresUserStore(pgDB, logger)
    tokenStore := store.NewPostgresTokenStore(pgDB)

    // 5. Initialize handlers (HTTP request handlers)
    workoutHandler := api.NewWorkoutHandler(workoutStore, logger)
    userHandler := api.NewUserHandler(userStore, logger)
    tokenHandler := api.NewTokenHandler(tokenStore, userStore, logger)

    // 6. Initialize middleware
    middlewareHandler := middleware.UserMiddleware{UserStore: userStore}

    // 7. Create and return Application struct
    app := &Application{
        Logger: logger,
        WorkoutHandler: workoutHandler,
        UserHandler: userHandler,
        TokenHandler: tokenHandler,
        Middleware: middlewareHandler,
        DB: pgDB,
    }
    return app, nil
}
```

## Application Structure

```
Application Struct
├── DB: *sql.DB                          (database connection)
├── Logger: *log.Logger                  (logging)
├── WorkoutHandler: *api.WorkoutHandler  (handles workout requests)
├── UserHandler: *api.UserHandler        (handles user registration)
├── TokenHandler: *api.TokenHandler      (handles authentication)
└── Middleware: middleware.UserMiddleware (handles auth on requests)
```

## Startup Flow Diagram

```
┌─────────────────────────────────────┐
│   Program Starts (main.go)          │
└────────────┬────────────────────────┘
             │
             ▼
┌─────────────────────────────────────┐
│   Parse CLI Flags (port: 8080)      │
└────────────┬────────────────────────┘
             │
             ▼
┌─────────────────────────────────────┐
│   NewApplication()                  │
├─────────────────────────────────────┤
│ 1. Connect to PostgreSQL            │
│ 2. Run Migrations                   │
│ 3. Create Logger                    │
│ 4. Initialize Stores                │
│ 5. Initialize Handlers              │
│ 6. Initialize Middleware            │
└────────────┬────────────────────────┘
             │
             ▼
┌─────────────────────────────────────┐
│   SetupRoutes()                     │
│   (Register all API endpoints)      │
└────────────┬────────────────────────┘
             │
             ▼
┌─────────────────────────────────────┐
│   Create HTTP Server                │
│   (Set timeouts, port, etc)         │
└────────────┬────────────────────────┘
             │
             ▼
┌─────────────────────────────────────┐
│   Server.ListenAndServe()           │
│   (Listen for requests)             │
└─────────────────────────────────────┘
```

## Running the Application

### Start with default port (8080)

```bash
go run main.go
```

### Start with custom port (9000)

```bash
go run main.go -port 9000
```

### Using Docker Compose

```bash
docker compose up --build
```

## Key Takeaways

1. **Entry Point**: `main.go` is where everything starts
2. **Initialization**: `NewApplication()` sets up all components
3. **CLI Flags**: Allow customization at startup (port)
4. **Timeouts**: Critical for production-ready applications
5. **Cleanup**: `defer` ensures resources are released properly
6. **Logging**: Tracks what the application is doing

## Next Steps

- Part 3: Learn about **Database & Migrations** - how data is persisted
- Part 4: Learn about **Authentication & Middleware** - securing requests

---

**Key Concept**: The entry point is like a recipe that says "First do this, then do that, then start cooking!" Everything is orchestrated here.
