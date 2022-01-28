package soapHandler

import (
	"bytes"
	"crypto/tls"
	"encoding/xml"
	"fmt"
	"io/ioutil"
	"net/http"
	"text/template"
)

type Request struct {
	//Values are set in below fields as per the request
	Codigo string
}

type Response struct {
	XMLName    xml.Name `xml:"http://schemas.xmlsoap.org/soap/envelope/ Envelope"`
	SoapHeader *SOAPHeaderResponse
	SoapBody   *SOAPBodyResponse
}

type SOAPHeaderResponse struct {
	XMLName xml.Name `xml:"Header"`
}

type SOAPBodyResponse struct {
	XMLName      xml.Name `xml:"Body"`
	Resp         *GetStockResponseBody
	FaultDetails *Fault
}

type Fault struct {
	XMLName     xml.Name `xml:"Fault"`
	Faultcode   string   `xml:"faultcode"`
	Faultstring string   `xml:"faultstring"`
}

type GetStockResponseBody struct {
	XMLName xml.Name `xml:"ZSDRFC_SKN_GET_STOCKResponse"`
	Result  *Return
	Stock   *Stock
}

type Return struct {
	XMLName    xml.Name `xml:"ET_RETURN"`
	ResultItem *ReturnItem
}

type ReturnItem struct {
	XMLName       xml.Name `xml:"item"`
	Type          string   `xml:"TYPE"`
	Code          string   `xml:"CODE"`
	ResultMessage string   `xml:"MESSAGE"`
	Log_No        string   `xml:"LOG_NO"`
	LOG_MSG_NO    string   `xml:"LOG_MSG_NO"`
	MESSAGE_V1    string   `xml:"MESSAGE_V1"`
	MESSAGE_V2    string   `xml:"MESSAGE_V2"`
	MESSAGE_V3    string   `xml:"MESSAGE_V3"`
	MESSAGE_V4    string   `xml:"MESSAGE_V4"`
}

type Stock struct {
	XMLName   xml.Name `xml:"ET_STOCK"`
	StockItem *StockItem
}

type StockItem struct {
	XMLName     xml.Name `xml:"item"`
	Description string   `xml:"ZSD_DCORTA"`
	Quantity    string   `xml:"ZSD_QSTUCO"`
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

	r, err := http.NewRequest(http.MethodPost, "https://servicioswebdex.alicorp.com.pe/nd1/sap/bc/srt/rfc/sap/zsdrfc_skn_get_stock/300/zsdrfc_skn_get_stock/zsdrfc_skn_get_stock", bytes.NewBuffer(doc.Bytes()))
	r.SetBasicAuth("nrodriguezv", "mikaela2013")
	r.Header.Add("Content-Type", "text/xml;charset=UTF-8")
	r.Header.Add("Accept-Encoding", "gzip,deflate")

	if err != nil {
		fmt.Printf("Error making a request. %s ", err.Error())
		return nil, err
	}

	return r, nil
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
	tr := &http.Transport{
		TLSClientConfig: &tls.Config{InsecureSkipVerify: true},
	}

	client := &http.Client{Transport: tr}
	resp, err := client.Do(req)

	if err != nil {
		return nil, err
	}

	body, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		return nil, err
	}

	fmt.Println(string(body))

	defer resp.Body.Close()

	r := &Response{}
	err = xml.Unmarshal(body, &r)

	fmt.Println()

	if err != nil {
		return nil, err
	}

	//if r.SoapBody.Resp.Status != "200" {
	//	return nil, err
	//}

	return r, nil
}
