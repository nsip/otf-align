package util

import (
	"crypto/rand"
	"io"
	"log"
	"math/big"
	"net"
	"net/http"
	"time"

	"github.com/nats-io/nuid"
	"github.com/pkg/errors"
	hashids "github.com/speps/go-hashids"
)

//
// generate a short useful unique name - hashid in this case
//
func GenerateName() string {

	name := "reader"

	// generate a random number
	number0, err := rand.Int(rand.Reader, big.NewInt(10000000))

	hd := hashids.NewData()
	hd.Salt = "otf-reader random name generator 2020"
	hd.MinLength = 5
	h, err := hashids.NewWithData(hd)
	if err != nil {
		log.Println("error auto-generating name: ", err)
		return name
	}
	e, err := h.EncodeInt64([]int64{number0.Int64()})
	if err != nil {
		log.Println("error encoding auto-generated name: ", err)
		return name
	}
	name = e

	return name

}

//
// generate a unique id - nuid in this case
//
func GenerateID() string {

	return nuid.Next()

}

// Fetch makes network calls using the method (POST/GET..), the URL // to hit, headers to add (if any), and the body of the request.
// Feel free to add more stuff to before/after making the actual n/w call!
func Fetch(method string, url string, header map[string]string, body io.Reader) (*http.Response, error) {
	// Create client with required custom parameters.
	// Options: Disable keep-alives, 30sec n/w call timeout.
	client := &http.Client{
		Transport: &http.Transport{
			DisableKeepAlives: true,
		},
		Timeout: time.Duration(10 * time.Second),
	}
	// Create request.
	req, _ := http.NewRequest(method, url, body)
	// Add any required headers.
	for key, value := range header {
		req.Header.Add(key, value)
	}
	// Perform said network call.
	res, err := client.Do(req)
	if err != nil {
		// glog.Error(err.Error()) // use glog it's amazing
		return nil, err
	}
	// If response from network call is not 200, return error too.
	if res.StatusCode != http.StatusOK {
		return res, errors.New("Network call did not return SUCCESS!")
	}
	return res, nil
}

//
// small utility function embedded in major ops
// to print a performance indicator.
//
func TimeTrack(start time.Time, name string) {
	elapsed := time.Since(start)
	log.Printf("%s took %s", name, elapsed.Truncate(time.Millisecond).String())

}

//
// find an available tcp port
//
func AvailablePort() (int, error) {

	listener, err := net.Listen("tcp", ":0")
	if err != nil {
		return 0, errors.Wrap(err, "cannot acquire a tcp port")
	}

	return listener.Addr().(*net.TCPAddr).Port, nil

}
