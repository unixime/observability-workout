package handler

import (
	"github.com/gin-gonic/gin"
	"github.com/swaggest/swgui/v5emb"
)

const (
	specURL  = "/v1/openapi.yaml"
	docsPath = "/docs/"
)

// RegisterSpec serves the OpenAPI definition read from specPath at request
// time, and a Swagger UI page at /docs that renders it.
func RegisterSpec(r *gin.Engine, specPath string) {
	r.GET(specURL, func(c *gin.Context) {
		c.Header("Content-Type", "application/yaml")
		c.File(specPath)
	})

	ui := gin.WrapH(v5emb.New("Condominium Manager API", specURL, docsPath))
	r.GET("/docs", func(c *gin.Context) { c.Redirect(301, docsPath) })
	r.GET(docsPath+"*any", ui)
}
