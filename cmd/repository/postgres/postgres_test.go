package postgres

import (
	"fmt"
	"os"
	"testing"

	. "github.com/brunovlucena/microservice/cmd/data"
	. "github.com/brunovlucena/microservice/cmd/utils"
	"github.com/stretchr/testify/assert"

	"github.com/sirupsen/logrus"
	"github.com/spf13/viper"
)

var (
	configs []map[string]interface{}
	repo    Repository
	dType   string
	dHost   string
	dPort   string
	dUser   string
	dPass   string
	dbName  string
)

func init() {
	// load json
	LoadJSON("postgres_test.json", &configs)

	// Setup Config path
	viper.SetConfigName("config.yaml")
	viper.AddConfigPath("../")
	err := viper.ReadInConfig()
	if err != nil {
		panic(fmt.Errorf("Fatal error config file: %s \n", err))
	}

	// Read Config
	viper.SetConfigType("yaml")
	if dType = os.Getenv("DATABASE_TYPE"); dType == "" {
		dType = viper.Get("dType").(string)
	}
	if dHost = os.Getenv("DATABASE_HOST"); dHost == "" {
		dHost = viper.Get("dHost").(string)
	}
	if dPort = os.Getenv("DATABASE_PORT"); dPort == "" {
		dPort = viper.Get("dPort").(string)
	}
	if dUser = os.Getenv("DATABASE_USER"); dUser == "" {
		dUser = viper.Get("dUser").(string)
	}
	if dPass = os.Getenv("DATABASE_PASS"); dPass == "" {
		dPass = viper.Get("dPass").(string)
	}
	if dbName = os.Getenv("DATABASE_NAME"); dbName == "" {
		dbName = viper.Get("dbName").(string)
	}

	logrus.WithFields(logrus.Fields{
		"dType": dType,
		"dHost": dHost,
		"dPort": dPort,
		"dUser": dUser,
		"dPass": dPass,
	}).Debug("Parameters Loaded")

	// initialize repository
	repo, err = NewRepository(dType, dHost, dPort, dUser, dPass, dbName)
	if err != nil {
		logrus.Info(err.Error())
	}

	logrus.Info("Successfully Loaded tests")
}

func TestCreate(t *testing.T) {
	for _, c := range configs {
		// create configs
		repo.Create(&Config{Data: c})
	}
}

func TestUpdate(t *testing.T) {
	// new changed metadata
	newData := map[string]interface{}{
		"name": "pod-2p",
		"metadata": map[string]interface{}{
			"monitoring": map[string]interface{}{
				"enabled": "true",
			},
		},
	}
	// update pod-2
	repo.Update(&Config{Data: newData})
}

func TestFind(t *testing.T) {
	// find pod-2
	config, _ := repo.Find("pod-2p")
	// compare metadata
	data := config.Data
	metadata := data["metadata"].(map[string]interface{})
	monitoring := metadata["monitoring"].(map[string]interface{})
	enabled := monitoring["enabled"].(string)
	// it becomes true because TestUpdate
	assert.Equal(t, "true", enabled)
}

func TestRemove(t *testing.T) {
	// remove all inserted
	for _, c := range configs {
		// create configs
		repo.Remove(c["name"].(string))
	}
	// remove pod-13-idonotexist
	_, err := repo.Remove("pod-13-idonotexit")
	if err != nil {
		assert.Equal(t, "sql: no rows in result set", err.Error())
	}
}

func TestFindAll(t *testing.T) {
	configs, _ := repo.FindAll()
	// compare metadata
	config := configs[0]
	data := config.Data
	metadata := data["metadata"].(map[string]interface{})
	monitoring := metadata["monitoring"].(map[string]interface{})
	enabled := monitoring["enabled"].(bool)
	// it becomes true because TestUpdate
	assert.Equal(t, false, enabled)
}

func TestSearch(t *testing.T) {

}
