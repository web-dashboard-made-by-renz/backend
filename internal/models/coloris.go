package models

import (
	"time"

	"go.mongodb.org/mongo-driver/bson/primitive"
)

type Coloris struct {
	ID                   primitive.ObjectID `json:"id" bson:"_id,omitempty"`
	Timestamp            time.Time          `json:"timestamp" bson:"timestamp"`
	Bulan                string             `json:"bulan" bson:"bulan"`
	Region               string             `json:"region" bson:"region"`
	Cabang               string             `json:"cabang" bson:"cabang"`
	Materi               string             `json:"materi" bson:"materi"`
	NamaAtasanLangsung   string             `json:"nama_atasan_langsung" bson:"nama_atasan_langsung"`
	NamaToko             string             `json:"nama_toko" bson:"nama_toko"`
	NamaLengkapSesuaiKTP string             `json:"nama_lengkap_sesuai_ktp" bson:"nama_lengkap_sesuai_ktp"`
	NilaiPG              float64            `json:"nilai_pg" bson:"nilai_pg"`
	NilaiAkhir           float64            `json:"nilai_akhir" bson:"nilai_akhir"`
	Total                float64            `json:"total" bson:"total"`
	CreatedAt            time.Time          `json:"created_at" bson:"created_at"`
	UpdatedAt            time.Time          `json:"updated_at" bson:"updated_at"`
}

type ColorisCreateRequest struct {
	Timestamp            string  `json:"timestamp" binding:"required"`
	Bulan                string  `json:"bulan" binding:"required"`
	Region               string  `json:"region" binding:"required"`
	Cabang               string  `json:"cabang" binding:"required"`
	Materi               string  `json:"materi" binding:"required"`
	NamaAtasanLangsung   string  `json:"nama_atasan_langsung" binding:"required"`
	NamaToko             string  `json:"nama_toko" binding:"required"`
	NamaLengkapSesuaiKTP string  `json:"nama_lengkap_sesuai_ktp" binding:"required"`
	NilaiPG              float64 `json:"nilai_pg" binding:"required"`
	NilaiAkhir           float64 `json:"nilai_akhir" binding:"required"`
	Total                float64 `json:"total" binding:"required"`
}

type ColorisListResponse struct {
	Data       []Coloris `json:"data"`
	Total      int64     `json:"total"`
	Page       int       `json:"page"`
	PerPage    int       `json:"per_page"`
	TotalPages int       `json:"total_pages"`
}
