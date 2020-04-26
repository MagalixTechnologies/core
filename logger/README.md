## **Logger Utils**

- This is an alternative to the log library where zap library is used. This library can be a standard logger library across the system. 
    All that is needed is to import the "github.com/MagalixTechnologies/core/logger" package and use the following method to initialize the logger:
     
    ```go
  import "github.com/MagalixTechnologies/core/logger"
  
  log := logger.New(logger.InfoLevel)    
    ```

- There are 4 levels that can be used; info, error, debug, and warning.