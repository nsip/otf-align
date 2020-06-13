package otfalign

import (
	"context"
	"fmt"
	"net/http"
	"os"
	"time"

	"github.com/labstack/echo"
	"github.com/labstack/gommon/log"
)

type OtfAlignService struct {
	e           *echo.Echo
	serviceName string
	serviceID   string
	serviceHost string
	servicePort int
	niasHost    string
	niasPort    int
	niasToken   string
	tcHost      string
	tcPort      int
}

//
// create a new service instance
//
func New(options ...Option) (*OtfAlignService, error) {

	srvc := OtfAlignService{}

	if err := srvc.setOptions(options...); err != nil {
		return nil, err
	}

	srvc.e = echo.New()
	srvc.e.Logger.SetLevel(log.INFO)
	// add pingable method to know we're up
	srvc.e.GET("/", func(c echo.Context) error {
		return c.JSON(http.StatusOK, "OK")
	})

	return &srvc, nil
}

//
// start the service running
//
func (s *OtfAlignService) Start() {

	address := fmt.Sprintf("%s:%d", s.serviceHost, s.servicePort)
	go func(addr string) {
		if err := s.e.Start(addr); err != nil {
			s.e.Logger.Info("error starting server: ", err, ", shutting down...")
			// attempt clean shutdown by raising sig int
			p, _ := os.FindProcess(os.Getpid())
			p.Signal(os.Interrupt)
		}
	}(address)

}

//
// shut the server down gracefully
//
func (s *OtfAlignService) Shutdown() {
	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()
	if err := s.e.Shutdown(ctx); err != nil {
		fmt.Println("could not shut down server cleanly: ", err)
		s.e.Logger.Fatal(err)
	}

}

func (s *OtfAlignService) PrintConfig() {

	fmt.Println("\n\tOTF-Align Service Configuration")
	fmt.Println("\t---------------------------------\n")

	s.printID()
	s.printNiasConfig()
	s.printClassifierConfig()

}

func (s *OtfAlignService) printID() {
	fmt.Println("\tservice name:\t\t", s.serviceName)
	fmt.Println("\tservice ID:\t\t", s.serviceID)
	fmt.Println("\tservice host:\t\t", s.serviceHost)
	fmt.Println("\tservice port:\t\t", s.servicePort)
}

func (s *OtfAlignService) printNiasConfig() {
	fmt.Println("\tnias n3w host:\t\t", s.niasHost)
	fmt.Println("\tnias n3w port:\t\t", s.niasPort)
	fmt.Println("\tnias token:\t\t", s.niasToken)
}

func (s *OtfAlignService) printClassifierConfig() {
	fmt.Println("\totf-class host:\t\t", s.tcHost)
	fmt.Println("\totf-class port:\t\t", s.tcPort)
}
