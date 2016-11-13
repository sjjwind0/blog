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

	// plugin
	ErrorPluginNotExist   = 2000
	ErrorPluginParseError = 2001
	ErrorPluginBuildError = 2002
	ErrorPluginLoadError  = 2003
	ErrorPluginPathError  = 2004

	// runtime
	ErrorRunTimeError = 3000

	// account
	ErrorAccountNotLogin  = 4000
	ErrorAccountAuthError = 5005

	// file
	ErrorFileNotExist = 5000
)
