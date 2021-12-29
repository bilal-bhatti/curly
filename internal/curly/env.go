/*
Copyright © 2021 Bilal Bhatti
*/

package curly

import (
	"io/ioutil"
	"os"
	"path"

	"github.com/pkg/errors"
	"gopkg.in/yaml.v3"
)

var efile = "env.yml"

type env struct {
	cwd  string
	f    []string
	Data interface{}
}

func Env(path string) (*env, error) {
	info, err := os.Stat(path)

	if errors.Is(err, os.ErrNotExist) {
		Tracef("env file not found, %s", path)
		return nil, err
	}

	if info.IsDir() {
		e := &env{cwd: path}

		home, err := os.UserHomeDir()
		if err != nil {
			return nil, errors.Errorf("home directory error, %v", err)
		}

		err = e.files(home, e.cwd)
		if err != nil {
			return nil, err
		}

		return e.load()
	}

	e := &env{cwd: path}
	e.f = append(e.f, path)
	return e.load()
}

func (e *env) load() (*env, error) {
	Tracef("read env files, %v", e.f)

	for i := len(e.f) - 1; i >= 0; i-- {

		yf, err := ioutil.ReadFile(e.f[i])
		if err != nil {
			return nil, err
		}

		var temp interface{}
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

func (e *env) files(home, cd string) error {
	yf := path.Join(cd, efile)
	if exists(yf) {
		e.f = append(e.f, yf)
	}

	if home == cd {
		return nil
	}

	return e.files(home, path.Dir(cd))
}

func exists(path string) bool {
	if _, err := os.Stat(path); errors.Is(err, os.ErrNotExist) {
		return false
	}

	return true
}
