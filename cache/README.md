# Cache
A middleware to cache requests based on  specified URLs & HTTP methods.

## Usage
1. Initiate the redis adapter if you want to use redis
    ```go
       rdb := cache.NewRedisAdapter(&cache.RedisOptions{
                Addr:     *RequestsCacheDSN,
                Password: *RequestsCachePassword,
                })
    ```
2. Initiate the cache client
    ```go
    		
   cacheClient, err := cache.NewClient(
            cache.ClientWithAdapter(rdb),
            cache.ClientWithTTL(time.Duration(*RequestsCacheExpireAfter)*time.Minute), //time to remove the data after
            cache.ClientWithRefreshKey("opn"),//key to clear the data for specifc resource
            cache.ClientWithMethods([]string{http.MethodGet, http.MethodPost}),//HTTP methods to cover
        )
    ```
3. use the middleware method
   ```go
   //router could by any method that follows this interface of  func(http.Handler) http.Handler
   router.Use(cacheClient.Middleware)
   ```
####Notes
To cache specific endpoints you could use the following.
```go
 allowedRoutes := []*cache.URLMatch{
   {
   Path:   "/api/v1/recommendations/stats",
   Method: http.MethodPost,
   },
   {
   Path:   "/api/v1/recommendations/:arg",
   Method: http.MethodPost,
   },
}
router.Use(cache.ApplyMiddlewareOnSpecificRoutes(CacheMiddleware, allowedRoutes...))
```
