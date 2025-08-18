package domain

import (
	"context"
	"prakarsa-app/transport/request"
	"prakarsa-app/transport/response"
)

type Institution struct {
	ID        string `json:"id"`
	Name      string `json:"name"`
	Alias     string `json:"alias"`
	Type      string `json:"type"`
	IsActive  *bool  `json:"is_active"`
	CreatedBy string `json:"created_by"`
	CreatedAt int64  `json:"created_at"`
	UpdatedBy string `json:"updated_by"`
	UpdatedAt int64  `json:"updated_at"`
	DeletedAt int64  `json:"deleted_at"`
}

// InstitutionRepository represent the institution repository contract
type InstitutionRepository interface {
	Create(ctx context.Context, institution *Institution) error
	Update(ctx context.Context, institution *Institution) error
	Delete(ctx context.Context, institution *Institution) (int64, error)
	GetList(ctx context.Context, request *request.GetListInstitutionReq) ([]response.GetListInstitutionRes, response.MetaRes, error)
	GetDetail(ctx context.Context, request *request.GetDetailInstitutionReq) (Institution, error)
}

// InstitutionUsecase represent the institution usecase contract
type InstitutionUsecase interface {
	Create(ctx context.Context, request *request.CreateInstitutionReq) (response.CreateInstitutionRes, error)
	Update(ctx context.Context, request *request.UpdateInstitutionReq) error
	Delete(ctx context.Context, request *request.DeleteInstitutionReq) (int64, error)
	GetList(ctx context.Context, request *request.GetListInstitutionReq) ([]response.GetListInstitutionRes, response.MetaRes, error)
	GetDetail(ctx context.Context, request *request.GetDetailInstitutionReq) (Institution, error)
}
