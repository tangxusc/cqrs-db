package config

import (
	"github.com/sirupsen/logrus"
	"github.com/spf13/cobra"
	"github.com/spf13/viper"
)

const debugArgName = "debug"

func InitLog() {
	if viper.GetBool(debugArgName) {
		logrus.SetLevel(logrus.DebugLevel)
		logrus.SetReportCaller(true)
		logrus.Debug("已开启debug模式...")
	} else {
		logrus.SetLevel(logrus.WarnLevel)
	}

	Instance.Debug = viper.GetBool(debugArgName)
}

func BindParameter(cmd *cobra.Command) {
	viper.SetEnvPrefix("DB")
	viper.AutomaticEnv()

	cmd.PersistentFlags().BoolVarP(&Instance.Debug, debugArgName, "v", false, "debug mod")
	cmd.PersistentFlags().StringVarP(&Instance.Db.Port, "db-port", "p", "3307", "数据库端口")
	cmd.PersistentFlags().StringVarP(&Instance.Db.Username, "db-Username", "u", "root", "用户名")
	cmd.PersistentFlags().StringVarP(&Instance.Db.Password, "db-Password", "d", "123456", "密码")

	cmd.PersistentFlags().StringVarP(&Instance.Proxy.Address, "proxy-address", "", "172.17.0.2", "proxy数据库连接地址")
	cmd.PersistentFlags().StringVarP(&Instance.Proxy.Port, "proxy-port", "", "3306", "proxy数据库端口")
	cmd.PersistentFlags().StringVarP(&Instance.Proxy.Database, "proxy-Database", "", "test", "proxy数据库实例")
	cmd.PersistentFlags().StringVarP(&Instance.Proxy.Username, "proxy-Username", "", "root", "proxy数据库用户名")
	cmd.PersistentFlags().StringVarP(&Instance.Proxy.Password, "proxy-Password", "", "123456", "proxy数据库密码")
	cmd.PersistentFlags().IntVarP(&Instance.Proxy.LifeTime, "proxy-LifeTime", "", 10, "proxy数据库连接最大连接周期(秒)")
	cmd.PersistentFlags().IntVarP(&Instance.Proxy.MaxOpen, "proxy-MaxOpen", "", 5, "proxy数据库最大连接数")
	cmd.PersistentFlags().IntVarP(&Instance.Proxy.MaxIdle, "proxy-MaxIdle", "", 5, "proxy数据库最大等待数量")

	_ = viper.BindPFlag(debugArgName, cmd.PersistentFlags().Lookup(debugArgName))
	_ = viper.BindPFlag("db-port", cmd.PersistentFlags().Lookup("db-port"))
	_ = viper.BindPFlag("db-Username", cmd.PersistentFlags().Lookup("db-Username"))
	_ = viper.BindPFlag("db-Password", cmd.PersistentFlags().Lookup("db-Password"))

	_ = viper.BindPFlag("proxy-address", cmd.PersistentFlags().Lookup("proxy-address"))
	_ = viper.BindPFlag("proxy-port", cmd.PersistentFlags().Lookup("proxy-port"))
	_ = viper.BindPFlag("proxy-Database", cmd.PersistentFlags().Lookup("proxy-Database"))
	_ = viper.BindPFlag("proxy-Username", cmd.PersistentFlags().Lookup("proxy-Username"))
	_ = viper.BindPFlag("proxy-Password", cmd.PersistentFlags().Lookup("proxy-Password"))
	_ = viper.BindPFlag("proxy-LifeTime", cmd.PersistentFlags().Lookup("proxy-LifeTime"))
	_ = viper.BindPFlag("proxy-MaxOpen", cmd.PersistentFlags().Lookup("proxy-MaxOpen"))
	_ = viper.BindPFlag("proxy-MaxIdle", cmd.PersistentFlags().Lookup("proxy-MaxIdle"))
}

type Config struct {
	Debug bool
	Db    *DbConfig
	Proxy *ProxyConfig
}

type DbConfig struct {
	Port     string
	Username string
	Password string
}

var Instance = &Config{
	Db:    &DbConfig{},
	Proxy: &ProxyConfig{},
}

type ProxyConfig struct {
	Address  string
	Port     string
	Database string
	Username string
	Password string

	LifeTime int
	MaxOpen  int
	MaxIdle  int
}
