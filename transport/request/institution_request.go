package request

import (
	validation "github.com/go-ozzo/ozzo-validation"
)

// CreateInstitutionReq represent create request body
type CreateInstitutionReq struct {
	Name  string `json:"name"`
	Alias string `json:"alias"`
	Type  string `json:"type"`
}

func (request CreateInstitutionReq) Validate() error {
	return validation.ValidateStruct(
		&request,
		validation.Field(&request.Name, validation.Required),
		validation.Field(&request.Alias, validation.Required),
	)
}

// Update request body
type UpdateInstitutionReq struct {
	ID       string `param:"id"`
	Name     string `json:"name"`
	Alias    string `json:"alias"`
	Type     string `json:"type"`
	IsActive *bool  `json:"is_active"`
}

func (request UpdateInstitutionReq) Validate() error {
	return validation.ValidateStruct(
		&request,
		validation.Field(&request.ID, validation.Required),
	)
}

// Delete request body
type DeleteInstitutionReq struct {
	ID string `param:"id"`
}

func (request DeleteInstitutionReq) Validate() error {
	return validation.ValidateStruct(
		&request,
		validation.Field(&request.ID, validation.Required),
	)
}

// GetList request body
type GetListInstitutionReq struct {
	Name     string `query:"name"`
	Type     string `query:"type"`
	Alias    string `query:"alias"`
	IsActive *bool  `query:"is_active"`
	PerPage  int64  `query:"per_page"`
	Page     int64  `query:"page"`
}

func (request GetListInstitutionReq) Validate() error {
	return validation.ValidateStruct(
		&request,
	)
}

// GetDetail request body
type GetDetailInstitutionReq struct {
	ID string `param:"id"`
}

func (request GetDetailInstitutionReq) Validate() error {
	return validation.ValidateStruct(
		&request,
	)
}
