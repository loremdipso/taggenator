package processor

import (
	"fmt"
	"internal/database"
)

func fix(self *QueryProcessor, args []string, db *database.Database) error {
	fmt.Println("# Fixing everything...")
	return db.FixAllEntries()
}
