package config

import (
	"fmt"
	"github.com/spf13/pflag"
	"github.com/spf13/viper"
	"path"
	"strings"
)

const (
	DBPostgres = "postgres"
	Dev        = "dev"
	Local      = "local"
)

type Config struct {
	viper.Viper
	ConfigName           string
	ConfigType           string
	ConfigPath           string
	LogLevel             int
	Env                  string
	DBtype               string
	DBhost               string
	DBport               int
	DBname               string
	DBuser               string
	DBpass               string
	MigrationPath        string
	APIPort              int
	APIKey               string
	UseProxy             bool
	UserProxy            string
	PassProxy            string
	SchemaProxy          string
	HostProxy            string
	PortProxy            string
	KeyManagerServiceURL string
	S3Region             string
	S3AccountKey         string
	S3Secret             string
	S3Endpoint           string
	S3Bucket             string
	PythonName           string
	KafkaBrokers         []string
}

func NewConfig() *Config {
	return &Config{
		Viper:                *viper.New(),
		ConfigName:           "config",
		ConfigType:           "json",
		ConfigPath:           "./config/",
		LogLevel:             -1,
		DBtype:               DBPostgres,
		DBhost:               "localhost",
		DBport:               6434,
		DBname:               "dev_app_db",
		DBuser:               "developer",
		DBpass:               "developer",
		MigrationPath:        "./migrations/postgres",
		APIPort:              8010,
		APIKey:               "test_api_key",
		Env:                  Local,
		KeyManagerServiceURL: "http://localhost:8003/api/v1/auth/google",
		UseProxy:             false,
		UserProxy:            "",
		PassProxy:            "",
		SchemaProxy:          "http",
		HostProxy:            "",
		PortProxy:            "",
		S3Region:             "",
		S3AccountKey:         "",
		S3Secret:             "",
		S3Endpoint:           "",
		S3Bucket:             "apks",
		PythonName:           "python",
		KafkaBrokers:         []string{"localhost:9092"},
	}
}

func init() {
	pflag.String("config", "", "config file path")
	pflag.Int("loglevel", 0, "debug level (debug:-1 .. fatal:5)")
	pflag.Int("apiport", 0, "API port")
	pflag.String("dbtype", "", "DB type")
	pflag.String("dbhost", "", "DB host")
	pflag.String("dbport", "", "DB port")
	pflag.String("dbname", "", "DB name")
	pflag.String("dbuser", "", "DB user")
	pflag.String("dbpass", "", "DB pass")
	pflag.String("migrationpath", "", "Migration DB path")
	pflag.String("env", "", "environment")
	pflag.String("keymanagerserviceurl", "", "key manager service url")
	pflag.Bool("useproxy", false, "Use Proxy")
	pflag.String("userproxy", "", "User proxy")
	pflag.String("passproxy", "", "Password proxy")
	pflag.String("schemaproxy", "", "Schema proxy (http or https)")
	pflag.String("hostproxy", "", "Host proxy")
	pflag.String("portproxy", "", "Port proxy")
	pflag.String("s3region", "", "S3 region")
	pflag.String("s3accountkey", "", "S3 account key")
	pflag.String("s3secret", "", "S3 secret")
	pflag.String("s3endpoint", "", "S3 endpoint")
	pflag.String("s3bucket", "", "S3 bucket")
	pflag.String("pythonname", "", "python command name")
	pflag.StringArray("kafkabrokers", []string{""}, "kafka brokers")
}

