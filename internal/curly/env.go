package curly

import (
	"io/ioutil"
	"log"
	"os"
	"path"

	"github.com/pkg/errors"
	"gopkg.in/yaml.v2"
)

var efile = "env.yml"

type env struct {
	f    []string
	Data interface{}
}

func Env() (*env, error) {
	e := new(env)
	return e.load()
}

func (e *env) load() (*env, error) {
	log.Println("initializing environment")

	home, err := os.UserHomeDir()
	if err != nil {
		return nil, errors.Errorf("home directory error, %v", err)
	}

	cwd, err := os.Getwd()
	if err != nil {
		return nil, errors.Errorf("current working directory error, %v", err)
	}

	err = e.files(home, cwd)
	if err != nil {
		return nil, err
	}

	for i := len(e.f) - 1; i >= 0; i-- {
		Tracef("read env file, %s", e.f[i])

		yf, err := ioutil.ReadFile(e.f[i])
		if err != nil {
			return nil, err
		}

		var temp interface{}
		err = yaml.Unmarshal(yf, &temp)
		if err != nil {
			return nil, err
		}

		temp = convert(temp)

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
	Tracef("exists? %s", path)
	if _, err := os.Stat(path); errors.Is(err, os.ErrNotExist) {
		return false
	}

	return true
}
