package main

import (
	"fmt"
	"log"
	"os"
	"path/filepath"
	"sort"
	"strings"

	"github.com/jmoiron/sqlx"
	"github.com/joho/godotenv"
	_ "github.com/lib/pq"
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

	direction := "up"
	if len(os.Args) > 1 {
		direction = os.Args[1]
	}

	pattern := fmt.Sprintf("migrations/*.%s.sql", direction)
	files, err := filepath.Glob(pattern)
	if err != nil {
		log.Fatalf("glob: %v", err)
	}

	sort.Strings(files)
	if direction == "down" {
		sort.Sort(sort.Reverse(sort.StringSlice(files)))
	}

	for _, f := range files {
		content, err := os.ReadFile(f)
		if err != nil {
			log.Fatalf("read %s: %v", f, err)
		}
		log.Printf("applying %s", f)
		if _, err := db.Exec(string(content)); err != nil {
			if strings.Contains(err.Error(), "already exists") {
				log.Printf("  skip (already exists)")
				continue
			}
			log.Fatalf("exec %s: %v", f, err)
		}
		log.Printf("  ok")
	}
	log.Println("migration done")
}

func getEnv(key, fallback string) string {
	if v := os.Getenv(key); v != "" {
		return v
	}
	return fallback
}
