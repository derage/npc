package github

import (
	"errors"
	"fmt"
	"os"

	"github.com/derage/npc/pkg/utils"
	"github.com/go-git/go-git/plumbing/transport"
	git "github.com/go-git/go-git/v5"
	"github.com/go-git/go-git/v5/plumbing"
	"github.com/go-git/go-git/v5/plumbing/transport/http"
	"github.com/go-git/go-git/v5/plumbing/transport/ssh"
	"github.com/mitchellh/mapstructure"
)

// In here we grab template files from github
var logger = utils.GetLogger()

func Initialize(repo map[string]string) (*Repo, error) {
	var repoStruct Repo
	err := mapstructure.Decode(repo, &repoStruct)
	if err != nil {
		return &repoStruct, err
	}
	repoStruct.GitUrl = fmt.Sprintf("https://github.com/%s/%s", repoStruct.Owner, repoStruct.Repo)
	return &repoStruct, nil
}

func (repo *Repo) GetTemplate(template string, location string) error {
	// Get template. We probably wont use template variable here, but
	// for some repos we will need the template name, so just doing it
	// for hte interface. Just call getrepo here
	if _, err := os.Stat(location); err == nil {
		logger.Infof("Repo already checked out at %s", location)
	} else if errors.Is(err, os.ErrNotExist) {
		logger.Infof("Repo not checkedout yet to location %s, checking out now", location)
		return repo.GetRepo(location)
	} else {
		logger.Warnf("Got an error when checking if template %s exists in location %s: %s", template, location, err.Error())
		return err
	}
	return nil

}

func (repo *Repo) UpdateTemplate(template string, location string) error {
	//Another method where we just need location, but other repo types
	// might pull just the template and not full repo so this is just to
	// satisfy interface
	return repo.PullLatest(location)
}

func (repo *Repo) TemplateExists(template string, location string) bool {
	//take location, add repo name from repo.Name onto it
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

// repoHttpUrl in form of https://github.com/go-git/go-billy
func (repo *Repo) GetRepo(location string) error {
	var auth transport.AuthMethod
	branch := repo.Branch
	repoHttpUrl := repo.GitUrl

	logger.Infof("Checking out repo %s branch %s to location %s", repoHttpUrl, branch, location)

	token, ok := os.LookupEnv("GITHUB_TOKEN")
	if ok {
		auth = authWithToken(token)
	} else {
		logger.Infof("variable GITHUB_TOKEN not set")
	}

	if auth == nil {
		sshKey, ok := os.LookupEnv("GITHUB_SSH_KEY")
		if ok {
			var err error
			auth, err = authWithSsh(sshKey)
			if err != nil {
				logger.Warnf("error reading ssh key %s, continuing with no auth: %s", sshKey, err.Error())
			}
		} else {
			logger.Warn("variable GITHUB_SSH_KEY not set ")
		}
	}

	if auth == nil {
		logger.Warnf("No auth method found for git checkout; continuing with no auth")
	}

	_, err := git.PlainClone(location, false, &git.CloneOptions{
		Auth:          auth,
		URL:           repoHttpUrl,
		Progress:      os.Stdout,
		ReferenceName: plumbing.ReferenceName("refs/heads/" + branch),
	})
	if err != nil {
		logger.Error(err.Error())
		return err
	}
	return nil
}

func authWithToken(token string) *http.BasicAuth {
	return &http.BasicAuth{
		Username: "npc_login", // yes, this can be anything except an empty string
		Password: token,
	}
}

func authWithSsh(privateKeyFile string) (*ssh.PublicKeys, error) {
	_, err := os.Stat(privateKeyFile)
	if err != nil {
		logger.Warnf("read file %s failed %s\n", privateKeyFile, err.Error())
		return nil, err
	}

	// Clone the given repository to the given directory
	publicKeys, err := ssh.NewPublicKeysFromFile("git", privateKeyFile, "")
	if err != nil {
		logger.Warnf("generate publickeys failed: %s\n", err.Error())
	}
	return publicKeys, err
}

func (repo *Repo) PullLatest(location string) error {
	branch := repo.Branch
	r, err := git.PlainOpen(location)
	if err != nil {
		logger.Warnf("error pulling latest; unable to open location as git repo: %s", err.Error())
		return err
	}

	// Get the working directory for the repository
	w, err := r.Worktree()
	if err != nil {
		logger.Warnf("error pulling latest; unable to get working tree: %s", err.Error())
		return err
	}

	err = w.Pull(&git.PullOptions{ReferenceName: plumbing.ReferenceName("refs/heads/" + branch), RemoteName: "origin"})
	if err != nil {
		logger.Warnf("error pulling latest; unable to pull from origin: %s", err.Error())
		return err
	}

	// Print the latest commit that was just pulled
	ref, err := r.Head()
	if err != nil {
		logger.Warnf("error pulling latest; unable get head commit: %s", err.Error())
		return err
	}

	commit, err := r.CommitObject(ref.Hash())
	if err != nil {
		logger.Warnf("error pulling latest; unable to get commit reference: %s", err.Error())
		return err
	}

	logger.Infof("Pulled latest commit of branch %s, commit %s", branch, commit)
	return nil
}
