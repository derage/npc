package npc

import (
	"fmt"
	"os"
	"path/filepath"
	"strings"
	"text/template"

	sprig "github.com/Masterminds/sprig/v3"
	githubRepo "github.com/derage/npc/pkg/repos/github"
	configFuncs "github.com/derage/npc/pkg/template_functions/config_funcs"
	"github.com/derage/npc/pkg/utils"

	"github.com/spf13/viper"
)

var logger = utils.GetLogger()

type Repo interface {
	GetRepo(location string) error
	GetTemplate(template string, location string) error
	TemplateExists(template string, location string) bool
	UpdateTemplate(template string, location string) error
}
type Config struct {
	Repo          Repo
	TemplateName  string
	RepoName      string
	BinaryViper   *viper.Viper
	TemplateViper *viper.Viper
	TemplateFuncs template.FuncMap
}

// Entrypoiny for npc, give it a viper config
// will use this to look up various variables needed
func Initialize(myViper *viper.Viper) (Config, error) {
	myConfig := Config{}
	var err error
	template := myViper.GetString("template")
	templatesplit := strings.Split(template, "/")
	repoConfig := myViper.GetStringMapString("repos." + templatesplit[0])
	switch os := repoConfig["type"]; os {
	case "github":
		myConfig.Repo, err = githubRepo.Initialize(repoConfig)
	default:
		return myConfig, fmt.Errorf("Repo type not supported yet")
	}
	myConfig.RepoName = templatesplit[0]
	myConfig.TemplateName = templatesplit[1]
	myConfig.BinaryViper = myViper

	return myConfig, err
}

func (config *Config) Bootstrap() error {
	var err error
	templateCache := config.BinaryViper.GetString("template-path") + "/" + config.RepoName
	template := config.TemplateName
	// First check if repo/template exist
	if !config.Repo.TemplateExists(template, templateCache) {
		err = config.Repo.GetTemplate(template, templateCache)
	} else if config.BinaryViper.GetBool("force-update") {
		err = config.Repo.UpdateTemplate(template, templateCache)
	}
	if err != nil {
		logger.Error(err.Error())
		return err
	}
	// Now we have the template, lets load the viper config
	err = config.LoadTemplateConfig()
	if err != nil {
		logger.Error(err.Error())
		return err
	}

	// Load templatefunctions
	err = config.LoadTemplateFunctions()
	if err != nil {
		logger.Error(err.Error())
		return err
	}

	// Once its loaded its time to boostrap
	return config.CopyTemplates()
}

func (config *Config) CopyTemplates() error {
	filePath := config.BinaryViper.GetString("template-path") + "/" + config.RepoName + "/" + config.TemplateName
	err := filepath.Walk(filePath,
		func(path string, info os.FileInfo, err error) error {
			var data interface{}
			if err != nil {
				return err
			}
			fileInfo, err := os.Stat(path)
			if err != nil {
				return err
			}
			if fileInfo.IsDir() {
				return nil
			}

			tmpl, err := template.New(path).Funcs(config.TemplateFuncs).ParseFiles(path)
			if err != nil {
				return err
			}

			file, err := os.Create(fmt.Sprintf("%s/%s", config.BinaryViper.GetString("directory"), fileInfo.Name()))
			if err != nil {
				return err
			}

			err = tmpl.Execute(file, data)
			if err != nil {
				return err
			}
			return nil
		})

	return err
}

func (config *Config) LoadTemplateConfig() error {
	templatePath := fmt.Sprintf("%s/%s/%s", config.BinaryViper.GetString("template-path"), config.RepoName, config.TemplateName)
	directory := config.BinaryViper.GetString("directory")

	config.TemplateViper = viper.New()

	// Search config in home directory with name ".npc" (without extension).
	// TODO: Figure out which one takes precidents. Directory should take precidents over teampltePath
	config.TemplateViper.AddConfigPath(templatePath)
	config.TemplateViper.AddConfigPath(directory)
	config.TemplateViper.SetConfigType("yaml")
	config.TemplateViper.SetConfigName("template-config")

	config.TemplateViper.AutomaticEnv() // read in environment variables that match

	// If a config file is found, read it in.
	if err := config.TemplateViper.ReadInConfig(); err == nil {
		fmt.Fprintln(os.Stderr, "Using config file:", config.TemplateViper.ConfigFileUsed())
	}

	return nil
}

func (config *Config) LoadTemplateFunctions() error {
	config.TemplateFuncs = sprig.TxtFuncMap()

	configFunctions := configFuncs.InitializeConfigFuncs(config.TemplateViper)

	for k, v := range configFunctions {
		config.TemplateFuncs[k] = v
	}

	return nil

}
