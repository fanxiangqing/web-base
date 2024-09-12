package env

import (
	"errors"
	"github.com/fanxiangqing/web-base/lib/utils"
	"os"
	"reflect"
	"strings"

	"github.com/sirupsen/logrus"

	"github.com/joho/godotenv"
)

const (
	EnvModeProd = "prod"
	EnvModeDev  = "dev"
)

// .env文件结构体
type env struct {
	RootPath string
	EnvMode  string
}

// Env 定义全局环境变量
var Env *env

func init() {
	Env = &env{}
	// 获取项目根目录
	Env.RootPath = getRootPath()
	// 读取系统环境变量
	readEnvironment()
	if Env.EnvMode == "" {
		loadEnvFile()
		if os.Getenv("EnvMode") != "" {
			v := reflect.ValueOf(Env).Elem()
			for i := 0; i < v.NumField(); i++ {
				f := v.Field(i)
				t := v.Type()
				if f.Interface() == "" {
					val := os.Getenv(t.Field(i).Name)
					if val == "" {
						logrus.Fatal(errors.New("env file loss key:" + t.Field(i).Name))
					}
					f.SetString(val)
				}
			}
		}
	}
}

// 获取项目根目录
func getRootPath() string {
	dir, err := os.Getwd()
	if err != nil {
		logrus.Fatal("failed to get path")
	}

	return strings.Replace(dir, "\\", "/", -1)
}

// 加载.env文件
func loadEnvFile() {
	file := Env.RootPath + "/.env"
	if utils.Exists(file) {
		err := godotenv.Load(file)
		if err != nil {
			logrus.Fatal("failed to load file .env")
		}
	} else {
		// logrus.Fatal("not exists file .env")
		Env.EnvMode = EnvModeDev
	}
}

// 读取系统环境变量
func readEnvironment() {
	getType1 := reflect.TypeOf(Env).Elem()
	getValue1 := reflect.ValueOf(Env).Elem()
	for i := 0; i < getType1.NumField(); i++ {
		tp := getType1.Field(i)
		val := getValue1.Field(i)
		osVal := os.Getenv(tp.Name)
		if osVal != "" {
			val.SetString(osVal)
		}
	}
}
