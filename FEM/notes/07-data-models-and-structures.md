# 07 - Data Models & Structures

## Introduction to Structs

A **struct** is a collection of fields grouped together. It's like a blueprint for data. The FEM project uses structs to represent:

- Database records (Workout, User, etc.)
- HTTP requests and responses
- Domain objects

## User Struct

Located in `internal/store/`:

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
```

**Database Mapping:**

```go
type User struct {
    ID           int            // users.id
    Username     string         // users.username
    Email        string         // users.email
    PasswordHash PasswordHash   // users.password_hash
    Bio          string         // users.bio
    CreatedAt    time.Time      // users.created_at
    UpdatedAt    time.Time      // users.updated_at
}
```

**JSON Tags:**
When converted to JSON, fields map like this:

```go
type User struct {
    ID           int       `json:"id"`
    Username     string    `json:"username"`
    Email        string    `json:"email"`
    // PasswordHash intentionally omitted from JSON!
    Bio          string    `json:"bio"`
    CreatedAt    time.Time `json:"created_at"`
    UpdatedAt    time.Time `json:"updated_at"`
}
```

**Example JSON:**

```json
{
  "id": 1,
  "username": "john_doe",
  "email": "john@example.com",
  "bio": "Fitness enthusiast",
  "created_at": "2024-01-15T10:30:00Z",
  "updated_at": "2024-01-15T10:30:00Z"
}
```

## Workout Struct

```go
type Workout struct {
    ID              int            `json:"id"`
    UserID          int            `json:"user_id"`
    Title           string         `json:"title"`
    Description     string         `json:"description"`
    DurationMinutes int            `json:"duration_minutes"`
    CaloriesBurned  int            `json:"calories_burned"`
    CreatedAt       time.Time      `json:"created_at"`
    UpdatedAt       time.Time      `json:"updated_at"`
    Entries         []WorkoutEntry `json:"entries"`
}
```

**Example:**

```json
{
  "id": 1,
  "user_id": 1,
  "title": "Morning Run",
  "description": "5K run in the park",
  "duration_minutes": 30,
  "calories_burned": 300,
  "created_at": "2024-01-15T10:30:00Z",
  "updated_at": "2024-01-15T10:30:00Z",
  "entries": [
    {
      "id": 1,
      "workout_id": 1,
      "exercise_name": "Running",
      "sets": 1,
      "reps": null,
      "duration_seconds": 1800,
      "weight": null,
      "notes": "Good pace",
      "order_index": 1
    }
  ]
}
```

## WorkoutEntry Struct

```go
type WorkoutEntry struct {
    ID              int      `json:"id"`
    WorkoutId       int      `json:"workout_id"`
    ExerciseName    string   `json:"exercise_name"`
    Sets            int      `json:"sets"`
    Reps            *int     `json:"reps"`              // Pointer!
    DurationSeconds *int     `json:"duration_seconds"` // Pointer!
    Weight          *float64 `json:"weight"`           // Pointer!
    Notes           string   `json:"notes"`
    OrderIndex      int      `json:"order_index"`
}
```

**Key Feature: Pointer Fields**

Notice some fields are `*int` and `*float64` (pointers) instead of `int` and `float64`:

```go
// Regular field
Reps int           // Always has a value (0 if not set)

// Pointer field
Reps *int          // Can be nil (null) or point to a value
```

**Why Pointers?**
In the check constraint: `(reps IS NOT NULL OR duration_seconds IS NOT NULL)`

We need to distinguish between:

- `nil` (null in database)
- `0` (actual value)

**Example: Strength Exercise vs Cardio**

```go
// Strength exercise (has reps, no duration)
entry1 := WorkoutEntry{
    ExerciseName: "Squats",
    Reps:         intPtr(12),        // Points to 12
    DurationSeconds: nil,            // null
    // ✓ Valid: one field has value, other is nil
}

// Cardio exercise (has duration, no reps)
entry2 := WorkoutEntry{
    ExerciseName: "Running",
    Reps:         nil,               // null
    DurationSeconds: intPtr(900),    // Points to 900
    // ✓ Valid: one field has value, other is nil
}

func intPtr(i int) *int { return &i }
```

**JSON Representation:**

Strength exercise:

```json
{
  "exercise_name": "Squats",
  "reps": 12,
  "duration_seconds": null,
  "weight": 100.0,
  "sets": 3
}
```

Cardio exercise:

```json
{
  "exercise_name": "Running",
  "reps": null,
  "duration_seconds": 900,
  "weight": null,
  "sets": 1
}
```

## Request Structs

Used to parse HTTP request bodies:

### User Registration Request

```go
type registerUserRequest struct {
    Username string `json:"username"`
    Email    string `json:"email"`
    Password string `json:"password"`
    Bio      string `json:"bio"`
}
```

**Decoding:**

```go
var req registerUserRequest
json.NewDecoder(r.Body).Decode(&req)
```

**Input JSON:**

```json
{
  "username": "john_doe",
  "email": "john@example.com",
  "password": "SecurePass123!",
  "bio": "I love fitness"
}
```

### Token Request

```go
type createTokenRequest struct {
    Username string `json:"username"`
    Password string `json:"password"`
}
```

**Input JSON:**

```json
{
  "username": "john_doe",
  "password": "SecurePass123!"
}
```

### Workout Update Request

```go
type updateWorkOutRequest struct {
    Title           *string              `json:"title"`
    Description     *string              `json:"description"`
    DurationMinutes *int                 `json:"duration_minutes"`
    CaloriesBurned  *int                 `json:"calories_burned"`
    Entries         []store.WorkoutEntry `json:"entries"`
}
```

**Why pointers?** To make all fields optional:

```json
{
  "title": "New Title"
  // Other fields omitted - they won't be updated!
}
```

If fields weren't pointers, we'd lose the difference between:

- "User didn't provide this field" (should preserve old value)
- "User provided zero value" (should update to zero)

## Response Structures

### Envelope Pattern

The API uses an `Envelope` for consistent responses:

```go
type Envelope map[string]interface{}
```

**Usage:**

```go
utils.WriteJson(w, http.StatusOK,
    Envelope{"workout": workout})

