# 06 - Workout Operations (CRUD)

## CRUD Operations Overview

**CRUD** stands for:

- **C**reate - Add new workouts
- **R**ead - Get workout data
- **U**pdate - Modify existing workouts
- **D**elete - Remove workouts

All workouts are owned by the user who created them. Users can only manage their own workouts.

## Workout Data Model

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

## CREATE - Adding a New Workout

### Handler Code

```go
func (h *WorkoutHandler) HandleCreateWorkout(w http.ResponseWriter, r *http.Request) {
    // Step 1: Decode request body
    var workout store.Workout
    err := json.NewDecoder(r.Body).Decode(&workout)
    if err != nil {
        h.logger.Printf("Error: decodeCreateWorkout : %v", err)
        utils.WriteJson(w, http.StatusBadRequest,
            utils.Envelope{"error": "Unable to Decode the workout"})
        return
    }

    // Step 2: Get authenticated user
    currentUser := middleware.GetUser(r)
    if currentUser == nil || currentUser == store.AnonymousUser {
        h.logger.Printf("Error: middleware GetUser : %v", err)
        utils.WriteJson(w, http.StatusBadRequest,
            utils.Envelope{"error": "you must be loggedIn"})
        return
    }

    // Step 3: Attach user ID to workout
    workout.UserID = currentUser.ID

    // Step 4: Store in database
    createdWorkout, err := h.store.CreateWorkout(&workout)
    if err != nil {
        h.logger.Printf("Error: createWorkout : %v", err)
        utils.WriteJson(w, http.StatusInternalServerError,
            utils.Envelope{"error": "unable to Store workout"})
        return
    }

    // Step 5: Return created workout
    utils.WriteJson(w, http.StatusOK,
        utils.Envelope{"CreatedWorkout": createdWorkout})
}
```

### Database Storage

```go
func (pgws *PostgresWorkoutStore) CreateWorkout(workout *Workout) (*Workout, error) {
    // Start transaction
    tx, err := pgws.db.Begin()
    if err != nil {
        return nil, err
    }
    defer tx.Rollback()

    // Step 1: Insert workout record
    query := `INSERT INTO workouts
              (title, description, duration_minutes, calories_burned, user_id)
              VALUES ($1, $2, $3, $4, $5) RETURNING id;`

    err = tx.QueryRow(query,
        workout.Title,
        workout.Description,
        workout.DurationMinutes,
        workout.CaloriesBurned,
        workout.UserID).Scan(&workout.ID)

    if err != nil {
        return nil, err
    }

    // Step 2: Insert workout entries
    for index := range workout.Entries {
        entry := &workout.Entries[index]

        query := `INSERT INTO workout_entries
                  (workout_id, exercise_name, sets, reps,
                   duration_seconds, weight, notes, order_index)
                  VALUES ($1, $2, $3, $4, $5, $6, $7, $8)
                  RETURNING id;`

        err := tx.QueryRow(query,
            workout.ID,
            entry.ExerciseName,
            entry.Sets,
            entry.Reps,
            entry.DurationSeconds,
            entry.Weight,
            entry.Notes,
            entry.OrderIndex).Scan(&entry.ID)

        if err != nil {
            return nil, err
        }
    }

    // Step 3: Commit transaction
    err = tx.Commit()
    if err != nil {
        return nil, err
    }

    return workout, nil
}
```

**Flow Diagram:**

```
POST /workouts
    │
    ▼
Inside transaction:
    ├─ INSERT workout row
    ├─ Get generated ID
    ├─ FOR EACH entry:
    │   └─ INSERT workout_entry
    │
    ▼
COMMIT transaction
    │
    ▼
Return 200 OK with created workout
```

### Example: Create Workout Request

**Request:**

```http
POST /workouts HTTP/1.1
Authorization: Bearer eyJhbGc...
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
      "notes": "Felt strong",
      "order_index": 1
    },
    {
      "exercise_name": "Bench Press",
      "sets": 4,
      "reps": 10,
      "weight": 80.0,
      "notes": "Good form",
      "order_index": 2
    }
  ]
}
```

**Response:**

```http
HTTP/1.1 200 OK
Content-Type: application/json

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
        "notes": "Felt strong",
        "order_index": 1
      },
      {
        "id": 2,
        "workout_id": 1,
        "exercise_name": "Bench Press",
        "sets": 4,
        "reps": 10,
        "duration_seconds": null,
        "weight": 80.0,
        "notes": "Good form",
        "order_index": 2
      }
    ]
  }
}
```

## READ - Getting a Workout

### Handler Code

