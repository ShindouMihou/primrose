package gen

import (
	"fmt"
	"github.com/kataras/golog"
	"os"
	"strings"
	"time"
)

func GenerateIndex(name string, collection string) []byte {
	bytes, err := os.ReadFile("./toolkit/gen/templates/index")
	if err != nil {
		golog.Fatal("[GEN] Cannot find the index template, are you in the correct root directory of the project?")
	}
	normalizedName := strings.Replace(strings.ToLower(name), " ", "_", -1)
	nonSpacedName := strings.Replace(name, " ", "", -1) + "Migration"

	contents := string(bytes)
	contents = strings.Replace(contents, "{IndexName}", nonSpacedName, -1)
	contents = strings.Replace(contents, "{IndexName$}", normalizedName, -1)
	contents = strings.Replace(contents, "{CollectionName}", collection, -1)
	return []byte(contents)
}

func WriteIndex(name string, collection string) {
	n := strings.ToLower(name)
	n = strings.Replace(n, " ", "_", -1)
	fName := fmt.Sprint("./toolkit/indexes/", time.Now().Unix(), "_", n, ".go")
	err := os.WriteFile(fName, GenerateIndex(name, collection), 0666)
	if err != nil {
		golog.Fatal("[GEN] Cannot write the file ", fName, ": ", err)
	}
	golog.Info("[GEN] Created an index under ", fName)
}
