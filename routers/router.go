/*
Package routers implements URL route of keyvault service
*/

package routers

import (
	_ "github.com/Drinkey/keyvault/docs"
	v1 "github.com/Drinkey/keyvault/routers/api/v1"
	"github.com/gin-gonic/gin"
	swaggerFiles "github.com/swaggo/files"
	ginSwagger "github.com/swaggo/gin-swagger"
)

/*
Classify Routers by endpoint required security level
- API:
	HTTP
- Secret:
	HTTPS + Verify client certificate
- Default(Maintenance):
	HTTPS + trusted location + user/pass, no client certificate verification
*/

const currentAPIVersion = "/api/v1"

// InitAPIRouter creates a gin handler and serve swagger API documentation
func InitAPIRouter() *gin.Engine {
	r := gin.New()
	r.GET("/swagger/*any", ginSwagger.WrapHandler(swaggerFiles.Handler))
	r.GET("/ping", v1.Ping)
	return r
}

// InitSecretRouter creates a gin handler and setup API routes for secret access
func InitSecretRouter() *gin.Engine {
	r := gin.New()
	r.Use(gin.Logger())
	r.Use(gin.Recovery())
	apiV1 := r.Group(currentAPIVersion)
	{
		apiV1.GET("/pings", v1.Ping)

		// only list namespace of vault
		apiV1.GET("/namespace", v1.ListNamespaces)
		apiV1.POST("/namespace", v1.CreateNamespace)
		// TODO: only name can be updated
		// apiV1.PUT("/namespace", v1.UpdateNamespace)
		// TODO: need authorization, it would be cost operation, rekey ns would
		// require all secrets with this ns updated.
		// apiV1.POST("/namespace/rekey", v1.RekeyNamespace)
		// respond to URL like /vault/gitlab/?q=k8s_password
		apiV1.GET("/vault/:namespace", v1.GetSecret)
		apiV1.POST("/vault/:namespace", v1.CreateSecret)
		// TODO: respond to URL like /vault/gitlab/?q=k8s_password
		// apiV1.DELETE("/vault/:namespace", v1.DeleteSecret)
		// TODO: only value can be updated, namespace and key can't
		// v1.PUT("/vault/:namespace", v1.UpdateSecret)
	}
	return r
}

// InitDefaultRouter creates a gin handler and setup API routes for maintenance functions
func InitDefaultRouter() *gin.Engine {
	r := gin.New()
	r.Use(gin.Logger())
	r.Use(gin.Recovery())
	apiV1 := r.Group(currentAPIVersion)
	{
		apiV1.GET("/ping", v1.Ping)
		// Certificate need extra works, auth or token
		// we may need this API exposed in HTTP, not HTTPS
		apiV1.POST("/cert/req", v1.CreateCertificateRequest)

		// TODO: read CA cert directly from a file
		apiV1.GET("/cert/ca", v1.GetCACertificate)

		// OU(name in certificate table) should be unique
		// respond to URL like /cert/?q=k8s_password
		apiV1.GET("/cert/", v1.GetCertificate)

		// TODO: only limited user should be able to access this API, how
		// apiV1.POST("/cert/issue", v1.CreateCertificateRequest)
	}
	return r
}
