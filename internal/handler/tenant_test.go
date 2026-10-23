package handler

import (
	"context"
	"bytes"
	"errors"
	"net/http"
	"net/http/httptest"
	"os"
	"path/filepath"
	"testing"

	"github.com/gin-gonic/gin"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"

	"condominium-manager/mocks"
	"condominium-manager/internal/model"
	"condominium-manager/internal/repository"
)

func setup(t *testing.T) (*gin.Engine, *mocks.MockTenantRepository) {
	t.Helper()
	gin.SetMode(gin.TestMode)
	repo := mocks.NewMockTenantRepository(t)
	r := gin.New()
	Register(r, repo)
	return r, repo
}

func do(r *gin.Engine, method, url, body string) *httptest.ResponseRecorder {
	req := httptest.NewRequest(method, url, bytes.NewBufferString(body))
	req.Header.Set("Content-Type", "application/json")
	w := httptest.NewRecorder()
	r.ServeHTTP(w, req)
	return w
}

func TestOpenAPI(t *testing.T) {
	gin.SetMode(gin.TestMode)
	path := filepath.Join(t.TempDir(), "spec.yml")
	assert.NoError(t, os.WriteFile(path, []byte("openapi: '3.0.3'\n"), 0o600))
	r := gin.New()
	RegisterSpec(r, path)

	w := do(r, http.MethodGet, "/v1/openapi.yaml", "")

	assert.Equal(t, http.StatusOK, w.Code)
	assert.Contains(t, w.Header().Get("Content-Type"), "application/yaml")
	assert.Contains(t, w.Body.String(), "openapi: '3.0.3'")

	docs := do(r, http.MethodGet, "/docs/", "")
	assert.Equal(t, http.StatusOK, docs.Code)
	assert.Contains(t, docs.Header().Get("Content-Type"), "text/html")
	assert.Contains(t, docs.Body.String(), "/v1/openapi.yaml")
}

func TestList(t *testing.T) {
	t.Run("ok", func(t *testing.T) {
		r, repo := setup(t)
		repo.EXPECT().List(mock.Anything).
			Return([]model.Tenant{{ID: 1, Name: "Ann", Lastname: "Lee", Age: 30}}, nil)

		w := do(r, http.MethodGet, "/v1/tenants", "")

		assert.Equal(t, http.StatusOK, w.Code)
		assert.JSONEq(t, `[{"id":1,"name":"Ann","lastname":"Lee","age":30}]`, w.Body.String())
	})

	t.Run("db error", func(t *testing.T) {
		r, repo := setup(t)
		repo.EXPECT().List(mock.Anything).Return(nil, errors.New("boom"))

		w := do(r, http.MethodGet, "/v1/tenants", "")

		assert.Equal(t, http.StatusInternalServerError, w.Code)
	})
}

func TestCreate(t *testing.T) {
	t.Run("created, client id ignored", func(t *testing.T) {
		r, repo := setup(t)
		repo.EXPECT().Create(mock.Anything, mock.MatchedBy(func(tn *model.Tenant) bool {
			return tn.ID == 0 && tn.Name == "Ann" && tn.Lastname == "Lee"
		})).Run(func(_ context.Context, tn *model.Tenant) { tn.ID = 7 }).Return(nil)

		w := do(r, http.MethodPost, "/v1/tenants", `{"id":99,"name":"Ann","lastname":"Lee","age":30}`)

		assert.Equal(t, http.StatusCreated, w.Code)
		assert.JSONEq(t, `{"id":7,"name":"Ann","lastname":"Lee","age":30}`, w.Body.String())
	})

	t.Run("missing required fields", func(t *testing.T) {
		r, _ := setup(t)

		assert.Equal(t, http.StatusBadRequest, do(r, http.MethodPost, "/v1/tenants", `{"name":"Ann"}`).Code)
		assert.Equal(t, http.StatusBadRequest, do(r, http.MethodPost, "/v1/tenants", `{"lastname":"Lee"}`).Code)
	})

	t.Run("invalid json", func(t *testing.T) {
		r, _ := setup(t)

		assert.Equal(t, http.StatusBadRequest, do(r, http.MethodPost, "/v1/tenants", `{`).Code)
	})

	t.Run("already exists", func(t *testing.T) {
		r, repo := setup(t)
		repo.EXPECT().Create(mock.Anything, mock.Anything).Return(repository.ErrAlreadyExists)

		w := do(r, http.MethodPost, "/v1/tenants", `{"name":"Ann","lastname":"Lee"}`)

		assert.Equal(t, http.StatusConflict, w.Code)
	})

	t.Run("db error", func(t *testing.T) {
		r, repo := setup(t)
		repo.EXPECT().Create(mock.Anything, mock.Anything).Return(errors.New("boom"))

		w := do(r, http.MethodPost, "/v1/tenants", `{"name":"Ann","lastname":"Lee"}`)

		assert.Equal(t, http.StatusInternalServerError, w.Code)
	})
}

