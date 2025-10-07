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

// NewServiceWithParams creates a new Service with XParams.
func NewServiceWithParams(name string, params XParams) *Service {
	core := NewCoreWithParams(name, params)
	return &Service{
		Core:   core,
		Crypto: NewCrypto(core.Cfg().ByteSliceVal(Key.SecEncryptionKey)),
	}
}

// NewServiceWithParamsAndOpts creates a new Service with XParams and additional options.
func NewServiceWithParamsAndOpts(name string, params XParams, opts ...Option) *Service {
	core := NewCoreWithParamsAndOpts(name, params, opts...)
	return &Service{
		Core:   core,
		Crypto: NewCrypto(core.Cfg().ByteSliceVal(Key.SecEncryptionKey)),
	}
}
