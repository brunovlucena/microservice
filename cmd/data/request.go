package data

import (
	"encoding/json"
	"errors"
	"net/http"
	"net/url"

	"github.com/sirupsen/logrus"
)

// ConfigRequest represents a Request
type ConfigRequest struct {
	Method string
	Path   *url.URL
	Config Config
}

// ConfigRequest implements Binder interface
func (cr *ConfigRequest) Bind(r *http.Request) error {
	var data map[string]interface{}
	err := json.NewDecoder(r.Body).Decode(&data)
	if err != nil {
		logrus.WithFields(logrus.Fields{
			"cmd": "bind",
		}).Error(err.Error())
		return err
	}
	// validate json
	if data == nil { // edge case where json payload is null
		err := errors.New("Null is not valid!")
		logrus.WithFields(logrus.Fields{
			"cmd": "bind",
		}).Error(err.Error())
		return err
	}
	cr.Method = r.Method
	cr.Path = r.URL
	cr.Config.Data = data
	return nil
}
