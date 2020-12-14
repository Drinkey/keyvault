package routers

import (
	v1 "github.com/Drinkey/keyvault/routers/api/v1"
	"github.com/gin-gonic/gin"
)

func InitRouter() *gin.Engine {
	r := gin.New()
	r.Use(gin.Logger())
	r.Use(gin.Recovery())
	apiv1 := r.Group("/api/v1")
	{
		apiv1.GET("/ping", v1.Ping)
		// only list namespace of vault
		apiv1.GET("/namespace", v1.ListNamespaces)
		apiv1.POST("/namespace", v1.CreateNamespace)
		// TODO: only name can be updated
		// apiv1.PUT("/namespace", v1.UpdateNamespace)
		// TODO: need authorization, it would be cost operation, rekey ns would
		// require all secrets with this ns updated.
		// apiv1.POST("/namespace/rekey", v1.RekeyNamespace)
		// respond to URL like /vault/gitlab/?q=k8s_password
		apiv1.GET("/vault/:namespace", v1.GetSecret)
		apiv1.POST("/vault/:namespace", v1.CreateSecret)
		// TODO: respond to URL like /vault/gitlab/?q=k8s_password
		// apiv1.DELETE("/vault/:namespace", v1.DeleteSecret)
		// TODO: only value can be updated, namespace and key can't
		// v1.PUT("/vault/:namespace", v1.UpdateSecret)

		// Certificate need extra works, auth or token
		// TODO: we may need this API exposed in HTTP, not HTTPS
		apiv1.POST("/cert/req", v1.CreateCertificateRequest)
		// TODO: read CA cert directly from a file
		// apiv1.GET("/cert/ca", v1.GetCACertificate)
		// TODO: OU(name in certificate table) should be unique
		// respond to URL like /cert/?q=k8s_password
		apiv1.GET("/cert/", v1.GetCertificate)
		// TODO: only limited user should be able to access this API, how
		// apiv1.POST("/cert/issue", v1.CreateCertificateRequest)

	}
	return r
}
