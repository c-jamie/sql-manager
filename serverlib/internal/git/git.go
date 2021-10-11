package git

import (
	"fmt"
	"io/ioutil"
	"os"

	"github.com/c-jamie/sql-manager/serverlib/log"
	"github.com/go-git/go-billy/v5"
	"github.com/go-git/go-billy/v5/memfs"
	"github.com/go-git/go-git/v5"
	"github.com/go-git/go-git/v5/plumbing/transport/http"
	"github.com/go-git/go-git/v5/storage/memory"
)

// Repo represents a git repository source
type Repo interface {
	// GetFile returns a given script from the Repo
	GetFile(file string) (string, error)
}

type repo struct {
	r    *git.Repository
	fs   billy.Filesystem
	Auth *http.BasicAuth
}

func New(repoLink string, username string, key string) (*repo, error) {
	var gitRepo repo

	gitRepo.fs = memfs.New()
	log.Info(repoLink, " | ", username, " | ", key)
	if repoLink == "" {
		return nil, fmt.Errorf("no git repo link to init")
	}
	r, err := git.Clone(memory.NewStorage(), gitRepo.fs, &git.CloneOptions{
		URL:      repoLink,
		Auth:     &http.BasicAuth{Username: username, Password: key},
		Progress: os.Stdout,
	})

	gitRepo.Auth = &http.BasicAuth{Username: username, Password: key}
	if err != nil {
		log.Error(err)
		return nil, fmt.Errorf("unable to init git: %w", err)
	}
	log.Info("git loaded")
	gitRepo.r = r
	return &gitRepo, nil
}

func (gt *repo) GetFile(file string) (string, error) {
	gt.update()
	log.Debug("grabbing: ", file)
	fo, err := gt.fs.Open(file)
	if err != nil {
		log.Error(fmt.Errorf("error grabbing git file %w", err))
		return "", err
	}
	buf, err := ioutil.ReadAll(fo)
	log.Debug("file is: ", string(buf))
	if err != nil {
		log.Error(fmt.Errorf("error reading file %w", err))
		return "", err
	}
	return string(buf), nil
}

func (gt *repo) update() {
	w, err := gt.r.Worktree()
	log.Debug("pulling remote")
	if err != nil {
		log.Error(fmt.Errorf("error loading worktree: %w", err))
	}
	err = w.Pull(&git.PullOptions{Force: true, RemoteName: "origin", Auth: gt.Auth})
	if err != nil {
		log.Error(fmt.Errorf("error pulling file: %w", err))
	}
}
