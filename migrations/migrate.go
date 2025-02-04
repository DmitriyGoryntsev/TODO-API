package migrations

import (
	"log"
	"os/exec"
)

func ApplyMigrations(dbURL string) {
	cmd := exec.Command("migrate", "-path", "migrations", "-database", dbURL, "-verbose", "up")
	output, err := cmd.CombinedOutput()
	if err != nil {
		log.Fatalf("migration failed: %v\n%s", err, output)
	}

	log.Println("migrations applied successfully")
}
