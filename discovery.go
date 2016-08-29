package main

import (
	"log"
	"net/http"

	"github.com/spf13/pflag"
	"github.com/spf13/viper"

	"github.com/coreos/go-systemd/activation"

	handling "github.com/coreos/discovery.etcd.io/http"
)

func init() {
	viper.SetEnvPrefix("disc")
	viper.AutomaticEnv()

	pflag.StringP("etcd", "e", "http://127.0.0.1:2379", "etcd endpoint location")
	pflag.StringP("host", "h", "https://discovery.etcd.io", "discovery url prefix")
	pflag.StringP("addr", "a", ":8087", "web service address")

	viper.BindPFlag("etcd", pflag.Lookup("etcd"))
	viper.BindPFlag("host", pflag.Lookup("host"))
	viper.BindPFlag("addr", pflag.Lookup("addr"))

	pflag.Parse()
}

func main() {
	log.SetFlags(0)
	handling.Setup()

	err := http.ListenAndServe(viper.GetString("addr"), nil)
	if err != nil {
		panic(err)
	}

	listeners, err := activation.Listeners(true)
	if err != nil {
		panic(err)
	}

	if len(listeners) != 1 {
		panic("Unexpected number of socket activation fds")
	}

	http.Serve(listeners[0], nil)
}
