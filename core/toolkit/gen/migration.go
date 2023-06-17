package gen

import (
	"fmt"
	"github.com/kataras/golog"
	"os"
	"strings"
	"time"
)

func GenerateMigration(name string, collection string) []byte {
	bytes, err := os.ReadFile("./toolkit/gen/templates/migration")
	if err != nil {
		golog.Fatal("[GEN] Cannot find the migration template, are you in the correct root directory of the project?")
	}
	normalizedName := strings.Replace(strings.ToLower(name), " ", "_", -1)
	nonSpacedName := strings.Replace(name, " ", "", -1) + "Migration"

	contents := string(bytes)
	contents = strings.Replace(contents, "{MigrationName}", nonSpacedName, -1)
	contents = strings.Replace(contents, "{MigrationName$}", normalizedName, -1)
	contents = strings.Replace(contents, "{CollectionName}", collection, -1)
	return []byte(contents)
}

func WriteMigration(name string, collection string) {
	n := strings.ToLower(name)
	n = strings.Replace(n, " ", "_", -1)
	fName := fmt.Sprint("./toolkit/migrations/", time.Now().Unix(), "_", n, ".go")
	err := os.WriteFile(fName, GenerateMigration(name, collection), 0666)
	if err != nil {
		golog.Fatal("[GEN] Cannot write the file ", fName, ": ", err)
	}
	golog.Info("[GEN] Created a migration under ", fName)
}
