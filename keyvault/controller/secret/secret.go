package secret

import (
	"fmt"
	"log"
	"net/http"

	"github.com/Drinkey/keyvault/internal"
	"github.com/Drinkey/keyvault/model"
	"github.com/gin-gonic/gin"
)

func Query(c *gin.Context) {
	namespace := c.Param("namespace")
	key := c.Query("q")

	var secret_model model.Secrets
	secret := secret_model.Get(key, namespace)

	if secret.IsEmpty() {
		c.JSON(http.StatusNotFound, gin.H{
			"error": fmt.Sprintf("Record Not Found: NameSpace=%s, Key=%s", namespace, key),
		})
		return
	}

	secret.Value = internal.Decrypt(secret.Value, secret.NameSpace.MasterKey)
	secret.NameSpace.MasterKey = internal.KeyMask

	c.JSON(http.StatusOK, secret)
}

func Create(c *gin.Context) {

	namespace := c.Param("namespace")
	log.Printf("Creating new secret under %s", namespace)

	var secret_data model.Secrets
	if err := c.ShouldBindJSON(&secret_data); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}
	// 1. Query namespace database for id, masterkey
	var ns_model model.Namespace
	ns, err := ns_model.Get(namespace)

	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}
	if ns.IsEmpty() {
		c.JSON(http.StatusBadRequest, gin.H{
			"error": fmt.Sprintf("namespace %s does not exist, create it first", namespace),
		})
		return
	}
	secret_data.NameSpace = ns
	// 2. encrypt secret value with master key
	secret_data.Value = internal.Encrypt(secret_data.Value, secret_data.NameSpace.MasterKey)
	fmt.Println(secret_data)
	// 3. save to database
	var secret_model model.Secrets
	err = secret_model.Create(secret_data)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}
	c.JSON(http.StatusCreated, gin.H{
		"message": fmt.Sprintf("Secret namespace=%s, key=%s created success", namespace, secret_data.Key),
	})
}

func Delete(c *gin.Context) {
	c.JSON(http.StatusCreated, gin.H{
		"message": fmt.Sprintf("delete"),
	})
}

func Update(c *gin.Context) {
	c.JSON(http.StatusCreated, gin.H{
		"message": fmt.Sprintf("update"),
	})
}
