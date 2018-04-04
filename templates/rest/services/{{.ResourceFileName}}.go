package services

import (
	"{{.PackagePath}}/app"
	"{{.PackagePath}}/models"
)

// {{unexported .ResourceName}}DAO specifies the interface of the {{unexported .ResourceName}} DAO needed by {{unexported .ResourceName}}Service.
type {{unexported .ResourceName}}DAO interface {
	// Get returns the {{unexported .ResourceName}} with the specified {{unexported .ResourceName}} ID.
	Get(rs app.RequestScope, id int) (*models.{{unexported .ResourceName}}, error)
	// Count returns the number of {{unexported .ResourceName}}s.
	Count(rs app.RequestScope) (int, error)
	// Query returns the list of {{unexported .ResourceName}}s with the given offset and limit.
	Query(rs app.RequestScope, offset, limit int) ([]models.{{unexported .ResourceName}}, error)
	// Create saves a new {{unexported .ResourceName}} in the storage.
	Create(rs app.RequestScope, {{unexported .ResourceName}} *models.{{unexported .ResourceName}}) error
	// Update updates the {{unexported .ResourceName}} with given ID in the storage.
	Update(rs app.RequestScope, id int, {{unexported .ResourceName}} *models.{{unexported .ResourceName}}) error
	// Delete removes the {{unexported .ResourceName}} with given ID from the storage.
	Delete(rs app.RequestScope, id int) error
}

// {{unexported .ResourceName}}Service provides services related with {{unexported .ResourceName}}s.
type {{unexported .ResourceName}}Service struct {
	dao {{unexported .ResourceName}}DAO
}

// New{{unexported .ResourceName}}Service creates a new {{unexported .ResourceName}}Service with the given {{unexported .ResourceName}} DAO.
func New{{unexported .ResourceName}}Service(dao {{unexported .ResourceName}}DAO) *{{unexported .ResourceName}}Service {
	return &{{unexported .ResourceName}}Service{dao}
}

// Get returns the {{unexported .ResourceName}} with the specified the {{unexported .ResourceName}} ID.
func (s *{{unexported .ResourceName}}Service) Get(rs app.RequestScope, id int) (*models.{{unexported .ResourceName}}, error) {
	return s.dao.Get(rs, id)
}

// Create creates a new {{unexported .ResourceName}}.
func (s *{{unexported .ResourceName}}Service) Create(rs app.RequestScope, model *models.{{unexported .ResourceName}}) (*models.{{unexported .ResourceName}}, error) {
	if err := model.Validate(); err != nil {
		return nil, err
	}
	if err := s.dao.Create(rs, model); err != nil {
		return nil, err
	}
	return s.dao.Get(rs, model.Id)
}

// Update updates the {{unexported .ResourceName}} with the specified ID.
func (s *{{unexported .ResourceName}}Service) Update(rs app.RequestScope, id int, model *models.{{unexported .ResourceName}}) (*models.{{unexported .ResourceName}}, error) {
	if err := model.Validate(); err != nil {
		return nil, err
	}
	if err := s.dao.Update(rs, id, model); err != nil {
		return nil, err
	}
	return s.dao.Get(rs, id)
}

// Delete deletes the {{unexported .ResourceName}} with the specified ID.
func (s *{{unexported .ResourceName}}Service) Delete(rs app.RequestScope, id int) (*models.{{unexported .ResourceName}}, error) {
	{{unexported .ResourceName}}, err := s.dao.Get(rs, id)
	if err != nil {
		return nil, err
	}
	err = s.dao.Delete(rs, id)
	return {{unexported .ResourceName}}, err
}

// Count returns the number of {{unexported .ResourceName}}s.
func (s *{{unexported .ResourceName}}Service) Count(rs app.RequestScope) (int, error) {
	return s.dao.Count(rs)
}

// Query returns the {{unexported .ResourceName}}s with the specified offset and limit.
func (s *{{unexported .ResourceName}}Service) Query(rs app.RequestScope, offset, limit int) ([]models.{{unexported .ResourceName}}, error) {
	return s.dao.Query(rs, offset, limit)
}
