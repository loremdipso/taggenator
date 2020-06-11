package main

import (
	"io"
	"io/ioutil"
	"os"
	"path/filepath"
	"strings"
)

// Reads all .json files in the current folder
// and encodes them as strings literals in jsonfiles.go
func main() {
	basedir := "settings"
	fs, _ := ioutil.ReadDir(filepath.Join(".", basedir))
	out, _ := os.Create(filepath.Join(".", "jsonfiles.go"))
	out.Write([]byte("package main \n\nconst (\n"))
	for _, f := range fs {
		if strings.HasSuffix(f.Name(), ".json") {
			out.Write([]byte(getVariableName(f.Name()) + " = `"))

			f, _ := os.Open(filepath.Join(".", basedir, f.Name()))
			io.Copy(out, f)

			out.Write([]byte("`\n"))
		}
	}
	out.Write([]byte(")\n"))
}

func getVariableName(name string) string {
	variableName := strings.TrimSuffix(name, ".json")
	variableName = strings.ReplaceAll(variableName, "_", " ")
	variableName = strings.Title(variableName)
	variableName = strings.ReplaceAll(variableName, " ", "")
	return variableName
}
