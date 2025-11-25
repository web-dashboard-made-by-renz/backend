package models

import (
	"time"

	"go.mongodb.org/mongo-driver/bson/primitive"
)

type Sellout struct {
	ID           primitive.ObjectID `json:"id" bson:"_id,omitempty"`
	Tahun        int                `json:"tahun" bson:"tahun"`
	Bulan        int                `json:"bulan" bson:"bulan"`
	Reg          string             `json:"reg" bson:"reg"`
	Cabang       string             `json:"cabang" bson:"cabang"`
	Outlet       string             `json:"outlet" bson:"outlet"`
	NamaColorist string             `json:"nama_colorist" bson:"nama_colorist"`
	NoReg        string             `json:"no_reg" bson:"no_reg"`
	CHL          string             `json:"chl" bson:"chl"`
	SelloutTT    float64            `json:"sellout_tt" bson:"sellout_tt"`
	SelloutRM    float64            `json:"sellout_rm" bson:"sellout_rm"`
	TotalSellout float64            `json:"total_sellout" bson:"total_sellout"`
	CreatedAt    time.Time          `json:"created_at" bson:"created_at"`
	UpdatedAt    time.Time          `json:"updated_at" bson:"updated_at"`
}

type SelloutCreateRequest struct {
	Tahun        int     `json:"tahun" binding:"required"`
	Bulan        int     `json:"bulan" binding:"required,min=1,max=12"`
	Reg          string  `json:"reg" binding:"required"`
	Cabang       string  `json:"cabang" binding:"required"`
	Outlet       string  `json:"outlet" binding:"required"`
	NamaColorist string  `json:"nama_colorist" binding:"required"`
	NoReg        string  `json:"no_reg" binding:"required"`
	CHL          string  `json:"chl" binding:"required"`
	SelloutTT    float64 `json:"sellout_tt" binding:"required"`
	SelloutRM    float64 `json:"sellout_rm" binding:"required"`
	TotalSellout float64 `json:"total_sellout" binding:"required"`
}

type SelloutListResponse struct {
	Data       []Sellout `json:"data"`
	Total      int64     `json:"total"`
	Page       int       `json:"page"`
	PerPage    int       `json:"per_page"`
	TotalPages int       `json:"total_pages"`
}
