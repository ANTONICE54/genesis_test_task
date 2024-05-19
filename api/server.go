package api

import (
	"database/sql"
	"genesis_tt/db"
	"genesis_tt/util"
	"log"
	"os"
	"os/signal"
	"sync"
	"syscall"

	"github.com/gin-gonic/gin"
	"github.com/golang-migrate/migrate/v4"
	_ "github.com/golang-migrate/migrate/v4/database/postgres"
	_ "github.com/golang-migrate/migrate/v4/source/file"
	"github.com/robfig/cron"
)

// Structure for server managing
type Server struct {
	router       *gin.Engine
	config       util.Config
	InfoLog      *log.Logger
	ErrorLog     *log.Logger
	store        db.PostgresDB
	wait         *sync.WaitGroup
	mailer       util.Mail
	cronOperator *cron.Cron
}

// NewServer initializes Server struct:
//
//	creates a new HTTP server and set up routing
//	establishes a connection with the database
//	initializes struct Mail for managing sending of emails
//	loads config
//	creates loggers
//
// creates Cron object that is responsible for sending emails once a day
func NewServer() (*Server, error) {

	infoLog := log.New(os.Stdout, "INFO\t", log.Ldate|log.Ltime)
	errorLog := log.New(os.Stdout, "ERROR\t", log.Ldate|log.Ltime|log.Lshortfile)

	//Loading config from .env file
	config, err := util.LoadConfig(".")
	if err != nil {
		errorLog.Fatal("cannot load config:", err)
	}

	server := &Server{
		InfoLog:      infoLog,
		ErrorLog:     errorLog,
		config:       config,
		wait:         &sync.WaitGroup{},
		cronOperator: cron.New(),
	}

	//Connecting to DB and run migration
	server.store = server.InitDB()
	server.runDBMigration()

	//Initializing of Mail struct
	server.mailer = server.createMailer()

	//Creating cron operation for sending emails once a day
	server.cronOperator.AddFunc(util.TimeToSendEmails(server.config.DailyEmailsTime), server.sendEmailsOncePerDay)
	server.cronOperator.Start()

	server.setUpRouter()
	return server, nil

}

func (server *Server) InitDB() db.PostgresDB {
	conn, err := sql.Open(server.config.DBDriver, server.config.DBSource)
	if err != nil {
		server.ErrorLog.Fatal("cannot connect to the db:", err)
	}

	return db.NewPostgresDB(conn)

}

func (server *Server) runDBMigration() {
	migration, err := migrate.New(server.config.MigrationURL, server.config.DBSource)
	if err != nil {
		server.ErrorLog.Fatal("cannot create migration:", err)
	}

	if err = migration.Up(); err != nil && err != migrate.ErrNoChange {
		server.ErrorLog.Fatal("failed to run migrate up:", err)
	}

	server.InfoLog.Println("DB migrated successfully")

}

func (server *Server) createMailer() util.Mail {
	return util.Mail{
		Domain:      server.config.MailerDomain,
		Host:        server.config.MailerHost,
		Port:        server.config.MailerPort,
		FromName:    "rateInfo",
		FromAddress: "rateInfo@example.com",
		Wait:        server.wait,
		ErrorChan:   make(chan error),
		MailerChan:  make(chan util.Message, 100),
		DoneChan:    make(chan bool),
	}
}

func (server *Server) setUpRouter() {
	router := gin.Default()

	router.GET("/rate", server.getRate)
	router.POST("/subscribe", server.subscribeEmail)
	router.POST("/sendEmails", server.sendEmails)

	server.router = router
}

// Start runs the HTTP server on a specific address.
func (server *Server) Start() error {
	return server.router.Run(server.config.ServerAddress)
}

// a function to listen for messages on the MailerChan
func (server *Server) ListenForMail() {
	for {
		select {
		case msg := <-server.mailer.MailerChan:
			go server.mailer.SendMail(msg, server.mailer.ErrorChan)
		case err := <-server.mailer.ErrorChan:
			server.ErrorLog.Println(err)
		case <-server.mailer.DoneChan:
			return

		}

	}

}

// This function listens for shutdown and does all the necessary work before exiting
func (server *Server) ListenForShutdown() {
	quit := make(chan os.Signal, 1)
	signal.Notify(quit, syscall.SIGINT, syscall.SIGTERM)
	<-quit

	server.shutdown()
	os.Exit(0)
}

func (server *Server) shutdown() {
	server.wait.Wait()
	server.mailer.DoneChan <- true
	server.cronOperator.Stop()
	close(server.mailer.MailerChan)
	close(server.mailer.ErrorChan)
	close(server.mailer.DoneChan)

}
