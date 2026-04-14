# 09 - Error Handling & Utilities

## Utility Functions

The FEM project uses helper functions to avoid code repetition and maintain consistency.

## JSON Response Writer

Located in `internal/utils/`:

### WriteJson Function

```go
type Envelope map[string]interface{}

func WriteJson(w http.ResponseWriter, status int, data Envelope) error {
    // Step 1: Set Content-Type header
    w.Header().Set("Content-Type", "application/json; charset=utf-8")

    // Step 2: Write HTTP status code
    w.WriteHeader(status)

    // Step 3: Encode data as JSON and write to response
    return json.NewEncoder(w).Encode(data)
}
```

**Simple Explanation:**

- Sets headers to indicate JSON response
- Writes HTTP status code
- Converts Go data to JSON
- Sends to client

### Usage Examples

**Success Response:**

```go
utils.WriteJson(w, http.StatusOK,
    Envelope{"workout": workout})

// Output:
// HTTP/1.1 200 OK
// Content-Type: application/json; charset=utf-8
//
// {"workout":{...}}
```

**Error Response:**

```go
utils.WriteJson(w, http.StatusBadRequest,
    Envelope{"error": "Invalid Email!"})

// Output:
// HTTP/1.1 400 Bad Request
// Content-Type: application/json; charset=utf-8
//
// {"error":"Invalid Email!"}
```

**Created Response:**

```go
utils.WriteJson(w, http.StatusCreated,
    Envelope{"auth_token": token})

// Output:
// HTTP/1.1 201 Created
// Content-Type: application/json; charset=utf-8
//
// {"auth_token":"eyJ..."}
```

## ID Parameter Reading

### ReadIDParam Function

```go
func ReadIDParam(r *http.Request) (int64, error) {
    // Step 1: Extract ID from URL path
    id := chi.URLParam(r, "id")

    if id == "" {
        return 0, errors.New("id parameter is required")
    }

    // Step 2: Convert string to int64
    idInt, err := strconv.ParseInt(id, 10, 64)
    if err != nil {
        return 0, errors.New("id must be a valid integer")
    }

    return idInt, nil
}
```

**Simple Explanation:**

- Extracts "id" parameter from URL
- Validates it's not empty
- Converts string to int64
- Returns error if invalid

### Usage in Handlers

```go
func (h *WorkoutHandler) HandleGetWorkoutById(w http.ResponseWriter, r *http.Request) {
    // Extract and validate ID
    workoutId, err := utils.ReadIDParam(r)

    if err != nil {
        h.logger.Printf("Error: readIDParam : %v", err)
        utils.WriteJson(w, http.StatusBadRequest,
            Envelope{"error": "invalid get workout id"})
        return
    }

    // Use ID
    workout, err := h.store.GetWorkoutById(workoutId)
}
```

## HTTP Status Codes

### Common Status Codes Used

```go
http.StatusOK                  // 200 - Success
http.StatusCreated              // 201 - Resource created
http.StatusBadRequest           // 400 - Invalid request
http.StatusUnauthorized         // 401 - Need authentication
http.StatusForbidden            // 403 - Access denied
http.StatusNotFound             // 404 - Resource not found
http.StatusInternalServerError  // 500 - Server error
```

### When to Use Each

| Status                    | When                        | Example                           |
| ------------------------- | --------------------------- | --------------------------------- |
| 200 OK                    | Request successful          | GET workout successful            |
| 201 Created               | Resource created            | POST user successful              |
| 400 Bad Request           | Invalid input               | Missing required field            |
| 401 Unauthorized          | Not authenticated           | No/invalid token                  |
| 403 Forbidden             | Authenticated but no access | User can't delete other's workout |
| 404 Not Found             | Resource doesn't exist      | Workout ID doesn't exist          |
| 500 Internal Server Error | Server error                | Database connection failed        |

## Error Handling Patterns

### Pattern 1: Validation Error

