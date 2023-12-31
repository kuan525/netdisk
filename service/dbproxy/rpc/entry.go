package rpc

import (
	"bytes"
	"context"
	"encoding/json"
	dbProxy "github.com/kuan525/netdisk/client/dbproxy/proto"
	"github.com/kuan525/netdisk/common/mapper"
	"github.com/kuan525/netdisk/common/orm"
)

// DBProxy 结构体
type DBProxy struct {
	dbProxy.UnimplementedDBProxyServiceServer
}

// ExecuteAction 请求执行sql函数
func (db *DBProxy) ExecuteAction(ctx context.Context, req *dbProxy.ReqExec) (res *dbProxy.RespExec, err error) {
	res = new(dbProxy.RespExec)
	resList := make([]orm.ExecResult, len(req.Action))

	// TODO: 检查 req.Sequence req.Transaction两个参数，执行不同的流程
	for idx, singleAction := range req.Action {
		var params []interface{}
		dec := json.NewDecoder(bytes.NewReader(singleAction.Params))
		dec.UseNumber()
		// 避免int/int32/int64等自动转换为float64
		if err := dec.Decode(&params); err != nil {
			resList[idx] = orm.ExecResult{
				Suc: false,
				Msg: "请求参数有误",
			}
			continue
		}

		for k, v := range params {
			if _, ok := v.(json.Number); ok {
				params[k], _ = v.(json.Number).Int64()
			}
		}

		execRes, err := mapper.FuncCall(singleAction.Name, params...)
		if err != nil {
			resList[idx] = orm.ExecResult{
				Suc: false,
				Msg: "函数调用有误",
			}
			continue
		}

		resList[idx] = execRes[0].Interface().(orm.ExecResult)
	}

	// TODO： 处理异常
	res.Data, err = json.Marshal(resList)
	return
}
