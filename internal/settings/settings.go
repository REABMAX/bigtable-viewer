package settings

import (
	"errors"
	"flag"
	"fmt"
	"os"
	"strings"
)

const (
	PROJECT_FLAG  = "project"
	INSTANCE_FLAG = "instance"
)

// TODO refactor
func All() (string, string, []error) {
	var errs []error

	projectFlag := flag.String(PROJECT_FLAG, "", "the gcp project id")
	instanceFlag := flag.String(INSTANCE_FLAG, "", "the bigtable instance name")
	flag.Parse()

	project := *projectFlag
	instance := *instanceFlag

	if project == "" {
		if e, ok := os.LookupEnv(strings.ToUpper(PROJECT_FLAG)); ok {
			project = e
		} else {
			errs = append(errs, errors.New(fmt.Sprintf("could not find flag or env var %s", PROJECT_FLAG)))
		}
	}

	if instance == "" {
		if e, ok := os.LookupEnv(strings.ToUpper(INSTANCE_FLAG)); ok {
			instance = e
		} else {
			errs = append(errs, errors.New(fmt.Sprintf("could not find flag or env var %s", INSTANCE_FLAG)))
		}
	}

	return project, instance, errs
}
