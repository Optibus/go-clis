package reader

import (
	"encoding/json"
	"io"
	"os"
	"strings"

	"gopkg.in/yaml.v2"
)

// The data struct for the decoded data
// Notice that all fields must be exportable!

type Step struct {
	Command string `json:"command"`
	Name    string `json:"name"`
}

type Data struct {
	ID    string `json:"id"`
	Steps []Step `json:"steps"`
}

type Workspace struct {
	Packages []string
}

func Read(path string) (Workspace, error) {
	// Open the file
	file, err := os.Open(path)
	if err != nil {
		return Workspace{}, err
	}
	defer file.Close()

	// Read the file
	byteValue, err := io.ReadAll(file)
	if err != nil {
		return Workspace{}, err
	}

	splitedStrings := strings.Split(path, ".")
	extension := splitedStrings[len(splitedStrings)-1]

	var payload Workspace
	switch extension {

	case "json":
		{
			// Decode the JSON content into the data structure
			err := json.Unmarshal(byteValue, &payload)
			if err != nil {
				return Workspace{}, err
			}

		}
	case "yaml", "yml":
		{
			err := yaml.Unmarshal(byteValue, &payload)
			if err != nil {
				return Workspace{}, err
			}
		}
	}
	return payload, nil
}
