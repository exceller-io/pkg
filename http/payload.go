package http

import (
	"encoding/json"
	"encoding/xml"
	"fmt"
	"io/ioutil"
	"net/http"

	"go.uber.org/zap"
	"gopkg.in/yaml.v2"
)

const (
	//ContentTypeJSON represents JSON content type
	ContentTypeJSON = "application/json; charset=UTF-8"
	//ContentTypeXML represents XML content type
	ContentTypeXML = "application/xml; charset=UTF-8"
	//ContentTypeYaml represents yaml
	ContentTypeYaml   = "application/x-yaml"
	contentTypeHeader = "Content-Type"
)

//Payload represents an generic interface to read and write data from http
type Payload interface {
	RequestReader
	ResponseReadWriter
}

//RequestReader request reader
type RequestReader interface {
	ReadRequest(contentType string, data interface{}, r *http.Request) error
}

//ResponseReadWriter response read writer
type ResponseReadWriter interface {
	WriteResponse(contentType string, status int, data interface{}, w http.ResponseWriter)
	ReadResponse(data interface{}, response *http.Response) error
}

type payloadImpl struct {
}

//NewPayload initializes a new payload
func NewPayload(log *zap.SugaredLogger) Payload {
	return &payloadImpl{}
}

//Write JSON data with status code to response
func (p *payloadImpl) WriteResponse(contentType string, status int, data interface{}, w http.ResponseWriter) {
	fmt.Printf("HTTP Status : %v Content-Type: %s", status, contentType)

	//write to response
	w.Header().Set(contentTypeHeader, contentType)
	w.WriteHeader(status)
	switch contentType {
	case ContentTypeJSON:
		if err := json.NewEncoder(w).Encode(data); err != nil {
			fmt.Printf("Error Writing JSON to response: %v", err)
			http.Error(w, err.Error(), http.StatusInternalServerError)
		}
	case ContentTypeXML:
		if err := xml.NewEncoder(w).Encode(data); err != nil {
			fmt.Printf("Error writing XML to response: %v", err)
			http.Error(w, err.Error(), http.StatusInternalServerError)
		}
	case ContentTypeYaml:
		if err := yaml.NewEncoder(w).Encode(data); err != nil {
			fmt.Printf("Error writing yaml to response: %v", err)
			http.Error(w, err.Error(), http.StatusInternalServerError)
		}
	default:
		fmt.Printf("Unknown content type: %s", contentType)
	}
}

func (p *payloadImpl) ReadResponse(data interface{}, response *http.Response) error {
	fmt.Println("Reading response body")
	defer response.Body.Close()

	payload, err := ioutil.ReadAll(response.Body)
	if err != nil {
		fmt.Printf("Error reading response from GitHub API call %v", err)
		return err
	}
	if err = json.Unmarshal(payload, &data); err != nil {
		fmt.Printf("Error unmarshalling golang struct")
		return err
	}
	return nil
}

func (p *payloadImpl) ReadRequest(contentType string, data interface{}, r *http.Request) error {
	//read headers from request
	defer r.Body.Close()
	payload, err := ioutil.ReadAll(r.Body)
	if err != nil {
		fmt.Printf("Error reading payload posted in http request: %v", err)
		return err
	}
	switch contentType {
	case ContentTypeJSON:
		if err := json.Unmarshal(payload, &data); err != nil {
			fmt.Printf("Error un-marshalling JSON payload: %v", err)
			return err
		}
		break
	case ContentTypeXML:
		if err := xml.Unmarshal(payload, &data); err != nil {
			fmt.Printf("Error unmarshalling XML payload: %v", err)
			return err
		}
	case ContentTypeYaml:
		if err := yaml.Unmarshal(payload, &data); err != nil {
			fmt.Printf("Error unmarshalling yaml payload: %v", err)
			return err
		}
	default:
		fmt.Printf("Unknown content type: %s", contentType)
	}
	return nil
}
