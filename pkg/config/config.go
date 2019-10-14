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
	cmd.PersistentFlags().IntVarP(&Instance.ServerDb.RecoveryInterval, "server-RecoveryInterval", "", 5, "恢复周期")
	cmd.PersistentFlags().IntVarP(&Instance.ServerDb.MaxEventToSnapshot, "server-MaxEventToSnapshot", "", 50, "最大事件转换为快照")

	cmd.PersistentFlags().BoolVarP(&Instance.Mysql.Enable, "mysql-enable", "", false, "启用mysql作为后端存储")
	cmd.PersistentFlags().StringVarP(&Instance.Mysql.Address, "mysql-address", "", "localhost", "mysql数据库连接地址")
	cmd.PersistentFlags().StringVarP(&Instance.Mysql.Port, "mysql-port", "", "3306", "mysql数据库端口")
	cmd.PersistentFlags().StringVarP(&Instance.Mysql.Database, "mysql-Database", "", "test", "mysql数据库实例")
	cmd.PersistentFlags().StringVarP(&Instance.Mysql.Username, "mysql-Username", "", "root", "mysql数据库用户名")
	cmd.PersistentFlags().StringVarP(&Instance.Mysql.Password, "mysql-Password", "", "123456", "mysql数据库密码")
	cmd.PersistentFlags().IntVarP(&Instance.Mysql.LifeTime, "mysql-LifeTime", "", 10, "mysql数据库连接最大连接周期(秒)")
	cmd.PersistentFlags().IntVarP(&Instance.Mysql.MaxOpen, "mysql-MaxOpen", "", 5, "mysql数据库最大连接数")
	cmd.PersistentFlags().IntVarP(&Instance.Mysql.MaxIdle, "mysql-MaxIdle", "", 5, "mysql数据库最大等待数量")

	cmd.PersistentFlags().StringVarP(&Instance.Pulsar.Url, "pulsar-url", "", "pulsar://localhost:6650", "pulsar消息中间件地址")
	cmd.PersistentFlags().StringVarP(&Instance.Pulsar.TopicName, "pulsar-topic-name", "", "cqrs-db", "pulsar消息中间件主题名称")

	cmd.PersistentFlags().BoolVarP(&Instance.Mongo.Enable, "mongo-enable", "", true, "启用mongo作为后端存储")
	cmd.PersistentFlags().StringVarP(&Instance.Mongo.Address, "mongo-address", "", "localhost", "mongo数据库连接地址")
	cmd.PersistentFlags().StringVarP(&Instance.Mongo.Port, "mongo-port", "", "27017", "mongo数据库端口")
	cmd.PersistentFlags().StringVarP(&Instance.Mongo.Username, "mongo-Username", "", "root", "数据库用户名")
	cmd.PersistentFlags().StringVarP(&Instance.Mongo.Password, "mongo-Password", "", "123456", "数据库密码")
	cmd.PersistentFlags().IntVarP(&Instance.Mongo.LocalThreshold, "mongo-LocalThreshold", "", 3, "本地阀值")
	cmd.PersistentFlags().IntVarP(&Instance.Mongo.MaxPoolSize, "mongo-MaxPoolSize", "", 10, "最大连接数")
	cmd.PersistentFlags().IntVarP(&Instance.Mongo.MaxConnIdleTime, "mongo-MaxConnIdleTime", "", 5, "最大等待时间")
	cmd.PersistentFlags().StringVarP(&Instance.Mongo.DbName, "mongo-DbName", "", "aggregate", "mongo数据库名称")
	cmd.PersistentFlags().StringVarP(&Instance.Mongo.EventCollectionName, "mongo-EventCollectionName", "", "event", "mongo event集合名称")
	cmd.PersistentFlags().StringVarP(&Instance.Mongo.SnapshotCollectionName, "mongo-SnapshotCollectionName", "", "snapshot", "mongo event集合名称")

	_ = viper.BindPFlag(debugArgName, cmd.PersistentFlags().Lookup(debugArgName))
	_ = viper.BindPFlag("mysql-enable", cmd.PersistentFlags().Lookup("mysql-enable"))
	_ = viper.BindPFlag("mysql-address", cmd.PersistentFlags().Lookup("mysql-address"))
	_ = viper.BindPFlag("mysql-port", cmd.PersistentFlags().Lookup("mysql-port"))
	_ = viper.BindPFlag("mysql-Database", cmd.PersistentFlags().Lookup("mysql-Database"))
	_ = viper.BindPFlag("mysql-Username", cmd.PersistentFlags().Lookup("mysql-Username"))
	_ = viper.BindPFlag("mysql-Password", cmd.PersistentFlags().Lookup("mysql-Password"))
	_ = viper.BindPFlag("mysql-LifeTime", cmd.PersistentFlags().Lookup("mysql-LifeTime"))
	_ = viper.BindPFlag("mysql-MaxOpen", cmd.PersistentFlags().Lookup("mysql-MaxOpen"))
	_ = viper.BindPFlag("mysql-MaxIdle", cmd.PersistentFlags().Lookup("mysql-MaxIdle"))

	_ = viper.BindPFlag("pulsar-url", cmd.PersistentFlags().Lookup("pulsar-url"))
	_ = viper.BindPFlag("pulsar-topic-name", cmd.PersistentFlags().Lookup("pulsar-topic-name"))

	_ = viper.BindPFlag("mongo-enable", cmd.PersistentFlags().Lookup("mongo-enable"))
	_ = viper.BindPFlag("mongo-address", cmd.PersistentFlags().Lookup("mongo-address"))
	_ = viper.BindPFlag("mongo-port", cmd.PersistentFlags().Lookup("mongo-port"))
	_ = viper.BindPFlag("mongo-Username", cmd.PersistentFlags().Lookup("mongo-Username"))
	_ = viper.BindPFlag("mongo-Password", cmd.PersistentFlags().Lookup("mongo-Password"))
	_ = viper.BindPFlag("mongo-LocalThreshold", cmd.PersistentFlags().Lookup("mongo-LocalThreshold"))
	_ = viper.BindPFlag("mongo-MaxPoolSize", cmd.PersistentFlags().Lookup("mongo-MaxPoolSize"))
	_ = viper.BindPFlag("mongo-MaxConnIdleTime", cmd.PersistentFlags().Lookup("mongo-MaxConnIdleTime"))
	_ = viper.BindPFlag("mongo-DbName", cmd.PersistentFlags().Lookup("mongo-DbName"))
	_ = viper.BindPFlag("mongo-EventCollectionName", cmd.PersistentFlags().Lookup("mongo-EventCollectionName"))
	_ = viper.BindPFlag("mongo-SnapshotCollectionName", cmd.PersistentFlags().Lookup("mongo-SnapshotCollectionName"))
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

	LocalThreshold         int
	MaxPoolSize            int
	MaxConnIdleTime        int
	DbName                 string
	EventCollectionName    string
	SnapshotCollectionName string
	Enable                 bool
}

type Config struct {
	Debug    bool
	ServerDb *ServerDbConfig
	Mysql    *MysqlConfig
	Pulsar   *PulsarConfig
	Mongo    *MongoConfig
}

type ServerDbConfig struct {
	Port               string
	Username           string
	Password           string
	RecoveryInterval   int
	MaxEventToSnapshot uint
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
	Enable   bool
}
