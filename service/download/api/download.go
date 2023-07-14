package api

import (
	"github.com/gin-gonic/gin"
	dbcli "github.com/kuan525/netdisk/dbclient"
)

func DownloadURLHandler(c *gin.Context) {
	filehash := c.Request.FormValue("filehash")

	dbResp, err := dbcli.
}
