# Vexal Key Server

A secure license key management system built with Go, featuring user authentication, invite system, and hardware ID (HWID) binding.

## 🚀 Features

- License key generation and validation
- User management system with admin controls
- Invite-based registration system
- Hardware ID (HWID) binding and reset functionality
- VPN detection for enhanced security
- PostgreSQL database integration

## 📋 Prerequisites

- Go 1.23.3 or higher
- PostgreSQL
- Docker (optional)

## 🛠️ Installation

1. Clone the repository
2. Create a `.env` file with the following variables:

```
DATABASE_URL=postgres://user:password@localhost:5432/dbname?sslmode=disable
PORT=3000
VERSION=1.0.0
```

3. Run with Docker:

```bash
docker-compose up
```

Or build manually:

```bash
make build_me
```

In case you want to build for every architecture:

```bash
make build
```

## 📚 API Documentation

### Authentication Endpoints

#### Verify License Key (HWID)

```
POST /api/auth/verify
```

Verifies a license key with HWID binding.

Request Body:

```json
{
  "username": "string",
  "license_key": "string",
  "hwid": "string",
  "nonce": "string",
  "timestamp": number
}
```

Response:

```json
{
  "success": boolean,
  "message": string,
  "result": boolean
}
```

#### Verify License Key (No HWID)

```
POST /api/auth/verifykey
```

Verifies a license key without HWID binding or additional security checks.

Request Body:

```json
{
  "username": "string",
  "license_key": "string"
}
```

Response:

```json
{
  "success": boolean,
  "message": string,
  "result": boolean
}
```

#### Login

```
POST /api/auth/login
```

Logs in a user with username and password.

Request Body:

```json
{
  "username": "string",
  "password": "string"
}
```

Response:

```json
{
  "success": boolean,
  "message": string,
  "result": boolean
}
```

#### Register

```
POST /api/auth/register
```

Registers a new user with an invite code.

Request Body:

```json
{
  "username": "string",
  "password": "string",
  "invite_code": "string"
}
```

Response:

```json
{
  "success": boolean,
  "message": string,
  "result": boolean
}
```

### User Management Endpoints

#### Get User Metadata

```
GET /api/@me/metadata
```

Retrieves user metadata.

Request Body:

```json
{
  "username": "string",
  "password": "string"
}
```

Response:

```json
{
  "message": string,
  "result": {
    "id": number,
    "created_at": string,
    "username": string,
    "password": string,
    "license_key": string | null,
    "admin": boolean,
    "active": boolean,
    "hwid": string | null,
    "invited_by": number
  },
  "success": boolean
}
```

#### Update User Metadata

```
PATCH /api/@me/metadata
```

Updates username or password.

Request Body:

```json
{
  "username": "string",
  "password": "string",
  "type": "username|password",
  "new_value": "string"
}
```

Response:

```json
{
  "success": boolean,
  "message": string,
  "result": boolean
}
```

#### Reset HWID

```
POST /api/@me/hwidreset
```

Resets the hardware ID with a 7-day cooldown.

Request Body:

```json
{
  "username": "string",
  "password": "string"
}
```

Response:

```json
{
  "success": boolean,
  "message": string,
  "result": boolean
}
```

#### List User Invites

```
GET /api/@me/invites
```

Lists all invite codes associated with the user.

Request Body:

```json
{
  "username": "string",
  "password": "string"
}
```

Response:

```json
{
  "message": string,
  "result": [
    {
      "id": number,
      "user_id": number,
      "invite_code": string,
      "created_at": string,
      "used_by": number | null,
      "used_at": string | null,
      "active": boolean
    }
  ],
  "success": boolean
}
```

#### Deactivate Account

```
DELETE /api/@me/deactivate
```

Deactivates user account.

Request Body:

```json
{
  "username": "string",
  "password": "string"
}
```

Response:

```json
{
  "success": boolean,
  "message": string,
  "result": boolean
}
```

### Admin Endpoints

### Login

```
POST /api/admin/auth/login
```

Logs in as an admin user.

Request Body:

```json
{
  "username": "string",
  "password": "string"
}
```

Response:

```json
{
  "success": boolean,
  "message": string,
  "result": boolean
}
```

#### Reset HWID

```
POST /api/admin/hwid/reset
```

Resets the hardware ID of a user.

Request Body:

```json
{
  "username": "string",
  "password": "string",
  "target_user": "string"
}
```

Response:

```json
{
  "success": boolean,
  "message": string,
  "result": boolean
}
```

#### Promote User

```
POST /api/admin/users/promote
```

Promotes a user to admin status.

Request Body:

```json
{
  "username": "string",
  "password": "string",
  "target_user": "string"
}
```

Response:

```json
{
  "success": boolean,
  "message": string,
  "result": boolean
}
```

#### Demote User

