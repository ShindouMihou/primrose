package gen

import (
	"fmt"
	"github.com/kataras/golog"
	"os"
	"strings"
	"time"
)

func GenerateCollection(name string) []byte {
	bytes, err := os.ReadFile("./toolkit/gen/templates/collection")
	if err != nil {
		golog.Fatal("[GEN] Cannot find the collection template, are you in the correct root directory of the project?")
	}
	normalizedName := strings.Replace(strings.ToLower(name), " ", "_", -1)
	nonSpacedName := strings.Replace(name, " ", "", -1) + "Migration"

	contents := string(bytes)
	contents = strings.Replace(contents, "{CollectionName}", nonSpacedName, -1)
	contents = strings.Replace(contents, "{CollectionName$}", normalizedName, -1)
	return []byte(contents)
}

func WriteCollection(name string) {
	n := strings.ToLower(name)
	n = strings.Replace(n, " ", "_", -1)
	fName := fmt.Sprint("./toolkit/collections/", time.Now().Unix(), "_", n, ".go")
	err := os.WriteFile(fName, GenerateCollection(name), 0666)
	if err != nil {
		golog.Fatal("[GEN] Cannot write the file ", fName, ": ", err)
	}
	golog.Info("[GEN] Created an collection migration under ", fName)
}
