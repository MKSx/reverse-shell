package util

import (
	"flag"
	"os"
	"reflect"
	"strconv"

	flags "github.com/jessevdk/go-flags"
)

func ParseArgs(options interface{}) {
	_, err := flags.Parse(options)
	if err != nil {
		if err.(*flags.Error).Type == flags.ErrHelp {
			os.Exit(0)
		} else {
			os.Exit(1)
		}
	}
}

func SetupLogging(options interface{}) {
	flag.Set("logtostderr", "true")
	verbose := reflect.ValueOf(options).FieldByName("Verbose")
	if verbose.IsValid() {
		flag.Set("v", strconv.Itoa(verbose.Interface().(int)))
	}
}

func SetupLogging2(leve int) {
	flag.Set("logtostderr", "true")
	flag.Set("v", strconv.Itoa(leve))
	flag.CommandLine.Parse([]string{})
}
