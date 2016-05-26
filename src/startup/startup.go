package startup

import (
	"controller"
	"framework/config"
	"framework/database"
	"framework/server"
	"model"
)

const defaultConfigPath = "default.conf"

func StartServer() {
	config := config.NewConfigFileManager(defaultConfigPath)
	localWebResourcePath := config.ReadConfig("resource.localpath").(string)
	port := int(config.ReadConfig("net.port").(int64))

	server.ShareServerMgrInstance().SetServerPort(port)
	server.ShareServerMgrInstance().RegisterController("/index", controller.NewIndexController())
	server.ShareServerMgrInstance().RegisterController("/blog", controller.NewBlogController())

	server.ShareServerMgrInstance().RegisterStaticFile("/js/", localWebResourcePath)
	server.ShareServerMgrInstance().RegisterStaticFile("/css/", localWebResourcePath)
	server.ShareServerMgrInstance().RegisterStaticFile("/img/", localWebResourcePath)
	server.ShareServerMgrInstance().RegisterStaticFile("/font/", localWebResourcePath)

	server.ShareServerMgrInstance().RegisterStaticFile("/blog/", localWebResourcePath+"/html")

	database.ShareDatabaseRunner().RegisterModel(model.ShareCommentModel())
	database.ShareDatabaseRunner().Start()

	server.ShareServerMgrInstance().StartServer()
}
