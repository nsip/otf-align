package otfalign

import (
	"context"
	"fmt"
	"net/http"
	"net/http/httputil"
	"os"
	"time"

	"github.com/labstack/echo"
	"github.com/labstack/gommon/log"
)

type OtfAlignService struct {
	// embedded web server to handle alignment requests
	e *echo.Echo
	// the unique name of this service when running multiple instances
	serviceName string
	// the unique id of this service when running multiple instances
	serviceID string
	// the host address this service instance is running on
	serviceHost string
	// the port that this service instance is running on
	servicePort int
	// the host address of the nias3 server used for data lookups
	niasHost string
	// the port the nias3 server is running on
	niasPort int
	// the jwt used to acess the nias service
	niasToken string
	// the host address of the text classifier service
	tcHost string
	// the port of the text classifier service
	tcPort int
}

//
// Query paramters sent to the
// web service.
// Params can be provided as json payload, via form components
// or as query params
//
type AlignRequest struct {
	//
	// method to be used for alignment one of...
	// prescribed: results in lookup/passthrough of NLP reference
	// mapped: maps from input token through known linkages such as Australian Curriculum to find link to NLP
	// inferred: uses text classifier lookup to try and identify desired NLP
	//
	AlignMethod string `json:"alignMethod" form:"alignMethod" query:"alignMethod"`
	//
	// parameter to guide chosen method...
	// prescribed: will typically be an NLP reference. Lookup may still occur to find full extent of GESDI block, or value may simply be passed through/back to user
	// mapped: will typically be a module or node reference in the providing system, which in turn will be looked up in avialable vendor maps to find link to NLP via (for example) a common Australian Curriculum link
	// inferred: will typically be a piece of free-form text such as a question or observation
	//
	AlignToken interface{} `json:"alignToken" form:"alignToken" query:"alignToken"`
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
	// add align method
	srvc.e.POST("/align", align)

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
// the main align method
// requires an input of request variables (in json)
// alignMethod: one of (prescribed|mapped|inferred)
// alignToken: string (reference such as an AC ref for mapped alignment,
// or the text to be used as input
// to the text classifier for inferred alignment)
// prescribed looks up full GESDI if necessary.
//
func align(c echo.Context) error {

	//
	// TODO: disable for production/release
	// show the full request
	//
	requestDump, err := httputil.DumpRequest(c.Request(), true)
	if err != nil {
		fmt.Println("req-dump error: ", err)
	}
	fmt.Println(string(requestDump))

	// check required params are in input
	ar := &AlignRequest{}
	if err := c.Bind(ar); err != nil {
		return echo.NewHTTPError(http.StatusBadRequest, err)
	}

	// token could be any json type so convert to string
	stringToken := fmt.Sprintf("%v", ar.AlignToken)

	if ar.AlignMethod == "" || stringToken == "" {
		fmt.Println("align binding failed")
		return echo.NewHTTPError(http.StatusBadRequest, "must supply values for alignMethod and alignToken")
	}

	// fmt.Printf("\ninput:\n%#v\n", ar)

	return c.JSON(http.StatusOK, ar)

	// switch based on method

	// call tc for inferred

	// call nias for mapped

	// call nias for gesdi

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
