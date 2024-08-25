package server

import (
	"database/sql"

	"github.com/gin-gonic/gin"
	"github.com/spf13/viper"
)

type HttpServer struct {
	config            *viper.Viper
	router            *gin.Engine
	runnersController *controllers.RunnersController
	resultController  *controllers.ResultController
}

func InitHttpServer(config *viper.Viper, dbHandler *sql.DB) HttpServer {

}

func (hs HttpServer) Start() {

}
