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

type ColorisRepository interface {
	Create(ctx context.Context, coloris *models.Coloris) error
	FindByID(ctx context.Context, id string) (*models.Coloris, error)
	FindAll(ctx context.Context, page, perPage int) ([]models.Coloris, int64, error)
	Update(ctx context.Context, id string, coloris *models.Coloris) error
	Delete(ctx context.Context, id string) error
	InsertMany(ctx context.Context, colorisData []models.Coloris) error
	FindWithFilters(ctx context.Context, filters bson.M, page, perPage int) ([]models.Coloris, int64, error)
}

type colorisRepository struct {
	collection *mongo.Collection
}

func NewColorisRepository(db *mongo.Database) ColorisRepository {
	return &colorisRepository{
		collection: db.Collection("coloris"),
	}
}

func (r *colorisRepository) Create(ctx context.Context, coloris *models.Coloris) error {
	coloris.ID = primitive.NewObjectID()
	coloris.CreatedAt = time.Now()
	coloris.UpdatedAt = time.Now()

	_, err := r.collection.InsertOne(ctx, coloris)
	return err
}

func (r *colorisRepository) FindByID(ctx context.Context, id string) (*models.Coloris, error) {
	objectID, err := primitive.ObjectIDFromHex(id)
	if err != nil {
		return nil, err
	}

	var coloris models.Coloris
	err = r.collection.FindOne(ctx, bson.M{"_id": objectID}).Decode(&coloris)
	if err != nil {
		return nil, err
	}

	return &coloris, nil
}

func (r *colorisRepository) FindAll(ctx context.Context, page, perPage int) ([]models.Coloris, int64, error) {
	skip := (page - 1) * perPage

	findOptions := options.Find()
	findOptions.SetSkip(int64(skip))
	findOptions.SetLimit(int64(perPage))
	findOptions.SetSort(bson.D{{Key: "created_at", Value: -1}})

	cursor, err := r.collection.Find(ctx, bson.M{}, findOptions)
	if err != nil {
		return nil, 0, err
	}
	defer cursor.Close(ctx)

	var colorisData []models.Coloris
	if err = cursor.All(ctx, &colorisData); err != nil {
		return nil, 0, err
	}

	total, err := r.collection.CountDocuments(ctx, bson.M{})
	if err != nil {
		return nil, 0, err
	}

	return colorisData, total, nil
}

func (r *colorisRepository) Update(ctx context.Context, id string, coloris *models.Coloris) error {
	objectID, err := primitive.ObjectIDFromHex(id)
	if err != nil {
		return err
	}

	coloris.UpdatedAt = time.Now()

	update := bson.M{
		"$set": bson.M{
			"timestamp":                coloris.Timestamp,
			"bulan":                    coloris.Bulan,
			"region":                   coloris.Region,
			"cabang":                   coloris.Cabang,
			"materi":                   coloris.Materi,
			"nama_atasan_langsung":     coloris.NamaAtasanLangsung,
			"nama_toko":                coloris.NamaToko,
			"nama_lengkap_sesuai_ktp":  coloris.NamaLengkapSesuaiKTP,
			"nilai_pg":                 coloris.NilaiPG,
			"nilai_akhir":              coloris.NilaiAkhir,
			"total":                    coloris.Total,
			"updated_at":               coloris.UpdatedAt,
		},
	}

	_, err = r.collection.UpdateOne(ctx, bson.M{"_id": objectID}, update)
	return err
}

func (r *colorisRepository) Delete(ctx context.Context, id string) error {
	objectID, err := primitive.ObjectIDFromHex(id)
	if err != nil {
		return err
	}

	_, err = r.collection.DeleteOne(ctx, bson.M{"_id": objectID})
	return err
}

func (r *colorisRepository) InsertMany(ctx context.Context, colorisData []models.Coloris) error {
	if len(colorisData) == 0 {
		return nil
	}

	documents := make([]interface{}, len(colorisData))
	for i, data := range colorisData {
		data.ID = primitive.NewObjectID()
		data.CreatedAt = time.Now()
		data.UpdatedAt = time.Now()
		documents[i] = data
	}

	_, err := r.collection.InsertMany(ctx, documents)
	return err
}

func (r *colorisRepository) FindWithFilters(ctx context.Context, filters bson.M, page, perPage int) ([]models.Coloris, int64, error) {
	skip := (page - 1) * perPage

	findOptions := options.Find()
	findOptions.SetSkip(int64(skip))
	findOptions.SetLimit(int64(perPage))
	findOptions.SetSort(bson.D{{Key: "created_at", Value: -1}})

	cursor, err := r.collection.Find(ctx, filters, findOptions)
	if err != nil {
		return nil, 0, err
	}
	defer cursor.Close(ctx)

	var colorisData []models.Coloris
	if err = cursor.All(ctx, &colorisData); err != nil {
		return nil, 0, err
	}

	total, err := r.collection.CountDocuments(ctx, filters)
	if err != nil {
		return nil, 0, err
	}

	return colorisData, total, nil
}
