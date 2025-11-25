package repository

import (
	"context"
	"time"

	"github.com/web-dashboard-made-by-renz/backend/internal/models"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/bson/primitive"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
)

type TrainingRepository interface {
	Create(ctx context.Context, training *models.Training) error
	FindByID(ctx context.Context, id string) (*models.Training, error)
	FindAll(ctx context.Context, page, perPage int) ([]models.Training, int64, error)
	Update(ctx context.Context, id string, training *models.Training) error
	Delete(ctx context.Context, id string) error
	InsertMany(ctx context.Context, trainings []models.Training) error
	FindWithFilters(ctx context.Context, filters bson.M, page, perPage int) ([]models.Training, int64, error)
}

type trainingRepository struct {
	collection *mongo.Collection
}

func NewTrainingRepository(db *mongo.Database) TrainingRepository {
	return &trainingRepository{
		collection: db.Collection("trainings"),
	}
}

func (r *trainingRepository) Create(ctx context.Context, training *models.Training) error {
	training.CreatedAt = time.Now()
	training.UpdatedAt = time.Now()

	_, err := r.collection.InsertOne(ctx, training)
	return err
}

func (r *trainingRepository) FindByID(ctx context.Context, id string) (*models.Training, error) {
	objectID, err := primitive.ObjectIDFromHex(id)
	if err != nil {
		return nil, err
	}

	var training models.Training
	err = r.collection.FindOne(ctx, bson.M{"_id": objectID}).Decode(&training)
	if err != nil {
		return nil, err
	}

	return &training, nil
}

func (r *trainingRepository) FindAll(ctx context.Context, page, perPage int) ([]models.Training, int64, error) {
	skip := (page - 1) * perPage

	opts := options.Find().
		SetSkip(int64(skip)).
		SetLimit(int64(perPage)).
		SetSort(bson.D{{Key: "timestamp", Value: -1}})

	cursor, err := r.collection.Find(ctx, bson.M{}, opts)
	if err != nil {
		return nil, 0, err
	}
	defer cursor.Close(ctx)

	var trainings []models.Training
	if err = cursor.All(ctx, &trainings); err != nil {
		return nil, 0, err
	}

	total, err := r.collection.CountDocuments(ctx, bson.M{})
	if err != nil {
		return nil, 0, err
	}

	return trainings, total, nil
}

func (r *trainingRepository) Update(ctx context.Context, id string, training *models.Training) error {
	objectID, err := primitive.ObjectIDFromHex(id)
	if err != nil {
		return err
	}

	training.UpdatedAt = time.Now()

	update := bson.M{
		"$set": bson.M{
			"timestamp":              training.Timestamp,
			"bulan":                  training.Bulan,
			"region":                 training.Region,
			"cabang_area":            training.CabangArea,
			"nama_atasan_langsung":   training.NamaAtasanLangsung,
			"materi_pelatihan":       training.MateriPelatihan,
			"nama_lengkap_sesuai_ktp": training.NamaLengkapSesuaiKTP,
			"jabatan":                training.Jabatan,
			"total_nilai":            training.TotalNilai,
			"nilai_essay":            training.NilaiEssay,
			"total":                  training.Total,
			"updated_at":             training.UpdatedAt,
		},
	}

	_, err = r.collection.UpdateOne(ctx, bson.M{"_id": objectID}, update)
	return err
}

func (r *trainingRepository) Delete(ctx context.Context, id string) error {
	objectID, err := primitive.ObjectIDFromHex(id)
	if err != nil {
		return err
	}

	_, err = r.collection.DeleteOne(ctx, bson.M{"_id": objectID})
	return err
}

func (r *trainingRepository) InsertMany(ctx context.Context, trainings []models.Training) error {
	docs := make([]interface{}, len(trainings))
	for i, training := range trainings {
		training.CreatedAt = time.Now()
		training.UpdatedAt = time.Now()
		docs[i] = training
	}

	_, err := r.collection.InsertMany(ctx, docs)
	return err
}

func (r *trainingRepository) FindWithFilters(ctx context.Context, filters bson.M, page, perPage int) ([]models.Training, int64, error) {
	skip := (page - 1) * perPage

	opts := options.Find().
		SetSkip(int64(skip)).
		SetLimit(int64(perPage)).
		SetSort(bson.D{{Key: "timestamp", Value: -1}})

	cursor, err := r.collection.Find(ctx, filters, opts)
	if err != nil {
		return nil, 0, err
	}
	defer cursor.Close(ctx)

	var trainings []models.Training
	if err = cursor.All(ctx, &trainings); err != nil {
		return nil, 0, err
	}

	total, err := r.collection.CountDocuments(ctx, filters)
	if err != nil {
		return nil, 0, err
	}

	return trainings, total, nil
}
