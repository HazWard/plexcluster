package loadbalancer

import (
	"fmt"
	"github.com/boltdb/bolt"
	"github.com/gorilla/mux"
	"log"
	"net/http"
	"os"
	"time"
)

const defaultDBPath = "/tmp/plexcluster.db"

type LoggingHandler struct {
	LoggerObj *log.Logger
}

func (l *LoggingHandler) HandlerMethod(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		t1 := time.Now()
		next.ServeHTTP(w, r)
		t2 := time.Now()
		l.LoggerObj.Printf("[%s] %q %v\n", r.Method, r.URL.String(), t2.Sub(t1))
	})
}


type Server struct {
	Database *bolt.DB
	Port int
	Log *log.Logger
}

func NewServer(port int, dbPath string, loggerObj *log.Logger) (*Server, error) {
	db, err := bolt.Open(dbPath, 0600, nil)
	if err != nil {
		return nil, err
	}
	return &Server{
		Database: db,
		Port: port,
		Log: loggerObj,
	}, nil
}

func (s *Server) TranscoderRegisterHandler(w http.ResponseWriter, r *http.Request) {

}

func (s *Server) TranscoderRemovalHandler(w http.ResponseWriter, r *http.Request) {

}

func (s *Server) TranscodeJobSubmissionHandler(w http.ResponseWriter, r *http.Request) {

}

func Run(port int) {
	logger := log.New(os.Stdout, "[PlexCluster] ", log.Ldate|log.Ltime)
	server, err := NewServer(port, defaultDBPath, logger)
	if err != nil {
		log.Fatalf("unable to create loadbalancer instance: %s", err)
	}

	loggingHandler := LoggingHandler{
		LoggerObj: logger,
	}

	router := mux.NewRouter()
	router.Use(loggingHandler.HandlerMethod)

	router.HandleFunc("/jobs", server.TranscodeJobSubmissionHandler).Methods("POST")
	router.HandleFunc("/transcoders/{id}", server.TranscoderRemovalHandler).Methods("DELETE")
	router.HandleFunc("/transcoders", server.TranscoderRegisterHandler).Methods("POST")

	logger.Printf("Listening on 0.0.0.0:%d", server.Port)
	srv := &http.Server{
		Addr:         fmt.Sprintf("0.0.0.0:%d", server.Port),
		Handler:      router,
		ReadTimeout:  15 * time.Second,
		WriteTimeout: 15 * time.Second,
		ErrorLog:     logger,
	}
	logger.Fatalln(srv.ListenAndServe())
}
