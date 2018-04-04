package apis

import (
	"strconv"

	"github.com/go-ozzo/ozzo-routing"
	"{{.PackagePath}}/app"
	"{{.PackagePath}}/models"
)

type (
	// {{unexported .ResourceName}}Service specifies the interface for the {{unexported .ResourceName}} service needed by {{unexported .ResourceName}}Resource.
	{{unexported .ResourceName}}Service interface {
		// Get returns the {{unexported .ResourceName}} with the specified {{unexported .ResourceName}}Id.
		Get(rs app.RequestScope, id int) (*models.{{.ResourceName}}, error)
		Query(rs app.RequestScope, offset, limit int) ([]models.{{.ResourceName}}, error)
		Count(rs app.RequestScope) (int, error)
		Create(rs app.RequestScope, model *models.{{.ResourceName}}) (*models.{{.ResourceName}}, error)
		Update(rs app.RequestScope, id int, model *models.{{.ResourceName}}) (*models.{{.ResourceName}}, error)
		Delete(rs app.RequestScope, id int) (*models.{{.ResourceName}}, error)
	}

	// {{unexported .ResourceName}}Resource defines the handlers for the CRUD APIs.
	{{unexported .ResourceName}}Resource struct {
		service {{unexported .ResourceName}}Service
	}
)

// Serve{{.ResourceName}}Resource sets up the routing of {{unexported .ResourceName}} endpoints and the corresponding handlers.
func Serve{{.ResourceName}}Resource(rg *routing.RouteGroup, service {{unexported .ResourceName}}Service) {
	r := &{{unexported .ResourceName}}Resource{service}
	rg.Get("/{{unexported .ResourceName}}s/<id>", r.get)
	rg.Get("/{{unexported .ResourceName}}s", r.query)
	rg.Post("/{{unexported .ResourceName}}s", r.create)
	rg.Put("/{{unexported .ResourceName}}s/<id>", r.update)
	rg.Delete("/{{unexported .ResourceName}}s/<id>", r.delete)
}

func (r *{{unexported .ResourceName}}Resource) get(c *routing.Context) error {
	id, err := strconv.Atoi(c.Param("id"))
	if err != nil {
		return err
	}

	response, err := r.service.Get(app.GetRequestScope(c), id)
	if err != nil {
		return err
	}

	return c.Write(response)
}

func (r *{{unexported .ResourceName}}Resource) query(c *routing.Context) error {
	rs := app.GetRequestScope(c)
	count, err := r.service.Count(rs)
	if err != nil {
		return err
	}
	paginatedList := getPaginatedListFromRequest(c, count)
	items, err := r.service.Query(app.GetRequestScope(c), paginatedList.Offset(), paginatedList.Limit())
	if err != nil {
		return err
	}
	paginatedList.Items = items
	return c.Write(paginatedList)
}

func (r *{{unexported .ResourceName}}Resource) create(c *routing.Context) error {
	var model models.{{unexported .ResourceName}}
	if err := c.Read(&model); err != nil {
		return err
	}
	response, err := r.service.Create(app.GetRequestScope(c), &model)
	if err != nil {
		return err
	}

	return c.Write(response)
}

func (r *{{unexported .ResourceName}}Resource) update(c *routing.Context) error {
	id, err := strconv.Atoi(c.Param("id"))
	if err != nil {
		return err
	}

	rs := app.GetRequestScope(c)

	model, err := r.service.Get(rs, id)
	if err != nil {
		return err
	}

	if err := c.Read(model); err != nil {
		return err
	}

	response, err := r.service.Update(rs, id, model)
	if err != nil {
		return err
	}

	return c.Write(response)
}

func (r *{{unexported .ResourceName}}Resource) delete(c *routing.Context) error {
	id, err := strconv.Atoi(c.Param("id"))
	if err != nil {
		return err
	}

	response, err := r.service.Delete(app.GetRequestScope(c), id)
	if err != nil {
		return err
	}

	return c.Write(response)
}
