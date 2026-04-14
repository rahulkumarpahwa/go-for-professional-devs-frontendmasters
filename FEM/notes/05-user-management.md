# 05 - User Management

## User System Overview

The user management system handles:

1. **User Registration** - Creating new accounts
2. **User Storage** - Persisting user data to database
3. **Password Hashing** - Securing passwords cryptographically
4. **User Validation** - Ensuring data quality

## User Data Model

In `internal/store/workout_store.go`:

```go
type User struct {
    ID           int
    Username     string
    Email        string
    PasswordHash PasswordHash  // Special type for security
    Bio          string
    CreatedAt    time.Time
    UpdatedAt    time.Time
}

// Anonymous user for unauthenticated requests
var AnonymousUser = &User{}

func (u *User) IsAnonymous() bool {
    return u.ID == 0
}
```

**Explanation:**

- `ID`: Unique identifier assigned by database
- `Username`: For login (must be unique)
- `Email`: For contact (must be unique)
- `PasswordHash`: Encrypted password (never stored plain!)
- `Bio`: Optional user biography
- `Timestamps`: Track creation and last update
- `IsAnonymous()`: Check if user is logged in

## User Handler

Located in `internal/api/user_handler.go`:

### Registration Request Structure

```go
type registerUserRequest struct {
    Username string `json:"username"`
    Email    string `json:"email"`
    Password string `json:"password"`
    Bio      string `json:"bio"`
}
```

**Simple Explanation:**

- Struct tags (`json:"..."`) map JSON fields to struct fields
- These come from the HTTP POST request body

### Validation

```go
func (h *UserHandler) validateRegisterRequest(req *registerUserRequest) error {
    // Check required fields
    if req.Username == "" {
        return errors.New("Username is required")
    }
    if req.Email == "" {
        return errors.New("Email is required")
    }
    if req.Password == "" {
        return errors.New("Password is required")
    }
    if req.Bio == "" {
        return errors.New("Bio is required")
    }

    // Check field lengths
    if len(req.Username) > 50 {
        return errors.New("Username can't be greater than 50 characters")
    }

    // Validate email format with regex
    emailRegex := regexp.MustCompile(`^[a-zA-Z0-9._%+-]+@[a-zA-Z0-9.-]+\.[a-zA-Z]{2,}$`)
    if !emailRegex.MatchString(req.Email) {
        return errors.New("Invalid Email!")
    }

    return nil
}
```

**Validation Rules:**

1. All fields required (no empty values)
2. Username max 50 characters
3. Email must match email format
4. Password required

### Password Hashing

```go
err = user.PasswordHash.Set(req.Password)
if err != nil {
    h.logger.Printf("ERROR : hashing password: %v", err)
    utils.WriteJson(w, http.StatusInternalServerError,
        utils.Envelope{"error": "Internal Server Error"})
    return
}
```

**What happens:**

1. Plaintext password received from client
2. `PasswordHash.Set()` hashes it using bcrypt
3. Hash stored in database (plaintext never stored!)
4. Original password discarded from memory

### Registration Handler

```go
func (h *UserHandler) HandleRegisterUser(w http.ResponseWriter, r *http.Request) {
    // Step 1: Decode JSON request
    var req registerUserRequest
    err := json.NewDecoder(r.Body).Decode(&req)
    if err != nil {
        h.logger.Printf("ERROR : decoding register request: %v", err)
        utils.WriteJson(w, http.StatusBadRequest,
            utils.Envelope{"error": "Invalid register request payload!"})
        return
    }

    // Step 2: Validate request
    err = h.validateRegisterRequest(&req)
    if err != nil {
        h.logger.Printf("ERROR : validating register request: %v", err)
        utils.WriteJson(w, http.StatusBadRequest,
            utils.Envelope{"error": err.Error()})
        return
    }

    // Step 3: Create user struct
    user := &store.User{
        Username: req.Username,
        Email:    req.Email,
    }

    if req.Bio != "" {
        user.Bio = req.Bio
    }

    // Step 4: Hash password
    err = user.PasswordHash.Set(req.Password)
    if err != nil {
        h.logger.Printf("ERROR : hashing password: %v", err)
        utils.WriteJson(w, http.StatusInternalServerError,
            utils.Envelope{"error": "Internal Server Error"})
        return
    }

    // Step 5: Store in database
    err = h.store.CreateUser(user)
    if err != nil {
        h.logger.Printf("ERROR : registering user: %v", err)
        utils.WriteJson(w, http.StatusInternalServerError,
            utils.Envelope{"error": "Internal Server Error"})
        return
    }

    // Step 6: Return success response
    utils.WriteJson(w, http.StatusOK,
        utils.Envelope{
            "message": "User Created Successfully!",
            "user": user,
        })
}
```

