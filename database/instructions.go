package database

import (
	"fmt"
	"log"
	"os"
	"vexal/auth"
)

var (
	RootUser     = "root" // Default root user
	rootPassword = ""     // Default root password
)

func hashRootPassword() {
	rootPassword = os.Getenv("ROOT_USER_PASSWORD")
	hash, err := auth.PasswordManagerInstance.HashPassword(rootPassword)
	if err != nil {
		log.Fatalf("Cannot hash root password: %v", err)
	}
	rootPassword = hash
}

func generateSetupInstructions() map[int]map[string][]string {
	hashRootPassword()

	return map[int]map[string][]string{
		0: {
			"users": {
				`CREATE TABLE IF NOT EXISTS users (
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
				)`,
				`CREATE INDEX idx_users_username ON users (username)`,
				`CREATE INDEX idx_users_license_key ON users (license_key)`,
				`CREATE INDEX idx_users_activated ON users (active)`,
				`CREATE INDEX idx_users_invited_by ON users (invited_by)`,
				fmt.Sprintf(`INSERT INTO users (username, password, admin) VALUES ('%s', '%s', true)`, RootUser, rootPassword),
			},
		},
		1: {
			"invites": {
				`CREATE TABLE IF NOT EXISTS invites (
					id SERIAL PRIMARY KEY,
					user_id INT REFERENCES users(id),
					invite_code VARCHAR(255) UNIQUE NOT NULL,
					created_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP,
					used_by INT REFERENCES users(id) DEFAULT NULL,
					used_at TIMESTAMP DEFAULT NULL,
					active BOOLEAN DEFAULT TRUE
				)`,
				`CREATE INDEX idx_invites_invite_code ON invites (invite_code)`,
				`CREATE INDEX idx_invites_user_id ON invites (user_id)`,
				`CREATE INDEX idx_invites_used_by ON invites (used_by)`,
			},
		},
		2: {
			"init": {
				`CREATE TABLE IF NOT EXISTS init (
					id SERIAL PRIMARY KEY,
					setup_complete BOOLEAN DEFAULT FALSE
				)`,
				`INSERT INTO init (setup_complete) VALUES (true)`,
			},
		},
	}
}
