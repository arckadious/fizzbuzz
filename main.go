// main package : read flags, load classes instances and run server
package main

import (
	"flag"
	"fmt"
	"os"
	"path/filepath"

	"github.com/arckadious/fizzbuzz/config"
	"github.com/arckadious/fizzbuzz/container"
	"github.com/arckadious/fizzbuzz/database"
	"github.com/arckadious/fizzbuzz/server"
	"github.com/arckadious/fizzbuzz/validator"
	"github.com/sirupsen/logrus"

	_ "github.com/go-sql-driver/mysql"
)

var paramFile string

func init() {
	p, _ := os.Getwd()

	flag.Usage = func() {
		fmt.Fprintln(os.Stderr, "usage: binary -config=[config]")
		flag.PrintDefaults()
		os.Exit(0)
	}

	flag.String("config", "", "JSON config file")
	flag.Parse()

	configDir := flag.Lookup("config")
	if configDir.Value.String() == configDir.DefValue {
		flag.Set("config", filepath.Join(p, "parameters", "parameters.json"))
	}

	paramFile = flag.Lookup("config").Value.String()
}

func main() {
	validator := validator.New()
	cf, err := config.New(paramFile, *validator)
	if err != nil {
		logrus.Fatal(err)
	}
	server.New(
		container.New(
			cf,
			validator,
			database.Connect(cf),
		),
	).Run()
}
