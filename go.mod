module github.com/Drinkey/keyvault

go 1.14

require (
	github.com/gin-gonic/gin v1.6.3
	github.com/mattn/go-sqlite3 v1.14.5
	github.com/smartystreets/goconvey v1.6.4
)

replace (
	github.com/Drinkey/keyvault/internal  => ./internal
	github.com/Drinkey/keyvault/certio  => ./certio
	github.com/Drinkey/keyvault/controller  => ./controller
	github.com/Drinkey/keyvault/namespace  => ./namespace
	github.com/Drinkey/keyvault/models  => ./models
	github.com/Drinkey/keyvault/secret  => ./secret
)