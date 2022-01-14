package unmarshal

import (
	"encoding/json"
	"fmt"
	"io/ioutil"
	"os"
	"path"
	"strings"

	"gopkg.in/yaml.v3"
)

type DataFormat uint8

const (
	DataFormatJSON DataFormat = iota
	DataFormatYAML
)

// FromFlag unmarshals a command line flag into an argument;
// if the command line argument starts with a '@' it is assumed to
// be a file on the loca filesystem, it is read into memory and then
// unmarshalled into the object struct, which must be appropriately
// annotated; if it does not start with '@', it is assumed to be an
// inline JSON representation ans is unmarshalled as such.
func FromFlag(value string, object interface{}) error {
	format := DataFormatJSON
	var content []byte
	if strings.HasPrefix(value, "@") {
		filename := strings.TrimPrefix(value, "@")
		info, err := os.Stat(filename)
		if os.IsNotExist(err) {
			return err
		}
		if info.IsDir() {
			return fmt.Errorf("%s is not a file", filename)
		}
		content, err = ioutil.ReadFile(filename)
		if err != nil {
			return err
		}
		switch path.Ext(filename) {
		case ".yaml", ".yml":
			format = DataFormatYAML
		case ".json":
			format = DataFormatJSON
		default:
			return fmt.Errorf("unsupported data format: %s", path.Ext(filename))
		}
	} else {
		content = []byte(value)
	}
	var err error

	if format == DataFormatJSON {
		err = json.Unmarshal(content, object)
	} else if format == DataFormatYAML {
		err = yaml.Unmarshal(content, object)
	}
	return err
}
