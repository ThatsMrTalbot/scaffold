# Scaffold
_Simple context aware router_

[![GoDoc](https://godoc.org/github.com/ThatsMrTalbot/scaffold?status.svg)](https://godoc.org/github.com/ThatsMrTalbot/scaffold) [![Coverage Status](https://coveralls.io/repos/github/ThatsMrTalbot/scaffold/badge.svg?branch=master)](https://coveralls.io/github/ThatsMrTalbot/scaffold?branch=master) [![Build Status](https://travis-ci.org/ThatsMrTalbot/scaffold.svg?branch=master)](https://travis-ci.org/ThatsMrTalbot/scaffold)

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