```go
func (h *WorkoutHandler) HandleGetWorkoutById(w http.ResponseWriter, r *http.Request) {
    // Step 1: Extract workout ID from URL
    workoutId, err := utils.ReadIDParam(r)
    if err != nil {
        h.logger.Printf("Error: readIDParam : %v", err)
        utils.WriteJson(w, http.StatusBadRequest,
            utils.Envelope{"error": "invalid get workout id"})
        return
    }

    // Step 2: Query database
    workout, err := h.store.GetWorkoutById(workoutId)
    if err != nil {
        h.logger.Printf("Error: getWorkoutByID: %v", err)
        utils.WriteJson(w, http.StatusInternalServerError,
            utils.Envelope{"error": "internal server error"})
        return
    }

    // Step 3: Return workout
    utils.WriteJson(w, http.StatusOK,
        utils.Envelope{"workout": workout})
}
```

### Database Query

```go
func (pgws *PostgresWorkoutStore) GetWorkoutById(id int64) (*Workout, error) {
    // Start transaction
    tx, err := pgws.db.Begin()
    if err != nil {
        return nil, err
    }
    defer tx.Rollback()

    // Step 1: Select workout
    query := `SELECT id, user_id, title, description, duration_minutes,
                     calories_burned, created_at, updated_at
              FROM workouts WHERE id = $1;`

    row := tx.QueryRow(query, id)
    workout := &Workout{}
    err = row.Scan(&workout.ID, &workout.UserID, &workout.Title,
        &workout.Description, &workout.DurationMinutes,
        &workout.CaloriesBurned, &workout.CreatedAt, &workout.UpdatedAt)

    if err == sql.ErrNoRows {
        return nil, err
    }
    if err != nil {
        return nil, err
    }

    // Step 2: Select workout entries
    entriesQuery := `SELECT id, workout_id, exercise_name, sets, reps,
                            duration_seconds, weight, notes, order_index
                     FROM workout_entries WHERE workout_id = $1
                     ORDER BY order_index;`

    rows, err := tx.Query(entriesQuery, id)
    if err != nil {
        return nil, err
    }
    defer rows.Close()

    var entries []WorkoutEntry
    for rows.Next() {
        var entry WorkoutEntry
        err := rows.Scan(&entry.ID, &entry.WorkoutId,
            &entry.ExerciseName, &entry.Sets, &entry.Reps,
            &entry.DurationSeconds, &entry.Weight,
            &entry.Notes, &entry.OrderIndex)
        if err != nil {
            return nil, err
        }
        entries = append(entries, entry)
    }

    workout.Entries = entries

    // Commit transaction
    err = tx.Commit()
    if err != nil {
        return nil, err
    }

    return workout, nil
}
```

## UPDATE - Modifying a Workout

### Handler Code

```go
func (h *WorkoutHandler) HandleUpdateWorkoutByID(w http.ResponseWriter, r *http.Request) {
    // Step 1: Get workout ID
    workoutId, err := utils.ReadIDParam(r)
    if err != nil {
        utils.WriteJson(w, http.StatusBadRequest,
            utils.Envelope{"error": "invalid update workout id"})
        return
    }

    // Step 2: Get existing workout
    existingWorkout, err := h.store.GetWorkoutById(workoutId)
    if err != nil || existingWorkout == nil {
        http.NotFound(w, r)
        return
    }

    // Step 3: Parse update request (using pointers for optional fields)
    var updateRequest struct {
        Title           *string              `json:"title"`
        Description     *string              `json:"description"`
        DurationMinutes *int                 `json:"duration_minutes"`
        CaloriesBurned  *int                 `json:"calories_burned"`
        Entries         []store.WorkoutEntry `json:"entries"`
    }
    err = json.NewDecoder(r.Body).Decode(&updateRequest)
    if err != nil {
        utils.WriteJson(w, http.StatusBadRequest,
            utils.Envelope{"error": "unable to decode workout"})
        return
    }

    // Step 4: Apply only provided updates (preserve other fields)
    if updateRequest.Title != nil {
        existingWorkout.Title = *updateRequest.Title
    }
    if updateRequest.Description != nil {
        existingWorkout.Description = *updateRequest.Description
    }
    if updateRequest.DurationMinutes != nil {
        existingWorkout.DurationMinutes = *updateRequest.DurationMinutes
    }
    if updateRequest.CaloriesBurned != nil {
        existingWorkout.CaloriesBurned = *updateRequest.CaloriesBurned
    }
    if updateRequest.Entries != nil {
        existingWorkout.Entries = updateRequest.Entries
    }

    // Step 5: Verify user owns this workout
    currentUser := middleware.GetUser(r)
    workoutOwner, err := h.store.GetWorkoutOwner(workoutId)

    if workoutOwner != currentUser.ID {
        utils.WriteJson(w, http.StatusForbidden,
            utils.Envelope{"error": "you are not authorized to update this workout!"})
        return
    }

    // Step 6: Update in database
    err = h.store.UpdateWorkout(existingWorkout)
    if err != nil {
        utils.WriteJson(w, http.StatusInternalServerError,
            utils.Envelope{"error": "unable to store the updated workout"})
        return
    }

    // Step 7: Return success
    utils.WriteJson(w, http.StatusOK,
        utils.Envelope{"Message": "Workout Updated Successfully!"})
}
```

