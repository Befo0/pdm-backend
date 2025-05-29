package services

import (
	"strconv"

	"github.com/gin-gonic/gin"
)

func GetFinanceId(c *gin.Context) (uint, error) {
	finanzaId := c.Query("finanza_id")

	if finanzaId != "" {
		id, err := strconv.ParseUint(finanzaId, 10, 64)
		if err != nil {
			return 0, err
		}
		finanzaId := uint(id)

		return finanzaId, nil
	}

	return 0, nil
}
