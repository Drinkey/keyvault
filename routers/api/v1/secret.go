package v1

import (
	"fmt"
	"log"
	"net/http"

	"github.com/Drinkey/keyvault/models"
	"github.com/Drinkey/keyvault/pkg/crypt"
	"github.com/Drinkey/keyvault/pkg/e"
	"github.com/gin-gonic/gin"
)

type Secret struct {
	ID          uint      `json:"id"`
	Key         string    `json:"key"`
	Value       string    `json:"value"`
	NamespaceID uint      `json:"namespace_id"`
	Namespace   Namespace `json:"namespace"`
}

func GetSecret(c *gin.Context) {

	namespace := c.Param("namespace")

	if !IsClientAuthorized(c.Request, namespace) {
		msg := fmt.Sprintf(`Client not authorized to access the namespace=%s`, namespace)
		c.JSON(http.StatusUnauthorized, gin.H{
			"code": e.NOT_AUTHORIED,
			"msg":  e.GetMsg(e.NOT_AUTHORIED),
			"data": msg,
		})
		return
	}

	key := c.Query("q")
	log.Printf("Query secret [%s] under namespace %s", key, namespace)

	s, err := models.GetSecret(crypt.Sha256Sum(key), namespace)

	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{
			"code": e.ERROR,
			"msg":  e.GetMsg(e.ERROR),
			"data": fmt.Sprintf("error when finding secret: NameSpace=%s, Key=%s", namespace, key),
		})
		return
	}

	if s.IsEmpty() {
		c.JSON(http.StatusNotFound, gin.H{
			"code": e.NOT_FOUND,
			"msg":  e.GetMsg(e.NOT_FOUND),
			"data": fmt.Sprintf("Record Not Found: NameSpace=%s, Key=%s", namespace, key),
		})
		return
	}

	cipherTextBytes, err := crypt.DecodeString(s.Value)
	if err != nil {
		log.Printf("failed to decode string %s", s.Value)
		c.JSON(http.StatusNotFound, gin.H{
			"code": e.DECODING_ERROR,
			"msg":  e.GetMsg(e.DECODING_ERROR),
			"data": fmt.Sprintf("failed to decode secret value string. NameSpace=%s, Key=%s", namespace, key),
		})
		return
	}
	nonceByte, err := crypt.DecodeString(s.Namespace.Nonce)
	masterKeyByte, err := crypt.DecodeString(s.Namespace.MasterKey)

	s.Value = crypt.Decrypt(cipherTextBytes, masterKeyByte, nonceByte)
	s.Namespace.MasterKey = crypt.KeyMask
	s.Namespace.Nonce = crypt.KeyMask

	c.JSON(http.StatusOK, gin.H{
		"code": e.SUCCESS,
		"msg":  e.GetMsg(e.SUCCESS),
		"data": s,
	})
}

// CreateSecret create a secret record in database, sensitive fields of info
// are encrypted or hashed.
func CreateSecret(c *gin.Context) {

	namespace := c.Param("namespace")

	if !IsClientAuthorized(c.Request, namespace) {
		msg := fmt.Sprintf("Client not authorized to create new key under the namespace=%s", namespace)
		c.JSON(http.StatusUnauthorized, gin.H{
			"code": e.NOT_AUTHORIED,
			"msg":  e.GetMsg(e.NOT_AUTHORIED),
			"data": msg,
		})
		return
	}

	log.Printf("Creating new secret under %s", namespace)

	var req Secret
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{
			"code": e.INVALID_PARAMS,
			"msg":  e.GetMsg(e.INVALID_PARAMS),
			"data": err.Error(),
		})
		return
	}
	// 1. Query namespace database for id, masterkey
	ns, err := models.GetNamespace(namespace)

	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{
			"code": e.INVALID_PARAMS,
			"msg":  e.GetMsg(e.INVALID_PARAMS),
			"data": fmt.Sprintf("Error when finding namespace %s: %s", namespace, err.Error()),
		})
		return
	}
	if ns.IsEmpty() {
		c.JSON(http.StatusNotFound, gin.H{
			"code": e.NOT_FOUND,
			"msg":  e.GetMsg(e.NOT_FOUND),
			"data": fmt.Sprintf("namespace %s does not exist, create it first", namespace),
		})
		return
	}

	// 2. encrypt secret value with master key
	nonceByte, err := crypt.DecodeString(ns.Nonce)
	masterKeyByte, err := crypt.DecodeString(ns.MasterKey)

	encryptDataBase64 := crypt.EncodeByte(
		crypt.Encrypt(req.Value, masterKeyByte, nonceByte),
	)
	hashKey := crypt.Sha256Sum(req.Key)

	// 3. save to database
	err = models.CreateSecret(hashKey, encryptDataBase64, namespace)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{
			"code": e.ERROR,
			"msg":  e.GetMsg(e.ERROR),
			"data": fmt.Sprintf("Error when creating secret: %s", err.Error()),
		})
		return
	}
	c.JSON(http.StatusCreated, gin.H{
		"code": e.SUCCESS,
		"msg":  e.GetMsg(e.SUCCESS),
		"data": fmt.Sprintf("Secret namespace=%s, key=%s created success", namespace, req.Key),
	})
}
