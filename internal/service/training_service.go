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

type TrainingService interface {
	CreateTraining(ctx context.Context, req *models.TrainingCreateRequest) error
	GetTrainingById(ctx context.Context, id string) (*models.Training, error)
	GetAllTraining(ctx context.Context, page, perPage int) (*models.TrainingListResponse, error)
	UpdateTraining(ctx context.Context, id string, req *models.TrainingCreateRequest) error
	DeleteTraining(ctx context.Context, id string) error
	ImportFromExcel(ctx context.Context, file *multipart.FileHeader) (int, error)
	ExportToExcel(ctx context.Context) (*excelize.File, error)
	GetTrainingWithFilters(ctx context.Context, filters map[string]string, page, perPage int) (*models.TrainingListResponse, error)
}

type trainingService struct {
	repo repository.TrainingRepository
}

func NewTrainingService(repo repository.TrainingRepository) TrainingService {
	return &trainingService{
		repo: repo,
	}
}

func (s *trainingService) CreateTraining(ctx context.Context, req *models.TrainingCreateRequest) error {
	timestamp, err := utils.ParseTimestamp(req.Timestamp)
	if err != nil {
		return fmt.Errorf("invalid timestamp format: %v", err)
	}

	training := &models.Training{
		Timestamp:            timestamp,
		Bulan:                req.Bulan,
		Region:               req.Region,
		CabangArea:           req.CabangArea,
		NamaAtasanLangsung:   req.NamaAtasanLangsung,
		MateriPelatihan:      req.MateriPelatihan,
		NamaLengkapSesuaiKTP: req.NamaLengkapSesuaiKTP,
		Jabatan:              req.Jabatan,
		TotalNilai:           req.TotalNilai,
		NilaiEssay:           req.NilaiEssay,
		Total:                req.Total,
	}

	return s.repo.Create(ctx, training)
}

func (s *trainingService) GetTrainingById(ctx context.Context, id string) (*models.Training, error) {
	return s.repo.FindByID(ctx, id)
}

func (s *trainingService) GetAllTraining(ctx context.Context, page, perPage int) (*models.TrainingListResponse, error) {
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

	return &models.TrainingListResponse{
		Data:       data,
		Total:      total,
		Page:       page,
		PerPage:    perPage,
		TotalPages: totalPages,
	}, nil
}

func (s *trainingService) UpdateTraining(ctx context.Context, id string, req *models.TrainingCreateRequest) error {
	timestamp, err := utils.ParseTimestamp(req.Timestamp)
	if err != nil {
		return fmt.Errorf("invalid timestamp format: %v", err)
	}

	training := &models.Training{
		Timestamp:            timestamp,
		Bulan:                req.Bulan,
		Region:               req.Region,
		CabangArea:           req.CabangArea,
		NamaAtasanLangsung:   req.NamaAtasanLangsung,
		MateriPelatihan:      req.MateriPelatihan,
		NamaLengkapSesuaiKTP: req.NamaLengkapSesuaiKTP,
		Jabatan:              req.Jabatan,
		TotalNilai:           req.TotalNilai,
		NilaiEssay:           req.NilaiEssay,
		Total:                req.Total,
	}

	return s.repo.Update(ctx, id, training)
}

func (s *trainingService) DeleteTraining(ctx context.Context, id string) error {
	return s.repo.Delete(ctx, id)
}

func (s *trainingService) ImportFromExcel(ctx context.Context, fileHeader *multipart.FileHeader) (int, error) {
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

	trainingData, err := utils.ParseExcelToTraining(excelFile)
	if err != nil {
		return 0, fmt.Errorf("failed to parse Excel data: %v", err)
	}

	if len(trainingData) == 0 {
		return 0, fmt.Errorf("no valid data found in Excel file")
	}

	err = s.repo.InsertMany(ctx, trainingData)
	if err != nil {
		return 0, fmt.Errorf("failed to insert data: %v", err)
	}

	return len(trainingData), nil
}

func (s *trainingService) ExportToExcel(ctx context.Context) (*excelize.File, error) {
	data, _, err := s.repo.FindAll(ctx, 1, 999999)
	if err != nil {
		return nil, fmt.Errorf("failed to fetch data: %v", err)
	}

	excelFile, err := utils.ExportTrainingToExcel(data)
	if err != nil {
		return nil, fmt.Errorf("failed to create Excel file: %v", err)
	}

	return excelFile, nil
}

func (s *trainingService) GetTrainingWithFilters(ctx context.Context, filters map[string]string, page, perPage int) (*models.TrainingListResponse, error) {
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

	return &models.TrainingListResponse{
		Data:       data,
		Total:      total,
		Page:       page,
		PerPage:    perPage,
		TotalPages: totalPages,
	}, nil
}
