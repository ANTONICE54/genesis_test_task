package api

import (
	"fmt"
	"genesis_tt/util"
	"net/http"

	"github.com/gin-gonic/gin"
)

// Handler for obtaining the exchange rate
func (server *Server) getRate(ctx *gin.Context) {

	rate, err := util.FetchRateData(server.config.RateAPIKey)

	if err != nil {
		ctx.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	desc := fmt.Sprintf("Сurrent dollar (USD) to hryvnia (UAH) exchange rate: %f", *rate)

	ctx.JSON(http.StatusOK, gin.H{"description": desc})

}
