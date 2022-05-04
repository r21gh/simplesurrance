package services

type key int

const (
	requestNumber         key = 0
	XRequestId                = "X-Request-Id"
	Empty                     = ""
	apiPath                   = "/api/v1/"
	CounterPath               = apiPath + "counter/"
	storeFile                 = "value.txt"
	ContentType               = "Content-Type"
	ContentTypeValue          = "text/plain; charset=utf-8"
	ServerName                = "simplesurrance"
	ServerPortWithColon       = ":9000"
	ServerUsage               = "server listen address"
	ServerWelcomeMessage      = "Server started ..."
	ServerShutdownMessage     = "Server is shutting down..."
	XContentTypeOptions       = "X-Content-Type-Options"
	NoSniff                   = "nosniff"
	window                    = 60.0 // 60 seconds in float64
)
