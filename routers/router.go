/*
Package routers implements URL route of keyvault service
*/

package routers

import (
	v1 "github.com/Drinkey/keyvault/routers/api/v1"
	"github.com/gin-gonic/gin"
)

// InitRouter creates a gin handler and setup API routes
func InitRouter() *gin.Engine {
	r := gin.New()
	r.Use(gin.Logger())
	r.Use(gin.Recovery())
	apiV1 := r.Group("/api/v1")
	{
		apiV1.GET("/ping", v1.Ping)
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
