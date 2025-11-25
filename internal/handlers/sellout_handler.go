package handlers

import (
	"fmt"
	"net/http"
	"strconv"

	"github.com/gin-gonic/gin"
	"github.com/web-dashboard-made-by-renz/backend/internal/models"
	"github.com/web-dashboard-made-by-renz/backend/internal/service"
)

type SelloutHandler struct {
	service service.SelloutService
}

func NewSelloutHandler(service service.SelloutService) *SelloutHandler {
	return &SelloutHandler{
		service: service,
	}
}

func (h *SelloutHandler) CreateSellout(c *gin.Context) {
	var req models.SelloutCreateRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	err := h.service.CreateSellout(c.Request.Context(), &req)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusCreated, gin.H{"message": "Data Sellout berhasil dibuat"})
}

func (h *SelloutHandler) GetSelloutById(c *gin.Context) {
	id := c.Param("id")

	sellout, err := h.service.GetSelloutById(c.Request.Context(), id)
	if err != nil {
		c.JSON(http.StatusNotFound, gin.H{"error": "Data tidak ditemukan"})
		return
	}

	c.JSON(http.StatusOK, gin.H{"data": sellout})
}

func (h *SelloutHandler) GetAllSellout(c *gin.Context) {
	page, _ := strconv.Atoi(c.DefaultQuery("page", "1"))
	perPage, _ := strconv.Atoi(c.DefaultQuery("per_page", "10"))

	tahun := c.Query("tahun")
	bulan := c.Query("bulan")
	cabang := c.Query("cabang")

	filters := make(map[string]string)
	if tahun != "" {
		filters["tahun"] = tahun
	}
	if bulan != "" {
		filters["bulan"] = bulan
	}
	if cabang != "" {
		filters["cabang"] = cabang
	}

	var response *models.SelloutListResponse
	var err error

	if len(filters) > 0 {
		response, err = h.service.GetSelloutWithFilters(c.Request.Context(), filters, page, perPage)
	} else {
		response, err = h.service.GetAllSellout(c.Request.Context(), page, perPage)
	}

	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusOK, response)
}

func (h *SelloutHandler) UpdateSellout(c *gin.Context) {
	id := c.Param("id")

	var req models.SelloutCreateRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	err := h.service.UpdateSellout(c.Request.Context(), id, &req)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusOK, gin.H{"message": "Data Sellout berhasil diupdate"})
}

func (h *SelloutHandler) DeleteSellout(c *gin.Context) {
	id := c.Param("id")

	err := h.service.DeleteSellout(c.Request.Context(), id)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusOK, gin.H{"message": "Data Sellout berhasil dihapus"})
}

func (h *SelloutHandler) ImportExcel(c *gin.Context) {
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

func (h *SelloutHandler) ExportExcel(c *gin.Context) {
	excelFile, err := h.service.ExportToExcel(c.Request.Context())
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}
	defer excelFile.Close()

	filename := fmt.Sprintf("sellout_data_%s.xlsx", c.Query("filename"))
	if c.Query("filename") == "" {
		filename = "sellout_data.xlsx"
	}

	c.Header("Content-Type", "application/vnd.openxmlformats-officedocument.spreadsheetml.sheet")
	c.Header("Content-Disposition", fmt.Sprintf("attachment; filename=%s", filename))

	if err := excelFile.Write(c.Writer); err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Gagal menulis file Excel"})
		return
	}
}
