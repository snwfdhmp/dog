package daos

import (
	"testing"

	"{{.PackagePath}}/app"
	"{{.PackagePath}}/models"
	"{{.PackagePath}}/testdata"
	"github.com/stretchr/testify/assert"
)

func Test{{.ResourceName}}DAO(t *testing.T) {
	db := testdata.ResetDB()
	dao := New{{.ResourceName}}DAO()

	{
		// Get
		testDBCall(db, func(rs app.RequestScope) {
			{{unexported .ResourceName}}, err := dao.Get(rs, 2)
			assert.Nil(t, err)
			if assert.NotNil(t, {{unexported .ResourceName}}) {
				assert.Equal(t, 2, {{unexported .ResourceName}}.Id)
			}
		})
	}

	{
		// Create
		testDBCall(db, func(rs app.RequestScope) {
			{{unexported .ResourceName}} := &models.{{.ResourceName}}{
				Id:   1000,
				Name: "tester",
			}
			err := dao.Create(rs, {{unexported .ResourceName}})
			assert.Nil(t, err)
			assert.NotEqual(t, 1000, {{unexported .ResourceName}}.Id)
			assert.NotZero(t, {{unexported .ResourceName}}.Id)
		})
	}

	{
		// Update
		testDBCall(db, func(rs app.RequestScope) {
			{{unexported .ResourceName}} := &models.{{.ResourceName}}{
				Id:   2,
				Name: "tester",
			}
			err := dao.Update(rs, {{unexported .ResourceName}}.Id, {{unexported .ResourceName}})
			assert.Nil(t, err)
		})
	}

	{
		// Update with error
		testDBCall(db, func(rs app.RequestScope) {
			{{unexported .ResourceName}} := &models.{{.ResourceName}}{
				Id:   2,
				Name: "tester",
			}
			err := dao.Update(rs, 99999, {{unexported .ResourceName}})
			assert.NotNil(t, err)
		})
	}

	{
		// Delete
		testDBCall(db, func(rs app.RequestScope) {
			err := dao.Delete(rs, 2)
			assert.Nil(t, err)
		})
	}

	{
		// Delete with error
		testDBCall(db, func(rs app.RequestScope) {
			err := dao.Delete(rs, 99999)
			assert.NotNil(t, err)
		})
	}

	{
		// Query
		testDBCall(db, func(rs app.RequestScope) {
			{{unexported .ResourceName}}s, err := dao.Query(rs, 1, 3)
			assert.Nil(t, err)
			assert.Equal(t, 3, len({{unexported .ResourceName}}s))
		})
	}

	{
		// Count
		testDBCall(db, func(rs app.RequestScope) {
			count, err := dao.Count(rs)
			assert.Nil(t, err)
			assert.NotZero(t, count)
		})
	}
}
