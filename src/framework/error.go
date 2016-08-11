package framework

const (
	// global
	ErrorOK                    = 0
	ErrorMethodError           = 1
	ErrorParamError            = 2
	ErrorSQLError              = 3
	ErrorRenderError           = 4
	ErrorNoSuchFileOrDirectory = 5

	// blog
	ErrorBlogExist = 1000
	ErrorEmptyBlog = 1001

	// runtime
	ErrorRunTimeError = 2000

	// account
	ErrorAccountNotLogin  = 3000
	ErrorAccountAuthError = 3005

	// file
	ErrorFileNotExist = 4000
)
