package borrower_handler

import (
	"errors"
	"net/http"
	"strconv"

	"github.com/gin-gonic/gin"
	"github.com/iwansofian0512/billing_service/internal/model"
	"github.com/iwansofian0512/billing_service/internal/service/borrower_service"
)

type BorrowerHandler struct {
	service borrower_service.BorrowerService
}

func NewBorrowerHandler(service borrower_service.BorrowerService) *BorrowerHandler {
	return &BorrowerHandler{
		service: service,
	}
}

func (h *BorrowerHandler) CreateBorrower(ctx *gin.Context) {
	var req model.CreateBorrowerRequest
	if err := ctx.ShouldBindJSON(&req); err != nil {
		ctx.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	borrower, err := h.service.CreateBorrower(ctx.Request.Context(), req.Name, req.Email)
	if err != nil {
		if errors.Is(err, borrower_service.ErrBorrowerEmailExists) {
			ctx.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
			return
		}
		ctx.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	ctx.JSON(http.StatusCreated, borrower)
}

func (h *BorrowerHandler) ListBorrowerLoans(ctx *gin.Context) {
	borrowerIDStr := ctx.Query("borrower_id")
	pageStr := ctx.Query("page")
	pageSizeStr := ctx.Query("page_size")

	var borrowerID int
	var err error
	if borrowerIDStr != "" {
		borrowerID, err = strconv.Atoi(borrowerIDStr)
		if err != nil {
			ctx.JSON(http.StatusBadRequest, gin.H{"error": "invalid borrower_id"})
			return
		}
	}

	page := 1
	if pageStr != "" {
		page, err = strconv.Atoi(pageStr)
		if err != nil || page < 1 {
			ctx.JSON(http.StatusBadRequest, gin.H{"error": "invalid page"})
			return
		}
	}

	pageSize := 10
	if pageSizeStr != "" {
		pageSize, err = strconv.Atoi(pageSizeStr)
		if err != nil || pageSize < 1 {
			ctx.JSON(http.StatusBadRequest, gin.H{"error": "invalid page_size"})
			return
		}
	}

	loans, err := h.service.ListBorrowerLoans(ctx.Request.Context(), borrowerID, page, pageSize)
	if err != nil {
		ctx.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	ctx.JSON(http.StatusOK, loans)
}