func TestSearch(t *testing.T) {
	t.Run("by all filters", func(t *testing.T) {
		r, repo := setup(t)
		id := int32(1)
		repo.EXPECT().Search(mock.Anything, repository.TenantFilter{ID: &id, Name: "Ann", Lastname: "Lee"}).
			Return([]model.Tenant{{ID: 1, Name: "Ann", Lastname: "Lee"}}, nil)

		w := do(r, http.MethodGet, "/v1/tenants/search?id=1&name=Ann&lastname=Lee", "")

		assert.Equal(t, http.StatusOK, w.Code)
		assert.JSONEq(t, `[{"id":1,"name":"Ann","lastname":"Lee","age":0}]`, w.Body.String())
	})

	t.Run("by lastname only", func(t *testing.T) {
		r, repo := setup(t)
		repo.EXPECT().Search(mock.Anything, repository.TenantFilter{Lastname: "Lee"}).
			Return([]model.Tenant{{ID: 1, Name: "Ann", Lastname: "Lee"}}, nil)

		assert.Equal(t, http.StatusOK, do(r, http.MethodGet, "/v1/tenants/search?lastname=Lee", "").Code)
	})

	t.Run("not found", func(t *testing.T) {
		r, repo := setup(t)
		repo.EXPECT().Search(mock.Anything, mock.Anything).Return([]model.Tenant{}, nil)

		assert.Equal(t, http.StatusNotFound, do(r, http.MethodGet, "/v1/tenants/search?name=Zed", "").Code)
	})

	t.Run("invalid id", func(t *testing.T) {
		r, _ := setup(t)

		assert.Equal(t, http.StatusBadRequest, do(r, http.MethodGet, "/v1/tenants/search?id=abc", "").Code)
	})

	t.Run("db error", func(t *testing.T) {
		r, repo := setup(t)
		repo.EXPECT().Search(mock.Anything, mock.Anything).Return(nil, errors.New("boom"))

		assert.Equal(t, http.StatusInternalServerError, do(r, http.MethodGet, "/v1/tenants/search?name=Ann", "").Code)
	})
}

func TestDelete(t *testing.T) {
	t.Run("deleted", func(t *testing.T) {
		r, repo := setup(t)
		repo.EXPECT().Delete(mock.Anything, int32(5)).Return(nil)

		w := do(r, http.MethodDelete, "/v1/tenants/5", "")

		assert.Equal(t, http.StatusNoContent, w.Code)
		assert.Empty(t, w.Body.String())
	})

	t.Run("not found", func(t *testing.T) {
		r, repo := setup(t)
		repo.EXPECT().Delete(mock.Anything, int32(5)).Return(repository.ErrNotFound)

		assert.Equal(t, http.StatusNotFound, do(r, http.MethodDelete, "/v1/tenants/5", "").Code)
	})

	t.Run("invalid id", func(t *testing.T) {
		r, _ := setup(t)

		assert.Equal(t, http.StatusBadRequest, do(r, http.MethodDelete, "/v1/tenants/abc", "").Code)
	})

	t.Run("db error", func(t *testing.T) {
		r, repo := setup(t)
		repo.EXPECT().Delete(mock.Anything, int32(5)).Return(errors.New("boom"))

		assert.Equal(t, http.StatusInternalServerError, do(r, http.MethodDelete, "/v1/tenants/5", "").Code)
	})
}

func TestUpdate(t *testing.T) {
	const body = `{"id":99,"name":"Ann","lastname":"Lee","age":31}`

	t.Run("updated, path id wins", func(t *testing.T) {
		r, repo := setup(t)
		repo.EXPECT().Update(mock.Anything, &model.Tenant{ID: 5, Name: "Ann", Lastname: "Lee", Age: 31}).Return(nil)

		w := do(r, http.MethodPut, "/v1/tenants/5", body)

		assert.Equal(t, http.StatusOK, w.Code)
		assert.JSONEq(t, `{"id":5,"name":"Ann","lastname":"Lee","age":31}`, w.Body.String())
	})

	t.Run("not found", func(t *testing.T) {
		r, repo := setup(t)
		repo.EXPECT().Update(mock.Anything, mock.Anything).Return(repository.ErrNotFound)

		assert.Equal(t, http.StatusNotFound, do(r, http.MethodPut, "/v1/tenants/5", body).Code)
	})

	t.Run("conflict", func(t *testing.T) {
		r, repo := setup(t)
		repo.EXPECT().Update(mock.Anything, mock.Anything).Return(repository.ErrAlreadyExists)

		assert.Equal(t, http.StatusConflict, do(r, http.MethodPut, "/v1/tenants/5", body).Code)
	})

	t.Run("invalid id", func(t *testing.T) {
		r, _ := setup(t)

		assert.Equal(t, http.StatusBadRequest, do(r, http.MethodPut, "/v1/tenants/abc", body).Code)
	})

	t.Run("invalid body", func(t *testing.T) {
		r, _ := setup(t)

		assert.Equal(t, http.StatusBadRequest, do(r, http.MethodPut, "/v1/tenants/5", `{"name":"Ann"}`).Code)
	})

	t.Run("db error", func(t *testing.T) {
		r, repo := setup(t)
		repo.EXPECT().Update(mock.Anything, mock.Anything).Return(errors.New("boom"))

		assert.Equal(t, http.StatusInternalServerError, do(r, http.MethodPut, "/v1/tenants/5", body).Code)
	})
}
