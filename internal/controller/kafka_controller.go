package controller

import (
	"hyperscal-go/internal/dto"
	"hyperscal-go/internal/service"
	"net/http"

	"github.com/gin-gonic/gin"
)

// KafkaController handles HTTP requests that publish events to Kafka.
type KafkaController struct {
	service *service.KafkaService
}

// NewKafkaController creates a new KafkaController.
func NewKafkaController(svc *service.KafkaService) *KafkaController {
	return &KafkaController{service: svc}
}

// PublishUserCreated handles POST /api/kafka/user-created
func (k *KafkaController) PublishUserCreated(ctx *gin.Context) {
	if k.service == nil {
		ctx.JSON(http.StatusServiceUnavailable, dto.ErrorResponse("Kafka is unavailable", "kafka service not initialized"))
		return
	}

	var req dto.PublishUserCreatedRequest
	if err := ctx.ShouldBindJSON(&req); err != nil {
		ctx.JSON(http.StatusBadRequest, dto.ErrorResponse("Invalid request body", err.Error()))
		return
	}

	if err := k.service.PublishUserCreated(ctx.Request.Context(), req.UserID, req.Email); err != nil {
		ctx.JSON(http.StatusInternalServerError, dto.ErrorResponse("Failed to publish event", err.Error()))
		return
	}

	ctx.JSON(http.StatusOK, dto.SuccessResponse("Event user.created published successfully", nil))
}

// PublishOrderPlaced handles POST /api/kafka/order-placed
func (k *KafkaController) PublishOrderPlaced(ctx *gin.Context) {
	if k.service == nil {
		ctx.JSON(http.StatusServiceUnavailable, dto.ErrorResponse("Kafka is unavailable", "kafka service not initialized"))
		return
	}

	var req dto.PublishOrderPlacedRequest
	if err := ctx.ShouldBindJSON(&req); err != nil {
		ctx.JSON(http.StatusBadRequest, dto.ErrorResponse("Invalid request body", err.Error()))
		return
	}

	if err := k.service.PublishOrderPlaced(ctx.Request.Context(), req.OrderID, req.UserID, req.Total); err != nil {
		ctx.JSON(http.StatusInternalServerError, dto.ErrorResponse("Failed to publish event", err.Error()))
		return
	}

	ctx.JSON(http.StatusOK, dto.SuccessResponse("Event order.placed published successfully", nil))
}
