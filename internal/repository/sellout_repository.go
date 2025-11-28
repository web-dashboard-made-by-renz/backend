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

type SelloutRepository interface {
	Create(ctx context.Context, sellout *models.Sellout) error
	FindByID(ctx context.Context, id string) (*models.Sellout, error)
	FindAll(ctx context.Context, page, perPage int) ([]models.Sellout, int64, error)
	Update(ctx context.Context, id string, sellout *models.Sellout) error
	Delete(ctx context.Context, id string) error
	InsertMany(ctx context.Context, sellouts []models.Sellout) error
	FindWithFilters(ctx context.Context, filters bson.M, page, perPage int) ([]models.Sellout, int64, error)
}

type selloutRepository struct {
	collection *mongo.Collection
}

func NewSelloutRepository(db *mongo.Database) SelloutRepository {
	return &selloutRepository{
		collection: db.Collection("sellouts"),
	}
}

func (r *selloutRepository) Create(ctx context.Context, sellout *models.Sellout) error {
	sellout.CreatedAt = time.Now()
	sellout.UpdatedAt = time.Now()

	_, err := r.collection.InsertOne(ctx, sellout)
	return err
}

func (r *selloutRepository) FindByID(ctx context.Context, id string) (*models.Sellout, error) {
	objectID, err := primitive.ObjectIDFromHex(id)
	if err != nil {
		return nil, err
	}

	var sellout models.Sellout
	err = r.collection.FindOne(ctx, bson.M{"_id": objectID}).Decode(&sellout)
	if err != nil {
		return nil, err
	}

	return &sellout, nil
}

func (r *selloutRepository) FindAll(ctx context.Context, page, perPage int) ([]models.Sellout, int64, error) {
	skip := (page - 1) * perPage

	opts := options.Find().
		SetSkip(int64(skip)).
		SetLimit(int64(perPage)).
		SetSort(bson.D{{Key: "tahun", Value: -1}, {Key: "bulan", Value: -1}})

	cursor, err := r.collection.Find(ctx, bson.M{}, opts)
	if err != nil {
		return nil, 0, err
	}
	defer cursor.Close(ctx)

	var sellouts []models.Sellout
	if err = cursor.All(ctx, &sellouts); err != nil {
		return nil, 0, err
	}

	total, err := r.collection.CountDocuments(ctx, bson.M{})
	if err != nil {
		return nil, 0, err
	}

	return sellouts, total, nil
}

func (r *selloutRepository) Update(ctx context.Context, id string, sellout *models.Sellout) error {
	objectID, err := primitive.ObjectIDFromHex(id)
	if err != nil {
		return err
	}

	sellout.UpdatedAt = time.Now()

	update := bson.M{
		"$set": bson.M{
			"tahun":             sellout.Tahun,
			"bulan":             sellout.Bulan,
			"reg":               sellout.Reg,
			"cabang":            sellout.Cabang,
			"outlet":            sellout.Outlet,
			"area_cover":        sellout.AreaCover,
			"mos_ss":            sellout.MosSs,
			"nama_colorist":     sellout.NamaColorist,
			"no_reg":            sellout.NoReg,
			"tanggal_bergabung": sellout.TanggalBergabung,
			"masa_kerja":        sellout.MasaKerja,
			"chl":               sellout.CHL,
			"wilayah":           sellout.Wilayah,
			"target_sellout":    sellout.TargetSellout,
			"sellout_tt":        sellout.SelloutTT,
			"sellout_rm":        sellout.SelloutRM,
			"primafix":          sellout.Primafix,
			"total_sellout":     sellout.TotalSellout,
			"updated_at":        sellout.UpdatedAt,
		},
	}

	_, err = r.collection.UpdateOne(ctx, bson.M{"_id": objectID}, update)
	return err
}

func (r *selloutRepository) Delete(ctx context.Context, id string) error {
	objectID, err := primitive.ObjectIDFromHex(id)
	if err != nil {
		return err
	}

	_, err = r.collection.DeleteOne(ctx, bson.M{"_id": objectID})
	return err
}

func (r *selloutRepository) InsertMany(ctx context.Context, sellouts []models.Sellout) error {
	docs := make([]interface{}, len(sellouts))
	for i, sellout := range sellouts {
		sellout.CreatedAt = time.Now()
		sellout.UpdatedAt = time.Now()
		docs[i] = sellout
	}

	_, err := r.collection.InsertMany(ctx, docs)
	return err
}

func (r *selloutRepository) FindWithFilters(ctx context.Context, filters bson.M, page, perPage int) ([]models.Sellout, int64, error) {
	skip := (page - 1) * perPage

	opts := options.Find().
		SetSkip(int64(skip)).
		SetLimit(int64(perPage)).
		SetSort(bson.D{{Key: "tahun", Value: -1}, {Key: "bulan", Value: -1}})

	cursor, err := r.collection.Find(ctx, filters, opts)
	if err != nil {
		return nil, 0, err
	}
	defer cursor.Close(ctx)

	var sellouts []models.Sellout
	if err = cursor.All(ctx, &sellouts); err != nil {
		return nil, 0, err
	}

	total, err := r.collection.CountDocuments(ctx, filters)
	if err != nil {
		return nil, 0, err
	}

	return sellouts, total, nil
}
