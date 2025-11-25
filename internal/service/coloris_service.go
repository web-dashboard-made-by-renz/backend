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

type ColorisService interface {
	CreateColoris(ctx context.Context, req *models.ColorisCreateRequest) error
	GetColorisById(ctx context.Context, id string) (*models.Coloris, error)
	GetAllColoris(ctx context.Context, page, perPage int) (*models.ColorisListResponse, error)
	UpdateColoris(ctx context.Context, id string, req *models.ColorisCreateRequest) error
	DeleteColoris(ctx context.Context, id string) error
	ImportFromExcel(ctx context.Context, file *multipart.FileHeader) (int, error)
	ExportToExcel(ctx context.Context) (*excelize.File, error)
	GetColorisWithFilters(ctx context.Context, filters map[string]string, page, perPage int) (*models.ColorisListResponse, error)
}

type colorisService struct {
	repo repository.ColorisRepository
}

func NewColorisService(repo repository.ColorisRepository) ColorisService {
	return &colorisService{
		repo: repo,
	}
}

func (s *colorisService) CreateColoris(ctx context.Context, req *models.ColorisCreateRequest) error {
	timestamp, err := utils.ParseTimestamp(req.Timestamp)
	if err != nil {
		return fmt.Errorf("invalid timestamp format: %v", err)
	}

	coloris := &models.Coloris{
		Timestamp:            timestamp,
		Bulan:                req.Bulan,
		Region:               req.Region,
		Cabang:               req.Cabang,
		Materi:               req.Materi,
		NamaAtasanLangsung:   req.NamaAtasanLangsung,
		NamaToko:             req.NamaToko,
		NamaLengkapSesuaiKTP: req.NamaLengkapSesuaiKTP,
		NilaiPG:              req.NilaiPG,
		NilaiAkhir:           req.NilaiAkhir,
		Total:                req.Total,
	}

	return s.repo.Create(ctx, coloris)
}

func (s *colorisService) GetColorisById(ctx context.Context, id string) (*models.Coloris, error) {
	return s.repo.FindByID(ctx, id)
}

func (s *colorisService) GetAllColoris(ctx context.Context, page, perPage int) (*models.ColorisListResponse, error) {
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

	return &models.ColorisListResponse{
		Data:       data,
		Total:      total,
		Page:       page,
		PerPage:    perPage,
		TotalPages: totalPages,
	}, nil
}

func (s *colorisService) UpdateColoris(ctx context.Context, id string, req *models.ColorisCreateRequest) error {
	timestamp, err := utils.ParseTimestamp(req.Timestamp)
	if err != nil {
		return fmt.Errorf("invalid timestamp format: %v", err)
	}

	coloris := &models.Coloris{
		Timestamp:            timestamp,
		Bulan:                req.Bulan,
		Region:               req.Region,
		Cabang:               req.Cabang,
		Materi:               req.Materi,
		NamaAtasanLangsung:   req.NamaAtasanLangsung,
		NamaToko:             req.NamaToko,
		NamaLengkapSesuaiKTP: req.NamaLengkapSesuaiKTP,
		NilaiPG:              req.NilaiPG,
		NilaiAkhir:           req.NilaiAkhir,
		Total:                req.Total,
	}

	return s.repo.Update(ctx, id, coloris)
}

func (s *colorisService) DeleteColoris(ctx context.Context, id string) error {
	return s.repo.Delete(ctx, id)
}

func (s *colorisService) ImportFromExcel(ctx context.Context, fileHeader *multipart.FileHeader) (int, error) {
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

	colorisData, err := utils.ParseExcelToColoris(excelFile)
	if err != nil {
		return 0, fmt.Errorf("failed to parse Excel data: %v", err)
	}

	if len(colorisData) == 0 {
		return 0, fmt.Errorf("no valid data found in Excel file")
	}

	err = s.repo.InsertMany(ctx, colorisData)
	if err != nil {
		return 0, fmt.Errorf("failed to insert data: %v", err)
	}

	return len(colorisData), nil
}

func (s *colorisService) ExportToExcel(ctx context.Context) (*excelize.File, error) {
	data, _, err := s.repo.FindAll(ctx, 1, 999999)
	if err != nil {
		return nil, fmt.Errorf("failed to fetch data: %v", err)
	}

	excelFile, err := utils.ExportColorisToExcel(data)
	if err != nil {
		return nil, fmt.Errorf("failed to create Excel file: %v", err)
	}

	return excelFile, nil
}

func (s *colorisService) GetColorisWithFilters(ctx context.Context, filters map[string]string, page, perPage int) (*models.ColorisListResponse, error) {
	if page < 1 {
		page = 1
	}
	if perPage < 1 {
		perPage = 10
	}

	bsonFilters := bson.M{}
	for key, value := range filters {
		if value != "" {
			bsonFilters[key] = bson.M{"$regex": value, "$options": "i"}
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

	return &models.ColorisListResponse{
		Data:       data,
		Total:      total,
		Page:       page,
		PerPage:    perPage,
		TotalPages: totalPages,
	}, nil
}
