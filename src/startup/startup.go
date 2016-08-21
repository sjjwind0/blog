package startup

import (
	"controller"
	"controller/personal"
	"framework/base/config"
	"framework/database"
	"framework/server"
	"model"
	"plugin"
)

func StartServer() {
	config := config.GetDefaultConfigJsonReader()
	localWebResourcePath := config.Get("resource.localpath").(string)
	pluginResourcePath := config.Get("resource.pluginpath").(string)
	port := int(config.Get("net.port").(int64))

	server.ShareServerMgrInstance().SetServerPort(port)

	// pubic api
	server.ShareServerMgrInstance().RegisterController(controller.NewIndexController())
	server.ShareServerMgrInstance().RegisterController(controller.NewBlogController())
	server.ShareServerMgrInstance().RegisterController(controller.NewArticleController())
	server.ShareServerMgrInstance().RegisterController(controller.NewAPIController())
	server.ShareServerMgrInstance().RegisterController(controller.NewLoginController())
	server.ShareServerMgrInstance().RegisterController(controller.NewAboutController())
	server.ShareServerMgrInstance().RegisterController(controller.NewPlayController())

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

	// plugin
	pluginsRunner := plugin.GetDefaultPluginManager().GetAllPluginRunner()
	for _, runner := range pluginsRunner {
		staticFiles := runner.ResourceHandler()
		for _, web := range staticFiles {
			server.ShareServerMgrInstance().RegisterStaticFile(web, pluginResourcePath)
		}
		normalControllers := runner.NormalHanlder()
		for _, controller := range normalControllers {
			server.ShareServerMgrInstance().RegisterController(controller.(server.NormalController))
		}
		websocketsControllers := runner.WebSocketHandler()
		for _, controller := range websocketsControllers {
			server.ShareServerMgrInstance().RegisterWebSocketController(
				controller.(server.WebSocketController))
		}
	}

	// 评论表
	database.ShareDatabaseRunner().RegisterModel(model.ShareCommentModel())
	// 博客表
	database.ShareDatabaseRunner().RegisterModel(model.ShareBlogModel())
	// 用户表
	database.ShareDatabaseRunner().RegisterModel(model.ShareUserModel())
	database.ShareDatabaseRunner().Start()

	server.ShareServerMgrInstance().StartServer()
}
