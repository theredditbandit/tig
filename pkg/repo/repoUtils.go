// this file contains all the code relating to anything that's got to do with repo creation
// and not a metod to the gitRepo struct
package repo

import (
	"errors"
	"fmt"
	"os"
	"path/filepath"
	"tig/utils"

	"github.com/charmbracelet/log"
	"github.com/spf13/viper"
)

var (
	ErrNotADir        = errors.New("not a directory")
	ErrNotAGitRepo    = errors.New("not a git repo")
	ErrNoConfigFile   = errors.New("config file does not exist")
	ErrInvalidVersion = errors.New("unsupported repositoryformatversion")
	ErrDirNotEmpty    = errors.New("This directory is not empty; it should be empty")
)

// initialize a new git repo on disk
func InitRepo(path string) (*gitRepo, error) {
	repo, err := getRepo(path, true)
	if err != nil {
		log.Error(err)
		return nil, err
	}
	// check that the path doesn't exist or is in an empty dir
	workTree, pathErr := os.Stat(repo.workTree)
	if os.IsNotExist(pathErr) {
		mkdirErr := os.MkdirAll(repo.workTree, 0755)
		if mkdirErr != nil {
			log.Error(mkdirErr)
			return nil, mkdirErr
		}
	} else { // path exists
		if !workTree.IsDir() {
			return nil, errors.Join(ErrNotADir, errors.New(repo.workTree))
		}
		gitDir, pathErr := os.Stat(repo.gitDir)
		if !os.IsNotExist(pathErr) { // when .git/ exists
			if !gitDir.IsDir() {
				return nil, errors.Join(ErrNotADir, errors.New(repo.gitDir))
			}
			entries, err := os.ReadDir(repo.gitDir)
			if err != nil {
				return nil, err
			}
			if len(entries) != 0 { // when the git dir path exists while initializing the repo , repo/.git must be empty
				return nil, errors.Join(ErrDirNotEmpty, errors.New(fmt.Sprintf("repo is %s/%s/", repo.workTree, repo.gitDir)))
			}
		}
	}

	_, errB := repo.repoDir(true, "branches")
	_, errO := repo.repoDir(true, "objects")
	_, errRT := repo.repoDir(true, "refs", "tags")
	_, errRH := repo.repoDir(true, "refs", "heads")

	if err := utils.CheckErrors(errB, errO, errRT, errRH); err != nil {
		log.Error(err)
		return nil, err
	}

	// write to .git/description
	desc, err := repo.repoFile(false, "description")
	if err != nil {
		return nil, err
	}
	descFile, err := os.Create(desc)
	defer descFile.Close()
	if err != nil {
		return nil, err
	}
	descFile.WriteString("Unnamed repository; edit this file 'description' to name the repository.\n")

	// write to .git/HEAD
	head, err := repo.repoFile(false, "HEAD")
	if err != nil {
		return nil, err
	}
	headFile, err := os.Create(head)
	defer headFile.Close()
	if err != nil {
		return nil, err
	}
	headFile.WriteString("ref: refs/heads/master\n")
	repo.setDefaultConfig()
	return repo, err
}

// getRepo returns the address to a gitRepo object and an error
//
// the path argument specifies where the git repo exists on the file system
// the force argument forces certain operations required while repo creation only
func getRepo(path string, force bool) (*gitRepo, error) {
	g := gitRepo{}
	g.conf = viper.GetViper()
	g.workTree = path
	g.gitDir = filepath.Join(path, ".git")
	_, err := os.Stat(g.gitDir)
	if !(force || !os.IsNotExist(err)) {
		return nil, errors.Join(ErrNotAGitRepo, errors.New(g.gitDir))
	}
	configPath, err := g.repoFile(false, "config") // cf should be /path/to/repo/.git/config
	if err != nil && !force {
		return nil, err
	}
	_, err = os.Stat(configPath)
	if err != nil && !force { // file does not exist
		return nil, errors.Join(ErrNoConfigFile, errors.New(configPath))
	}
	// when config file does exist , read config file
	err = g.readConfig(configPath, force)
	if err != nil {
		return nil, err
	}

	if !force {
		vers := g.conf.GetInt("core.repositoryformatversion")
		if vers != 0 { // invalid version
			return nil, errors.Join(ErrInvalidVersion, errors.New(fmt.Sprint(vers)))
		}
	}
	return &g, nil
}

// returns a gitRepo object if any of the parent directory has the .git directory
//
// usually path can be set to an empty string, then it automatically takes in the current directory.
// require bool
func FindGitRepo(path string, required bool) (*gitRepo, error) {
	if path == "" {
		path = "."
	}
	absPath, err := filepath.Abs(path)
	if err != nil {
		return nil, err
	}
	gitDirPath := filepath.Join(absPath, ".git")
	gitDir, err := os.Stat(gitDirPath)
	if os.IsNotExist(err) { // cwd doesn't have git dir , check parent
		parent := filepath.Dir(absPath)
		if parent == absPath { // this happens when parent == absPath == /
			if required {
				return nil, ErrNotAGitRepo
			} else {
				return nil, nil /// ??? why
			}
		} else {
			return FindGitRepo(parent, required) // check above directories recursively
		}
	}
	if gitDir.IsDir() {
		return getRepo(absPath, false)
	}
	return nil, ErrNotAGitRepo
}
