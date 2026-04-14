# 08 - Routing System

## What is Routing?

**Routing** is the process of matching incoming HTTP requests to handler functions based on the URL path and HTTP method.

```
Request comes in
    │
    ▼
┌─────────────────────────────────────┐
│   Routing System                    │
│                                     │
│   Match URL path → Find handler     │
│   Match HTTP method → Find function │
└─────────────────────────────────────┘
    │
    ▼
Call appropriate handler
    │
    ▼
Send response
```

## Chi Router

The FEM project uses **Chi**, a lightweight HTTP routing library for Go.

**Installation:**

```go
import "github.com/go-chi/chi/v5"
```

## Routes Setup

Located in `internal/routes/routes.go`:

```go
func SetupRoutes(app *app.Application) *chi.Mux {
    r := chi.NewRouter()

    // Protected routes group
    r.Group(func(r chi.Router) {
        r.Use(app.Middleware.Authenticate)

        r.Get("/workouts/{id}",
            app.Middleware.RequireUser(app.WorkoutHandler.HandleGetWorkoutById))
        r.Post("/workouts",
            app.Middleware.RequireUser(app.WorkoutHandler.HandleCreateWorkout))
        r.Put("/workouts/{id}",
            app.Middleware.RequireUser(app.WorkoutHandler.HandleUpdateWorkoutByID))
        r.Delete("/workouts/{id}",
            app.Middleware.RequireUser(app.WorkoutHandler.HandleDeleteWorkoutByID))
    })

    // Public routes
    r.Get("/health", app.HealthCheck)
    r.Post("/users", app.UserHandler.HandleRegisterUser)
    r.Post("/tokens/authentication", app.TokenHandler.HandleCreateToken)

    return r
}
```

## Route Groups

### Concept

Route groups allow applying middleware to multiple routes:

```go
r.Group(func(r chi.Router) {
    // Apply middleware to all routes in this group
    r.Use(app.Middleware.Authenticate)

    // All routes here have Authenticate middleware
    r.Get("/workouts/{id}", ...)
    r.Post("/workouts", ...)
    r.Put("/workouts/{id}", ...)
    r.Delete("/workouts/{id}", ...)
})
```

**Benefits:**

- Organize related routes
- Avoid repeating middleware declarations
- Clear separation of concerns

### Middleware Application

```
Request arrives
    │
    ▼
┌──────────────────────────────────────┐
│   Check if in Group with middleware  │
│   Apply Group middleware             │
└──────────────────────────────────────┘
    │
    ▼
┌──────────────────────────────────────┐
│   Route Matching                     │
│   Find matching HTTP method + path   │
└──────────────────────────────────────┘
    │
    ▼
┌──────────────────────────────────────┐
│   Check route-specific middleware    │
│   Apply route middleware             │
└──────────────────────────────────────┘
    │
    ▼
Call handler function
```

## HTTP Methods and Routes

### GET - Retrieve Data

```go
r.Get("/workouts/{id}", app.Middleware.RequireUser(...))
```

**Request:**

```http
GET /workouts/1 HTTP/1.1
Authorization: Bearer token...
```

**Response:**

```http
HTTP/1.1 200 OK
{workout data}
```

### POST - Create Data

```go
r.Post("/workouts", app.Middleware.RequireUser(...))
```

**Request:**

```http
POST /workouts HTTP/1.1
Authorization: Bearer token...
Content-Type: application/json

{workout data}
```

**Response:**

```http
HTTP/1.1 200 OK
{created workout}
```

### PUT - Replace/Update Data

```go
r.Put("/workouts/{id}", app.Middleware.RequireUser(...))
```

**Request:**

```http
PUT /workouts/1 HTTP/1.1
Authorization: Bearer token...
Content-Type: application/json

{updated fields}
```

**Response:**

```http
HTTP/1.1 200 OK
{success message}
```

### DELETE - Remove Data

```go
r.Delete("/workouts/{id}", app.Middleware.RequireUser(...))
```

**Request:**

```http
DELETE /workouts/1 HTTP/1.1
Authorization: Bearer token...
```

**Response:**

```http
HTTP/1.1 200 OK
{success message}
```

## URL Parameters

### Path Parameters

Extracted from URL path using `{name}` syntax:

```go
r.Get("/workouts/{id}", handler)
```

**Request:**

```
GET /workouts/123
```

**In Handler:**

```go
func (h *WorkoutHandler) HandleGetWorkoutById(w http.ResponseWriter, r *http.Request) {
    // Extract ID from URL
    workoutId, err := utils.ReadIDParam(r)  // Gets "123"

    // Use ID to query database
    workout, err := h.store.GetWorkoutById(workoutId)
}
```

### Example URLs

| Route            | Example URL     | ID  |
| ---------------- | --------------- | --- |
| `/workouts/{id}` | `/workouts/1`   | 1   |
| `/workouts/{id}` | `/workouts/999` | 999 |
| `/users/{id}`    | `/users/42`     | 42  |

## Route Matching Examples

```go
// Exact match
GET /health
├─ /health          ✓
├─ /health/check    ✗

// Path parameter
GET /workouts/{id}
├─ /workouts/1      ✓
├─ /workouts/abc    ✓ (but will fail validation)
├─ /workouts        ✗

// Multiple parameters (not used in FEM)
GET /users/{userId}/workouts/{workoutId}
├─ /users/1/workouts/5  ✓
├─ /users/1/workouts    ✗
```

