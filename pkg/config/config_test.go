package config

import (
	"testing"
)

// func TestIntegratedLoadCfgFile(t *testing.T) {
// 	cwd, _ := os.Getwd()
// 	result, err := LoadCfgFile(
// 		fmt.Sprintf("%s/examples/v1demoapi.yml", cwd),
// 		&Yaml{}
// 	)

// 	if err != nil {
// 		t.Fatalf("error during LoadCfgFile: %s", err.Error())
// 	}

// 	if result.Version != 1.0 {
// 		t.Fatalf("expected result.version == 1.0, got: %f", result.Version)
// 	}
// }

func TestLoadCfgFile(t *testing.T) {
	mockYml := &MockYaml{
		ReadFromFileObjSetter: SuccesfulRead,
	}

	result, err := LoadCfgFile("", mockYml)

	if err != nil {
		t.Fatalf("error during LoadCfgFile: %s", err.Error())
	}

	if !mockYml.ReadFromFileCalled {
		t.Fatal("ReadFromFile not called")
	}

	if result.Version != 1.0 {
		t.Fatalf("expected result.version == 1.0, got: %f", result.Version)
	}
}
