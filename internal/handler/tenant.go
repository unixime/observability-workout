package handler

import (
	"errors"
	"net/http"
	"strconv"

	"github.com/gin-gonic/gin"

	"condominium-manager/internal/model"
	"condominium-manager/internal/repository"
)

type TenantHandler struct {
	repo repository.TenantRepository
}

func Register(r *gin.Engine, repo repository.TenantRepository) {
	h := &TenantHandler{repo: repo}
	r.GET("/v1/tenants", h.list)
	r.POST("/v1/tenants", h.create)
	r.GET("/v1/tenants/search", h.search)
	r.DELETE("/v1/tenants/:id", h.delete)
	r.PUT("/v1/tenants/:id", h.update)
}

func (h *TenantHandler) update(c *gin.Context) {
	id, err := strconv.ParseInt(c.Param("id"), 10, 32)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "id must be an integer"})
		return
	}

	var t model.Tenant
	if err := c.ShouldBindJSON(&t); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}
	// id comes from the path only
	t.ID = int32(id)

	if err := h.repo.Update(c.Request.Context(), &t); err != nil {
		switch {
		case errors.Is(err, repository.ErrNotFound):
			c.JSON(http.StatusNotFound, gin.H{"error": err.Error()})
		case errors.Is(err, repository.ErrAlreadyExists):
			c.JSON(http.StatusConflict, gin.H{"error": err.Error()})
		default:
			c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		}
		return
	}
	c.JSON(http.StatusOK, t)
}

func (h *TenantHandler) delete(c *gin.Context) {
	id, err := strconv.ParseInt(c.Param("id"), 10, 32)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "id must be an integer"})
		return
	}

	if err := h.repo.Delete(c.Request.Context(), int32(id)); err != nil {
		if errors.Is(err, repository.ErrNotFound) {
			c.JSON(http.StatusNotFound, gin.H{"error": err.Error()})
			return
		}
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}
	c.Status(http.StatusNoContent)
}

func (h *TenantHandler) list(c *gin.Context) {
	tenants, err := h.repo.List(c.Request.Context())
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}
	c.JSON(http.StatusOK, tenants)
}

func (h *TenantHandler) create(c *gin.Context) {
	var t model.Tenant
	if err := c.ShouldBindJSON(&t); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}
	// id is read-only: always assigned by the database
	t.ID = 0

	if err := h.repo.Create(c.Request.Context(), &t); err != nil {
		if errors.Is(err, repository.ErrAlreadyExists) {
			c.JSON(http.StatusConflict, gin.H{"error": err.Error()})
			return
		}
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}
	c.JSON(http.StatusCreated, t)
}

func (h *TenantHandler) search(c *gin.Context) {
	f := repository.TenantFilter{
		Name:     c.Query("name"),
		Lastname: c.Query("lastname"),
	}
	if v := c.Query("id"); v != "" {
		id, err := strconv.ParseInt(v, 10, 32)
		if err != nil {
			c.JSON(http.StatusBadRequest, gin.H{"error": "id must be an integer"})
			return
		}
		id32 := int32(id)
		f.ID = &id32
	}

	tenants, err := h.repo.Search(c.Request.Context(), f)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}
	if len(tenants) == 0 {
		c.JSON(http.StatusNotFound, gin.H{"error": "no tenant found"})
		return
	}
	c.JSON(http.StatusOK, tenants)
}
