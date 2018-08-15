package consul2ssh

import (
	"fmt"
	"log"
	"os"
)

func checkErrCMD(err error) {
	if err != nil {
		log.Print(err)
		os.Exit(1)
	}
}

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
