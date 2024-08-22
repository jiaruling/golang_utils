package lib

type BytesTo struct{}

func NewBytesTo() *BytesTo {
	return &BytesTo{}
}

// 将字节转换为KB
func (bt *BytesTo) ByteToKB(bytes uint64) uint64 {
	return bytes / 1024
}

// 将字节转换为MB
func (bt *BytesTo) ByteToMB(bytes uint64) uint64 {
	return bytes / 1024 / 1024
}

// 将字节转换为GB
func (bt *BytesTo) ByteToGB(bytes uint64) uint64 {
	return bytes / 1024 / 1024 / 1024
}

// 将KB转换为字节
func (bt *BytesTo) KBToByte(kb uint64) uint64 {
	return kb * 1024
}

// 将MB转换为字节
func (bt *BytesTo) MBToByte(mb uint64) uint64 {
	return mb * 1024 * 1024
}

// 将GB转换为字节
func (bt *BytesTo) GBToByte(gb uint64) uint64 {
	return gb * 1024 * 1024 * 1024
}
