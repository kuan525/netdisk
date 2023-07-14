package api

import (
	"github.com/gin-gonic/gin"
	"github.com/kuan525/netdisk/common"
	dbcli "github.com/kuan525/netdisk/dbclient"
	"net/http"
)

func DownloadURLHandler(c *gin.Context) {
	filehash := c.Request.FormValue("filehash")

	dbResp, err := dbcli.GetFileMeta(filehash)
	if err != nil {
		c.JSON(
			http.StatusOK,
			gin.H{
				"code": common.StatusServerError,
				"msg":  "server error",
			})
		return
	}

	tblFile := dbcli.ToTableFile(dbResp.Data)

}
