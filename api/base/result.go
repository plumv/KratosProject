package base

var SUCCESS = &ResultResp{Code: 0, Msg: "成功"}
var ERROR = &ResultResp{Code: -1, Msg: "错误"}
var TIP = &ResultResp{Code: 1, Msg: "提示:"}

func (x *ResultResp) FillMsg(s string) *ResultResp {
	return &ResultResp{Code: x.Code, Msg: x.Msg + s}
}
