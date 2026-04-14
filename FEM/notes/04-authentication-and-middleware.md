# 04 - Authentication & Middleware

## What is Authentication?

**Authentication** is the process of verifying that a user is who they claim to be. In the FEM project, users prove their identity using username and password, then receive a **token** to use for future requests.

## What is Middleware?

**Middleware** is code that runs before or after a request is handled. In the FEM project, middleware:

1. Checks for authentication tokens
2. Validates tokens
3. Extracts user information from tokens
4. Prevents unauthorized access

## The Authentication Flow

```
┌─────────────────────────┐
│   User Registers        │
│   (username/password)   │
└────────────┬────────────┘
             │
             ▼
┌─────────────────────────┐
│   User Logs In          │
│   (sends credentials)   │
└────────────┬────────────┘
             │
             ▼
┌─────────────────────────┐
│   Server Validates      │
│   (checks password)     │
└────────────┬────────────┘
             │
             ▼
┌─────────────────────────┐
│   Generate Token        │
│   (valid for 24 hours)  │
└────────────┬────────────┘
             │
             ▼
┌─────────────────────────┐
│   Return Token          │
│   (send to client)      │
└────────────┬────────────┘
             │
             ▼
┌─────────────────────────┐
│   User Makes Request    │
│   (includes token)      │
└────────────┬────────────┘
             │
             ▼
┌─────────────────────────┐
│   Middleware Validates  │
│   (checks token)        │
└────────────┬────────────┘
             │
             ▼
┌─────────────────────────┐
│   Access Granted        │
│   (request proceeds)    │
└─────────────────────────┘
```

## Middleware Code

Located in `internal/middleware/middleware.go`:

### Context Setup

```go
type contextKey string
const UserContextKey = contextKey("user")

func SetUser(r *http.Request, user *store.User) *http.Request {
    ctx := context.WithValue(r.Context(), UserContextKey, user)
    return r.WithContext(ctx)
}

func GetUser(r *http.Request) *store.User {
    user, ok := r.Context().Value(UserContextKey).(*store.User)
    if !ok {
        panic("missing user in request!")
    }
    return user
}
```

**Simple Explanation:**

- **Context**: A way to store data associated with a specific request
- **SetUser**: Stores user info in the request context
- **GetUser**: Retrieves user info from the request context
- Each request gets its own context with user info

### Authenticate Middleware

```go
func (um *UserMiddleware) Authenticate(next http.Handler) http.Handler {
    return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
        // Set Vary header for caching
        w.Header().Add("Vary", "Authorization")

        // Get Authorization header
        authHeader := r.Header.Get("Authorization")

        // If no token, set anonymous user
        if authHeader == "" {
            r = SetUser(r, store.AnonymousUser)
            next.ServeHTTP(w, r)
            return
        }

        // Parse the Authorization header
        headerParts := strings.Split(authHeader, " ") // ["Bearer", "TOKEN"]

        if len(headerParts) != 2 || headerParts[0] != "Bearer" {
            utils.WriteJson(w, http.StatusUnauthorized,
                utils.Envelope{"error": "Invalid Authorization Header"})
            return
        }

        token := headerParts[1]

        // Get user associated with this token
        user, err := um.UserStore.GetUserToken(tokens.ScopeAuth, token)
        if err != nil {
            utils.WriteJson(w, http.StatusUnauthorized,
                utils.Envelope{"error": "Invalid Token"})
            return
        }

        if user == nil {
            utils.WriteJson(w, http.StatusUnauthorized,
                utils.Envelope{"error": "Invalid User of Token"})
            return
        }

        // Set user in context and continue
        SetUser(r, user)
        next.ServeHTTP(w, r)
    })
}
```

**Step by Step:**

1. **Check Authorization Header**

   ```
   Request header: "Authorization: Bearer <token>"
   ```

2. **Parse Header**

   ```
   Split by space: ["Bearer", "<token>"]
   Verify format: Must have 2 parts and start with "Bearer"
   ```

3. **Validate Token**

   ```
   Query database for token
   Get associated user
   Verify user exists
   ```

4. **Store User in Context**
   ```
   Next handlers can access user via GetUser(r)
   ```

### RequireUser Middleware

```go
func (um *UserMiddleware) RequireUser(next http.HandlerFunc) http.HandlerFunc {
    return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
        user := GetUser(r)

        // Check if user is anonymous
        if user.IsAnonymous() {
            utils.WriteJson(w, http.StatusUnauthorized,
                utils.Envelope{"error": "You must be loggedin to access the route!"})
            return
        }

        // User is authenticated, continue
        next.ServeHTTP(w, r)
    })
}
```

**Simple Explanation:**

- **RequireUser**: Ensures only authenticated users can access this route
- If user is anonymous → Return 401 Unauthorized
- If user is authenticated → Continue to handler

## Middleware in Routes

In `internal/routes/routes.go`:

