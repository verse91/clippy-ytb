package response

// HTTP Status Code for team's project

// Custom business codes (2xx/4xx/5xx/other)
const (
	SuccessCode         = 200001 // success
	MissUserIDErrCode   = 200002 // missing user id
	ParamInvalidErrCode = 200003 // email is invalid
)

// Client error codes (400xxx)
const (
	ErrInvalidRequestBody = 400001 // invalid request body
	ErrURLRequired        = 400002 // url is required
	ErrDownloadIDRequired = 400003 // download id is required
)

// Server error codes (500xxx)
const (
	ErrDownloadStartFailed = 500001 // failed to start download
	ErrSerializeResponse   = 500002 // failed to serialize response
	ErrSerializeStatus     = 500003 // failed to serialize status
)

// Not found error codes (404xxx)
const (
	ErrDownloadNotFound = 404001 // download not found
)

// Unauthorized error codes (401xxx)
const (
	ErrUnauthorized = 401001 // unauthorized access
)

const (
    ErrTooManyRequests = 429001 // too many requests
)

// message

var msg = map[int]string{
	SuccessCode:            "Success",
	MissUserIDErrCode:      "Missing user id",
	ParamInvalidErrCode:    "Email is invalid",
	ErrInvalidRequestBody:  "Invalid request body",
	ErrURLRequired:         "URL is required",
	ErrDownloadIDRequired:  "Download ID is required",
	ErrDownloadStartFailed: "Failed to start download",
	ErrSerializeResponse:   "Failed to serialize response",
	ErrSerializeStatus:     "Failed to serialize status",
	ErrDownloadNotFound:    "Download not found",
	ErrUnauthorized:        "Unauthorized access",
    ErrTooManyRequests:    "Too many requests",
}
