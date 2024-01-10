package exception

var (
	EmptyTokenMsg   string = "EMPTY_TOKEN"   // Frontend will redirect to login page
	InvalidTokenMsg string = "INVALID_TOKEN" // Frontend will redirect to login page
	ExpiredTokenMsg string = "EXPIRED_TOKEN" // Frontend will request a new access token to /api/auth/refresh
)
