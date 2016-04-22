# Scaffold
_Simple context aware router_

Scaffold is a simple router with middleware support.
It allows for alternate dispatchers.

The default dispatcher can parse basic patterns in the form /segment/:param/segment. Parameters can be accessed in the Handler/Middleware using GetParam. 

Example:
```go
dispatcher := scaffold.DefaultDispatcher()
router := scaffold.New(dispatcher)

// More specific routes have precedence
router.Host("example.com").Get("", somehandler)
router.Get("", somehandler)

http.ListenAndServe(":8080", dispatcher)
```

