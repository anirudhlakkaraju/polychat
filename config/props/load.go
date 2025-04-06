package props

import (
	"fmt"
	"os"
	"sync"

	"github.com/magiconair/properties"
)

var (
	onceInit        = new(sync.Once)
	implementations = make(map[string]interface{})
	propsKey        = "propsKey"
)

// Init loads properties files and initializes the properties only once.
// dirPath: Directory path for the properties files.
// Returns an error if loading fails.
func Init(dirPath string) error {
	var err error
	onceInit.Do(func() {
		if implementations[propsKey] == nil {
			props := loadProps(dirPath)
			implementations[propsKey] = props
		}
	})
	return err
}

func loadProps(dirPath string) *properties.Properties {
	files := []string{dirPath + "/app.properties"}
	if env, ok := os.LookupEnv("configEnvironment"); ok {
		files = append(files, dirPath+"/app."+env+".properties")
	}

	// Loads the app properties and current env specific properties
	// and also the fills in any env vars of the form ${ENV_VAR} from os
	p, err := properties.LoadFiles(files, properties.UTF8, false)
	if err != nil {
		fmt.Printf("Empty properties returned. Error: %v", err)
		return properties.NewProperties()
	}
	return p
}

// GetProps returns the loaded properties object
func GetProps() *properties.Properties {
	v := implementations[propsKey]
	return v.(*properties.Properties)
}
