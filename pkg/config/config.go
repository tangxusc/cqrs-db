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
	viper.SetEnvPrefix("cqrs")
	viper.AutomaticEnv()

	cmd.PersistentFlags().BoolVarP(&Instance.Debug, debugArgName, "v", false, "debug mod")
	cmd.PersistentFlags().StringVarP(&Instance.ServerDb.Port, "server-port", "p", "3307", "数据库端口")
	cmd.PersistentFlags().StringVarP(&Instance.ServerDb.Username, "server-Username", "u", "root", "用户名")
	cmd.PersistentFlags().StringVarP(&Instance.ServerDb.Password, "server-Password", "d", "123456", "密码")

	cmd.PersistentFlags().StringVarP(&Instance.Mysql.Address, "mysql-address", "", "172.17.0.2", "mysql数据库连接地址")
	cmd.PersistentFlags().StringVarP(&Instance.Mysql.Port, "mysql-port", "", "3306", "mysql数据库端口")
	cmd.PersistentFlags().StringVarP(&Instance.Mysql.Database, "mysql-Database", "", "test", "mysql数据库实例")
	cmd.PersistentFlags().StringVarP(&Instance.Mysql.Username, "mysql-Username", "", "root", "mysql数据库用户名")
	cmd.PersistentFlags().StringVarP(&Instance.Mysql.Password, "mysql-Password", "", "123456", "mysql数据库密码")
	cmd.PersistentFlags().IntVarP(&Instance.Mysql.LifeTime, "mysql-LifeTime", "", 10, "mysql数据库连接最大连接周期(秒)")
	cmd.PersistentFlags().IntVarP(&Instance.Mysql.MaxOpen, "mysql-MaxOpen", "", 5, "mysql数据库最大连接数")
	cmd.PersistentFlags().IntVarP(&Instance.Mysql.MaxIdle, "mysql-MaxIdle", "", 5, "mysql数据库最大等待数量")

	cmd.PersistentFlags().StringVarP(&Instance.Pulsar.Url, "pulsar-url", "", "pulsar://localhost:6650", "pulsar消息中间件地址")
	cmd.PersistentFlags().StringVarP(&Instance.Pulsar.TopicName, "pulsar-topic-name", "", "cqrs-db", "pulsar消息中间件主题名称")

	cmd.PersistentFlags().StringVarP(&Instance.Mongo.Address, "mongo-address", "", "172.17.0.2", "mongo数据库连接地址")
	cmd.PersistentFlags().StringVarP(&Instance.Mongo.Port, "mongo-port", "", "27017", "mongo数据库端口")
	cmd.PersistentFlags().StringVarP(&Instance.Mongo.Username, "mongo-Username", "", "root", "数据库用户名")
	cmd.PersistentFlags().StringVarP(&Instance.Mongo.Password, "mongo-Password", "", "123456", "数据库密码")
	cmd.PersistentFlags().IntVarP(&Instance.Mongo.LocalThreshold, "mongo-LocalThreshold", "", 3, "本地阀值")
	cmd.PersistentFlags().IntVarP(&Instance.Mongo.MaxPoolSize, "mongo-MaxPoolSize", "", 10, "最大连接数")
	cmd.PersistentFlags().IntVarP(&Instance.Mongo.MaxConnIdleTime, "mongo-MaxConnIdleTime", "", 5, "最大等待时间")
	cmd.PersistentFlags().StringVarP(&Instance.Mongo.DbName, "mongo-DbName", "", "game", "mongo数据库名称")
	cmd.PersistentFlags().StringVarP(&Instance.Mongo.CollectionName, "mongo-CollectionName", "", "game", "mongo集合名称")

	//TODO:重写绑定
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

	_ = viper.BindPFlag("pulsar-url", cmd.PersistentFlags().Lookup("pulsar-url"))
	_ = viper.BindPFlag("pulsar-topic-name", cmd.PersistentFlags().Lookup("pulsar-topic-name"))
	_ = viper.BindPFlag("mongo-port", cmd.PersistentFlags().Lookup("mongo-port"))
}

type PulsarConfig struct {
	Url       string
	TopicName string
}

type MongoConfig struct {
	Address  string
	Port     string
	Username string
	Password string

	LocalThreshold  int
	MaxPoolSize     int
	MaxConnIdleTime int
	DbName          string
	CollectionName  string
}

type Config struct {
	Debug    bool
	ServerDb *ServerDbConfig
	Mysql    *MysqlConfig
	Pulsar   *PulsarConfig
	Mongo    *MongoConfig
}

type ServerDbConfig struct {
	Port     string
	Username string
	Password string
}

var Instance = &Config{
	ServerDb: &ServerDbConfig{},
	Mysql:    &MysqlConfig{},
	Pulsar:   &PulsarConfig{},
	Mongo:    &MongoConfig{},
}

type MysqlConfig struct {
	Address  string
	Port     string
	Database string
	Username string
	Password string

	LifeTime int
	MaxOpen  int
	MaxIdle  int
}
