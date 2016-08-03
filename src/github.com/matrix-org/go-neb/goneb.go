package main

import (
	log "github.com/Sirupsen/logrus"
	"github.com/matrix-org/go-neb/auth"
	"github.com/matrix-org/go-neb/clients"
	"github.com/matrix-org/go-neb/database"
	"github.com/matrix-org/go-neb/server"
	_ "github.com/matrix-org/go-neb/services/echo"
	_ "github.com/matrix-org/go-neb/services/github"
	_ "github.com/mattn/go-sqlite3"
	"net/http"
	_ "net/http/pprof"
	"os"
)

func main() {
	bindAddress := os.Getenv("BIND_ADDRESS")
	databaseType := os.Getenv("DATABASE_TYPE")
	databaseURL := os.Getenv("DATABASE_URL")

	db, err := database.Open(databaseType, databaseURL)
	if err != nil {
		log.Panic(err)
	}

	clients := clients.New(db)
	if err := clients.Start(); err != nil {
		log.Panic(err)
	}

	auth.RegisterModules(db)

	http.Handle("/test", server.MakeJSONAPI(&heartbeatHandler{}))
	http.Handle("/admin/configureClient", server.MakeJSONAPI(&configureClientHandler{db: db, clients: clients}))
	http.Handle("/admin/configureService", server.MakeJSONAPI(&configureServiceHandler{db: db, clients: clients}))
	http.Handle("/admin/configureAuth", server.MakeJSONAPI(&configureAuthHandler{db: db}))
	wh := &webhookHandler{db: db}
	http.HandleFunc("/services/hooks/", wh.handle)

	http.ListenAndServe(bindAddress, nil)
}
