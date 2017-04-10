package http

type MessageCode int

const (
	SUCCESS          MessageCode = 0
	NOTAUTHORIZED                = 400
	NOAUTHENTICATION             = 401
	NOSERVER                     = 402
	MALFORMEDJSON                = 403
	NOFILE                       = 404
	NOSERVERID                   = 405
	INVALIDTIME                  = 406
	UNKNOWN                      = 999
)