**Key Features:**

1. **Optional Fields**: Uses pointers (`*string`, `*int`) for optional updates
2. **Partial Updates**: Only updates provided fields
3. **Authorization Check**: Ensures user owns workout
4. **Entry Replacement**: Replaces all entries (delete old, insert new)

### Example: Update Workout Request

**Request:**

```http
PUT /workouts/1 HTTP/1.1
Authorization: Bearer eyJhbGc...
Content-Type: application/json

{
  "title": "Afternoon Strength Training",
  "duration_minutes": 70
}
```

**Response:**

```http
HTTP/1.1 200 OK
{"Message": "Workout Updated Successfully!"}
```

## DELETE - Removing a Workout

### Handler Code

```go
func (h *WorkoutHandler) HandleDeleteWorkoutByID(w http.ResponseWriter, r *http.Request) {
    // Step 1: Get workout ID
    workoutId, err := utils.ReadIDParam(r)
    if err != nil {
        utils.WriteJson(w, http.StatusBadRequest,
            utils.Envelope{"error": "invalid delete workout id"})
        return
    }

    // Step 2: Verify user is authenticated
    currentUser := middleware.GetUser(r)
    if currentUser == nil || currentUser == store.AnonymousUser {
        utils.WriteJson(w, http.StatusBadRequest,
            utils.Envelope{"error": "you must be loggedIn to delete"})
        return
    }

    // Step 3: Verify user owns this workout
    workoutOwner, err := h.store.GetWorkoutOwner(workoutId)
    if err != nil {
        if errors.Is(err, sql.ErrNoRows) {
            utils.WriteJson(w, http.StatusNotFound,
                utils.Envelope{"error": "unable to get the workout owner"})
            return
        }
        utils.WriteJson(w, http.StatusInternalServerError,
            utils.Envelope{"error": "internal server error"})
        return
    }

    if workoutOwner != currentUser.ID {
        utils.WriteJson(w, http.StatusForbidden,
            utils.Envelope{"error": "you are not authorized to delete this workout!"})
        return
    }

    // Step 4: Delete from database
    err = h.store.DeleteWorkout(workoutId)
    if err != nil {
        utils.WriteJson(w, http.StatusInternalServerError,
            utils.Envelope{"error": "Unable to Delete Workout"})
        return
    }

    // Step 5: Return success
    utils.WriteJson(w, http.StatusOK,
        utils.Envelope{"message": "Workout Deleted Successfully!"})
}
```

### Database Deletion

```go
func (pgws *PostgresWorkoutStore) DeleteWorkout(id int64) error {
    tx, err := pgws.db.Begin()
    if err != nil {
        return err
    }
    defer tx.Rollback()

    // Delete workout (ON DELETE CASCADE automatically deletes entries)
    query := `DELETE FROM workouts WHERE id = $1`
    result, err := tx.Exec(query, id)
    if err != nil {
        return err
    }

    rowsAffected, err := result.RowsAffected()
    if err != nil {
        return err
    }

    if rowsAffected == 0 {
        return errors.New("No Rows Affected!")
    }

    return tx.Commit()
}
```

**Note:** Thanks to `ON DELETE CASCADE` in the workout_entries table, deleting a workout automatically deletes all its entries!

## Authorization Pattern

All modifying operations follow this pattern:

```
1. Get resource from database
2. Check if exists (404 if not)
3. Get current user from context
4. Check if user owns resource
5. ├─ No → Return 403 Forbidden
6. └─ Yes → Perform operation
7. Return result
```

## CRUD HTTP Methods

| Operation | Method | Endpoint       | Status |
| --------- | ------ | -------------- | ------ |
| Create    | POST   | /workouts      | 200    |
| Read      | GET    | /workouts/{id} | 200    |
| Update    | PUT    | /workouts/{id} | 200    |
| Delete    | DELETE | /workouts/{id} | 200    |

## Transaction Pattern

All database operations use transactions:

```go
tx, err := db.Begin()
defer tx.Rollback()

// Do operations
tx.QueryRow()
tx.Exec()

err = tx.Commit()
```

**Why transactions?**

- All-or-nothing: If any operation fails, all rollback
- Consistency: Data never in partial state
- Atomicity: Operation succeeds completely or not at all

## Key Takeaways

1. **Authorization**: Always check user owns workout
2. **Transactions**: Bundle related operations
3. **Optional Fields**: Use pointers for partial updates
4. **Cascading Deletes**: Database enforces cleanup
5. **Error Handling**: Return appropriate status codes
6. **Entries**: Always managed with workouts

## Next Steps

- Part 7: Learn about **Data Models & Structures** - request/response formats
- Part 8: Learn about **Routing System** - how requests are routed

---

**Key Concept**: CRUD operations always follow the pattern: parse → validate → authorize → execute → respond!