## User Store (Database Layer)

The `UserStore` interface defines database operations:

```go
type UserStore interface {
    CreateUser(*User) error
    GetUserByUsername(username string) (*User, error)
    UpdateUser(*User) error
    DeleteUser(id int) error
}
```

### CreateUser Flow

```
Request arrives with credentials
        │
        ▼
Decode JSON
        │
        ▼
Validate all fields
        │
        ▼
Hash password with bcrypt
        │
        ▼
INSERT into users table
        ├─ Check username unique ✓
        └─ Check email unique ✓
        │
        ▼
Get new user ID from database
        │
        ▼
Return success response
```

## Registration Flow Example

### Request:

```http
POST /users HTTP/1.1
Content-Type: application/json

{
  "username": "john_doe",
  "email": "john@example.com",
  "password": "SecurePassword123!",
  "bio": "Fitness enthusiast"
}
```

### Response (Success):

```http
HTTP/1.1 200 OK
Content-Type: application/json

{
  "message": "User Created Successfully!",
  "user": {
    "id": 1,
    "username": "john_doe",
    "email": "john@example.com",
    "bio": "Fitness enthusiast",
    "created_at": "2024-01-15T10:30:00Z",
    "updated_at": "2024-01-15T10:30:00Z"
  }
}
```

### Response (Validation Error):

```http
HTTP/1.1 400 Bad Request
Content-Type: application/json

{
  "error": "Invalid Email!"
}
```

### Response (Email Already Exists):

```http
HTTP/1.1 500 Internal Server Error
Content-Type: application/json

{
  "error": "Internal Server Error"
}
```

## Authentication Handler

In `internal/api/token_handler.go`:

### Login Request

```go
type createTokenRequest struct {
    Username string `json:"username"`
    Password string `json:"password"`
}
```

### Token Creation Handler

```go
func (h *TokenHandler) HandleCreateToken(w http.ResponseWriter, r *http.Request) {
    // Step 1: Decode login request
    var req createTokenRequest
    err := json.NewDecoder(r.Body).Decode(&req)
    if err != nil {
        h.logger.Printf("ERROR: createTokenRequest : %v", err)
        utils.WriteJson(w, http.StatusBadRequest,
            utils.Envelope{"error ": "Invalid request payload"})
        return
    }

    // Step 2: Get user by username
    user, err := h.userStore.GetUserByUsername(req.Username)
    if err != nil || user == nil {
        h.logger.Printf("ERROR: getUserByUsername : %v", err)
        utils.WriteJson(w, http.StatusInternalServerError,
            utils.Envelope{"error ": "Internal Server Error"})
        return
    }

    // Step 3: Check password matches
    passwordDoMatch, err := user.PasswordHash.Matches(req.Password)
    if err != nil {
        h.logger.Printf("ERROR: passwordHash Match : %v", err)
        utils.WriteJson(w, http.StatusInternalServerError,
            utils.Envelope{"error ": "Internal Server Error"})
        return
    }

    if !passwordDoMatch {
        utils.WriteJson(w, http.StatusUnauthorized,
            utils.Envelope{"error ": "Invalid Credentials"})
        return
    }

    // Step 4: Create token for user
    token, err := h.tokenStore.CreateNewToken(
        user.ID,
        time.Hour*24,           // Valid for 24 hours
        tokens.ScopeAuth,       // Scope: authentication
    )
    if err != nil {
        h.logger.Printf("ERROR: CreatingToken : %v", err)
        utils.WriteJson(w, http.StatusInternalServerError,
            utils.Envelope{"error ": "Internal Server Error"})
        return
    }

    // Step 5: Return token
    utils.WriteJson(w, http.StatusCreated,
        utils.Envelope{"auth_token": token})
}
```

