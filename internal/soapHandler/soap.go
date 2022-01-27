package soaphandler

import (
	"bytes"
	"encoding/xml"
	"fmt"
	"io/ioutil"
	"net/http"
	"text/template"
)

type Request struct {
	//Values are set in below fields as per the request
	Blz string
}

func generateSOAPRequest(req *Request) (*http.Request, error) {
	// Using the var getTemplate to construct request
	template, err := template.New("InputRequest").Parse(getTemplate)
	if err != nil {
		fmt.Printf("Error while marshling object. %s ", err.Error())
		return nil, err
	}

	doc := &bytes.Buffer{}
	// Replacing the doc from template with actual req values
	err = template.Execute(doc, req)
	if err != nil {
		fmt.Printf("template.Execute error. %s ", err.Error())
		return nil, err
	}

	buffer := &bytes.Buffer{}
	encoder := xml.NewEncoder(buffer)
	err = encoder.Encode(doc.String())
	if err != nil {
		fmt.Printf("encoder.Encode error. %s ", err.Error())
		return nil, err
	}

	r, err := http.NewRequest(http.MethodPost, "http://www.thomas-bayer.com/axis2/services/BLZService", bytes.NewBuffer(doc.Bytes()))
	if err != nil {
		fmt.Printf("Error making a request. %s ", err.Error())
		return nil, err
	}

	return r, nil
}

type Response struct {
	XMLName  xml.Name `xml:"http://schemas.xmlsoap.org/soap/envelope/ Envelope"`
	SoapBody *SOAPBodyResponse
}

type SOAPBodyResponse struct {
	XMLName      xml.Name `xml:"Body"`
	Resp         *GetBankResponseBody
	FaultDetails *Fault
}

type Fault struct {
	XMLName     xml.Name `xml:"Fault"`
	Faultcode   string   `xml:"faultcode"`
	Faultstring string   `xml:"faultstring"`
}

type GetBankResponseBody struct {
	XMLName  xml.Name `xml:"getBankResponse"`
	Response *Details
}

type Details struct {
	XMLName     xml.Name `xml:"details"`
	Bezeichnung string   `xml:"bezeichnung"`
	Bic         string   `xml:"bic"`
	Ort         string   `xml:"ort"`
	Plz         string   `xml:"plz"`
}

func CallSOAPClientSteps(req *Request) (*Response, error) {

	httpReq, err := generateSOAPRequest(req)
	if err != nil {
		fmt.Println("Some problem occurred in request generation")
	}

	response, err := soapCall(httpReq)
	if err != nil {
		fmt.Println("Problem occurred in making a SOAP call")
	}

	return response, err

}

func soapCall(req *http.Request) (*Response, error) {
	client := &http.Client{}
	resp, err := client.Do(req)

	if err != nil {
		return nil, err
	}

	body, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		return nil, err
	}
	defer resp.Body.Close()

	r := &Response{}
	err = xml.Unmarshal(body, &r)

	if err != nil {
		return nil, err
	}

	//if r.SoapBody.Resp.Status != "200" {
	//	return nil, err
	//}

	return r, nil
}