```go
func SetupRoutes(app *app.Application) *chi.Mux {
    r := chi.NewRouter()

    // Apply Authenticate middleware to all routes
    r.Group(func(r chi.Router) {
        r.Use(app.Middleware.Authenticate)

        // These routes require authentication
        r.Get("/workouts/{id}",
            app.Middleware.RequireUser(app.WorkoutHandler.HandleGetWorkoutById))
        r.Post("/workouts",
            app.Middleware.RequireUser(app.WorkoutHandler.HandleCreateWorkout))
        r.Put("/workouts/{id}",
            app.Middleware.RequireUser(app.WorkoutHandler.HandleUpdateWorkoutByID))
        r.Delete("/workouts/{id}",
            app.Middleware.RequireUser(app.WorkoutHandler.HandleDeleteWorkoutByID))
    })

    // Public routes (no authentication required)
    r.Get("/health", app.HealthCheck)
    r.Post("/users", app.UserHandler.HandleRegisterUser)
    r.Post("/tokens/authentication", app.TokenHandler.HandleCreateToken)

    return r
}
```

**Middleware Application:**

```
Request comes in
    │
    ▼
Authenticate middleware (runs for ALL requests)
    ├─ If no token ──────────────► r = SetUser(r, AnonymousUser)
    └─ If token exists ──────────► Validate & extract user
    │
    ▼
Route handler
    ├─ Public routes ────────────► Accessed by anyone
    └─ Protected routes
        ├─ RequireUser checks
        │   ├─ If anonymous ─────► Return 401
        │   └─ If authenticated ─► Proceed
        │
        ▼
    Handler executes
```

## Request Flow with Authentication

### Example 1: Create a Workout (Protected Route)

**Request:**

```http
POST /workouts HTTP/1.1
Authorization: Bearer eyJhbGc...
Content-Type: application/json

{
  "title": "Morning Run",
  "duration_minutes": 30
}
```

**Processing:**

```
1. Authenticate middleware runs
   ├─ Reads "Authorization: Bearer eyJhbGc..."
   ├─ Extracts token: "eyJhbGc..."
   ├─ Queries database
   ├─ Finds user with this token
   └─ Sets user in context

2. RequireUser middleware runs
   ├─ Gets user from context
   ├─ Checks if anonymous
   └─ Allows request to continue

3. Handler executes
   ├─ Gets user from context: currentUser = middleware.GetUser(r)
   ├─ Creates workout with currentUser.ID
   └─ Returns created workout

Response: 200 OK with workout data
```

### Example 2: Get Health (Public Route)

**Request:**

```http
GET /health HTTP/1.1
```

**Processing:**

```
1. Authenticate middleware runs
   ├─ No Authorization header
   └─ Sets AnonymousUser in context

2. Handler executes directly
   ├─ No RequireUser check
   └─ Returns "Status is Available!"

Response: 200 OK
```

### Example 3: Try Protected Route Without Token

**Request:**

```http
GET /workouts/1 HTTP/1.1
```

**Processing:**

```
1. Authenticate middleware runs
   ├─ No Authorization header
   └─ Sets AnonymousUser in context

2. RequireUser middleware runs
   ├─ Gets user from context (AnonymousUser)
   ├─ Checks IsAnonymous() ──► TRUE
   └─ Returns 401 Unauthorized

Response: 401 {"error": "You must be loggedin..."}
```

## User Context

The `store.User` struct is stored in context:

```go
type User struct {
    ID           int
    Username     string
    Email        string
    PasswordHash PasswordHash
    Bio          string
    CreatedAt    time.Time
    UpdatedAt    time.Time
}

var AnonymousUser = &User{} // User with all zero values

func (u *User) IsAnonymous() bool {
    return u.ID == 0
}
```

**Checking User:**

```go
func (h *WorkoutHandler) HandleCreateWorkout(w http.ResponseWriter, r *http.Request) {
    var workout store.Workout
    json.NewDecoder(r.Body).Decode(&workout)

    // Get current user from context
    currentUser := middleware.GetUser(r)

    // Set workout's user ID
    workout.UserID = currentUser.ID

    // Save workout
    createdWorkout, h.store.CreateWorkout(&workout)

    // ... return response
}
```

## Token System

Tokens are created with:

- **User ID**: Who the token belongs to
- **Expiration**: When token expires (24 hours from creation)
- **Scope**: What the token is for (auth, password-reset, etc.)

```go
token, err := h.tokenStore.CreateNewToken(
    user.ID,           // User ID
    time.Hour*24,      // Valid for 24 hours
    tokens.ScopeAuth,  // Scope: authentication
)
```

## Security Considerations

1. **Tokens in Headers**: Not in URL (URLs logged, headers not)
2. **Bearer Prefix**: Standard format for tokens
3. **Expiration**: Tokens expire to limit damage if stolen
4. **Password Hashing**: Passwords hashed with bcrypt
5. **Context**: User data passed safely through request lifecycle

## Key Takeaways

1. **Middleware runs in order** before handlers
2. **Authenticate** runs for all requests
3. **RequireUser** enforces authentication
4. **Context** stores user data for the request
5. **Tokens** enable stateless authentication
6. **Anonymous user** for unauthenticated requests

## Next Steps

- Part 5: Learn about **User Management** - registration and user operations
- Part 6: Learn about **Workout Operations** - CRUD operations

---

**Key Concept**: Middleware is like a security checkpoint. Authenticate checks your credentials, RequireUser checks it's not fake, then you're allowed through!
