package local

import (
	"errors"
	"os"
	"reflect"

	"github.com/derage/npc/pkg/utils"
	"github.com/mitchellh/mapstructure"
)

var logger = utils.GetLogger()

func Initialize(repo map[string]string) (*Repo, error) {
	var repoStruct Repo

	err := mapstructure.Decode(repo, &repoStruct)

	if err != nil {
		return &repoStruct, err
	}

	return &repoStruct, nil
}

func (repo *Repo) GetRepo(location string) error {
	logger.Info("Repo is already on filesystem")
	return nil
}

func (repo *Repo) GetTemplate(template string, location string) error {
	logger.Info("Template is already on filesystem")
	return nil
}

func (repo *Repo) ReadProperty(property string) string {
	// Gets a property defined in repo.<REPONAME> from .npc.yaml
	value := reflect.ValueOf(repo)
	valueE := value.Elem()

	for _, name := range []string{property} {
		field := valueE.FieldByName(name)

		if field == (reflect.Value{}) {
			logger.Warnf("field %s not exist in struct", name)
			return ""

		}
	}
	return reflect.Indirect(value).FieldByName(property).String()
}

func (repo *Repo) TemplateExists(template string, location string) bool {
	// take location, add repo name from repo.Name onto it
	// then add template onto it, and make sure it all exists under
	// location/RepoName/template, return true or false
	templateLocation := location + "/" + template

	if _, err := os.Stat(templateLocation); err == nil {
		logger.Infof("Template %s exist at %s", template, templateLocation)
		return true

	} else if errors.Is(err, os.ErrNotExist) {
		logger.Infof("Template %s does not exist at %s", template, templateLocation)
		return false

	} else {
		logger.Warnf("Got an error when checking if template %s exists in location %s: %s", template, location, err.Error())
	}

	return false

}

func (repo *Repo) UpdateTemplate(template string, location string) error {
	logger.Info("Repo is configured as local, there is nothing to pull")
	return nil
}
