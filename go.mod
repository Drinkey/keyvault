module github.com/Drinkey/keyvault

go 1.14

require (
	github.com/gin-gonic/gin v1.6.3
	github.com/mattn/go-sqlite3 v1.14.5
	github.com/smartystreets/goconvey v1.6.4
	gorm.io/driver/sqlite v1.1.4
	gorm.io/gorm v1.20.8
)

replace (
	github.com/Drinkey/keyvault/certio => ./certio
	github.com/Drinkey/keyvault/pkg/crypt => ./pkg/crypt
	github.com/Drinkey/keyvault/pkg/utils => ./pkt/utils
	github.com/Drinkey/keyvault/models => ./models
	github.com/Drinkey/keyvault/routers => ./routers
	github.com/Drinkey/keyvault/pkg/e => ./pkg/e
)
