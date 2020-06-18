## **Logger Utils**

- This is an alternative to the log library where zap library is used. This library can be a standard logger library across the system. 
    All that is needed is to import the "github.com/MagalixTechnologies/core/logger" package and use the following method to initialize the logger:
     
    ```go
  import "github.com/MagalixTechnologies/core/logger"
  
  log := logger.New(logger.InfoLevel)    
    ```

- There are 4 levels that can be used; info, error, debug, and warning.

#### Examples in Code

```go
    mgx_middleware "github.com/MagalixTechnologies/core/middleware"
    var logLevel mgx_logger.Level
    r := chi.NewRouter()
    r.Use(middleware.Timeout(time.Duration(timeout) * time.Second))
    r.Route("/api/v1", func(r chi.Router) {
        r.Use(mgx_middleware.Log(logLevel))
        r.Mount("/", api.HandlerCustom(advisorServer))
    })
```