package settings

import (
	"errors"
	"flag"
	"os"
	"strings"
)

// TODO refactor
func All() (string, string, []error) {
	var errs []error

	project := *flag.String("project", "", "the gcp project id")
	instance := *flag.String("instance", "", "the bigtable instance name")
	flag.Parse()

	if project == "" {
		if e, ok := os.LookupEnv(strings.ToUpper("project")); ok {
			project = e
		} else {
			errs = append(errs, errors.New("could not find flag or env var project"))
		}
	}

	if instance == "" {
		if e, ok := os.LookupEnv(strings.ToUpper("instance")); ok {
			instance = e
		} else {
			errs = append(errs, errors.New("could not find flag or env var project"))
		}
	}

	return project, instance, errs
}
