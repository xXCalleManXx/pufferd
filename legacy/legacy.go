package legacy

import "github.com/gin-gonic/gin"
import "github.com/PufferPanel/pufferd/environments/system"

func GetServerInfo(c *gin.Context) {
}

func CreateServer(c *gin.Context) {
}

func UpdateServerInfo(c *gin.Context) {
}

func DeleteServer(c *gin.Context) {
}

func ServerPower(c *gin.Context) {
    system.StartServer(c)
}

func ServerConsole(c *gin.Context) {
}

func GetServerLog(c *gin.Context) {
}

func GetFile(c *gin.Context) {
}

func UpdateFile(c *gin.Context) {
}

func DeleteFile(c *gin.Context) {
}

func DownloadFile(c *gin.Context) {
}

func GetDirectory(c *gin.Context) {
}

func ReinstallServer(c *gin.Context) {
}

func ResetPassword(c *gin.Context) {
}
