package path

import (
	"path/filepath"
	"strings"

	homedir "github.com/mitchellh/go-homedir"
)

func Abs(path string) (string, error) {
	if strings.HasPrefix(path, "~") {
		home, err := homedir.Dir()
		if err != nil {
			return "", err
		}
		path = filepath.Join(home, strings.TrimPrefix(path, "~"))
	}
	return filepath.Abs(path)
}
