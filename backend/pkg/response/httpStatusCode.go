package response

// HTTP Status Code for team's project

const (
	SuccessCode         = 200001 // success
	MissUserIDErrCode   = 200002 // missing user id
	ParamInvalidErrCode = 200003 // email is invalid
)

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

// Added new error code for unauthorized access (401xxx)
const (
	ErrUnauthorized = 401001 // unauthorized access
)

// message

var msg = map[int]string{
	SuccessCode:         "Success",
	MissUserIDErrCode:   "Missing user id",
	ParamInvalidErrCode: "Email is invalid",
	ErrInvalidRequestBody: "Invalid request body",
	ErrURLRequired:        "URL is required",
	ErrDownloadIDRequired: "Download ID is required",
	ErrDownloadStartFailed: "Failed to start download",
	ErrSerializeResponse:   "Failed to serialize response",
	ErrSerializeStatus:     "Failed to serialize status",
	ErrDownloadNotFound:    "Download not found",
	ErrUnauthorized:        "Unauthorized access",
}
