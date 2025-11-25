package models

import (
	"time"

	"go.mongodb.org/mongo-driver/bson/primitive"
)

type Training struct {
	ID                   primitive.ObjectID `json:"id" bson:"_id,omitempty"`
	Timestamp            time.Time          `json:"timestamp" bson:"timestamp"`
	Bulan                string             `json:"bulan" bson:"bulan"`
	Region               string             `json:"region" bson:"region"`
	CabangArea           string             `json:"cabang_area" bson:"cabang_area"`
	NamaAtasanLangsung   string             `json:"nama_atasan_langsung" bson:"nama_atasan_langsung"`
	MateriPelatihan      string             `json:"materi_pelatihan" bson:"materi_pelatihan"`
	NamaLengkapSesuaiKTP string             `json:"nama_lengkap_sesuai_ktp" bson:"nama_lengkap_sesuai_ktp"`
	Jabatan              string             `json:"jabatan" bson:"jabatan"`
	TotalNilai           float64            `json:"total_nilai" bson:"total_nilai"`
	NilaiEssay           float64            `json:"nilai_essay" bson:"nilai_essay"`
	Total                float64            `json:"total" bson:"total"`
	CreatedAt            time.Time          `json:"created_at" bson:"created_at"`
	UpdatedAt            time.Time          `json:"updated_at" bson:"updated_at"`
}

type TrainingCreateRequest struct {
	Timestamp            string  `json:"timestamp" binding:"required"`
	Bulan                string  `json:"bulan" binding:"required"`
	Region               string  `json:"region" binding:"required"`
	CabangArea           string  `json:"cabang_area" binding:"required"`
	NamaAtasanLangsung   string  `json:"nama_atasan_langsung" binding:"required"`
	MateriPelatihan      string  `json:"materi_pelatihan" binding:"required"`
	NamaLengkapSesuaiKTP string  `json:"nama_lengkap_sesuai_ktp" binding:"required"`
	Jabatan              string  `json:"jabatan" binding:"required"`
	TotalNilai           float64 `json:"total_nilai" binding:"required"`
	NilaiEssay           float64 `json:"nilai_essay" binding:"required"`
	Total                float64 `json:"total" binding:"required"`
}

type TrainingListResponse struct {
	Data       []Training `json:"data"`
	Total      int64      `json:"total"`
	Page       int        `json:"page"`
	PerPage    int        `json:"per_page"`
	TotalPages int        `json:"total_pages"`
}
