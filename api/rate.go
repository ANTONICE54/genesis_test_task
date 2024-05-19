package api

import (
	"genesis_tt/util"
	"net/http"
	"strconv"

	"github.com/gin-gonic/gin"
)

// Handler for obtaining the exchange rate
func (server *Server) getRate(ctx *gin.Context) {

	rate, err := util.FetchRateData(server.config.RateAPIKey)

	if err != nil {
		ctx.Status(http.StatusBadRequest)
		return
	}

	ctx.String(http.StatusOK, strconv.FormatFloat(*rate, 'f', -1, 64))

}
