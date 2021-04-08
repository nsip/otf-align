package main

import (
	"flag"
	"fmt"
	"os"
	"os/signal"
	"syscall"

	otfal "github.com/nsip/otf-align"
	"github.com/peterbourgon/ff/v3"
)

func main() {

	fs := flag.NewFlagSet("otf-reader", flag.ExitOnError)
	var (
		_           = fs.String("config", "", "config file (optional), json format.")
		serviceName = fs.String("name", "", "name for this alignment service instance")
		serviceID   = fs.String("id", "", "id for this alignment service instance, leave blank to auto-generate a unique id")
		serviceHost = fs.String("host", "localhost", "name/address of host for this service")
		servicePort = fs.Int("port", 0, "port to run service on, if not specified will assign an available port automatically")
		niasHost    = fs.String("niasHost", "localhost", "host name/address of nias3 (n3w) web service")
		niasPort    = fs.Int("niasPort", 1323, "port that nias3 web (n3w) service is running on")
		niasToken   = fs.String("niasToken", "", "access token for nias server when making queries")
		tcHost      = fs.String("tcHost", "localhost", "host name/address of text classification server")
		tcPort      = fs.Int("tcPort", 1576, "port that text classification server is running on")
	)

	ff.Parse(fs, os.Args[1:],
		ff.WithConfigFileFlag("config"),
		ff.WithConfigFileParser(ff.JSONParser),
		ff.WithEnvVarPrefix("OTF_ALIGN_SRVC"),
	)

	opts := []otfal.Option{
		otfal.Name(*serviceName),
		otfal.ID(*serviceID),
		otfal.Host(*serviceHost),
		otfal.Port(*servicePort),
		otfal.NiasHost(*niasHost),
		otfal.NiasPort(*niasPort),
		otfal.NiasToken(*niasToken),
		otfal.TcHost(*tcHost),
		otfal.TcPort(*tcPort),
	}

	srvc, err := otfal.New(opts...)
	if err != nil {
		fmt.Printf("\nCannot create otf-align service:\n%s\n\n", err)
		return
	}

	srvc.PrintConfig()

	// signal handler for shutdown
	closed := make(chan struct{})
	c := make(chan os.Signal)
	signal.Notify(c, os.Interrupt, syscall.SIGTERM)
	go func() {
		<-c
		fmt.Println("\notf-align shutting down")
		srvc.Shutdown()
		fmt.Println("otf-align closed")
		close(closed)
	}()

	srvc.Start()

	// block until shutdown by sig-handler
	<-closed

}
