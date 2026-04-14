# 07 - Data Models & Structures

## Understanding Data Models

A **data model** is a structure that represents how data is organized. In Go, we use `struct` types to define these models. The FEM project has several key models.

## Workout Structure

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

### Field Explanations

| Field             | Type           | Purpose               | Example              |
| ----------------- | -------------- | --------------------- | -------------------- |
| `ID`              | int            | Unique identifier     | 1                    |
| `UserID`          | int            | Owner's user ID       | 42                   |
| `Title`           | string         | Workout name          | "Morning Run"        |
| `Description`     | string         | Details about workout | "5km run in park"    |
| `DurationMinutes` | int            | How long (minutes)    | 30                   |
| `CaloriesBurned`  | int            | Energy expended       | 300                  |
| `CreatedAt`       | time.Time      | When created          | 2024-01-15T10:30:00Z |
| `UpdatedAt`       | time.Time      | Last modification     | 2024-01-15T10:30:00Z |
| `Entries`         | []WorkoutEntry | List of exercises     | [...]                |

### JSON Tags

```go
`json:"title"`
```

**What it means:**

- Maps Go field to JSON field name
- When converting to/from JSON, Go uses these names
- Example: `Title` (Go) → `"title"` (JSON)

### Example JSON

```json
{
  "id": 1,
  "user_id": 42,
  "title": "Morning Strength Training",
  "description": "Full body strength training",
  "duration_minutes": 60,
  "calories_burned": 500,
  "created_at": "2024-01-15T10:30:00Z",
  "updated_at": "2024-01-15T10:30:00Z",
  "entries": [...]
}
```

## WorkoutEntry Structure

```go
type WorkoutEntry struct {
    ID              int      `json:"id"`
    WorkoutId       int      `json:"workout_id"`
    ExerciseName    string   `json:"exercise_name"`
    Sets            int      `json:"sets"`
    Reps            *int     `json:"reps"`
    DurationSeconds *int     `json:"duration_seconds"`
    Weight          *float64 `json:"weight"`
    Notes           string   `json:"notes"`
    OrderIndex      int      `json:"order_index"`
}
```

### Field Explanations

| Field             | Type      | Purpose                  | Example       |
| ----------------- | --------- | ------------------------ | ------------- |
| `ID`              | int       | Unique identifier        | 1             |
| `WorkoutId`       | int       | Which workout            | 5             |
| `ExerciseName`    | string    | Exercise name            | "Squats"      |
| `Sets`            | int       | Number of sets           | 3             |
| `Reps`            | \*int     | Repetitions (optional)   | 12            |
| `DurationSeconds` | \*int     | Duration (optional)      | 300           |
| `Weight`          | \*float64 | Weight lifted (optional) | 100.5         |
| `Notes`           | string    | Notes about exercise     | "Felt strong" |
| `OrderIndex`      | int       | Position in workout      | 1             |

### Pointer Types

Notice `Reps`, `DurationSeconds`, and `Weight` are **pointers** (indicated by `*`):

```go
Reps            *int        // Pointer to int (can be nil)
DurationSeconds *int        // Pointer to int (can be nil)
Weight          *float64    // Pointer to float64 (can be nil)
```

**Why pointers?**

Some exercises are strength-based (use reps), others are cardio (use duration). We need to allow these to be optional.

**Examples:**

```go
// Strength exercise
entry1 := WorkoutEntry{
    ExerciseName:    "Squats",
    Reps:            &12,              // 12 reps
    DurationSeconds: nil,              // No duration
    Weight:          &100.5,           // 100.5 kg
}

// Cardio exercise
entry2 := WorkoutEntry{
    ExerciseName:    "Running",
    Reps:            nil,              // No reps
    DurationSeconds: &300,             // 300 seconds
    Weight:          nil,              // No weight
}

// Invalid (violates constraint)
entry3 := WorkoutEntry{
    Reps:            &12,
    DurationSeconds: &300,             // ERROR: Both set!
}

// Invalid (violates constraint)
entry4 := WorkoutEntry{
    Reps:            nil,
    DurationSeconds: nil,              // ERROR: Neither set!
}
```

### Example JSON

```json
{
  "id": 1,
  "workout_id": 5,
  "exercise_name": "Squats",
  "sets": 3,
  "reps": 12,
  "duration_seconds": null,
  "weight": 100.5,
  "notes": "Felt strong today",
  "order_index": 1
}
```

## User Structure

```go
type User struct {
    ID           int
    Username     string
    Email        string
    PasswordHash PasswordHash  // Custom type for security
    Bio          string
    CreatedAt    time.Time
    UpdatedAt    time.Time
}
```

### Example

```go
user := User{
    ID:       1,
    Username: "john_doe",
    Email:    "john@example.com",
    Bio:      "Fitness enthusiast",
    CreatedAt: time.Now(),
    UpdatedAt: time.Now(),
}
```

## Request/Response Structures

### User Registration Request

```go
type registerUserRequest struct {
    Username string `json:"username"`
    Email    string `json:"email"`
    Password string `json:"password"`
    Bio      string `json:"bio"`
}
```

### Token Creation Request