## Complete Route Tree

```
Router (chi.Mux)
├── [Group with Authenticate middleware]
│   ├── GET    /workouts/{id}    [RequireUser] → HandleGetWorkoutById
│   ├── POST   /workouts         [RequireUser] → HandleCreateWorkout
│   ├── PUT    /workouts/{id}    [RequireUser] → HandleUpdateWorkoutByID
│   └── DELETE /workouts/{id}    [RequireUser] → HandleDeleteWorkoutByID
│
├── GET    /health                            → HealthCheck
├── POST   /users                             → HandleRegisterUser
└── POST   /tokens/authentication             → HandleCreateToken
```

## Request Processing Flow

### 1. GET /workouts/1 (Protected Route)

```
GET /workouts/1
Authorization: Bearer TOKEN
    │
    ▼
Route matching:
    ├─ Method: GET ✓
    └─ Path: /workouts/{id} ✓
    │
    ▼
Apply Authenticate middleware
    ├─ Parse Authorization header
    ├─ Validate token
    └─ Extract user → Set in context
    │
    ▼
Apply RequireUser middleware
    ├─ Check if user is anonymous
    ├─ Is anonymous? → Return 401
    └─ Not anonymous? → Continue
    │
    ▼
Call HandleGetWorkoutById
    ├─ Extract ID from URL: "1"
    ├─ Query database: GetWorkoutById(1)
    └─ Return workout
    │
    ▼
Response: 200 OK + workout data
```

### 2. GET /health (Public Route)

```
GET /health
    │
    ▼
Route matching:
    ├─ Method: GET ✓
    └─ Path: /health ✓
    │
    ▼
No middleware applied
    │
    ▼
Call HealthCheck
    └─ Return "Status is Available!"
    │
    ▼
Response: 200 OK + status text
```

### 3. POST /workouts Without Token (Protected Route)

```
POST /workouts
Content-Type: application/json
{workout data}
    │
    ▼
Route matching:
    ├─ Method: POST ✓
    └─ Path: /workouts ✓
    │
    ▼
Apply Authenticate middleware
    ├─ Check Authorization header
    └─ No header found → Set AnonymousUser
    │
    ▼
Apply RequireUser middleware
    ├─ Check if user is anonymous
    ├─ Is anonymous? → Return 401
    └─ STOP HERE
    │
    ▼
Response: 401 Unauthorized
```

## Middleware Chain

Middlewares execute in order:

```
Request
    │
    ▼
middleware1
    │
    ▼
middleware2
    │
    ▼
middleware3
    │
    ▼
handler function
    │
    ▼
Response (flows back through middlewares)
```

**Example in FEM:**

```go
r.Group(func(r chi.Router) {
    r.Use(app.Middleware.Authenticate)  // First

    r.Get("/workouts/{id}",
        app.Middleware.RequireUser(...))  // Second
})
```

**Execution order:**

1. Authenticate middleware
2. RequireUser middleware
3. Handler function

## Handler Signature

All handlers follow this signature:

```go
func(w http.ResponseWriter, r *http.Request)
```

**Parameters:**

- `w`: Write responses to client
- `r`: Read request data

**Example:**

```go
func (h *WorkoutHandler) HandleGetWorkoutById(w http.ResponseWriter, r *http.Request) {
    // w - used to send response
    // r - used to read request
}
```

## Query Parameters (Not Used in FEM)

Chi also supports query parameters:

```go
r.Get("/workouts", handler)

// Request: GET /workouts?limit=10&offset=0
// In handler:
limit := r.URL.Query().Get("limit")      // "10"
offset := r.URL.Query().Get("offset")    // "0"
```

## Method Not Allowed

If request method doesn't match:

```
GET /workouts

But only route is:
POST /workouts

Response: 405 Method Not Allowed
```

## Not Found

If path doesn't match any route:

```
GET /invalid-path

Response: 404 Not Found
```

## Testing Routes

View all registered routes:

```go
chi.Walk(router, func(method string, route string, handler http.Handler, middlewares ...func(http.Handler) http.Handler) error {
    fmt.Printf("%s %s\n", method, route)
    return nil
})
```

**Output:**

```
GET /workouts/{id}
POST /workouts
PUT /workouts/{id}
DELETE /workouts/{id}
GET /health
POST /users
POST /tokens/authentication
```

## Best Practices

1. **Grouping**: Use groups for shared middleware
2. **Naming**: Consistent URL patterns meaningful
3. **Methods**: Use correct HTTP method (GET for read, POST for create, etc.)
4. **Parameters**: Extract from URL, validate in handler
5. **Middleware**: Apply at appropriate level (global vs route)

## Key Takeaways

1. **Routes map URLs to handlers**
2. **Route groups apply middleware to multiple routes**
3. **Path parameters extracted from URL using {name}**
4. **Method determines the action (GET, POST, PUT, DELETE)**
5. **Middleware runs before handlers**
6. **Chi is lightweight and intuitive**

## Next Steps

- Part 9: Learn about **Error Handling & Utils** - helper functions and utilities
- Part 10: Learn about **Integration & Complete Flow** - how everything connects

---

**Key Concept**: Routing is matching. URL path + HTTP method → Handler function. Middleware runs before the handler!
