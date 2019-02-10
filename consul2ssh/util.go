package consul2ssh

import (
	"log"
	"os"
	"path/filepath"
	"os/user"
	"strings"
)

func checkErrCMD(err error) {
	if err != nil {
		log.Print(err)
		os.Exit(1)
	}
}

// GetEnvKey - Look up for environment variables and set default value if not exist.
func GetEnvKey(key, defaultVal string) string {
	if value, isSet := os.LookupEnv(key); isSet {
		return value
	}
	return defaultVal
}

func getFilePath(file string) string {
        // Expand tilde to home directory or just return the file path.
        user, _ := user.Current()
        userHome := user.HomeDir
        var filePath string
        if file == "~" {
                filePath = userHome
        } else if strings.HasPrefix(file, "~/") {
                filePath = filepath.Join(userHome, file[2:])
        } else {
                filePath = file
        }
	return filePath
}

func mergeMaps(src, dest mapInterface) {
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
