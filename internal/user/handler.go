package user

import (
	"net/http"
	"strconv"

	"github.com/gin-gonic/gin"
)

type Handler struct {
    Service *Service
}

func NewHandler(service *Service) *Handler {
    return &Handler{Service: service}
}


func (h *Handler) GetUserAssets(c *gin.Context) {
    userID, _ := strconv.Atoi(c.Param("userId"))

    folders, folderShares, noteShares, err := h.Service.GetUserAssets(uint(userID))
    if err != nil {
        c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
        return
    }

    c.JSON(http.StatusOK, gin.H{
        "owned_folders":  folders,
        "shared_folders": folderShares,
        "shared_notes":   noteShares,
    })
}



