package main

import (
	"fmt"
	"log"
	"net/http"
	"net/url"
	"os"

	"github.com/spf13/pflag"
	"github.com/spf13/viper"

	"github.com/coreos/go-systemd/activation"

	handling "github.com/coreos/discovery.etcd.io/http"
)

func fail(err string) {
	log.Print(err)
	pflag.PrintDefaults()
	os.Exit(2) // default go flag error code
}

func getValidUrlOrExit(givenUrl string) string {
	u, err := url.Parse(givenUrl)

	if err != nil {
		fail(fmt.Sprintf("Invalid url given: %v", err))
	}

	if len(u.Path) > 1 || (len(u.Path) == 1 && u.Path != "/") {
		fail(fmt.Sprintf("Expected url without path (%v)", u.Path))
	}

	if u.RawQuery != "" {
		fail(fmt.Sprintf("Expected url without query (?%v)", u.RawQuery))
	}

	if u.Fragment != "" {
		fail(fmt.Sprintf("Expected url without fragment (%v)", u.Fragment))
	}

	if u.Host == "" {
		fail(fmt.Sprint("Expected hostname (none given)"))
	}

	return u.Scheme + "://" + u.Host
}

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
	viper.Set("etcd", getValidUrlOrExit(viper.GetString("etcd")))
	viper.Set("host", getValidUrlOrExit(viper.GetString("host")))

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