```
POST /api/admin/users/demote
```

Demotes an admin user to regular status.
Note: The root user cannot be demoted.

Request Body:

```json
{
  "username": "string",
  "password": "string",
  "target_user": "string"
}
```

Response:

```json
{
  "success": boolean,
  "message": string,
  "result": boolean
}
```

#### Generate License Key

```
POST /api/admin/keygen
```

Generates a new license key for a user.

Request Body:

```json
{
  "username": "string",
  "password": "string",
  "target_user": "string"
}
```

Response:

```json
{
  "success": boolean,
  "message": string,
  "result": string // License key
}
```

#### List Users

```
GET /api/admin/users/list
```

Lists all users.

Request Body:

```json
{
  "username": "string",
  "password": "string"
}
```

Response:

```json
{
  "message": string,
  "result": [
    {
      "id": number,
      "created_at": string,
      "username": string,
      "password": string,
      "license_key": string | null,
      "admin": boolean,
      "active": boolean,
      "hwid": string | null,
      "invited_by": number
    }
  ],
  "success": boolean
}
```

#### Create Invite

```
POST /api/admin/invites/create
```

Creates a new invite code.

Request Body:

```json
{
  "username": "string",
  "password": "string"
}
```

Response:

```json
{
  "success": boolean,
  "message": string,
  "result": string // Invite code
}
```

#### Invite Wave

```
POST /api/admin/invites/wave
```

Launches a wave of invite codes, given to every active user.

Request Body:

```json
{
  "username": "string",
  "password": "string"
}
```

Response:

```json
{
  "success": boolean,
  "message": string,
  "result": boolean
}
```

#### List Invites

```
GET /api/admin/invites/list
```

Lists all invite codes.

Request Body:

```json
{
  "username": "string",
  "password": "string"
}
```

Response:

```json
{
  "message": string,
  "result": [
    {
      "id": number,
      "user_id": number,
      "invite_code": string,
      "created_at": string,
      "used_by": number | null,
      "used_at": string | null,
      "active": boolean
    }
  ],
  "success": boolean
}
```

#### Delete Invite

```
DELETE /api/admin/invites/delete
```

Deletes an invite code.

Request Body:

```json
{
  "username": "string",
  "password": "string",
  "invite_code": "string"
}
```

Response:

```json
{
  "success": boolean,
  "message": string,
  "result": boolean
}
```

### Health Endpoints

#### Health Check

```
GET /api/health
```

Returns API health status.

Response:

```json
{
  "message": string,
  "result": string,
  "success": boolean
}
```

#### Version Check

```
GET /api/version
```

Returns current API version.

Response:

```json
{
  "message": string,
  "result": string,
  "success": boolean
}
```

## 🔒 Security Features

1. Password Hashing

- Bcrypt hashing with configurable cost factor
- Secure password storage and verification

2. VPN Detection

- IP intelligence check
- VPN and proxy detection
- Configurable threshold values

3. Nonce-based Request Validation

- Prevents replay attacks
- Time-based validation
- Unique nonce requirement

## 🗄️ Database Schema

### Users Table

```sql
CREATE TABLE IF NOT EXISTS users (
	id SERIAL PRIMARY KEY,
	username VARCHAR(255) UNIQUE,
	password VARCHAR(255),
	license_key VARCHAR(255) UNIQUE DEFAULT NULL,
	license_key_cooldown TIMESTAMP DEFAULT NULL,
	admin BOOLEAN DEFAULT FALSE,
	active BOOLEAN DEFAULT TRUE,
	hwid VARCHAR(255) UNIQUE DEFAULT NULL,
	hwid_cooldown TIMESTAMP DEFAULT NULL,
	created_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP,
	invited_by INT REFERENCES users(id) DEFAULT NULL
)
```

### Invites Table

```sql
CREATE TABLE invites (
    id SERIAL PRIMARY KEY,
    user_id INT REFERENCES users(id),
    invite_code VARCHAR(255) UNIQUE NOT NULL,
    created_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP,
    used_by INT REFERENCES users(id) DEFAULT NULL,
    used_at TIMESTAMP DEFAULT NULL,
    active BOOLEAN DEFAULT TRUE
)
```

## 🛡️ Default Configuration

- Default root user: "root"
- Password hashing cost: 12
- HWID reset cooldown: 7 days
- License key cooldown: 7 days
- Request time window: 15 seconds
- VPN detection enabled

## 🏗️ Building

The project includes a Makefile for building various architectures:

```bash
make build              # Build all architectures
make build_amd64        # Build for AMD64
make build_arm64        # Build for ARM64
make debug              # Build with debug symbols
make clean             # Clean build directory
```

## 🐳 Docker Support

The project includes Docker and Docker Compose configurations for easy deployment. See `docker-compose.yml` and `Dockerfile` for details.

## 📝 License
