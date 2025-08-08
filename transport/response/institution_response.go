package response

// CreateInstitutionRes represent create request body
type CreateInstitutionRes struct {
	ID string `json:"id"`
}

// Get List response
type GetListInstitutionRes struct {
	ID       string `json:"id"`
	Name     string `json:"name"`
	Alias    string `json:"alias"`
	Type     string `json:"type"`
	IsActive bool   `json:"is_active"`
}

// Get Detail response
type GetDetailInstitutionRes struct {
	ID        string `json:"id"`
	Name      string `json:"name"`
	Alias     string `json:"alias"`
	Type      string `json:"type"`
	IsActive  bool   `json:"is_active"`
	CreatedAt int64  `json:"created_at"`
	CreatedBy string `json:"created_by"`
	UpdatedAt int64  `json:"updated_at"`
	UpdatedBy string `json:"updated_by"`
}
