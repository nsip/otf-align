package otfalign

import (
	util "github.com/nsip/otf-util"
)

type Option func(*OtfAlignService) error

//
// apply all supplied options to the service
// returns any error encountered while applying the options
//
func (srvc *OtfAlignService) setOptions(options ...Option) error {
	for _, opt := range options {
		if err := opt(srvc); err != nil {
			return err
		}
	}
	return nil
}

//
// set the name of this instance of the service
// used for audit tracing purposes if multiple
// instances of the service are active.
// If no name provided a hashid-style unique
// short name will be generated
//
func Name(name string) Option {
	return func(s *OtfAlignService) error {
		if name != "" {
			s.serviceName = name
			return nil
		}
		s.serviceName = util.GenerateName("otf-align")
		return nil
	}
}

//
// create a unique id for this service instance, if none
// provided a nuid will be generated by default
//
func ID(id string) Option {
	return func(s *OtfAlignService) error {
		if id != "" {
			s.serviceID = id
			return nil
		}
		s.serviceID = util.GenerateID()
		return nil
	}
}

//
// set the hostname/address of this service
//
func Host(hname string) Option {
	return func(s *OtfAlignService) error {
		if hname != "" {
			s.serviceHost = hname
			return nil
		}
		s.serviceHost = "localhost"
		return nil
	}
}

//
// set the port to run this service on.
// if 0 then acquire available port from OS
//
func Port(port int) Option {
	return func(s *OtfAlignService) error {
		if port != 0 {
			s.servicePort = port
			return nil
		}
		osPort, err := util.AvailablePort()
		if err != nil {
			return err
		}
		s.servicePort = osPort
		return nil
	}
}

//
// set the hostname/address of the nias3 web server
// defaults to localhost if no host given
//
func NiasHost(hname string) Option {
	return func(s *OtfAlignService) error {
		if hname != "" {
			s.niasHost = hname
			return nil
		}
		s.niasHost = "localhost"
		return nil
	}
}

//
// set the port of the nias3 web server
// defaults to 1323 (n3w defalt port)
//
func NiasPort(port int) Option {
	return func(s *OtfAlignService) error {
		if port != 0 {
			s.niasPort = port
			return nil
		}
		s.niasPort = 1323
		return nil
	}
}

//
// set the access token of the nias3 web server
// defaults to demo otf token if not given
//
func NiasToken(tkn string) Option {
	return func(s *OtfAlignService) error {
		if tkn != "" {
			s.niasToken = tkn
			return nil
		}
		s.niasToken = "Bearer eyJhbGciOiJIUzI1NiIsInR5cCI6IkpXVCJ9.eyJhdWQiOiJkZW1vIiwiY25hbWUiOiJhbGlnbm1lbnRNYXBzIiwidW5hbWUiOiJuc2lwT3RmIn0.Hxepr1xqGpUC6amoO8eTlszM-M2sakLhtwBYHSi-Cig"
		return nil
	}
}

//
// set the hostname/address of the text classifier (otf-classifier) web service
// defaults to localhost if no host given
//
func TcHost(hname string) Option {
	return func(s *OtfAlignService) error {
		if hname != "" {
			s.tcHost = hname
			return nil
		}
		s.tcHost = "localhost"
		return nil
	}
}

//
// set the port of the text classifier web service
// defaults to 1576 (otf-classifier default port)
//
func TcPort(port int) Option {
	return func(s *OtfAlignService) error {
		if port != 0 {
			s.tcPort = port
			return nil
		}
		s.tcPort = 1576
		return nil
	}
}
