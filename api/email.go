package api

import (
	"context"
	"fmt"
	"genesis_tt/db"
	"genesis_tt/util"
	"net/http"

	"github.com/gin-gonic/gin"
	"github.com/lib/pq"
)

type subscribeEmailRequest struct {
	Email string `json:"email" binding:"email"`
}

// Handler for email subscription to receive information on course changes
func (server *Server) subscribeEmail(ctx *gin.Context) {
	var req subscribeEmailRequest

	if err := ctx.ShouldBindJSON(&req); err != nil {
		ctx.Status(http.StatusBadRequest)
		return
	}

	_, err := server.store.AddEmail(ctx, req.Email)

	if err != nil {
		if pqErr, ok := err.(*pq.Error); ok && pqErr.Code.Name() == "unique_violation" {
			ctx.Status(http.StatusConflict)
			return

		}
		ctx.Status(http.StatusInternalServerError)
		return
	}

	ctx.Status(http.StatusOK)
}

// Handler for sending exchange rate information to emails
func (server *Server) sendEmails(ctx *gin.Context) {

	emailsList, err := server.store.ListEmails(ctx)

	if err != nil {
		ctx.Status(http.StatusInternalServerError)
		return
	}

	rate, err := util.FetchRateData(server.config.RateAPIKey)

	if err != nil {
		ctx.Status(http.StatusBadRequest)
		return
	}

	server.sendEmailsToEveryone(*rate, emailsList)

	ctx.Status(http.StatusOK)

}

// Function for sending messages to the MailerChan channel, where these messages will be concurrently sent via e-mail
func (server *Server) sendEmailsToEveryone(rate float64, emailsList []db.Email) {

	msg := util.Message{
		Subject: "Daily rate info",
		Data:    fmt.Sprintf("Сurrent dollar (USD) to hryvnia (UAH) exchange rate: %f", rate),
	}

	for _, emailInfo := range emailsList {
		msg.To = emailInfo.Email
		server.wait.Add(1)
		server.mailer.MailerChan <- msg

	}

}

// Function used to send emails once a day
func (server *Server) sendEmailsOncePerDay() {

	rate, err := util.FetchRateData(server.config.RateAPIKey)
	if err != nil {
		server.ErrorLog.Println(err)
		return
	}

	emailsList, err := server.store.ListEmails(context.Background())
	if err != nil {
		server.ErrorLog.Println(err)
		return
	}

	msg := util.Message{
		Subject: "Daily rate info",
		Data:    fmt.Sprintf("Сurrent dollar (USD) to hryvnia (UAH) exchange rate: %f", *rate),
	}

	for _, emailInfo := range emailsList {
		msg.To = emailInfo.Email
		server.wait.Add(1)
		server.mailer.MailerChan <- msg

	}

}
