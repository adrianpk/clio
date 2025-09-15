package am

type Service struct {
	Core
	Crypto *Crypto
}

func NewService(name string, opts ...Option) *Service {
	core := NewCore(name, opts...)
	return &Service{
		Core:   core,
		Crypto: NewCrypto(core.Cfg().ByteSliceVal(Key.SecEncryptionKey)),
	}
}
