package postgres

import (
	"database/sql"
	"fmt"
	"net/url"
	"strings"

	// import driver for database/sql
	_ "github.com/lib/pq"

	. "github.com/brunovlucena/microservice/cmd/data"
	. "github.com/brunovlucena/microservice/cmd/utils"
	"github.com/sirupsen/logrus"
)

const (
	appName = "repositoryApp"
)

// Represents a Connection.
type Postgres struct {
	host     string
	port     string
	user     string
	password string
	dbname   string
	dbconn   *sql.DB
}

// NewPostgres creates a connection with the database.
func NewPostgres(host, port, user, password, dbname string) (*Postgres, error) {
	dbconn, err := connect(host, port, user, password, dbname)
	return &Postgres{host, port, user, password, dbname, dbconn}, err
}

// Connect stablishes a connection with the database
func connect(host, port, user, password, dbname string) (*sql.DB, error) {
	// info: sslmode - disabled
	psqlInfo := fmt.Sprintf("host=%s port=%s user=%s password=%s dbname=%s sslmode=disable", host, port, user, password, dbname)

	// opens connection
	db, err := sql.Open("postgres", psqlInfo)

	// error checking
	if err != nil {
		logrus.WithFields(logrus.Fields{
			"cmd":      "connect",
			"host":     host,
			"port":     port,
			"database": dbname,
		}).Error("Cannot connect!")

		return nil, err
	} else {
		// setup for highter performance
		db.SetMaxOpenConns(50)
		db.SetMaxIdleConns(20)
		// SQL.Open only creates the DB object, but dies not open
		//any connections to the database.
		err = db.Ping()
		if err != nil {
			logrus.WithFields(logrus.Fields{
				"cmd":      "connect",
				"host":     host,
				"port":     port,
				"database": dbname,
			}).Error("Failed to ping database!")

			return nil, err
		}
	}

	// Success
	logrus.WithFields(logrus.Fields{
		"cmd":             "connect",
		"host":            host,
		"port":            port,
		"database":        dbname,
		"max_connections": db.Stats().MaxOpenConnections,
	}).Infoln("Successfully connected to Postgres!")

	return db, nil
}

// Create creates a Record in the database.
func (p *Postgres) Create(config *Config) (*Config, error) {
	sqlStatement := `INSERT INTO configs (data) VALUES ($1) RETURNING id`
	var id int
	dataMap := config.Data

	// number of open connections
	openConn := p.dbconn.Stats().OpenConnections

	// query
	row := p.dbconn.QueryRow(sqlStatement, dataMap)
	err := row.Scan(&id)
	if err != nil {
		LogErr("Create", sqlStatement, openConn, err)
		return nil, err
	}

	// log
	LogInfo("Create", sqlStatement, "Record created!", openConn)
	logrus.WithFields(logrus.Fields{
		"cmd": "Create",
	}).Infoln("New record ID is", id)

	return config, nil
}

// Find finds a Record in the database.
func (p *Postgres) Find(name string) (*Config, error) {
	sqlStatement := `SELECT data FROM configs WHERE data->>'name' = $1;`
	row := p.dbconn.QueryRow(sqlStatement, name)
	config := new(Config)

	// number of open connections
	openConn := p.dbconn.Stats().OpenConnections

	var err error
	switch err = row.Scan(&config.Data); err {
	case sql.ErrNoRows:
		LogInfo("Find", sqlStatement, "No Record found!", openConn)
	case nil:
		LogInfo("Find", sqlStatement, "Record found!", openConn)
	}

	return config, err
}

// Update updates a Record in the database.
func (p *Postgres) Update(config *Config) (*Config, error) {
	sqlStatement := `UPDATE configs SET data = $1 WHERE data->>'name' = $2;`
	dataMap := config.Data

	// number of open connections
	openConn := p.dbconn.Stats().OpenConnections

	// query
	_, err := p.dbconn.Exec(sqlStatement, dataMap, dataMap["name"])
	if err != nil {
		LogErr("Update", sqlStatement, openConn, err)
		return nil, err
	}

	// log
	LogInfo("Update", sqlStatement, "Record updated!", openConn)

	return config, nil
}

// Remove removes a Record from database.
func (p *Postgres) Remove(name string) (*Config, error) {
	sqlStatement := `DELETE FROM configs WHERE data->>'name' = $1;`
	config := &Config{}

	// number of open connections
	openConn := p.dbconn.Stats().OpenConnections

	// query
	row := p.dbconn.QueryRow(sqlStatement, name)
	if row != nil {
		err := row.Scan(&config.Data)
		if err != nil {
			if err.Error() != "sql: no rows in result set" {
				LogErr("Remove", sqlStatement, openConn, err)
				return nil, err
			}
		}
	}

	// log
	LogInfo("Remove", sqlStatement, "Record removed!", openConn)

	return config, nil
}

// FindAll returns all Records in the database.
func (p *Postgres) FindAll() ([]*Config, error) {
	// return value
	var configs []*Config

	// query
	sqlStatement := `SELECT data FROM configs;`
	rows, err := p.dbconn.Query(sqlStatement)
	if rows != nil {
		defer rows.Close()
	}

	// number of open connections
	openConn := p.dbconn.Stats().OpenConnections

	// check for errors
	if err != nil {
		LogErr("FindAll", sqlStatement, openConn, err)
		return nil, err
	}

	// Loop through the data
	for rows.Next() {
		var m DataMap
		err := rows.Scan(&m)
		// check for errors
		if err != nil {
			LogErr("FindAll", sqlStatement, openConn, err)
			return nil, err
		}
		// append results
		configs = append(configs, &Config{Data: m})
	}

	// log
	LogInfo("FindAll", sqlStatement, "Record(s) found!", openConn)

	return configs, err
}

// Search searchs a records in the Database.
// Query: 	/search?metadata.{key}={value}
func (p *Postgres) Search(params url.Values) ([]*Config, error) {
	var configs []*Config
	// parse params: map[/search?metadata.limits.cpu.value:[120m]]
	// result: `SELECT FROM configs WHERE data->>'name' = $1;`
	var sqlStatement strings.Builder
	sqlStatement.WriteString(`SELECT data FROM configs WHERE data->`)

	// build statement
	var val string
	for k, v := range params {
		// k = /search?metadata.limits.cpu.value&&name->>''
		// v = 120m
		sk := strings.Split(string(k), "?")
		si := strings.Split(sk[1], ".")
		// [metadata, limits, cpu, value]
		for j, i := range si {
			if j == len(si)-1 {
				sqlStatement.WriteString("'" + i + "'" + "=")
			} else if j == len(si)-2 {
				sqlStatement.WriteString("'" + i + "'" + "->>")
			} else {
				sqlStatement.WriteString("'" + i + "'" + "->")
			}
		}
		val = strings.Join(v, "")
		sqlStatement.WriteString("'" + val + "'" + ";")
	}

	// log
	openConn := p.dbconn.Stats().OpenConnections
	LogInfo("search", sqlStatement.String(), "checking database", openConn)

	// query database
	rows, err := p.dbconn.Query(sqlStatement.String())
	if rows != nil {
		defer rows.Close()
	}

	// check for errors
	if err != nil {
		LogErr("search", sqlStatement.String(), openConn, err)
		return nil, err
	}

	// Loop through the data
	for rows.Next() {
		var m DataMap
		err := rows.Scan(&m)
		// check for errors
		if err != nil {
			LogErr("search", "check rows", openConn, err)
			return nil, err
		}

		// append results
		configs = append(configs, &Config{Data: m})
	}

	return configs, nil
}
