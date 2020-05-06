## **Logger Middleware**

- The logger in the middleware package has two main methods that can be used by the user.
    - The first and main one is the Log() method which is a wrapper for handling requests that wraps a request with a structured log that contains the request id and details of the request. 
    The log is implemented using the **`zap`** package. 
    The Log() function takes only one parameter which is the Level of the log. 
    The levels are constants also inside the logger file. The use can choose from 4 levels; info, warning, error, or debug.
    Here is an example of how to use the Log() function:
     
        ```go
      var handler http.Handler = mux 
      { 
          handler = mw.Log(mw.InfoLevel)(handler) 
      }
        ```
    - The second method that can be used after applying the Log() method to the request is the GetLoggerFromContext() method.
    In the Log() method the **`zap`** logger is added to the context.
    Therefore, you can get the logger from context to apply anywhere else other than the default logger for the request. 
    After implementing this method, you can use any of the methods of type zap.SugaredLogger and the same request id for the request will be printed with your log.
    Here is an example of how to apply this:
        ```go
      logger, ok := mw.GetLoggerFromContext(context)
      if ok {
        logger.Infow("Add", "a", p.A, "b", p.B)  
      }
        ```

     - If I apply both these methods I will get an output like this:
         ```json
          {"level":"info","timestamp":"2020-04-16T01:40:13.092+0200","caller":"calc/calc.go:24","msg":"Add","requestId":"oabGdo0U","a":123,"b":24}
          {"level":"info","timestamp":"2020-04-16T01:40:13.092+0200","caller":"mw/logger.go:71","msg":"Default Log","requestId":"oabGdo0U","method":"GET","url":"/add/123/24","status":200,"from":"127.0.0.1","duration":"112.34µs"}
         ```

