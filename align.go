package otfalign

import (
	"fmt"

	"github.com/labstack/echo"
)

type OtfAlignService struct {
	srvr        *echo.Echo
	serviceName string
	serviceID   string
	serviceHost string
	servicePort int
	niasHost    string
	niasPort    int
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

	return &srvc, nil
}

func (s *OtfAlignService) PrintConfig() {

	fmt.Println("\n\tOTF-Align Service Configuration")
	fmt.Println("\t---------------------------------\n")

	s.printID()
	// rdr.printNiasConfig()
	// rdr.printClassifierConfig()

}

func (s *OtfAlignService) printID() {
	fmt.Println("\talign service name:\t\t", s.serviceName)
	fmt.Println("\talign service ID:\t\t", s.serviceID)
}