func (conf *Config) bindAllEnv() {
	_ = conf.BindEnv("loglevel", "AM_SERVER_LOGLEVEL")
	_ = conf.BindEnv("apiport", "AM_SERVER_APIPORT")
	_ = conf.BindEnv("dbtype", "AM_SERVER_DBTYPE")
	_ = conf.BindEnv("dbhost", "AM_SERVER_DBHOST")
	_ = conf.BindEnv("dbport", "AM_SERVER_DBPORT")
	_ = conf.BindEnv("dbname", "AM_SERVER_DBNAME")
	_ = conf.BindEnv("dbuser", "AM_SERVER_DBUSER")
	_ = conf.BindEnv("dbpass", "AM_SERVER_DBPASS")
	_ = conf.BindEnv("migrationpath", "AM_MIGRATION_PATH")
	_ = conf.BindEnv("env", "AM_SERVER_ENV")
	_ = conf.BindEnv("keymanagerserviceurl", "AM_SERVICE")
	_ = conf.BindEnv("useproxy", "AM_USE_PROXY")
	_ = conf.BindEnv("userproxy", "AM_USER_PROXY")
	_ = conf.BindEnv("passproxy", "AM_PASS_PROXY")
	_ = conf.BindEnv("schemaproxy", "AM_SCHEMA_PROXY")
	_ = conf.BindEnv("hostproxy", "AM_HOST_PROXY")
	_ = conf.BindEnv("portproxy", "AM_PORT_PROXY")
	_ = conf.BindEnv("s3region", "AM_S3_REGION")
	_ = conf.BindEnv("s3accountkey", "AM_S3_ACCOUNT_KEY")
	_ = conf.BindEnv("s3secret", "AM_S3_SECRET")
	_ = conf.BindEnv("s3endpoint", "AM_S3_ENDPOINT")
	_ = conf.BindEnv("s3bucket", "AM_S3_BUCKET")
	_ = conf.BindEnv("pythonname", "AM_PYTHON_NAME")
	_ = conf.BindEnv("kafkabrokers", "AM_KAFKA_BROKERS")
}

func (conf *Config) bindAllFlags() {
	pflag.Parse()
	_ = conf.BindPFlag("config", pflag.Lookup("config"))
	_ = conf.BindPFlag("apiport", pflag.Lookup("apiport"))
	_ = conf.BindPFlag("loglevel", pflag.Lookup("loglevel"))
	_ = conf.BindPFlag("dbtype", pflag.Lookup("dbtype"))
	_ = conf.BindPFlag("dbhost", pflag.Lookup("dbhost"))
	_ = conf.BindPFlag("dbport", pflag.Lookup("dbport"))
	_ = conf.BindPFlag("dbname", pflag.Lookup("dbname"))
	_ = conf.BindPFlag("dbuser", pflag.Lookup("dbuser"))
	_ = conf.BindPFlag("dbpass", pflag.Lookup("dbpass"))
	_ = conf.BindPFlag("migrationpath", pflag.Lookup("migrationpath"))
	_ = conf.BindPFlag("env", pflag.Lookup("env"))
	_ = conf.BindPFlag("keymanagerserviceurl", pflag.Lookup("keymanagerserviceurl"))
	_ = conf.BindPFlag("useproxy", pflag.Lookup("useproxy"))
	_ = conf.BindPFlag("userproxy", pflag.Lookup("userproxy"))
	_ = conf.BindPFlag("passproxy", pflag.Lookup("passproxy"))
	_ = conf.BindPFlag("schemaproxy", pflag.Lookup("schemaproxy"))
	_ = conf.BindPFlag("hostproxy", pflag.Lookup("hostproxy"))
	_ = conf.BindPFlag("portproxy", pflag.Lookup("portproxy"))
	_ = conf.BindPFlag("s3region", pflag.Lookup("s3region"))
	_ = conf.BindPFlag("s3accountkey", pflag.Lookup("s3accountkey"))
	_ = conf.BindPFlag("s3secret", pflag.Lookup("s3secret"))
	_ = conf.BindPFlag("s3endpoint", pflag.Lookup("s3endpoint"))
	_ = conf.BindPFlag("s3bucket", pflag.Lookup("s3bucket"))
	_ = conf.BindPFlag("pythonname", pflag.Lookup("pythonname"))
	_ = conf.BindPFlag("kafkabrokers", pflag.Lookup("kafkabrokers"))
}

