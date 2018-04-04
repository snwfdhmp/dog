package apis

import (
	"net/http"
	"testing"

	"{{.PackagePath}}/daos"
	"{{.PackagePath}}/services"
	"{{.PackagePath}}/testdata"
)

func Test{{.ResourceName}}(t *testing.T) {
	testdata.ResetDB()
	router := newRouter()
	Serve{{.ResourceName}}Resource(&router.RouteGroup, services.New{{unexported .ResourceName}}Service(daos.New{{.ResourceName}}DAO()))

	notFoundError := `{"error_code":"NOT_FOUND", "message":"NOT_FOUND"}`
	nameRequiredError := `{"error_code":"INVALID_DATA","message":"INVALID_DATA","details":[{"field":"name","error":"cannot be blank"}]}`

	runAPITests(t, router, []apiTestCase{
		{"t1 - get an {{unexported .ResourceName}}", "GET", "/{{unexported .ResourceName}}s/2", "", http.StatusOK, `{"id":2,"name":"Accept"}`},
		{"t2 - get a nonexisting {{unexported .ResourceName}}", "GET", "/{{unexported .ResourceName}}s/99999", "", http.StatusNotFound, notFoundError},
		{"t3 - create an {{unexported .ResourceName}}", "POST", "/{{unexported .ResourceName}}s", `{"name":"Qiang"}`, http.StatusOK, `{"id": 276, "name":"Qiang"}`},
		{"t4 - create an {{unexported .ResourceName}} with validation error", "POST", "/{{unexported .ResourceName}}s", `{"name":""}`, http.StatusBadRequest, nameRequiredError},
		{"t5 - update an {{unexported .ResourceName}}", "PUT", "/{{unexported .ResourceName}}s/2", `{"name":"Qiang"}`, http.StatusOK, `{"id": 2, "name":"Qiang"}`},
		{"t6 - update an {{unexported .ResourceName}} with validation error", "PUT", "/{{unexported .ResourceName}}s/2", `{"name":""}`, http.StatusBadRequest, nameRequiredError},
		{"t7 - update a nonexisting {{unexported .ResourceName}}", "PUT", "/{{unexported .ResourceName}}s/99999", "{}", http.StatusNotFound, notFoundError},
		{"t8 - delete an {{unexported .ResourceName}}", "DELETE", "/{{unexported .ResourceName}}s/2", ``, http.StatusOK, `{"id": 2, "name":"Qiang"}`},
		{"t9 - delete a nonexisting {{unexported .ResourceName}}", "DELETE", "/{{unexported .ResourceName}}s/99999", "", http.StatusNotFound, notFoundError},
		{"t10 - get a list of {{unexported .ResourceName}}s", "GET", "/{{unexported .ResourceName}}s?page=3&per_page=2", "", http.StatusOK, `{"page":3,"per_page":2,"page_count":138,"total_count":275,"items":[{"id":6,"name":"Antônio Carlos Jobim"},{"id":7,"name":"Apocalyptica"}]}`},
	})
}
