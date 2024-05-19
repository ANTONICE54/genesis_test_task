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
		ctx.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	_, err := server.store.AddEmail(ctx, req.Email)

	if err != nil {
		if pqErr, ok := err.(*pq.Error); ok && pqErr.Code.Name() == "unique_violation" {
			ctx.JSON(http.StatusConflict, gin.H{"error": err.Error()})
			return

		}
		ctx.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	ctx.JSON(http.StatusOK, gin.H{"description": "E-mail was successfully added"})
}

// Handler for sending exchange rate information to emails
func (server *Server) sendEmails(ctx *gin.Context) {

	emailsList, err := server.store.ListEmails(ctx)

	if err != nil {
		ctx.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	rate, err := util.FetchRateData(server.config.RateAPIKey)

	if err != nil {
		ctx.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	server.sendEmailsToEveryone(*rate, emailsList)

	ctx.JSON(http.StatusOK, gin.H{"description": "E-mails have been sent successfully"})

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
