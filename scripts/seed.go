package main

import (
	"fmt"
	"log"
	"os"

	"github.com/jmoiron/sqlx"
	"github.com/joho/godotenv"
	_ "github.com/lib/pq"
	"golang.org/x/crypto/bcrypt"
)

func main() {
	_ = godotenv.Load()

	dsn := fmt.Sprintf(
		"host=%s port=%s user=%s password=%s dbname=%s sslmode=%s",
		getEnv("DB_HOST", "localhost"),
		getEnv("DB_PORT", "5432"),
		getEnv("DB_USER", "postgres"),
		getEnv("DB_PASSWORD", "postgres"),
		getEnv("DB_NAME", "orchestralog"),
		getEnv("DB_SSLMODE", "disable"),
	)

	db, err := sqlx.Connect("postgres", dsn)
	if err != nil {
		log.Fatalf("connect: %v", err)
	}
	defer db.Close()

	seedUsers(db)
	seedClusters(db)
	log.Println("seed complete")
}

func seedUsers(db *sqlx.DB) {
	users := []struct {
		email     string
		password  string
		firstName string
		lastName  string
		role      string
		dept      string
	}{
		{"admin@orchestralog.com", "Admin1234!", "Admin", "User", "admin", "Engineering"},
		{"operator@orchestralog.com", "Oper1234!", "Operator", "User", "operator", "Data Engineering"},
		{"viewer@orchestralog.com", "View1234!", "Viewer", "User", "viewer", "Analytics"},
	}

	for _, u := range users {
		var exists bool
		db.Get(&exists, `SELECT EXISTS(SELECT 1 FROM users WHERE email=$1)`, u.email)
		if exists {
			log.Printf("user %s already exists, skipping", u.email)
			continue
		}

		hash, _ := bcrypt.GenerateFromPassword([]byte(u.password), bcrypt.DefaultCost)
		_, err := db.Exec(`
			INSERT INTO users (email, password_hash, first_name, last_name, role, department, is_active)
			VALUES ($1, $2, $3, $4, $5, $6, true)
		`, u.email, string(hash), u.firstName, u.lastName, u.role, u.dept)
		if err != nil {
			log.Printf("seed user %s: %v", u.email, err)
			continue
		}
		log.Printf("created user: %s (%s)", u.email, u.role)
	}
}

func seedClusters(db *sqlx.DB) {
	var adminID string
	if err := db.Get(&adminID, `SELECT id FROM users WHERE email='admin@orchestralog.com' LIMIT 1`); err != nil {
		log.Printf("admin user not found, skipping clusters: %v", err)
		return
	}

	clusters := []struct {
		name       string
		status     string
		region     string
		k8sVersion string
		nodes      int
	}{
		{"production-cluster-01", "healthy", "eu-west-1", "1.28.4", 12},
		{"staging-cluster-01", "healthy", "eu-central-1", "1.28.2", 6},
		{"dev-cluster-01", "warning", "us-east-1", "1.27.8", 3},
	}

	for _, c := range clusters {
		var exists bool
		db.Get(&exists, `SELECT EXISTS(SELECT 1 FROM clusters WHERE name=$1)`, c.name)
		if exists {
			log.Printf("cluster %s already exists, skipping", c.name)
			continue
		}

		_, err := db.Exec(`
			INSERT INTO clusters (name, status, region, k8s_version, nodes_total, nodes_ready, created_by)
			VALUES ($1, $2, $3, $4, $5, $5, $6)
		`, c.name, c.status, c.region, c.k8sVersion, c.nodes, adminID)
		if err != nil {
			log.Printf("seed cluster %s: %v", c.name, err)
			continue
		}
		log.Printf("created cluster: %s (%s)", c.name, c.region)
	}
}

func getEnv(key, fallback string) string {
	if v := os.Getenv(key); v != "" {
		return v
	}
	return fallback
}
