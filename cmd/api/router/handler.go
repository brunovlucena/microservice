package router

import (
	"context"
	"encoding/json"
	"fmt"
	"io/ioutil"
	"net/http"
	"os"
	"strings"
	"time"

	. "github.com/brunovlucena/microservice/cmd/data"
	. "github.com/brunovlucena/microservice/cmd/messaging"
	"github.com/streadway/amqp"

	"github.com/go-chi/chi"
	"github.com/go-chi/render"
	"github.com/prometheus/client_golang/prometheus"
	"github.com/prometheus/client_golang/prometheus/promhttp"
	"github.com/sirupsen/logrus"
	"github.com/spf13/viper"
	"golang.org/x/net/http2"
)

const (
	appName = "apiApp"
)

var (
	histogram        *prometheus.HistogramVec
	serverAddr       string
	client           IMessagingClient
	connectionString string
)

func init() {
	// Setup Config path
	viper.SetConfigName("config.yaml")
	viper.AddConfigPath("/etc/appname/")
	viper.AddConfigPath("$HOME/.appname")
	viper.AddConfigPath(".")
	viper.AddConfigPath("../")
	err := viper.ReadInConfig()
	if err != nil {
		panic(fmt.Errorf("Fatal error config file: %s \n", err))
	}

	// Read Config
	viper.SetConfigType("yaml")
	if serverAddr = os.Getenv("SERVER_ADDR"); serverAddr == "" {
		serverAddr = viper.Get("serverAddr").(string)
	}

	// Broker
	if connectionString = os.Getenv("AMQP_ADDR"); connectionString == "" {
		connectionString = viper.Get("amqpAddr").(string)
	}
	client = &MessagingClient{}
	client.ConnectToBroker(connectionString)

	// prometheus
	histogram = prometheus.NewHistogramVec(prometheus.HistogramOpts{
		Name: "task_duration_seconds",
		Help: "Time taken to performa a task",
	}, []string{"code", "function"})
	prometheus.Register(histogram)
}

// StartWebServerHTTP starts the App
func (r *MyRouter) StartWebServerHTTP() {

	srv := &http.Server{
		Addr:         serverAddr,
		Handler:      r.Mux,
		ReadTimeout:  5 * time.Second,
		WriteTimeout: 10 * time.Second,
	}

	// Use HTTP2
	err := http2.ConfigureServer(srv, &http2.Server{})

	// log error
	if err != nil {
		logrus.WithFields(logrus.Fields{
			"cmd": "StartWebServerHTTP",
		}).Error(err.Error())
	}

	// setup
	r.setupRoutes()

	// start listening
	logrus.Infof("Starting %v on 0.0.0.0%s", appName, serverAddr)
	logrus.Fatalln(srv.ListenAndServe())
}

func (r *MyRouter) setupRoutes() {
	// add healthcheck
	r.Mux.Get("/", func(w http.ResponseWriter, r *http.Request) {
		w.Header().Set("Content-Type", "text/html")
		w.WriteHeader(http.StatusOK)
		b, err := ioutil.ReadFile("router/welcome.html")
		if err != nil {
			logrus.Error(err.Error())
		}
		welcome := string(b)
		w.Write([]byte(welcome))
	})

	// add healthcheck
	r.Mux.Get("/ping", func(w http.ResponseWriter, r *http.Request) {
		w.Header().Set("Content-Type", "text/plain")
		w.WriteHeader(http.StatusOK)
		w.Write([]byte("pong"))
	})

	// prometheus
	r.Mux.Handle("/metrics", promhttp.Handler())

	// add routes
	r.Mux.Route("/configs", func(r chi.Router) {
		r.Get("/", FindAll)
		r.Post("/", Create)
		r.Route("/{name}", func(r chi.Router) {
			r.Use(ConfigCtx) // Loads a config in the request's Context
			r.Get("/", Find)
			r.Put("/", Update)
			r.Patch("/", Update)
			r.Delete("/", Delete)
		})
	})
	r.Mux.Get("/search", Search) // GET /search?metadata.{key}={value}
}

// ConfigCtx middleware is used to load an Config object from
// the URL parameters passed through as the request. In case
// the Config could not be found.
func ConfigCtx(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		// start tracking
		start := time.Now()

		sk := strings.Split(r.URL.Path, "/")      // E.g route: /configs/foo/bar/
		configName := sk[2]                       // configName: foo
		logrus.Info("GET /configs/" + configName) //TODO: remove

		// publish to requests
		config := Config{
			Data: DataMap{"name": configName},
		}
		cr := &ConfigRequest{
			Method: r.Method,
			Path:   r.URL,
			Config: config,
		}

		jBytes, err := json.Marshal(&cr)
		if err != nil {
			logrus.WithFields(logrus.Fields{
				"cmd": "ConfigCtx",
			}).Debug(err.Error())
		}
		err = client.Publish(jBytes, "api", "fanout", "requests")
		if err != nil {
			logrus.WithFields(logrus.Fields{
				"cmd": "ConfigCtx",
			}).Debug(err.Error())
		}

		// consume from responses
		// RabbitMQ
		exchangeName := "api"
		exchangeType := "fanout"
		queueName := "responses"

		client.ConnectToBroker(connectionString)
		client.Subscribe(exchangeName, exchangeType, appName, queueName, func(d amqp.Delivery) {
			var c *Config
			err := json.Unmarshal(d.Body, &c)
			duration := time.Since(start)
			// end tracking

			// render errors
			if err != nil {
				render.Render(w, r, ErrRender(err))
				return
			} else {
				// prometheus: observe error
				code := http.StatusUnprocessableEntity
				observe(duration, code, "findall")

				// log
				logrus.WithFields(logrus.Fields{
					"cmd":      "Find",
					"duration": duration,
					"code":     code,
				}).Debug("Records found!")
			}

			// prometheus
			code := http.StatusFound
			observe(duration, code, "find")

			// save config
			ctx := context.WithValue(r.Context(), "config", c)
			next.ServeHTTP(w, r.WithContext(ctx))
		})
	})
}