```go
type createTokenRequest struct {
    Username string `json:"username"`
    Password string `json:"password"`
}
```

### Workout Update Request

```go
var updateWorkOutRequest struct {
    Title           *string              `json:"title"`
    Description     *string              `json:"description"`
    DurationMinutes *int                 `json:"duration_minutes"`
    CaloriesBurned  *int                 `json:"calories_burned"`
    Entries         []store.WorkoutEntry `json:"entries"`
}
```

**Note:** Pointers make fields optional - only update fields that are provided!

## Envelope Pattern

Responses wrap data in an `Envelope`:

```go
type Envelope map[string]interface{}
```

**Usage:**

```go
utils.Envelope{"workout": workout}
utils.Envelope{"error": "Not found"}
utils.Envelope{"message": "Success"}
```

**JSON Response:**

```json
{
  "workout": {
    "id": 1,
    "title": "...",
    ...
  }
}
```

## Working with Pointers

### Checking if Pointer is Nil

```go
if entry.Reps != nil {
    repsValue := *entry.Reps  // Dereference to get value
    fmt.Println(repsValue)
}
```

### Creating Pointer Values

```go
// Create new int, get pointer to it
reps := 12
entry.Reps = &reps

// Or directly
entry.Reps = new(int)
*entry.Reps = 12
```

### JSON Marshaling with Pointers

Go automatically handles pointers in JSON:

```go
// Go struct
entry := WorkoutEntry{
    Reps:            &12,
    DurationSeconds: nil,
}

// JSON output
{
  "reps": 12,
  "duration_seconds": null
}
```

## Type Hierarchy

```
Application
├── Workout
│   ├── ID
│   ├── UserID (references User)
│   ├── Title
│   ├── Entries (WorkoutEntry slice)
│   └── ...
├── WorkoutEntry
│   ├── ID
│   ├── WorkoutId (references Workout)
│   ├── ExerciseName
│   ├── Reps (optional)
│   ├── DurationSeconds (optional)
│   └── ...
├── User
│   ├── ID
│   ├── Username
│   ├── Email
│   ├── PasswordHash
│   └── ...
└── Token
    ├── User ID (references User)
    ├── Expiration
    └── Scope
```

## Data Conversion Flow

```
HTTP Request (JSON)
    │
    ▼
json.NewDecoder().Decode()
    (JSON → Go struct)
    │
    ▼
Go struct (handler uses it)
    │
    ▼
json.NewEncoder().Encode()
    (Go struct → JSON)
    │
    ▼
HTTP Response (JSON)
```

## Example: From JSON to Go

### Request JSON

```json
POST /workouts
{
  "title": "Morning Run",
  "description": "5km run",
  "duration_minutes": 30,
  "calories_burned": 250,
  "entries": [
    {
      "exercise_name": "Warm-up jog",
      "sets": 1,
      "reps": null,
      "duration_seconds": 300,
      "weight": null,
      "notes": "Easy pace",
      "order_index": 1
    }
  ]
}
```

### Handler Code

```go
var workout store.Workout
json.NewDecoder(r.Body).Decode(&workout)

// Now workout contains:
// - Title: "Morning Run"
// - Description: "5km run"
// - DurationMinutes: 30
// - CaloriesBurned: 250
// - Entries[0].ExerciseName: "Warm-up jog"
// - Entries[0].DurationSeconds: &300
// - Entries[0].Reps: nil
```

## Example: From Go to JSON

### Handler Code

```go
workout := Workout{
    ID: 1,
    Title: "Morning Run",
    Entries: []WorkoutEntry{
        {
            ID: 1,
            ExerciseName: "Warm-up",
            Reps: nil,
            DurationSeconds: &300,
        },
    },
}

utils.WriteJson(w, http.StatusOK, Envelope{"workout": workout})
```

### Response JSON

```json
{
  "workout": {
    "id": 1,
    "title": "Morning Run",
    "entries": [
      {
        "id": 1,
        "exercise_name": "Warm-up",
        "reps": null,
        "duration_seconds": 300
      }
    ]
  }
}
```

## Null vs Zero Values

### Go (before JSON)

```go
var reps *int        // nil (no value)
var duration *int    // nil (no value)
```

### JSON (after encoding)

```json
{
  "reps": null,
  "duration_seconds": null
}
```

### Go (before decoding from JSON)

```json
{
  "reps": 12,
  "duration_seconds": null
}
```

### JSON (after decoding)

```go
reps := &12           // Pointer to 12
duration := (*int)(nil)  // nil
```

## Key Takeaways

1. **Struct Tags**: Map Go fields to JSON fields
2. **Pointers**: Allow optional fields
3. **Type Safety**: Go types ensure valid data
4. **JSON Conversion**: Automatic with encoding/json
5. **Nil Handling**: Use pointers for truly optional data
6. **Envelope Pattern**: Wraps responses for consistency

## Next Steps

- Part 8: Learn about **Routing System** - URL patterns and method handling
- Part 9: Learn about **Error Handling & Utils** - common utilities

---

**Key Concept**: Data models are the blueprint for how information flows through the system. Tags and pointers make JSON handling flexible!
