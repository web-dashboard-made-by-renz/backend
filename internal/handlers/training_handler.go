package handlers

import (
	"fmt"
	"net/http"
	"strconv"

	"github.com/gin-gonic/gin"
	"github.com/web-dashboard-made-by-renz/backend/internal/models"
	"github.com/web-dashboard-made-by-renz/backend/internal/service"
)

type TrainingHandler struct {
	service service.TrainingService
}

func NewTrainingHandler(service service.TrainingService) *TrainingHandler {
	return &TrainingHandler{
		service: service,
	}
}

func (h *TrainingHandler) CreateTraining(c *gin.Context) {
	var req models.TrainingCreateRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	err := h.service.CreateTraining(c.Request.Context(), &req)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusCreated, gin.H{"message": "Data Training berhasil dibuat"})
}

func (h *TrainingHandler) GetTrainingById(c *gin.Context) {
	id := c.Param("id")

	training, err := h.service.GetTrainingById(c.Request.Context(), id)
	if err != nil {
		c.JSON(http.StatusNotFound, gin.H{"error": "Data tidak ditemukan"})
		return
	}

	c.JSON(http.StatusOK, gin.H{"data": training})
}

func (h *TrainingHandler) GetAllTraining(c *gin.Context) {
	page, _ := strconv.Atoi(c.DefaultQuery("page", "1"))
	perPage, _ := strconv.Atoi(c.DefaultQuery("per_page", "10"))

	region := c.Query("region")
	cabangArea := c.Query("cabang_area")
	bulan := c.Query("bulan")

	filters := make(map[string]string)
	if region != "" {
		filters["region"] = region
	}
	if cabangArea != "" {
		filters["cabang_area"] = cabangArea
	}
	if bulan != "" {
		filters["bulan"] = bulan
	}

	var response *models.TrainingListResponse
	var err error

	if len(filters) > 0 {
		response, err = h.service.GetTrainingWithFilters(c.Request.Context(), filters, page, perPage)
	} else {
		response, err = h.service.GetAllTraining(c.Request.Context(), page, perPage)
	}

	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusOK, response)
}

func (h *TrainingHandler) UpdateTraining(c *gin.Context) {
	id := c.Param("id")

	var req models.TrainingCreateRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	err := h.service.UpdateTraining(c.Request.Context(), id, &req)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusOK, gin.H{"message": "Data Training berhasil diupdate"})
}

func (h *TrainingHandler) DeleteTraining(c *gin.Context) {
	id := c.Param("id")

	err := h.service.DeleteTraining(c.Request.Context(), id)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusOK, gin.H{"message": "Data Training berhasil dihapus"})
}

func (h *TrainingHandler) ImportExcel(c *gin.Context) {
	file, err := c.FormFile("file")
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "File tidak ditemukan"})
		return
	}

	if file.Header.Get("Content-Type") != "application/vnd.openxmlformats-officedocument.spreadsheetml.sheet" &&
		file.Header.Get("Content-Type") != "application/vnd.ms-excel" {
		ext := file.Filename[len(file.Filename)-5:]
		if ext != ".xlsx" && ext != ".xls" {
			c.JSON(http.StatusBadRequest, gin.H{"error": "File harus berformat Excel (.xlsx atau .xls)"})
			return
		}
	}

	count, err := h.service.ImportFromExcel(c.Request.Context(), file)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"message": "Import berhasil",
		"count":   count,
	})
}

func (h *TrainingHandler) ExportExcel(c *gin.Context) {
	excelFile, err := h.service.ExportToExcel(c.Request.Context())
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}
	defer excelFile.Close()

	filename := fmt.Sprintf("training_data_%s.xlsx", c.Query("filename"))
	if c.Query("filename") == "" {
		filename = "training_data.xlsx"
	}

	c.Header("Content-Type", "application/vnd.openxmlformats-officedocument.spreadsheetml.sheet")
	c.Header("Content-Disposition", fmt.Sprintf("attachment; filename=%s", filename))

	if err := excelFile.Write(c.Writer); err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Gagal menulis file Excel"})
		return
	}
}