```go
func (h *UserHandler) HandleRegisterUser(w http.ResponseWriter, r *http.Request) {
    // Step 1: Decode request
    var req registerUserRequest
    err := json.NewDecoder(r.Body).Decode(&req)
    if err != nil {
        h.logger.Printf("ERROR : decoding: %v", err)
        utils.WriteJson(w, http.StatusBadRequest,
            Envelope{"error": "Invalid request payload!"})
        return  // ← Important: stop here
    }

    // Step 2: Validate
    err = h.validateRegisterRequest(&req)
    if err != nil {
        h.logger.Printf("ERROR : validating: %v", err)
        utils.WriteJson(w, http.StatusBadRequest,
            Envelope{"error": err.Error()})
        return  // ← Stop here
    }

    // Step 3: Process
    err = h.store.CreateUser(user)
    if err != nil {
        h.logger.Printf("ERROR : creating user: %v", err)
        utils.WriteJson(w, http.StatusInternalServerError,
            Envelope{"error": "Internal Server Error"})
        return
    }

    // Step 4: Success
    utils.WriteJson(w, http.StatusOK,
        Envelope{"message": "User Created Successfully!", "user": user})
}
```

**Flow:**

```
Check inputs → Validate → Store → Return
     │            │        │       │
     ├─Error ─────×        │       │
     │                     ├─Error─×
     │                     │
     │                     └─Success → Return 200
```

### Pattern 2: Authorization Check

```go
func (h *WorkoutHandler) HandleDeleteWorkoutByID(w http.ResponseWriter, r *http.Request) {
    // Step 1: Get workout ID
    workoutId, err := utils.ReadIDParam(r)
    if err != nil {
        utils.WriteJson(w, http.StatusBadRequest, ...)
        return
    }

    // Step 2: Get current user
    currentUser := middleware.GetUser(r)
    if currentUser == nil || currentUser == store.AnonymousUser {
        utils.WriteJson(w, http.StatusUnauthorized,
            Envelope{"error": "You must be logged in!"})
        return
    }

    // Step 3: Check ownership
    workoutOwner, err := h.store.GetWorkoutOwner(workoutId)
    if err != nil {
        if errors.Is(err, sql.ErrNoRows) {
            utils.WriteJson(w, http.StatusNotFound, ...)
            return
        }
        utils.WriteJson(w, http.StatusInternalServerError, ...)
        return
    }

    if workoutOwner != currentUser.ID {
        utils.WriteJson(w, http.StatusForbidden,
            Envelope{"error": "Not authorized!"})
        return
    }

    // Step 4: Delete
    err = h.store.DeleteWorkout(workoutId)
    if err != nil {
        utils.WriteJson(w, http.StatusInternalServerError, ...)
        return
    }

    utils.WriteJson(w, http.StatusOK,
        Envelope{"message": "Deleted!"})
}
```

### Pattern 3: Database Error Handling

```go
func (h *TokenHandler) HandleCreateToken(w http.ResponseWriter, r *http.Request) {
    // ... decode request ...

    // Get user
    user, err := h.userStore.GetUserByUsername(req.Username)
    if err != nil {
        // SQL error occurred
        h.logger.Printf("ERROR: database error: %v", err)
        utils.WriteJson(w, http.StatusInternalServerError,
            Envelope{"error": "Internal Server Error"})
        return
    }

    if user == nil {
        // User doesn't exist (not an error, just not found)
        utils.WriteJson(w, http.StatusBadRequest,
            Envelope{"error": "User not found"})
        return
    }

    // Continue...
}
```

## Error Type Checking

### Check for Specific Database Errors

```go
import "database/sql"

func someHandler(w http.ResponseWriter, r *http.Request) {
    workout, err := store.GetWorkoutById(id)

    if err == sql.ErrNoRows {
        // Resource not found
        http.NotFound(w, r)
        return
    }

    if err != nil {
        // Some other database error
        utils.WriteJson(w, http.StatusInternalServerError,
            Envelope{"error": "Internal Server Error"})
        return
    }

    // Success
    utils.WriteJson(w, http.StatusOK, Envelope{"workout": workout})
}
```

### Check for Error Equality

```go
import "errors"

err := someOperation()

if errors.Is(err, sql.ErrNoRows) {
    // Handle "no rows" error
}

if err != nil && !errors.Is(err, sql.ErrNoRows) {
    // Handle other errors
}
```

## Logging

The application logs errors for debugging:

```go
h.logger.Printf("ERROR : decoding register request: %v", err)
h.logger.Printf("Error: getWorkoutByID: %v", err)
h.logger.Printf("ERROR: authenticating: Invalid Token")
```

**Log Output:**

```
2024/01/15 10:30:45 ERROR : decoding register request: EOF
2024/01/15 10:30:46 Error: getWorkoutByID: sql: no rows in result set
2024/01/15 10:30:47 App Started Running on PORT 8080!
```

