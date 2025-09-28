/*
Copyright © 2021 Bilal Bhatti
*/

package curly

import (
	"log"
	"os"
	"path"

	"github.com/pkg/errors"
	"gopkg.in/yaml.v3"
)

var efile = "env.yml"

type env struct {
	cwd  string
	f    []string
	Data any
}

func Env(path string) (*env, error) {
	info, err := os.Stat(path)

	if errors.Is(err, os.ErrNotExist) {
		Tracef("environment settings file not found, %s", path)
		return nil, err
	}

	if info.IsDir() {
		e := &env{cwd: path}

		home, err := os.UserHomeDir()
		if err != nil {
			return nil, errors.Errorf("home directory error, %v", err)
		}

		e.files(home, e.cwd)

		return e.load()
	}

	e := &env{cwd: path}
	e.f = append(e.f, path)
	return e.load()
}

func (e *env) load() (*env, error) {
	for i := len(e.f) - 1; i >= 0; i-- {
		if Verbose {
			log.Printf("* settings file, %s\n", e.f[i])
		}

		yf, err := os.ReadFile(e.f[i])
		if err != nil {
			return nil, err
		}

		var temp any
		err = yaml.Unmarshal(yf, &temp)
		if err != nil {
			return nil, err
		}

		temp = MapI2MapS(temp)

		if i == len(e.f)-1 {
			e.Data = temp
		} else {
			err = Merge(e.Data, temp)
			if err != nil {
				return nil, err
			}
		}

	}

	return e, nil
}

func (e *env) files(home, cd string) {
	yf := path.Join(cd, efile)
	if exists(yf) {
		e.f = append(e.f, yf)
	}

	if home == cd {
		return
	}

	e.files(home, path.Dir(cd))
}

func exists(path string) bool {
	if _, err := os.Stat(path); errors.Is(err, os.ErrNotExist) {
		return false
	}

	return true
}
