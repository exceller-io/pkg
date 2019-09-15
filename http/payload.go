// Copyright ©  2019 AppsByRam authors.
//
// Licensed under the Apache License, Version 2.0 (the "License");
// you may not use this file except in compliance with the License.
// You may obtain a copy of the License at
//
//    http://www.apache.org/licenses/LICENSE-2.0
//
// Unless required by applicable law or agreed to in writing, software
// distributed under the License is distributed on an "AS IS" BASIS,
// WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
// See the License for the specific language governing permissions and
// limitations under the License.

package http

import (
	"encoding/json"
	"encoding/xml"
	"fmt"
	"io/ioutil"
	"net/http"

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

//Payload represents an generic interface to read and write data from HTTP Request or Response
type Payload interface {
	RequestReader
	ResponseReadWriter
}

//RequestReader represents a generic interface to read from HTTP Request
type RequestReader interface {
	ReadRequest(contentType string, data interface{}, r *http.Request) error
}

//ResponseReadWriter represents a generic interface to read from and write to HTTP Response
type ResponseReadWriter interface {
	WriteResponse(contentType string, status int, data interface{}, w http.ResponseWriter)
	ReadResponse(contentType string, data interface{}, response *http.Response) error
}

type payloadImpl struct {
}

//NewPayload returns a new instance of Payload
func NewPayload() Payload {
	return &payloadImpl{}
}

//WriteResponse used to send some payload over the wire as HTTP response
//WriteResponse accepts Content-Type (json|yaml|xml), http status code and instance of any struct that represents
//payload you want to send over the wire via HTTP Response
// Example:
// type SomeModel struct {
//    //fields
//}
//
//func getHandler() http.HandlerFunc {
//    return func (w http.ResponseWriter, r *http.Request) {
//        var model SomeModel
//
//        //retrieve data from DB
//
//        //Setup your model
//
//        //Write to response
//        p := ws.NewPayload()
//        err := p.WriteResponse(ws.ContentTypeJSON, http.StatusOK, &mode, w)
//    }
//}
func (p *payloadImpl) WriteResponse(contentType string, status int, data interface{}, w http.ResponseWriter) {
	fmt.Printf("HTTP Status : %v Content-Type: %s", status, contentType)

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

//ReadResponse is used to read data from HTTP response and return an instance of struct that represents the data
//received via HTTP response.
//Example:
// type SomeModel struct {
//    //fields
//}
//
//func postHandler() http.HandlerFunc {
//    return func (w http.ResponseWriter, r *http.Request) {
//        var model SomeModel
//
//        //setup request
//
//        //setup http client
//
//        //make the call
//        resp, err := client.Do(req)
//
//        //Handle error if required
//
//        //read response
//        p := ws.NewPayload()
//        err := p.ReadResponse(ws.ContentTypeJSON, &model, resp)
//        //Handle error if required
//    }
//}
func (p *payloadImpl) ReadResponse(contentType string, data interface{}, response *http.Response) error {
	fmt.Println("Reading response body")
	defer response.Body.Close()

	payload, err := ioutil.ReadAll(response.Body)
	if err != nil {
		fmt.Printf("Error reading response from GitHub API call %v", err)
		return err
	}
	switch contentType {
	case ContentTypeJSON:
		if err = json.Unmarshal(payload, &data); err != nil {
			fmt.Printf("Error unmarshalling payload")
			return err
		}
		break
	case ContentTypeXML:
		if err = xml.Unmarshal(payload, &data); err != nil {
			fmt.Printf("Error unmarshalling payload")
			return err
		}
		break
	case ContentTypeYaml:
		if err = yaml.Unmarshal(payload, &data); err != nil {
			fmt.Printf("Error unmarshalling payload")
			return err
		}
		break
	default:
		fmt.Printf("Unknown content type: %s", contentType)
	}

	return nil
}

//ReadRequest reads payloaded posted via HTTP request and returns instance of struct that represents the payload
//Example:
// type SomeModel struct {
//    //fields
//}
//
//func postHandler() http.HandlerFunc {
//    return func (w http.ResponseWriter, r *http.Request) {
//        var model SomeModel
//
//        p := ws.NewPayload()
//        err := p.ReadRequest(ws.ContentTypeJSON, &model, r)
//
//        //Handle error if required
//    }
//}
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
