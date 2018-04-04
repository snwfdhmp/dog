package daos

import (
	"{{.PackagePath}}/app"
	"{{.PackagePath}}/models"
)

// {{.ResourceName}}DAO persists {{unexported .ResourceName}} data in database
type {{.ResourceName}}DAO struct{}

// New{{.ResourceName}}DAO creates a new {{.ResourceName}}DAO
func New{{.ResourceName}}DAO() *{{.ResourceName}}DAO {
	return &{{.ResourceName}}DAO{}
}

// Get reads the {{unexported .ResourceName}} with the specified ID from the database.
func (dao *{{.ResourceName}}DAO) Get(rs app.RequestScope, id int) (*models.{{.ResourceName}}, error) {
	var {{unexported .ResourceName}} models.{{.ResourceName}}
	err := rs.Tx().Select().Model(id, &{{unexported .ResourceName}})
	return &{{unexported .ResourceName}}, err
}

// Create saves a new {{unexported .ResourceName}} record in the database.
// The {{.ResourceName}}.Id field will be populated with an automatically generated ID upon successful saving.
func (dao *{{.ResourceName}}DAO) Create(rs app.RequestScope, {{unexported .ResourceName}} *models.{{.ResourceName}}) error {
	{{unexported .ResourceName}}.Id = 0
	return rs.Tx().Model({{unexported .ResourceName}}).Insert()
}

// Update saves the changes to an {{unexported .ResourceName}} in the database.
func (dao *{{.ResourceName}}DAO) Update(rs app.RequestScope, id int, {{unexported .ResourceName}} *models.{{.ResourceName}}) error {
	if _, err := dao.Get(rs, id); err != nil {
		return err
	}
	{{unexported .ResourceName}}.Id = id
	return rs.Tx().Model({{unexported .ResourceName}}).Exclude("Id").Update()
}

// Delete deletes an {{unexported .ResourceName}} with the specified ID from the database.
func (dao *{{.ResourceName}}DAO) Delete(rs app.RequestScope, id int) error {
	{{unexported .ResourceName}}, err := dao.Get(rs, id)
	if err != nil {
		return err
	}
	return rs.Tx().Model({{unexported .ResourceName}}).Delete()
}

// Count returns the number of the {{unexported .ResourceName}} records in the database.
func (dao *{{.ResourceName}}DAO) Count(rs app.RequestScope) (int, error) {
	var count int
	err := rs.Tx().Select("COUNT(*)").From("{{unexported .ResourceName}}").Row(&count)
	return count, err
}

// Query retrieves the {{unexported .ResourceName}} records with the specified offset and limit from the database.
func (dao *{{.ResourceName}}DAO) Query(rs app.RequestScope, offset, limit int) ([]models.{{.ResourceName}}, error) {
	{{unexported .ResourceName}}s := []models.{{.ResourceName}}{}
	err := rs.Tx().Select().OrderBy("id").Offset(int64(offset)).Limit(int64(limit)).All(&{{unexported .ResourceName}}s)
	return {{unexported .ResourceName}}s, err
}
