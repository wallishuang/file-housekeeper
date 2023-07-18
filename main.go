package main

import (
	"fmt"
	"io/ioutil"
	"os"
	"strings"
	"time"

	"github.com/robfig/cron/v3"
	"github.com/spf13/viper"
)

type Config struct {
	MonitorFolders       string
	DeleteExpirationDays int
	Scheduler            string
}

func getConf() *Config {
	viper.AddConfigPath(".")
	viper.SetConfigName("config")
	err := viper.ReadInConfig()

	if err != nil {
		fmt.Printf("read config toml error: %v", err)
	}

	conf := &Config{}
	err = viper.Unmarshal(conf)
	if err != nil {
		fmt.Printf("unable to decode into config struct, %v", err)
	}

	return conf
}

func main() {
	GoWork()
	config := getConf()
	c := cron.New()
	c.AddFunc(config.Scheduler, func() { GoWork() })
	c.Run()

}

func GoWork() {
	config := getConf()
	fmt.Printf("check monitor folders :%s \n", config.MonitorFolders)
	fmt.Printf("delete files before %d days ago \n", config.DeleteExpirationDays)

	monitorFolders := strings.Split(config.MonitorFolders, ";")
	expirationHours := config.DeleteExpirationDays * 24

	for _, folder := range monitorFolders {
		files, err := ioutil.ReadDir(folder)
		if err != nil {
			fmt.Println(err)
		}

		for _, file := range files {
			expirationTime := time.Now().Add(-time.Hour * time.Duration(expirationHours))
			if file.ModTime().Before(expirationTime) {
				err = os.Remove(folder + "/" + file.Name())
				if err != nil {
					fmt.Println(err)
				} else {
					fmt.Println("file [" + file.Name() + "] has been deleted")
				}
			}
		}
	}
}
