package utils

import (
	"encoding/json"
	"io/ioutil"
	"net"
	"os"

	"github.com/sirupsen/logrus"
)

// LoadJSON loads from a json file.
func LoadJSON(filePath string, configs *[]map[string]interface{}) {
	// Open our jsonFile
	jsonArrayFile, err := os.Open(filePath)
	defer jsonArrayFile.Close()
	if err != nil {
		logrus.Infoln(err)
	}
	// read our opened json
	byteValue, _ := ioutil.ReadAll(jsonArrayFile)
	if err := json.Unmarshal(byteValue, configs); err != nil {
		logrus.Info(err.Error())
	}
	logrus.Info(configs)
}

// helper logger
func LogInfo(cmd, topic, message string, connections int) {
	logrus.WithFields(logrus.Fields{
		"cmd":              cmd,
		"topic":            topic,
		"open_connections": connections,
	}).Info(message)
}

// helper logger
func LogErr(cmd, msg string, connections int, err error) {
	if err != nil {
		logrus.WithFields(logrus.Fields{
			"cmd":              cmd,
			"msg":              msg,
			"open_connections": connections,
		}).Error(err.Error())
	}
}

// GetIP returns the ipv4 address from the server
func GetIP() string {
	addrs, err := net.InterfaceAddrs()
	if err != nil {
		if err != nil {
			logrus.Error(err)
		}
		return "error"
	}
	for _, address := range addrs {
		// check the address type and if it is not a loopback the display it
		if ipnet, ok := address.(*net.IPNet); ok && !ipnet.IP.IsLoopback() {
			if ipnet.IP.To4() != nil {
				return ipnet.IP.String()
			}
		}
	}
	panic("Unable to determine local IP address (non loopback). Exiting.")
}
