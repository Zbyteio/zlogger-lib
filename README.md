# zlogger-lib
- logger instances created are of 2 types - [DEBUG]/[RELEASE]
- instances created depends on gin mode of the code.

## Default log levels
- Error
- Warn
- Info
- Debug

## Usage

Fetch the latest version of zlogger-lib
```
GOPROXY=https://proxy.golang.org GO111MODULE=on go get github.com/Zbyteio/zlogger-lib@v1.1.0
```
### Create an app logger
- use this logger to log throughout the library


```
appLogger := zlogger.NewZlogger(gin.DebugMode)

............

ZBlocksAppDebugLogger.Debugf("%s", "success print debug via applogger[DEBUG]")
```



### Create a gin logger
- use this logger as middleware for gin route logging

```  
ginLogger := zlogger.NewGinLogger(gin.DebugMode)

............
ginEng := gin.New()


ginEng.Use(
    ZBlocksGinDebugLogger.GinRequestLoggerMiddleware(),
)
ginEng.GET("/abc", func(c *gin.Context) {
    c.String(http.StatusOK, "Welcome Gin Server")
})
```


### Create a gorm logger
- use this for replacing gorm trace logging
- use this for logging at database level inside code

```  
gormLogger := zlogger.NewGormLogger(gin.DebugMode)

............

ZBlocksGormDebugLogger.SetAsDefault()
db, err := gorm.Open(postgres.Open(dsn), &gorm.Config{
  Logger: ZBlocksGormDebugLogger,
})
```

## TODO 
- [X] remove setting of gin mode, maintain local state instead
- [ ] Add custom log levels if required later
- [ ] Add basic configurability
