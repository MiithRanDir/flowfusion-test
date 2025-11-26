# API Documentation

## Authentication Endpoints

### Register User
Register a new user account.

- **URL**: `/auth/register`
- **Method**: `POST`
- **Auth Required**: No

#### Request Body
```json
{
  "email": "user@example.com",
  "password": "password123",
  "name": "John Doe"
}
```

#### Success Response (201 Created)
```json
{
  "message": "User registered successfully"
}
```

#### Error Response (400 Bad Request)
```json
{
  "error": "Email and password are required"
}
```

---

### Login
Authenticate a user and return access and refresh tokens.

- **URL**: `/auth/login`
- **Method**: `POST`
- **Auth Required**: No

#### Request Body
```json
{
  "email": "user@example.com",
  "password": "password123"
}
```

#### Success Response (200 OK)
```json
{
  "access_token": "eyJhbGciOiJIUzI1NiIs...",
  "refresh_token": "eyJhbGciOiJIUzI1NiIs..."
}
```

#### Error Response (401 Unauthorized)
```json
{
  "error": "Invalid credentials"
}
```

---

### Refresh Token
Get a new access token using a valid refresh token.

- **URL**: `/auth/refresh`
- **Method**: `POST`
- **Auth Required**: No

#### Request Body
```json
{
  "refresh_token": "eyJhbGciOiJIUzI1NiIs..."
}
```

#### Success Response (200 OK)
```json
{
  "access_token": "eyJhbGciOiJIUzI1NiIs...",
  "refresh_token": "eyJhbGciOiJIUzI1NiIs..."
}
```

#### Error Response (401 Unauthorized)
```json
{
  "error": "Invalid refresh token"
}
```

---

### Logout
Invalidate the refresh token and log the user out.

- **URL**: `/auth/logout`
- **Method**: `POST`
- **Auth Required**: Yes (Bearer Token)

#### Request Body
```json
{
  "refresh_token": "eyJhbGciOiJIUzI1NiIs..."
}
```

#### Success Response (200 OK)
```json
{
  "message": "Logged out successfully"
}
```

#### Error Response (500 Internal Server Error)
```json
{
  "error": "Logout failed"
}
```

---

## User Endpoints

### Get Current User
Get the profile of the currently authenticated user.

- **URL**: `/me`
- **Method**: `GET`
- **Auth Required**: Yes (Bearer Token)

#### Success Response (200 OK)
```json
{
  "id": 1,
  "email": "user@example.com",
  "name": "John Doe",
  "created_at": "2023-10-27T10:00:00Z",
  "updated_at": "2023-10-27T10:00:00Z"
}
```

#### Error Response (404 Not Found)
```json
{
  "error": "User not found"
}
```
