package settings

import (
	"errors"
	"flag"
	"os"
	"strings"
)

func All() (string, string, []error) {
	var errors []error

	project, err := Get("project")
	if err != nil {
		errors = append(errors, err)
	}

	instance, err := Get("instance")
	if err != nil {
		errors = append(errors, err)
	}

	return project, instance, errors
}

func Get(name string) (string, error) {
	s := *flag.String(name, "", "the gcp project id")
	if s != "" {
		return s, nil
	}

	if e, ok := os.LookupEnv(strings.ToUpper(name)); ok {
		return e, nil
	}

	return "", errors.New("no project id found. Either set it by cli flag -project or by env var PROJECT")
}
