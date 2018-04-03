package services

import (
	"errors"
	"testing"

	"{{.PackagePath}}/app"
	"{{.PackagePath}}/models"
	"github.com/stretchr/testify/assert"
)

func TestNew{{.ResourceName}}Service(t *testing.T) {
	dao := newMock{{.ResourceName}}DAO()
	s := New{{.ResourceName}}Service(dao)
	assert.Equal(t, dao, s.dao)
}

func Test{{.ResourceName}}Service_Get(t *testing.T) {
	s := New{{.ResourceName}}Service(newMock{{.ResourceName}}DAO())
	{{little .ResourceName}}, err := s.Get(nil, 1)
	if assert.Nil(t, err) && assert.NotNil(t, {{little .ResourceName}}) {
		assert.Equal(t, "aaa", {{little .ResourceName}}.Name)
	}

	{{little .ResourceName}}, err = s.Get(nil, 100)
	assert.NotNil(t, err)
}

func Test{{.ResourceName}}Service_Create(t *testing.T) {
	s := New{{.ResourceName}}Service(newMock{{.ResourceName}}DAO())
	{{little .ResourceName}}, err := s.Create(nil, &models.{{.ResourceName}}{
		Name: "ddd",
	})
	if assert.Nil(t, err) && assert.NotNil(t, {{little .ResourceName}}) {
		assert.Equal(t, 4, {{little .ResourceName}}.Id)
		assert.Equal(t, "ddd", {{little .ResourceName}}.Name)
	}

	// dao error
	_, err = s.Create(nil, &models.{{.ResourceName}}{
		Id:   100,
		Name: "ddd",
	})
	assert.NotNil(t, err)

	// validation error
	_, err = s.Create(nil, &models.{{.ResourceName}}{
		Name: "",
	})
	assert.NotNil(t, err)
}

func Test{{.ResourceName}}Service_Update(t *testing.T) {
	s := New{{.ResourceName}}Service(newMock{{.ResourceName}}DAO())
	{{little .ResourceName}}, err := s.Update(nil, 2, &models.{{.ResourceName}}{
		Name: "ddd",
	})
	if assert.Nil(t, err) && assert.NotNil(t, {{little .ResourceName}}) {
		assert.Equal(t, 2, {{little .ResourceName}}.Id)
		assert.Equal(t, "ddd", {{little .ResourceName}}.Name)
	}

	// dao error
	_, err = s.Update(nil, 100, &models.{{.ResourceName}}{
		Name: "ddd",
	})
	assert.NotNil(t, err)

	// validation error
	_, err = s.Update(nil, 2, &models.{{.ResourceName}}{
		Name: "",
	})
	assert.NotNil(t, err)
}

func Test{{.ResourceName}}Service_Delete(t *testing.T) {
	s := New{{.ResourceName}}Service(newMock{{.ResourceName}}DAO())
	{{little .ResourceName}}, err := s.Delete(nil, 2)
	if assert.Nil(t, err) && assert.NotNil(t, {{little .ResourceName}}) {
		assert.Equal(t, 2, {{little .ResourceName}}.Id)
		assert.Equal(t, "bbb", {{little .ResourceName}}.Name)
	}

	_, err = s.Delete(nil, 2)
	assert.NotNil(t, err)
}

func Test{{.ResourceName}}Service_Query(t *testing.T) {
	s := New{{.ResourceName}}Service(newMock{{.ResourceName}}DAO())
	result, err := s.Query(nil, 1, 2)
	if assert.Nil(t, err) {
		assert.Equal(t, 2, len(result))
	}
}

func newMock{{.ResourceName}}DAO() {{little .ResourceName}}DAO {
	return &mock{{.ResourceName}}DAO{
		records: []models.{{.ResourceName}}{
			{Id: 1, Name: "aaa"},
			{Id: 2, Name: "bbb"},
			{Id: 3, Name: "ccc"},
		},
	}
}

type mock{{.ResourceName}}DAO struct {
	records []models.{{.ResourceName}}
}

func (m *mock{{.ResourceName}}DAO) Get(rs app.RequestScope, id int) (*models.{{.ResourceName}}, error) {
	for _, record := range m.records {
		if record.Id == id {
			return &record, nil
		}
	}
	return nil, errors.New("not found")
}

func (m *mock{{.ResourceName}}DAO) Query(rs app.RequestScope, offset, limit int) ([]models.{{.ResourceName}}, error) {
	return m.records[offset : offset+limit], nil
}

func (m *mock{{.ResourceName}}DAO) Count(rs app.RequestScope) (int, error) {
	return len(m.records), nil
}

func (m *mock{{.ResourceName}}DAO) Create(rs app.RequestScope, {{little .ResourceName}} *models.{{.ResourceName}}) error {
	if {{little .ResourceName}}.Id != 0 {
		return errors.New("Id cannot be set")
	}
	{{little .ResourceName}}.Id = len(m.records) + 1
	m.records = append(m.records, *{{little .ResourceName}})
	return nil
}

func (m *mock{{.ResourceName}}DAO) Update(rs app.RequestScope, id int, {{little .ResourceName}} *models.{{.ResourceName}}) error {
	{{little .ResourceName}}.Id = id
	for i, record := range m.records {
		if record.Id == id {
			m.records[i] = *{{little .ResourceName}}
			return nil
		}
	}
	return errors.New("not found")
}

func (m *mock{{.ResourceName}}DAO) Delete(rs app.RequestScope, id int) error {
	for i, record := range m.records {
		if record.Id == id {
			m.records = append(m.records[:i], m.records[i+1:]...)
			return nil
		}
	}
	return errors.New("not found")
}
