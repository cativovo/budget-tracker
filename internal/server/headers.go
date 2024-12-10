package server

// https://github.com/labstack/echo/blob/3b017855b4d331002e2b8b28e903679b875ae3e9/echo.go#L201

const (
	headerAccept         = "Accept"
	headerAcceptEncoding = "Accept-Encoding"
	// HeaderAllow is the name of the "Allow" header field used to list the set of methods
	// advertised as supported by the target resource. Returning an Allow header is mandatory
	// for status 405 (method not found) and useful for the OPTIONS method in responses.
	// See RFC 7231: https://datatracker.ietf.org/doc/html/rfc7231#section-7.4.1
	headerAllow               = "Allow"
	headerAuthorization       = "Authorization"
	headerContentDisposition  = "Content-Disposition"
	headerContentEncoding     = "Content-Encoding"
	headerContentLength       = "Content-Length"
	headerContentType         = "Content-Type"
	headerCookie              = "Cookie"
	headerSetCookie           = "Set-Cookie"
	headerIfModifiedSince     = "If-Modified-Since"
	headerLastModified        = "Last-Modified"
	headerLocation            = "Location"
	headerRetryAfter          = "Retry-After"
	headerUpgrade             = "Upgrade"
	headerVary                = "Vary"
	headerWWWAuthenticate     = "WWW-Authenticate"
	headerXForwardedFor       = "X-Forwarded-For"
	headerXForwardedProto     = "X-Forwarded-Proto"
	headerXForwardedProtocol  = "X-Forwarded-Protocol"
	headerXForwardedSsl       = "X-Forwarded-Ssl"
	headerXUrlScheme          = "X-Url-Scheme"
	headerXHTTPMethodOverride = "X-HTTP-Method-Override"
	headerXRealIP             = "X-Real-Ip"
	headerXRequestID          = "X-Request-Id"
	headerXCorrelationID      = "X-Correlation-Id"
	headerXRequestedWith      = "X-Requested-With"
	headerServer              = "Server"
	headerOrigin              = "Origin"
	headerCacheControl        = "Cache-Control"
	headerConnection          = "Connection"

	// Access control
	headerAccessControlRequestMethod    = "Access-Control-Request-Method"
	headerAccessControlRequestHeaders   = "Access-Control-Request-Headers"
	headerAccessControlAllowOrigin      = "Access-Control-Allow-Origin"
	headerAccessControlAllowMethods     = "Access-Control-Allow-Methods"
	headerAccessControlAllowHeaders     = "Access-Control-Allow-Headers"
	headerAccessControlAllowCredentials = "Access-Control-Allow-Credentials"
	headerAccessControlExposeHeaders    = "Access-Control-Expose-Headers"
	headerAccessControlMaxAge           = "Access-Control-Max-Age"

	// Security
	headerStrictTransportSecurity         = "Strict-Transport-Security"
	headerXContentTypeOptions             = "X-Content-Type-Options"
	headerXXSSProtection                  = "X-XSS-Protection"
	headerXFrameOptions                   = "X-Frame-Options"
	headerContentSecurityPolicy           = "Content-Security-Policy"
	headerContentSecurityPolicyReportOnly = "Content-Security-Policy-Report-Only"
	headerXCSRFToken                      = "X-CSRF-Token"
	headerReferrerPolicy                  = "Referrer-Policy"
)
