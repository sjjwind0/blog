package startup

import (
	"controller"
	"controller/personal"
	"framework/config"
	"framework/database"
	"framework/server"
	"model"
)

func StartServer() {
	config := config.GetDefaultConfigFileManager()
	localWebResourcePath := config.ReadConfig("resource.localpath").(string)
	port := int(config.ReadConfig("net.port").(int64))

	server.ShareServerMgrInstance().SetServerPort(port)

	// pubic api
	server.ShareServerMgrInstance().RegisterController(controller.NewIndexController())
	server.ShareServerMgrInstance().RegisterController(controller.NewBlogController())
	server.ShareServerMgrInstance().RegisterController(controller.NewArticleController())
	server.ShareServerMgrInstance().RegisterController(controller.NewAPIController())
	server.ShareServerMgrInstance().RegisterController(controller.NewLoginController())
	server.ShareServerMgrInstance().RegisterController(controller.NewAboutController())
	server.ShareServerMgrInstance().RegisterController(controller.NewNotImplController())

	// personal api
	server.ShareServerMgrInstance().RegisterController(personal.NewSyncController())
	server.ShareServerMgrInstance().RegisterController(personal.NewPersonalAuthController())
	server.ShareServerMgrInstance().RegisterController(personal.NewPersonalFetchController())
	server.ShareServerMgrInstance().RegisterController(personal.NewPersonalFileController())
	server.ShareServerMgrInstance().RegisterController(personal.NewPersonalDeleteController())

	// staitc file
	server.ShareServerMgrInstance().RegisterStaticFile("/js/", localWebResourcePath)
	server.ShareServerMgrInstance().RegisterStaticFile("/css/", localWebResourcePath)
	server.ShareServerMgrInstance().RegisterStaticFile("/img/", localWebResourcePath)
	server.ShareServerMgrInstance().RegisterStaticFile("/font/", localWebResourcePath)

	// 评论表
	database.ShareDatabaseRunner().RegisterModel(model.ShareCommentModel())
	// 博客表
	database.ShareDatabaseRunner().RegisterModel(model.ShareBlogModel())
	// 用户表
	database.ShareDatabaseRunner().RegisterModel(model.ShareUserModel())
	database.ShareDatabaseRunner().Start()

	server.ShareServerMgrInstance().StartServer()
}
