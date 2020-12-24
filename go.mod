module github.com/Drinkey/keyvault

go 1.14

require (
	github.com/alecthomas/template v0.0.0-20190718012654-fb15b899a751
	github.com/cpuguy83/go-md2man/v2 v2.0.0 // indirect
	github.com/gin-gonic/gin v1.6.3
	github.com/go-openapi/spec v0.20.0 // indirect
	github.com/go-playground/validator/v10 v10.4.1 // indirect
	github.com/golang/protobuf v1.4.3 // indirect
	github.com/json-iterator/go v1.1.10 // indirect
	github.com/leodido/go-urn v1.2.1 // indirect
	github.com/mattn/go-sqlite3 v1.14.5
	github.com/russross/blackfriday/v2 v2.1.0 // indirect
	github.com/smartystreets/goconvey v1.6.4
	github.com/swaggo/files v0.0.0-20190704085106-630677cd5c14
	github.com/swaggo/gin-swagger v1.3.0
	github.com/swaggo/swag v1.7.0
	github.com/ugorji/go v1.2.2 // indirect
	golang.org/x/crypto v0.0.0-20201221181555-eec23a3978ad // indirect
	golang.org/x/net v0.0.0-20201216054612-986b41b23924 // indirect
	golang.org/x/sync v0.0.0-20201020160332-67f06af15bc9
	golang.org/x/sys v0.0.0-20201221093633-bc327ba9c2f0 // indirect
	golang.org/x/tools v0.0.0-20201221201019-196535612888 // indirect
	google.golang.org/protobuf v1.25.0 // indirect
	gorm.io/driver/sqlite v1.1.4
	gorm.io/gorm v1.20.8
)

replace (
	github.com/Drinkey/keyvault/certio => ./certio
	github.com/Drinkey/keyvault/models => ./models
	github.com/Drinkey/keyvault/pkg/crypt => ./pkg/crypt
	github.com/Drinkey/keyvault/pkg/e => ./pkg/e
	github.com/Drinkey/keyvault/pkg/utils => ./pkt/utils
	github.com/Drinkey/keyvault/routers => ./routers
	github.com/Drinkey/keyvault/pkg/app => ./pkg/app
	github.com/Drinkey/keyvault/pkg/server => ./pkg/server
	github.com/Drinkey/keyvault/pkg/settings => ./pkg/settings
	github.com/Drinkey/keyvault/services/certificate_service => ./services/certificate_service
	github.com/Drinkey/keyvault/services/namespace_service => ./services/namespace_service
	github.com/Drinkey/keyvault/services/secret_service => ./services/secret_service
)
