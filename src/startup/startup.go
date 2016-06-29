package startup

import (
	"controller"
	"controller/personal"
	"framework/config"
	"framework/database"
	"framework/server"
	"model"
)

const defaultConfigPath = "default.conf"

func StartServer() {
	config := config.GetConfigFileManager(defaultConfigPath)
	localWebResourcePath := config.ReadConfig("resource.localpath").(string)
	port := int(config.ReadConfig("net.port").(int64))

	server.ShareServerMgrInstance().SetServerPort(port)
	server.ShareServerMgrInstance().RegisterController(controller.NewIndexController())
	server.ShareServerMgrInstance().RegisterController(controller.NewBlogController())
	server.ShareServerMgrInstance().RegisterController(controller.NewAPIController())
	server.ShareServerMgrInstance().RegisterController(controller.NewLoginController())
	server.ShareServerMgrInstance().RegisterController(personal.NewSyncController())

	server.ShareServerMgrInstance().RegisterStaticFile("/js/", localWebResourcePath)
	server.ShareServerMgrInstance().RegisterStaticFile("/css/", localWebResourcePath)
	server.ShareServerMgrInstance().RegisterStaticFile("/img/", localWebResourcePath)
	server.ShareServerMgrInstance().RegisterStaticFile("/font/", localWebResourcePath)

	database.ShareDatabaseRunner().RegisterModel(model.ShareCommentModel())
	database.ShareDatabaseRunner().RegisterModel(model.ShareBlogModel())
	database.ShareDatabaseRunner().Start()

	server.ShareServerMgrInstance().StartServer()
}
