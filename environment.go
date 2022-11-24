package zlogger

import "github.com/gin-gonic/gin"


const (
  Debug LogEnvironment = gin.DebugMode
  Release LogEnvironment = gin.ReleaseMode
)

type LogEnvironment string