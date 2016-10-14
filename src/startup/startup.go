package startup

import (
	"controller"
	"controller/personal"
	"fmt"
	"framework/base/config"
	"framework/database"
	"framework/server"
	"model"
	"path/filepath"
	// "plugin"
)

func StartServer() {
	config := config.GetDefaultConfigJsonReader()
	localWebResourcePath := config.GetString("storage.file.res")
	fmt.Println("localWebResourcePath: ", localWebResourcePath)
	//pluginResourcePath := config.GetString("resource.pluginpath")
	port := config.GetInteger("net.port")

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
	server.ShareServerMgrInstance().RegisterStaticFile("js", filepath.Join(localWebResourcePath, "js"))
	server.ShareServerMgrInstance().RegisterStaticFile("css", filepath.Join(localWebResourcePath, "css"))
	server.ShareServerMgrInstance().RegisterStaticFile("img", filepath.Join(localWebResourcePath, "img"))
	server.ShareServerMgrInstance().RegisterStaticFile("font", filepath.Join(localWebResourcePath, "font"))

	/*
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
	*/
	// 评论表
	database.ShareDatabaseRunner().RegisterModel(model.ShareCommentModel())
	// 博客表
	database.ShareDatabaseRunner().RegisterModel(model.ShareBlogModel())
	// 用户表
	database.ShareDatabaseRunner().RegisterModel(model.ShareUserModel())
	// 插件表
	database.ShareDatabaseRunner().RegisterModel(model.SharePluginModel())

	database.ShareDatabaseRunner().Start()

	// // plugin
	// plugin.SharePluginMgrInstance().LoadPlugin(1)

	server.ShareServerMgrInstance().StartServer()
}
