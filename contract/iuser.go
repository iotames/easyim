package contract

// 定义服务接口
type IUser interface {
	ReceiveData([]byte)
	SendData([]byte)
}
