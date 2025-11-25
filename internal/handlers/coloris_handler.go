package handlers

import (
	"fmt"
	"net/http"
	"strconv"

	"github.com/gin-gonic/gin"
	"github.com/web-dashboard-made-by-renz/backend/internal/models"
	"github.com/web-dashboard-made-by-renz/backend/internal/service"
)

type ColorisHandler struct {
	service service.ColorisService
}

func NewColorisHandler(service service.ColorisService) *ColorisHandler {
	return &ColorisHandler{
		service: service,
	}
}

func (h *ColorisHandler) CreateColoris(c *gin.Context) {
	var req models.ColorisCreateRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	err := h.service.CreateColoris(c.Request.Context(), &req)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusCreated, gin.H{"message": "Data Coloris berhasil dibuat"})
}

func (h *ColorisHandler) GetColorisById(c *gin.Context) {
	id := c.Param("id")

	coloris, err := h.service.GetColorisById(c.Request.Context(), id)
	if err != nil {
		c.JSON(http.StatusNotFound, gin.H{"error": "Data tidak ditemukan"})
		return
	}

	c.JSON(http.StatusOK, gin.H{"data": coloris})
}

func (h *ColorisHandler) GetAllColoris(c *gin.Context) {
	page, _ := strconv.Atoi(c.DefaultQuery("page", "1"))
	perPage, _ := strconv.Atoi(c.DefaultQuery("per_page", "10"))

	region := c.Query("region")
	cabang := c.Query("cabang")
	bulan := c.Query("bulan")

	filters := make(map[string]string)
	if region != "" {
		filters["region"] = region
	}
	if cabang != "" {
		filters["cabang"] = cabang
	}
	if bulan != "" {
		filters["bulan"] = bulan
	}

	var response *models.ColorisListResponse
	var err error

	if len(filters) > 0 {
		response, err = h.service.GetColorisWithFilters(c.Request.Context(), filters, page, perPage)
	} else {
		response, err = h.service.GetAllColoris(c.Request.Context(), page, perPage)
	}

	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusOK, response)
}

func (h *ColorisHandler) UpdateColoris(c *gin.Context) {
	id := c.Param("id")

	var req models.ColorisCreateRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	err := h.service.UpdateColoris(c.Request.Context(), id, &req)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusOK, gin.H{"message": "Data Coloris berhasil diupdate"})
}

func (h *ColorisHandler) DeleteColoris(c *gin.Context) {
	id := c.Param("id")

	err := h.service.DeleteColoris(c.Request.Context(), id)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusOK, gin.H{"message": "Data Coloris berhasil dihapus"})
}

func (h *ColorisHandler) ImportExcel(c *gin.Context) {
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

func (h *ColorisHandler) ExportExcel(c *gin.Context) {
	excelFile, err := h.service.ExportToExcel(c.Request.Context())
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}
	defer excelFile.Close()

	filename := fmt.Sprintf("coloris_data_%s.xlsx", c.Query("filename"))
	if c.Query("filename") == "" {
		filename = "coloris_data.xlsx"
	}

	c.Header("Content-Type", "application/vnd.openxmlformats-officedocument.spreadsheetml.sheet")
	c.Header("Content-Disposition", fmt.Sprintf("attachment; filename=%s", filename))

	if err := excelFile.Write(c.Writer); err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Gagal menulis file Excel"})
		return
	}
}
