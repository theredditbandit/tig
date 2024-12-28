// this file contains all the code relating to the gitRepo struct and any associated methods
package repo

import (
	"errors"
	"os"
	"path/filepath"

	"github.com/spf13/viper"
)

// GitRepo struct represents a git repo
type GitRepo struct {
	workTree string // workTree is the project repo
	gitDir   string // git dir is the .git inside the repo i.e. /path/to/repo/.git
	conf     *viper.Viper
}

// compute path under the repo's git dir and optionally create it if not exists.
//
// # Takes gitDir field implicitly i.e. /path/to/repo/.git is present
//
// eg RepoPath("config") will return /path/to/repo/.git/config
func (g *GitRepo) RepoPath(pathArgs ...string) string {
	return filepath.Join(g.gitDir, filepath.Join(pathArgs...))
}

// same as repoPath optionally create it if not exists.
//
// # Takes gitDir field implicitly i.e. /path/to/repo/.git is present
//
// eg RepoFile(True,"refs","remotes","origin","HEAD") will create .git/refs/remotes/origin
func (g *GitRepo) RepoFile(mkdir bool, pathArgs ...string) (string, error) {
	_, err := g.RepoDir(mkdir, pathArgs[:len(pathArgs)-1]...)
	if err == nil {
		return g.RepoPath(pathArgs...), nil
	} else {
		return "", err
	}
}

// return the path to the last directory in a path and optionally create it if not exists.
//
// # Takes gitDir field implicitly i.e. /path/to/repo/.git is present
func (g *GitRepo) RepoDir(mkdir bool, pathArgs ...string) (string, error) {
	path := g.RepoPath(pathArgs...)
	file, pathErr := os.Stat(path)
	if pathErr == nil { // when path exists
		if file.IsDir() {
			return path, nil
		} else {
			return "", errors.Join(ErrNotADir, errors.New(path))
		}
	}

	if mkdir {
		return path, os.MkdirAll(path, 0755)
	}

	return "", pathErr
}

func (g *GitRepo) readConfig(configPath string, force bool) error {
	g.conf.SetConfigName("config")
	g.conf.AddConfigPath(configPath)
	g.conf.SetConfigType("ini")
	if !force {
		err := g.conf.ReadInConfig()
		if err != nil {
			return err
		}
	}
	return nil
}

func (g *GitRepo) setDefaultConfig() {
	g.conf.Set("core.repositoryformatversion", 0)
	g.conf.Set("core.filemode", false)
	g.conf.Set("core.bare", false)
	g.conf.WriteConfigAs(g.RepoPath("config"))
}