## Complete User Flow

### Step 1: Registration

```
Client                          Server
  │                               │
  ├─ POST /users ─────────────────►
  │   + username: "john_doe"      │
  │   + email: "john@..."         │
  │   + password: "abc123"        │
  │                               │
  │                    ┌──────────────────────┐
  │                    │ Validate request     │
  │                    │ Hash password        │
  │                    │ Store in DB          │
  │                    └──────────────────────┘
  │                               │
  │◄─ 200 OK ────────────────────┤
  │   + user data               │
```

### Step 2: Login

```
Client                          Server
  │                               │
  ├─ POST /tokens/authentication──►
  │   + username: "john_doe"      │
  │   + password: "abc123"        │
  │                               │
  │                    ┌──────────────────────┐
  │                    │ Get user by username │
  │                    │ Check password match │
  │                    │ Create token         │
  │                    └──────────────────────┘
  │                               │
  │◄─ 201 Created ────────────────┤
  │   + auth_token: "eyJ..."     │
```

### Step 3: Use Token

```
Client                          Server
  │                               │
  ├─ POST /workouts ─────────────►
  │   + Authorization: Bearer...  │
  │   + workout data              │
  │                               │
  │                    ┌──────────────────────┐
  │                    │ Authenticate         │
  │                    │ RequireUser          │
  │                    │ Create workout       │
  │                    │ Link to user         │
  │                    └──────────────────────┘
  │                               │
  │◄─ 200 OK ────────────────────┤
  │   + created workout          │
```

## Error Scenarios

### Duplicate Username

**Request:**

```json
{
  "username": "john_doe",
  "email": "different@example.com",
  "password": "password",
  "bio": "bio"
}
```

**Database Error:**

```
ERROR: duplicate key value violates unique constraint "users_username_key"
```

**Response:**

```http
HTTP/1.1 500 Internal Server Error
{"error": "Internal Server Error"}
```

### Invalid Password During Login

**Request:**

```json
{
  "username": "john_doe",
  "password": "wrongpassword"
}
```

**Response:**

```http
HTTP/1.1 401 Unauthorized
{"error": "Invalid Credentials"}
```

## Security Practices

1. **Password Hashing**: Passwords hashed with bcrypt
2. **No Plaintext**: Passwords never stored in plain text
3. **Validation**: Both client and server validation
4. **Unique Fields**: Username and email must be unique
5. **Token Expiration**: Tokens expire after 24 hours
6. **HTTPS**: In production, use HTTPS (not shown here)

## Database Operations

When registering, a transaction runs:

```sql
BEGIN TRANSACTION;

INSERT INTO users (username, email, password_hash, bio)
VALUES ('john_doe', 'john@example.com', '$2a$10$...', 'bio')
RETURNING id;

COMMIT;
```

Password hash example (bcrypt):

```
Plain text:  abc123
Hashed:      $2a$10$N9qo8ucoqwgaqHZF8v/...
```

## Key Takeaways

1. **Validation First**: Always validate before storing
2. **Hash Passwords**: Never store plain passwords
3. **Token System**: Users get tokens to authenticate
4. **Unique Email/Username**: Database enforces uniqueness
5. **Error Handling**: Return appropriate HTTP status codes
6. **User Context**: User info stored in request context

## Next Steps

- Part 6: Learn about **Workout Operations** - CRUD for workouts
- Part 7: Learn about **Data Models** - request/response structures

---

**Key Concept**: User management is the foundation of authentication. Registration creates the user, login creates the token, token is used for all other operations!
