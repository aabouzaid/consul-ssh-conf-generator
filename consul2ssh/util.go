package consul2ssh

import (
	"os"
	"fmt"
	"log"
)

func GetEnvKey(key, defaultVal string) string {
	if value, isSet := os.LookupEnv(key); isSet {
		return value
	}
	return defaultVal
}

func getFilePath(filePath string) string {

	// Get working dir.
	workingDir, err := os.Getwd()
	if err != nil {
		log.Print(err)
		return ""
	}
	fullFilePath := fmt.Sprintf("%s/%s", workingDir, filePath)

	return fullFilePath
}

func mergeMaps(src, dest MapInterface) {
	for key, value := range src {
		dest[key] = value
	}
}

func setStrVal(value, defaultVal string) string {
	if value == "" {
		return defaultVal
	}
	return value
}
