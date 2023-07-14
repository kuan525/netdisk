package mapper

import (
	"dbclient/orm"
	"errors"
	"reflect"
)

var funcs = map[string]interface{}{
	"/file/OnFileUploadFinished": orm.OnFileUploadFinished,
	"/file/GetFileMeta":          orm.GetFileMeta,
	"/file/GetFileMetaList":      orm.GetFileMetaList,
	"/file/UpdateFileLocation":   orm.UpdateFileLocation,

	"/user/UserSignup":  orm.UserSignup,
	"/user/UserSignin":  orm.UserSignin,
	"/user/UpdateToken": orm.UpdateToken,
	"/user/GetUserInfo": orm.GetUserInfo,
	"/user/UserExist":   orm.UserExist,

	"/ufile/OnUserFileUploadFinished": orm.OnUserFileUploadFinished,
	"/ufile/QueryUserFileMetas":       orm.QueryUserFileMetas,
	"/ufile/DeleteUserFile":           orm.DeleteUserFile,
	"/ufile/RenameFileName":           orm.RenameFileName,
	"/ufile/QueryUserFileMeta":        orm.QueryUserFileMeta,
}

func FuncCall(name string, params ...interface{}) (result []reflect.Value, err error) {
	if _, ok := funcs[name]; !ok {
		err = errors.New("函数名不存在")
		return
	}

	// 通过反射可以动态调用对象的导出方法
	f := reflect.ValueOf(funcs[name])
	if len(params) != f.Type().NumIn() {
		err = errors.New("传入参数数量与调用方法要求的数量不一致")
		return
	}

	// 构造一个Value的slice，用作Call方法的传入参数
	in := make([]reflect.Value, len(params))
	for k, param := range params {
		in[k] = reflect.ValueOf(param)
	}

	// 执行方法f，并将方法结果复制给result
	result = f.Call(in)
	return
}