func (conf *Config) setDefaults() {
	conf.SetDefault("config", fmt.Sprintf("%s/%s.%s", strings.TrimSuffix(conf.ConfigPath, "/"), conf.ConfigName, conf.ConfigType))
	conf.SetDefault("apiport", conf.APIPort)
	conf.SetDefault("apikey", conf.APIKey)
	conf.SetDefault("loglevel", conf.LogLevel)
	conf.SetDefault("dbtype", conf.DBtype)
	conf.SetDefault("dbhost", conf.DBhost)
	conf.SetDefault("dbport", conf.DBport)
	conf.SetDefault("dbname", conf.DBname)
	conf.SetDefault("dbuser", conf.DBuser)
	conf.SetDefault("dbpass", conf.DBpass)
	conf.SetDefault("migrationpath", conf.MigrationPath)
	conf.SetDefault("env", conf.Env)
	conf.SetDefault("keymanagerserviceurl", conf.KeyManagerServiceURL)
	conf.SetDefault("useproxy", conf.UseProxy)
	conf.SetDefault("userproxy", conf.UserProxy)
	conf.SetDefault("passproxy", conf.PassProxy)
	conf.SetDefault("schemaproxy", conf.SchemaProxy)
	conf.SetDefault("hostproxy", conf.HostProxy)
	conf.SetDefault("portproxy", conf.PortProxy)
	conf.SetDefault("s3region", conf.S3Region)
	conf.SetDefault("s3accountkey", conf.S3AccountKey)
	conf.SetDefault("s3secret", conf.S3Secret)
	conf.SetDefault("s3endpoint", conf.S3Endpoint)
	conf.SetDefault("s3bucket", conf.S3Bucket)
	conf.SetDefault("pythonname", conf.PythonName)
	conf.SetDefault("kafkabrokers", conf.KafkaBrokers)
}

// ReadSettings ...
// viper precedence order:
// 1 explicit call to Set
// 2 flag
// 3 env
// 4 config
// 5 key/value store
// 6 default
func (conf *Config) ReadAllSettings() error {
	conf.setDefaults()
	conf.bindAllEnv()
	conf.bindAllFlags()

	flagConfig := conf.GetString("config")
	if flagConfig != "" {
		conf.ConfigPath = path.Dir(flagConfig)
		conf.ConfigType = strings.TrimPrefix(path.Ext(flagConfig), ".")
		if conf.ConfigType != "" {
			conf.ConfigName = strings.TrimSuffix(path.Base(flagConfig), "."+conf.ConfigType)
		} else {
			conf.ConfigName = path.Base(flagConfig)
		}
	}

	conf.SetConfigName(conf.ConfigName) // read config
	conf.SetConfigType(conf.ConfigType)
	conf.AddConfigPath(conf.ConfigPath)
	if err := conf.ReadInConfig(); err != nil {
		if errW := conf.WriteConfigAs(fmt.Sprintf("%s/%s.%s", conf.ConfigPath, conf.ConfigName, conf.ConfigType)); errW != nil {
			return err
		}
	}

	conf.APIPort = conf.GetInt("apiport")
	conf.APIKey = conf.GetString("apikey")
	conf.LogLevel = conf.GetInt("loglevel")
	conf.DBtype = conf.GetString("dbtype")
	conf.DBhost = conf.GetString("dbhost")
	conf.DBport = conf.GetInt("dbport")
	conf.DBname = conf.GetString("dbname")
	conf.DBuser = conf.GetString("dbuser")
	conf.DBpass = conf.GetString("dbpass")
	conf.MigrationPath = conf.GetString("migrationpath")
	conf.Env = conf.GetString("env")
	conf.KeyManagerServiceURL = conf.GetString("keymanagerserviceurl")
	conf.UseProxy = conf.GetBool("useproxy")
	conf.UserProxy = conf.GetString("userproxy")
	conf.PassProxy = conf.GetString("passproxy")
	conf.SchemaProxy = conf.GetString("schemaproxy")
	conf.HostProxy = conf.GetString("hostproxy")
	conf.PortProxy = conf.GetString("portproxy")
	conf.S3Region = conf.GetString("s3region")
	conf.S3AccountKey = conf.GetString("s3accountkey")
	conf.S3Secret = conf.GetString("s3secret")
	conf.S3Endpoint = conf.GetString("s3endpoint")
	conf.S3Bucket = conf.GetString("s3bucket")
	conf.PythonName = conf.GetString("pythonname")
	conf.KafkaBrokers = conf.GetStringSlice("kafkabrokers")

	return nil
}
