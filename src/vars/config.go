package vars

type GenMiniConfig struct {
	GenMiniConfigItem `yaml:"GenMini"`
}

type GenMiniConfigItem struct {
	AppKey  string `yaml:"AppKey"`
	BaseUrl string `yaml:"BaseUrl"`
}

func (g GenMiniConfig) String() string {
	return "AppKey: " + g.GenMiniConfigItem.AppKey + "\n" + "BaseUrl: " + g.GenMiniConfigItem.BaseUrl
}
