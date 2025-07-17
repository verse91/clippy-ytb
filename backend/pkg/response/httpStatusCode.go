package response

// HTTP Status Code for team's project

const (
	SuccessCode         = 200001 // success
	MissUserIDErrCode   = 200002 // missing user id
	ParamInvalidErrCode = 200003 // email is invalid
)

// message

var msg = map[int]string{
	SuccessCode:         "Success",
	MissUserIDErrCode:   "Missing user id",
	ParamInvalidErrCode: "Email is invalid",
}
