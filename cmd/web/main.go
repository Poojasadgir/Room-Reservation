package main

import (
	"encoding/gob"
	"flag"
	"fmt"
	"log"
	"net/http"
	"os"
	"time"

	"github.com/Poojasadgir/room-reservation/internal/config"
	"github.com/Poojasadgir/room-reservation/internal/driver"
	"github.com/Poojasadgir/room-reservation/internal/handlers"
	"github.com/Poojasadgir/room-reservation/internal/helpers"
	"github.com/Poojasadgir/room-reservation/internal/models"
	"github.com/Poojasadgir/room-reservation/internal/render"
	"github.com/alexedwards/scs/v2"
)

var portNumber = ":1023"

var app config.AppConfig
var session *scs.SessionManager
var infoLog, errorLog *log.Logger

// main is the entry point of the web application.
// It initializes the database connection, starts the mail listener,
// and sets up the HTTP server to listen on the specified port.
func main() {
	db, err := run()
	if err != nil {
		log.Fatal(err)
	}
	defer db.SQL.Close()
	defer close(app.MailChannel)
	listenForMail()

	fmt.Println("Starting mail listener...")
	fmt.Printf("Starting port number on port %s\n", portNumber)

	srv := &http.Server{
		Addr:    portNumber,
		Handler: routes(&app),
	}
	err = srv.ListenAndServe()
	log.Fatal(err)
}

// run function initializes the application by setting up the session, logging, database connection, and other configurations.
// It also reads command line flags to determine the application's environment and settings.
// It returns a pointer to the database and an error if any.
func run() (*driver.DB, error) {
	// What is going to be put into the session
	gob.Register(models.Reservation{})
	gob.Register(models.User{})
	gob.Register(models.Room{})
	gob.Register(models.Restriction{})
	gob.Register(map[string]int{})

	// read cmd flags
	inProduction := flag.Bool("production", true, "Application is in production")
	useCache := flag.Bool("cache", true, "Use template cache")
	dbHost := flag.String("dbhost", "localhost", "Database host")
	dbName := flag.String("dbname", "postgres", "Database name")
	dbUser := flag.String("dbuser", "postgres", "Database user")
	dbPass := flag.String("dbpass", "Pooja@2706", "Database password")
	dbPort := flag.Int("dbport", 5432, "Database port")
	dbSSL := flag.String("dbssl", "disable", "Database SSL settings (disable, prefer, require)")

	flag.Parse()

	if *dbName == "" || *dbUser == "" {
		fmt.Println("missing required flags")
		os.Exit(1)
	}

	if *inProduction {
		portNumber = ":8080"
	}

	// Channel for sending and receiving mail
	mailChannel := make(chan models.MailData)
	app.MailChannel = mailChannel

	// Change this to true when in production
	app.InProduction = *inProduction
	app.UseCache = *useCache

	// Logging and error handling
	infoLog = log.New(os.Stdout, "INFO\t", log.Ldate|log.Ltime)
	app.InfoLog = infoLog
	errorLog = log.New(os.Stdout, "ERROR\t", log.Ldate|log.Ltime|log.Lshortfile)
	app.ErrorLog = errorLog

	// Session info
	session = scs.New()
	session.Lifetime = 24 * time.Hour
	session.Cookie.Persist = true
	session.Cookie.SameSite = http.SameSiteLaxMode
	session.Cookie.Secure = app.InProduction
	app.Session = session

	// Connect to database
	log.Println("Connecting to database...")
	connestionString := fmt.Sprintf("host=%s port=%d dbname=%s user=%s password=%s sslmode=%s", *dbHost, *dbPort, *dbName, *dbUser, *dbPass, *dbSSL)
	db, err := driver.ConnectSQL(connestionString)
	if err != nil {
		log.Fatal("cannot connect to database.")
	}

	tc, err := render.CreateTemplateCache()
	if err != nil {
		log.Fatal("cannot create template cache")
		return nil, err
	}
	app.TemplateCache = tc

	repo := handlers.NewRepo(&app, db)
	handlers.NewHandlers(repo)
	render.NewRenderer(&app)
	helpers.NewHelpers(&app)

	return db, nil
}
