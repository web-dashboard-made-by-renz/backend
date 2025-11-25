package service

import (
	"context"
	"fmt"
	"mime/multipart"

	"github.com/web-dashboard-made-by-renz/backend/internal/models"
	"github.com/web-dashboard-made-by-renz/backend/internal/repository"
	"github.com/web-dashboard-made-by-renz/backend/pkg/utils"
	"github.com/xuri/excelize/v2"
	"go.mongodb.org/mongo-driver/bson"
)

type SelloutService interface {
	CreateSellout(ctx context.Context, req *models.SelloutCreateRequest) error
	GetSelloutById(ctx context.Context, id string) (*models.Sellout, error)
	GetAllSellout(ctx context.Context, page, perPage int) (*models.SelloutListResponse, error)
	UpdateSellout(ctx context.Context, id string, req *models.SelloutCreateRequest) error
	DeleteSellout(ctx context.Context, id string) error
	ImportFromExcel(ctx context.Context, file *multipart.FileHeader) (int, error)
	ExportToExcel(ctx context.Context) (*excelize.File, error)
	GetSelloutWithFilters(ctx context.Context, filters map[string]string, page, perPage int) (*models.SelloutListResponse, error)
}

type selloutService struct {
	repo repository.SelloutRepository
}

func NewSelloutService(repo repository.SelloutRepository) SelloutService {
	return &selloutService{
		repo: repo,
	}
}

func (s *selloutService) CreateSellout(ctx context.Context, req *models.SelloutCreateRequest) error {
	sellout := &models.Sellout{
		Tahun:        req.Tahun,
		Bulan:        req.Bulan,
		Reg:          req.Reg,
		Cabang:       req.Cabang,
		Outlet:       req.Outlet,
		NamaColorist: req.NamaColorist,
		NoReg:        req.NoReg,
		CHL:          req.CHL,
		SelloutTT:    req.SelloutTT,
		SelloutRM:    req.SelloutRM,
		TotalSellout: req.TotalSellout,
	}

	return s.repo.Create(ctx, sellout)
}

func (s *selloutService) GetSelloutById(ctx context.Context, id string) (*models.Sellout, error) {
	return s.repo.FindByID(ctx, id)
}

func (s *selloutService) GetAllSellout(ctx context.Context, page, perPage int) (*models.SelloutListResponse, error) {
	if page < 1 {
		page = 1
	}
	if perPage < 1 {
		perPage = 10
	}

	data, total, err := s.repo.FindAll(ctx, page, perPage)
	if err != nil {
		return nil, err
	}

	totalPages := int(total) / perPage
	if int(total)%perPage != 0 {
		totalPages++
	}

	return &models.SelloutListResponse{
		Data:       data,
		Total:      total,
		Page:       page,
		PerPage:    perPage,
		TotalPages: totalPages,
	}, nil
}

func (s *selloutService) UpdateSellout(ctx context.Context, id string, req *models.SelloutCreateRequest) error {
	sellout := &models.Sellout{
		Tahun:        req.Tahun,
		Bulan:        req.Bulan,
		Reg:          req.Reg,
		Cabang:       req.Cabang,
		Outlet:       req.Outlet,
		NamaColorist: req.NamaColorist,
		NoReg:        req.NoReg,
		CHL:          req.CHL,
		SelloutTT:    req.SelloutTT,
		SelloutRM:    req.SelloutRM,
		TotalSellout: req.TotalSellout,
	}

	return s.repo.Update(ctx, id, sellout)
}

func (s *selloutService) DeleteSellout(ctx context.Context, id string) error {
	return s.repo.Delete(ctx, id)
}

func (s *selloutService) ImportFromExcel(ctx context.Context, fileHeader *multipart.FileHeader) (int, error) {
	file, err := fileHeader.Open()
	if err != nil {
		return 0, fmt.Errorf("failed to open file: %v", err)
	}
	defer file.Close()

	excelFile, err := excelize.OpenReader(file)
	if err != nil {
		return 0, fmt.Errorf("failed to read Excel file: %v", err)
	}
	defer excelFile.Close()

	selloutData, err := utils.ParseExcelToSellout(excelFile)
	if err != nil {
		return 0, fmt.Errorf("failed to parse Excel data: %v", err)
	}

	if len(selloutData) == 0 {
		return 0, fmt.Errorf("no valid data found in Excel file")
	}

	err = s.repo.InsertMany(ctx, selloutData)
	if err != nil {
		return 0, fmt.Errorf("failed to insert data: %v", err)
	}

	return len(selloutData), nil
}

func (s *selloutService) ExportToExcel(ctx context.Context) (*excelize.File, error) {
	data, _, err := s.repo.FindAll(ctx, 1, 999999)
	if err != nil {
		return nil, fmt.Errorf("failed to fetch data: %v", err)
	}

	excelFile, err := utils.ExportSelloutToExcel(data)
	if err != nil {
		return nil, fmt.Errorf("failed to create Excel file: %v", err)
	}

	return excelFile, nil
}

func (s *selloutService) GetSelloutWithFilters(ctx context.Context, filters map[string]string, page, perPage int) (*models.SelloutListResponse, error) {
	if page < 1 {
		page = 1
	}
	if perPage < 1 {
		perPage = 10
	}

	bsonFilters := bson.M{}
	for key, value := range filters {
		if value != "" {
			// For tahun and bulan, convert to int
			if key == "tahun" || key == "bulan" {
				// Try to parse as int for exact match
				bsonFilters[key] = value
			} else {
				bsonFilters[key] = bson.M{"$regex": value, "$options": "i"}
			}
		}
	}

	data, total, err := s.repo.FindWithFilters(ctx, bsonFilters, page, perPage)
	if err != nil {
		return nil, err
	}

	totalPages := int(total) / perPage
	if int(total)%perPage != 0 {
		totalPages++
	}

	return &models.SelloutListResponse{
		Data:       data,
		Total:      total,
		Page:       page,
		PerPage:    perPage,
		TotalPages: totalPages,
	}, nil
}
