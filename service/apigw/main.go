package apigw

import "github.com/kuan525/netdisk/service/apigw/route"

func main() {
	r := route.Router()
	r.Run(":8080")
}