// FindsAll returns all configs.
func FindAll(w http.ResponseWriter, r *http.Request) {
	// start tracking
	start := time.Now()

	// get configs
	var configs []*Config
	var err error

	duration := time.Since(start)
	// end tracking

	// render errors
	if err != nil {
		// prometheus: observe error
		code := http.StatusUnprocessableEntity
		observe(duration, code, "findall")

		// log
		logrus.WithFields(logrus.Fields{
			"cmd":      "FindAll",
			"duration": duration,
			"code":     code,
		}).Debug("Records not found!")

		// render
		render.Render(w, r, ErrRender(err))
		return
	}

	// prometheus
	code := http.StatusFound
	observe(duration, code, "findall")

	// log
	logrus.WithFields(logrus.Fields{
		"cmd":      "FindAll",
		"duration": duration,
		"code":     code,
	}).Debug("Records found!")

	// render
	render.Status(r, code)
	render.RenderList(w, r, NewConfigListResponse(configs))
}

// Create creates a new config.
func Create(w http.ResponseWriter, r *http.Request) {
	// start tracking
	start := time.Now()

	// convert post data into json format
	cr := ConfigRequest{}
	err := cr.Bind(r)

	// check error
	if err != nil {
		// return error
		render.Render(w, r, ErrRender(err))
		return
	}

	// publish to request

	// consume from responses
	c := &Config{}

	duration := time.Since(start)
	// end tracking

	// render errors
	if err != nil {
		// prometheus: observe error
		code := http.StatusUnprocessableEntity
		observe(duration, code, "create")

		// render
		render.Render(w, r, ErrRender(err))
		return
	}

	// prometheus
	code := http.StatusCreated
	observe(duration, code, "create")

	// log
	logrus.WithFields(logrus.Fields{
		"cmd":      "Create",
		"duration": duration,
		"code":     code,
	}).Info("Record created!")

	// render
	render.Status(r, code)
	render.Render(w, r, NewConfigResponse(c))
}

// Find returns the specified config.
func Find(w http.ResponseWriter, r *http.Request) {
	// get from context
	config := r.Context().Value("config").(*Config)

	// render
	render.Status(r, http.StatusFound)
	render.Render(w, r, NewConfigResponse(config))
}

// Update updates the specified Config.
func Update(w http.ResponseWriter, r *http.Request) {
	// start tracking
	start := time.Now()

	// convert post data into json format
	cr := ConfigRequest{}
	err := cr.Bind(r)

	// check error
	if err != nil {
		// return error
		render.Render(w, r, ErrRender(err))
		return
	}

	// publish to request

	// consume from responses
	c := &Config{}

	duration := time.Since(start)
	// end tracking

	// render errors
	if err != nil {
		// prometheus: observe error
		code := http.StatusUnprocessableEntity
		observe(duration, code, "update")

		// render
		render.Render(w, r, ErrRender(err))
		return
	}

	// prometheus
	code := http.StatusFound
	observe(duration, code, "update")

	// log
	logrus.WithFields(logrus.Fields{
		"cmd":      "Update",
		"duration": duration,
		"code":     code,
	}).Info("Record updated!")

	// render
	render.Status(r, code)
	render.Render(w, r, NewConfigResponse(c))
}

// Delete removes the specified Config.
func Delete(w http.ResponseWriter, r *http.Request) {
	// start tracking
	start := time.Now()

	// get config from context
	config := r.Context().Value("config").(*Config)

	// removes from database
	var err error

	duration := time.Since(start)
	// end tracking

	// render errors
	if err != nil {
		// prometheus: observe error
		code := http.StatusUnprocessableEntity
		observe(duration, code, "delete")

		// render
		render.Render(w, r, ErrRender(err))
		return
	}

	// prometheus
	code := http.StatusFound
	observe(duration, code, "delete")

	// log
	logrus.WithFields(logrus.Fields{
		"cmd":      "Delete",
		"duration": duration,
		"code":     code,
	}).Info("Record removed!")

	// render
	render.Status(r, code)
	render.Render(w, r, NewConfigResponse(config))
}

// Search returns the Configs data for a matching config.
// Query:   /metadata.{key}={value}
func Search(w http.ResponseWriter, r *http.Request) {
	// start tracking
	start := time.Now()
	// get params
	//params, _ := url.ParseQuery(r.URL.String())

	// get configs
	var configs []*Config
	var err error

	duration := time.Since(start)
	// end tracking

	// check errors
	if err != nil {
		// prometheus: observe error
		code := http.StatusUnprocessableEntity
		observe(duration, code, "search")

		// render error
		render.Render(w, r, ErrRender(err))
		return
	}

	// prometheus
	code := http.StatusFound
	observe(duration, code, "search")

	// log
	logrus.WithFields(logrus.Fields{
		"cmd":      "Search",
		"duration": duration,
		"code":     code,
	}).Info("Record(s) found!")

	// render result
	render.RenderList(w, r, NewConfigListResponse(configs))
}

// Helper Funcion for Prometheus
func observe(duration time.Duration, code int, task string) {
	histogram.WithLabelValues(fmt.Sprintf("%d", code), task).Observe(duration.Seconds())
}
