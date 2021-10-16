package config

type MockYaml struct {
	ReadFromFileCalled    bool
	ReadFromFileResult    error
	ReadFromFileObjSetter func(interface{}) error
}

func (y *MockYaml) ReadFromFile(path string, obj interface{}) error {
	y.ReadFromFileCalled = true
	return y.ReadFromFileObjSetter(obj)
}

func SuccesfulRead(obj interface{}) error {
	asserted := obj.(*Config)
	asserted.Version = 1.0
	return nil
}
