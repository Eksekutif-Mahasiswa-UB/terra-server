# Terra Server - Authentication API Documentation

**Base URL:** `http://localhost:8080/api/v1`

This document covers all authentication endpoints for Issue #3: Authentication System.

---

## Table of Contents

1. [Register User (Email/Password)](#1-register-user-emailpassword)
2. [Login (Email/Password)](#2-login-emailpassword)
3. [Login (Google Account)](#3-login-google-account)
4. [Forgot Password](#4-forgot-password)
5. [Reset Password](#5-reset-password)
6. [Refresh Token (Remember Me)](#6-refresh-token-remember-me)
7. [Logout](#7-logout)

---

## 1. Register User (Email/Password)

- **Endpoint:** `POST /auth/register`
- **Description:** Registers a new user with their email, password, and full name. The system will hash the password before storing it in the database.

### Request Body:
```json
{
  "full_name": "John Doe",
  "email": "john.doe@example.com",
  "password": "password123"
}
```

### Success Response (201 Created):
```json
{
  "message": "User registered successfully",
  "data": {
    "id": "550e8400-e29b-41d4-a716-446655440000",
    "full_name": "John Doe",
    "email": "john.doe@example.com",
    "role": "user",
    "auth_method": "email",
    "created_at": "2025-11-04T10:30:00Z",
    "updated_at": "2025-11-04T10:30:00Z"
  }
}
```

### Error Responses:

**400 Bad Request - Validation Error:**
```json
{
  "error": "Validation error",
  "message": "Email cannot be empty"
}
```

**409 Conflict - Email Already Registered:**
```json
{
  "error": "Registration failed",
  "message": "email is registered"
}
```

**500 Internal Server Error:**
```json
{
  "error": "Registration failed",
  "message": "Internal server error"
}
```

---

## 2. Login (Email/Password)

- **Endpoint:** `POST /auth/login`
- **Description:** Authenticates a user with email and password. Returns both an access token (15 minutes expiry) and a refresh token (7 days expiry).

### Request Body:
```json
{
  "email": "john.doe@example.com",
  "password": "password123"
}
```

### Success Response (200 OK):
```json
{
  "message": "Login successful",
  "data": {
    "access_token": "eyJhbGciOiJIUzI1NiIsInR5cCI6IkpXVCJ9.eyJ1c2VyX2lkIjoiNTUwZTg0MDAtZTI5Yi00MWQ0LWE3MTYtNDQ2NjU1NDQwMDAwIiwicm9sZSI6InVzZXIiLCJwdXJwb3NlIjoiYWNjZXNzIiwiZXhwIjoxNzMwNzI1ODAwLCJpYXQiOjE3MzA3MjQ5MDAsImp0aSI6ImFiYzEyMzQ1In0.signature",
    "refresh_token": "eyJhbGciOiJIUzI1NiIsInR5cCI6IkpXVCJ9.eyJ1c2VyX2lkIjoiNTUwZTg0MDAtZTI5Yi00MWQ0LWE3MTYtNDQ2NjU1NDQwMDAwIiwicm9sZSI6InVzZXIiLCJwdXJwb3NlIjoicmVmcmVzaCIsImV4cCI6MTczMTMyOTcwMCwiaWF0IjoxNzMwNzI0OTAwLCJqdGkiOiJ4eXo2Nzg5MCJ9.signature"
  }
}
```

### Error Responses:

**400 Bad Request - Invalid Request:**
```json
{
  "error": "Invalid request",
  "message": "Key: 'LoginRequest.Email' Error:Field validation for 'Email' failed on the 'required' tag"
}
```

**400 Bad Request - Google Account User:**
```json
{
  "error": "Login failed",
  "message": "this account is registered with Google. Please use Google login"
}
```

**401 Unauthorized - Invalid Credentials:**
```json
{
  "error": "Login failed",
  "message": "email or password is incorrect"
}
```

**500 Internal Server Error:**
```json
{
  "error": "Login failed",
  "message": "failed to generate tokens"
}
```

---

## 3. Login (Google Account)

- **Endpoint:** `POST /auth/login/google`
- **Description:** Authenticates a user using Google OAuth. If the user doesn't exist, a new account will be created automatically. Returns both access and refresh tokens.

### Request Body:
```json
{
  "credential": "eyJhbGciOiJSUzI1NiIsImtpZCI6IjY4YTE1..."
}
```

**Note:** The `credential` is the Google ID token obtained from Google Sign-In on the frontend.

### Success Response (200 OK):
```json
{
  "message": "Google login successful",
  "data": {
    "access_token": "eyJhbGciOiJIUzI1NiIsInR5cCI6IkpXVCJ9...",
    "refresh_token": "eyJhbGciOiJIUzI1NiIsInR5cCI6IkpXVCJ9..."
  }
}
```

### Error Responses:

**400 Bad Request - Email Account User:**
```json
{
  "error": "Google login failed",
  "message": "please log in using email and password"
}
```

**400 Bad Request - Validation Error:**
```json
{
  "error": "Validation error",
  "message": "Google credential cannot be empty"
}
```

**401 Unauthorized - Invalid Google Token:**
```json
{
  "error": "Google login failed",
  "message": "invalid Google token"
}
```

**500 Internal Server Error:**
```json
{
  "error": "Google login failed",
  "message": "failed to create user account"
}
```

---

## 4. Forgot Password

- **Endpoint:** `POST /auth/forgot-password`
- **Description:** Initiates a password reset process. If the email exists and is registered with email authentication (not Google), a password reset email will be sent with a token valid for 15 minutes. Always returns success to prevent email enumeration attacks.

### Request Body:
```json
{
  "email": "john.doe@example.com"
}
```

### Success Response (200 OK):
```json
{
  "message": "If your email is registered, you will receive a password reset link."
}
```

**Note:** This endpoint always returns a success message, regardless of whether the email exists or not. This is a security measure to prevent email enumeration attacks.

### Error Responses:

**400 Bad Request - Validation Error:**
```json
{
  "error": "Validation error",
  "message": "Email cannot be empty"
}
```

---

## 5. Reset Password

- **Endpoint:** `POST /auth/reset-password`
- **Description:** Resets a user's password using the token received via email. The token is valid for 15 minutes. Requires password confirmation.

### Request Body:
```json
{
  "token": "eyJhbGciOiJIUzI1NiIsInR5cCI6IkpXVCJ9...",
  "password": "newPassword123",
  "confirm_password": "newPassword123"
}
```

### Success Response (200 OK):
```json
{
  "message": "Password has been reset successfully."
}
```

### Error Responses:

**400 Bad Request - Validation Error:**
```json
{
  "error": "Validation error",
  "message": "Passwords do not match"
}
```

**400 Bad Request - Password Too Short:**
```json
{
  "error": "Validation error",
  "message": "Password must be at least 8 characters"
}
```

**401 Unauthorized - Invalid/Expired Token:**
```json
{
  "error": "Password reset failed",
  "message": "invalid or expired token"
}
```

**401 Unauthorized - Invalid Token Purpose:**
```json
{
  "error": "Password reset failed",
  "message": "invalid token purpose"
}
```

**500 Internal Server Error:**
```json
{
  "error": "Password reset failed",
  "message": "failed to update password"
}
```

---

## 6. Refresh Token (Remember Me)

- **Endpoint:** `POST /auth/refresh`
- **Description:** Obtains a new access token using a valid refresh token. The refresh token must have `purpose: "refresh"` in its claims and must exist in the database. Returns a new access token valid for 15 minutes.

### Request Body:
```json
{
  "refresh_token": "eyJhbGciOiJIUzI1NiIsInR5cCI6IkpXVCJ9..."
}
```

### Success Response (200 OK):
```json
{
  "message": "Token refreshed successfully",
  "data": {
    "access_token": "eyJhbGciOiJIUzI1NiIsInR5cCI6IkpXVCJ9.eyJ1c2VyX2lkIjoiNTUwZTg0MDAtZTI5Yi00MWQ0LWE3MTYtNDQ2NjU1NDQwMDAwIiwicm9sZSI6InVzZXIiLCJwdXJwb3NlIjoiYWNjZXNzIiwiZXhwIjoxNzMwNzI2NzAwLCJpYXQiOjE3MzA3MjU4MDAsImp0aSI6ImRlZjQ1Njc4In0.signature"
  }
}
```

### Error Responses:

**400 Bad Request - Validation Error:**
```json
{
  "error": "Validation error",
  "message": "Refresh token cannot be empty"
}
```

**401 Unauthorized - Invalid Token:**
```json
{
  "error": "Token refresh failed",
  "message": "invalid or expired token"
}
```

**401 Unauthorized - Not a Refresh Token:**
```json
{
  "error": "Token refresh failed",
  "message": "invalid token: not a refresh token"
}
```

**401 Unauthorized - Revoked Token:**
```json
{
  "error": "Token refresh failed",
  "message": "token is invalid or has been revoked"
}
```

**401 Unauthorized - Expired Token:**
```json
{
  "error": "Token refresh failed",
  "message": "refresh token has expired"
}
```

**500 Internal Server Error:**
```json
{
  "error": "Token refresh failed",
  "message": "failed to generate access token"
}
```

---

## 7. Logout

- **Endpoint:** `POST /auth/logout`
- **Description:** Logs out a user by invalidating their refresh token. The refresh token is deleted from the database, preventing it from being used for future authentication.

### Request Body:
```json
{
  "refresh_token": "eyJhbGciOiJIUzI1NiIsInR5cCI6IkpXVCJ9..."
}
```

### Success Response (200 OK):
```json
{
  "message": "Logged out successfully"
}
```

### Error Responses:

**400 Bad Request - Validation Error:**
```json
{
  "error": "Validation error",
  "message": "Refresh token cannot be empty"
}
```

**500 Internal Server Error:**
```json
{
  "error": "Logout failed",
  "message": "failed to logout"
}
```

---

## Authentication Flow

### Initial Login Flow:
1. User logs in via `/auth/login` or `/auth/login/google`
2. Server returns both `access_token` (15 min) and `refresh_token` (7 days)
3. Client stores both tokens securely
4. Client uses `access_token` for authenticated API requests

### Token Refresh Flow (When Access Token Expires):
1. Client detects access token expiry (401 error)
2. Client calls `/auth/refresh` with `refresh_token`
3. Server validates refresh token and returns new `access_token`
4. Client continues making authenticated requests

### Logout Flow:
1. Client calls `/auth/logout` with `refresh_token`
2. Server deletes refresh token from database
3. Token is permanently invalidated
4. Client clears all stored tokens

---

## Security Features

1. **Token Purpose Validation**: Access and refresh tokens have distinct `purpose` claims to prevent misuse
2. **Stateful Refresh Tokens**: Refresh tokens are stored in the database and can be invalidated
3. **Password Hashing**: All passwords are hashed using bcrypt before storage
4. **Email Enumeration Prevention**: Forgot password always returns success
5. **Auth Method Validation**: Prevents users from mixing Google and email authentication
6. **Token Expiry**: Short-lived access tokens (15 min) and longer refresh tokens (7 days)
7. **Secure Logout**: Tokens are deleted server-side, not just cleared client-side

---

## HTTP Status Codes

| Code | Description |
|------|-------------|
| 200 | Success (OK) |
| 201 | Resource Created Successfully |
| 400 | Bad Request (Validation Error) |
| 401 | Unauthorized (Invalid Credentials/Token) |
| 409 | Conflict (Resource Already Exists) |
| 500 | Internal Server Error |

---

## Common Error Response Format

All error responses follow this structure:

```json
{
  "error": "Error Category",
  "message": "Detailed error message"
}
```

---

## Testing with cURL

### Register:
```bash
curl -X POST http://localhost:8080/api/v1/auth/register \
  -H "Content-Type: application/json" \
  -d '{
    "full_name": "John Doe",
    "email": "john.doe@example.com",
    "password": "password123"
  }'
```

### Login:
```bash
curl -X POST http://localhost:8080/api/v1/auth/login \
  -H "Content-Type: application/json" \
  -d '{
    "email": "john.doe@example.com",
    "password": "password123"
  }'
```

### Refresh Token:
```bash
curl -X POST http://localhost:8080/api/v1/auth/refresh \
  -H "Content-Type: application/json" \
  -d '{
    "refresh_token": "YOUR_REFRESH_TOKEN_HERE"
  }'
```

### Logout:
```bash
curl -X POST http://localhost:8080/api/v1/auth/logout \
  -H "Content-Type: application/json" \
  -d '{
    "refresh_token": "YOUR_REFRESH_TOKEN_HERE"
  }'
```

---

## Notes

- All endpoints accept `Content-Type: application/json`
- All responses return `Content-Type: application/json`
- Access tokens should be sent in the `Authorization` header as `Bearer <token>` for protected routes
- Refresh tokens should only be sent to the `/auth/refresh` endpoint
- Store tokens securely on the client (e.g., httpOnly cookies or secure storage)

---

**Last Updated:** November 4, 2025  
**API Version:** v1  
**Issue:** #3 - Authentication User