utils.WriteJson(w, http.StatusOK,
    Envelope{"message": "Success!", "user": user})

utils.WriteJson(w, http.StatusBadRequest,
    Envelope{"error": "Invalid request"})
```

**Success Response:**

```json
{
  "workout": {
    "id": 1,
    "title": "Morning Run",
    "duration_minutes": 30,
    "entries": [...]
  }
}
```

**Error Response:**

```json
{
  "error": "Invalid Email!"
}
```

## Struct Tags

### JSON Tags

```go
type Workout struct {
    // Maps to "id" in JSON
    ID int `json:"id"`

    // Maps to "title" in JSON
    Title string `json:"title"`

    // Omitted from JSON entirely
    InternalField string `json:"-"`
}
```

**Common Tags:**

```go
type Example struct {
    // Regular JSON mapping
    Name string `json:"name"`

    // JSON only if not zero value
    Age int `json:"age,omitempty"`

    // Flatten nested struct into parent
    Address Address `json:"address,inline"`
}
```

## Type Conversion

### Scans from Database

```go
// When reading from database
var user User
rows.Scan(
    &user.ID,           // int from BIGINT
    &user.Username,     // string from VARCHAR
    &user.Email,        // string from VARCHAR
    &user.PasswordHash, // PasswordHash from VARCHAR
    &user.Bio,          // string from TEXT
    &user.CreatedAt,    // time.Time from TIMESTAMP
    &user.UpdatedAt,    // time.Time from TIMESTAMP
)
```

### Query Parameters

```go
type WorkoutEntry struct {
    OrderIndex int `json:"order_index"`
}

// User sends: {"order_index": 1}
// Go converts: 1 (string) → 1 (int)
// Database stores as: 1 (INTEGER)
```

## Nullable Fields in Database

### Option 1: Pointers (Used in FEM)

```go
type WorkoutEntry struct {
    Reps *int  // Can be null
}

// In database
reps INT               -- can be NULL
```

**Go Code:**

```go
// Null value
entry.Reps = nil

// Has value
reps := 12
entry.Reps = &reps

// Check if present
if entry.Reps != nil {
    fmt.Println("Reps:", *entry.Reps)
}
```

### Option 2: sql.NullInt64

```go
type WorkoutEntry struct {
    Reps sql.NullInt64
}

// Check if present
if entry.Reps.Valid {
    fmt.Println("Reps:", entry.Reps.Int64)
}
```

## Nested Structures

```go
type Workout struct {
    ID      int            `json:"id"`
    Title   string         `json:"title"`
    // Nested slice of entries
    Entries []WorkoutEntry `json:"entries"`
}
```

**JSON Output:**

```json
{
  "id": 1,
  "title": "Morning Workout",
  "entries": [
    {
      "id": 101,
      "exercise_name": "Squats",
      "sets": 3
    },
    {
      "id": 102,
      "exercise_name": "Bench Press",
      "sets": 4
    }
  ]
}
```

## Table: Data Type Mappings

| Go Type     | JSON        | PostgreSQL | Example          |
| ----------- | ----------- | ---------- | ---------------- |
| `int`       | number      | INT        | 42               |
| `string`    | string      | VARCHAR    | "hello"          |
| `*int`      | number/null | INT        | null or 42       |
| `bool`      | boolean     | BOOLEAN    | true             |
| `time.Time` | string      | TIMESTAMP  | "2024-01-15T..." |
| `[]struct`  | array       | —          | [{...}, {...}]   |
| `float64`   | number      | DECIMAL    | 100.50           |
| `*float64`  | number/null | DECIMAL    | null or 100.50   |

## Key Takeaways

1. **Struct Tags**: Map Go fields to JSON/database names
2. **Pointers**: Enable nil/null values for optional fields
3. **Envelope**: Consistent response wrapper
4. **Nesting**: Support complex hierarchical data
5. **Type Safety**: Database types match Go types
6. **Validation**: Request structs validate input

## Next Steps

- Part 8: Learn about **Routing System** - how requests are mapped to handlers
- Part 9: Learn about **Error Handling & Utils** - helper functions and error patterns

---

**Key Concept**: Structs are blueprints. JSON tags tell Go how to convert between structs and JSON. Pointers allow optional fields!
