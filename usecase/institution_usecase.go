package usecase

import (
	"context"
	"prakarsa-app/transport/response"
	"time"

	"prakarsa-app/domain"
	"prakarsa-app/repository/redis"
	"prakarsa-app/transport/request"

	"github.com/google/uuid"
)

type InstitutionUsecase struct {
	institutionRepo domain.InstitutionRepository
	redisRepo       redis.RedisRepository
	ctxTimeout      time.Duration
}

// NewInstitutionUsecase will create new an tagUsecase object representation of ThreadUsecase interface
func NewInstitutionUsecase(institutionRepo domain.InstitutionRepository, redisRepo redis.RedisRepository, ctxTimeout time.Duration) *InstitutionUsecase {
	return &InstitutionUsecase{
		institutionRepo: institutionRepo,
		redisRepo:       redisRepo,
		ctxTimeout:      ctxTimeout,
	}
}

func (u *InstitutionUsecase) Create(c context.Context, request *request.CreateInstitutionReq) (res response.CreateInstitutionRes, err error) {
	ctx, cancel := context.WithTimeout(c, u.ctxTimeout)
	defer cancel()

	// Create Payload
	institutionUUID := uuid.NewString()
	t := true
	institutionPayload := &domain.Institution{
		ID:        institutionUUID,
		Name:      request.Name,
		Alias:     request.Alias,
		Type:      request.Type,
		IsActive:  &t,
		CreatedBy: request.UserID,
		CreatedAt: time.Now().Unix(),
	}

	// Response Payload
	res.ID = institutionUUID

	err = u.institutionRepo.Create(ctx, institutionPayload)
	return
}

func (u *InstitutionUsecase) Update(c context.Context, request *request.UpdateInstitutionReq) (err error) {
	ctx, cancel := context.WithTimeout(c, u.ctxTimeout)
	defer cancel()

	// Update Payload
	institutionPayload := &domain.Institution{
		ID:        request.ID,
		UpdatedBy: request.UserID,
		UpdatedAt: time.Now().Unix(),
	}

	if request.Name != "" {
		institutionPayload.Name = request.Name
	}

	if request.Alias != "" {
		institutionPayload.Alias = request.Alias
	}

	if request.Type != "" {
		institutionPayload.Type = request.Type
	}

	if request.IsActive != nil {
		institutionPayload.IsActive = request.IsActive
	}

	err = u.institutionRepo.Update(ctx, institutionPayload)
	return
}

func (u *InstitutionUsecase) Delete(c context.Context, request *request.DeleteInstitutionReq) (rowsAffected int64, err error) {
	ctx, cancel := context.WithTimeout(c, u.ctxTimeout)
	defer cancel()

	threadPayload := &domain.Institution{
		ID: request.ID,
	}

	rowsAffected, err = u.institutionRepo.Delete(ctx, threadPayload)
	return
}

func (u *InstitutionUsecase) GetList(c context.Context, request *request.GetListInstitutionReq) (res []response.GetListInstitutionRes, meta response.MetaRes, err error) {
	ctx, cancel := context.WithTimeout(c, u.ctxTimeout)
	defer cancel()

	res, meta, err = u.institutionRepo.GetList(ctx, request)
	return
}

func (u *InstitutionUsecase) GetDetail(c context.Context, request *request.GetDetailInstitutionReq) (res domain.Institution, err error) {
	ctx, cancel := context.WithTimeout(c, u.ctxTimeout)
	defer cancel()

	res, err = u.institutionRepo.GetDetail(ctx, request)
	return
}
