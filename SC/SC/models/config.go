package models

// 配置文件对应结构体对象
// 对应yaml的mysql  (Mysql--转成小写对应), 如果有不符合的就用标签去指定 `mapstructure:"db"`
type Config struct {
	// 匿名结构体
	// 数据库信息
	MySQL struct {
		Host     string `yaml:"host"`
		Port     int    `yaml:"port"`
		UserName string `yaml:"username"`
		PassWord string `yaml:"password"`
		DBName   string `yaml:"dbname"`
	} `mapstructure:"mysql"`

	// SenseCity 相关信息
	SenseCity struct {
		SCIP         string `yaml:"scip"`
		Token        string `yaml:"token"`
		UserFilePath string `yaml:"userfilepath"`
		UserName     string `yaml:"username"`
		PassWord     string `yaml:"password"`
		GrantType    string `yaml:"grant_type"`
	} `mapstructure:"sensecity"`

	// 日志相关目录
	Log struct {
		Log_path    string `yaml:"log_path"`
		Max_size    int    `yaml:"max_size"`
		Max_backups int    `yaml:"max_backups"`
		Compress    bool   `yaml:"compress"`
	} `mapstructure:"log"`

	// 调用身份库接口参数
	Portrait struct {
		Page          int64    `yaml:"page"`
		SenseType     int64    `yaml:"sensetype"`
		SageSize      int64    `yaml:"pagesize"`
		TarLibSerials []string `yaml:"tarlibserials"`
	} `mapstructure:"portrait"`
}
