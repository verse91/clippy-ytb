package errors

// Client error codes (400xxx)
const (
	ErrInvalidRequestBody = 400001
	ErrURLRequired        = 400002
	ErrDownloadIDRequired = 400003
)

// Server error codes (500xxx)
const (
	ErrDownloadStartFailed = 500001
	ErrSerializeResponse   = 500002
	ErrSerializeStatus     = 500003
)

// Not found error codes (404xxx)
const (
	ErrDownloadNotFound = 404001
)
