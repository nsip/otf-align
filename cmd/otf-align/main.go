package main

import (
	"flag"
	"fmt"
	"os"

	otfal "github.com/nsip/otf-align"
	"github.com/peterbourgon/ff"
)

func main() {

	fs := flag.NewFlagSet("otf-reader", flag.ExitOnError)
	var (
		_           = fs.String("config", "", "config file (optional), json format.")
		serviceName = fs.String("name", "", "name for this alignment service instance")
		serviceID   = fs.String("id", "", "id for this alignment service instance, leave blank to auto-generate a unique id")
		// serviceHost = fs.String("host", "localhost", "name/address of host for this service")
		// servicePort = fs.Int("port", 0, "port to run service on, if not specified will assign an available port automatically")
		// niasHost    = fs.String("niasHost", "localhost", "host name/address of nias3 data server")
		// niasPort    = fs.Int("niasPort", 0, "port that nias3 server is running on")
		// tcHost      = fs.String("tcHost", "localhost", "host name/address of text classification server")
		// tcPort      = fs.Int("tcPort", 0, "port that text classification server is running on")
	)

	ff.Parse(fs, os.Args[1:],
		ff.WithConfigFileFlag("config"),
		ff.WithConfigFileParser(ff.JSONParser),
		ff.WithEnvVarPrefix("OTF_ALIGN_SRVC"),
	)

	opts := []otfal.Option{
		otfal.Name(*serviceName),
		otfal.ID(*serviceID),
		// otfal.NiasPort(*niasPort),
		// otfal.NiasHostName(*niasHost),
	}

	srvc, err := otfal.New(opts...)
	if err != nil {
		fmt.Printf("\nCannot create otf-align service:\n%s\n\n", err)
		return
	}

	srvc.PrintConfig()

}