## Error Response Format

### Consistent Format

All errors follow the same pattern:

```json
{
  "error": "Error message describing what went wrong"
}
```

### Examples

**Missing Field:**

```json
{ "error": "Username is required" }
```

**Invalid Format:**

```json
{ "error": "Invalid Email!" }
```

**Authentication Failed:**

```json
{ "error": "Invalid Credentials" }
```

**Not Authorized:**

```json
{ "error": "You must be loggedin to access the route!" }
```

**Not Found:**

```json
{ "error": "unable to get the workout owner" }
```

## Common Error Scenarios

### 1. Decode JSON Error

```go
err := json.NewDecoder(r.Body).Decode(&req)
if err != nil {
    utils.WriteJson(w, http.StatusBadRequest,
        Envelope{"error": "Invalid JSON format"})
    return
}
```

**Response:**

```http
HTTP/1.1 400 Bad Request
{"error": "Invalid JSON format"}
```

### 2. Missing Required Field

```go
if req.Username == "" {
    utils.WriteJson(w, http.StatusBadRequest,
        Envelope{"error": "Username is required"})
    return
}
```

**Response:**

```http
HTTP/1.1 400 Bad Request
{"error": "Username is required"}
```

### 3. Database Query Failed

```go
workout, err := h.store.GetWorkoutById(workoutId)
if err != nil {
    h.logger.Printf("Error: database: %v", err)
    utils.WriteJson(w, http.StatusInternalServerError,
        Envelope{"error": "internal server error"})
    return
}
```

**Response:**

```http
HTTP/1.1 500 Internal Server Error
{"error": "internal server error"}
```

### 4. Resource Not Found

```go
if existingWorkout == nil {
    http.NotFound(w, r)
    return
}

// OR

utils.WriteJson(w, http.StatusNotFound,
    Envelope{"error": "Workout not found"})
```

**Response:**

```http
HTTP/1.1 404 Not Found
{"error": "Workout not found"}
```

### 5. Not Authenticated

```go
user := middleware.GetUser(r)
if user.IsAnonymous() {
    utils.WriteJson(w, http.StatusUnauthorized,
        Envelope{"error": "You must be loggedin!"})
    return
}
```

**Response:**

```http
HTTP/1.1 401 Unauthorized
{"error": "You must be loggedin!"}
```

### 6. Not Authorized (Don't Own Resource)

```go
if workoutOwner != currentUser.ID {
    utils.WriteJson(w, http.StatusForbidden,
        Envelope{"error": "Not authorized!"})
    return
}
```

**Response:**

```http
HTTP/1.1 403 Forbidden
{"error": "Not authorized!"}
```

## Error Handling Best Practices

1. **Validate Early**: Check inputs before processing
2. **Log Errors**: Always log for debugging
3. **Return Appropriate Status**: 400 vs 401 vs 403 vs 500
4. **Don't Expose Internals**: Generic message to client, details in logs
5. **Stop on Error**: Use `return` to prevent further processing
6. **Consistency**: All errors same format

## Examples of Good Error Handling

### ✓ Good

```go
// Validates, logs, returns specific status
err := json.NewDecoder(r.Body).Decode(&req)
if err != nil {
    h.logger.Printf("ERROR: decode: %v", err)  // Log detail
    utils.WriteJson(w, http.StatusBadRequest,   // Specific status
        Envelope{"error": "Invalid request"})   // Generic message
    return
}
```

### ✗ Bad

```go
// Doesn't log, generic status, leaks details
err := json.NewDecoder(r.Body).Decode(&req)
if err != nil {
    utils.WriteJson(w, http.StatusInternalServerError,
        Envelope{"error": err.Error()})  // Exposes internals!
}
// No return - continues executing!
```

## Key Takeaways

1. **Envelope**: Consistent response wrapper
2. **WriteJson**: Helper to send JSON responses
3. **ReadIDParam**: Helper to extract/validate URL parameters
4. **Status Codes**: Use the right code for the situation
5. **Error Logging**: Always log for debugging
6. **Validation**: Check inputs early
7. **Stop on Error**: Use return to prevent further processing

## Next Steps

- Part 10: Learn about **Integration & Complete Flow** - how all parts work together

---

**Key Concept**: Handle errors gracefully. Log details for debugging, return generic messages to client, use correct HTTP status codes!
