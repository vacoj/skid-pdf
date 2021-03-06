package main

import (
	"bytes"
	"fmt"
	"log"
	"os"
	"os/exec"
	"path"
)

const wkhtmltopdfCmd = "wkhtmltopdf"

// WkOrientationLandscape - if passed, sets orientation to landscape
var WkOrientationLandscape = []string{"-O", "Landscape"}

// WkHeader - prefixes header key pairs
var WkHeader = []string{"--custom-header"}

// WkPostParam - prefixes post key pairs
var WkPostParam = []string{"--post"}

//WkGrayscale - If passed, created PDF will be grayscale
var WkGrayscale = []string{"-g"}

// GenerateWKPDF creates PDF from target URL.  First argument is the URL to
// be converted, other arguments should be pulled from the constants above.
// Additionally, arbitrary params can be passed as long as the string is
// buffered by spaces on both sides.  The result will look something
// like this:
// $> wkhtmltopdf http://someurl/something.html - --your-arguments
func generateWKPDF(targetURL string, params []string) []byte {
	var result bytes.Buffer
	wkCommand := []string{}

	for _, param := range params {
		wkCommand = append(wkCommand, param)
	}

	wkCommand = append(wkCommand, targetURL)
	wkCommand = append(wkCommand, "-")

	cmd := exec.Command(wkhtmltopdfCmd, wkCommand...)

	// for testing
	fmt.Println(wkhtmltopdfCmd, wkCommand)

	cmd.Stdout = &result
	err := cmd.Run()

	if err != nil {
		log.Print(err)
	}
	return result.Bytes()
}

func generateFromPDFRequest(p *pdfRequest) []byte {
	if p.Data != "" {
		p.URL += p.Data
	}
	extraParams := []string{}
	if p.Grayscale {
		extraParams = append(extraParams, WkGrayscale...)
	}
	if p.Landscape {
		extraParams = append(extraParams, WkOrientationLandscape...)
	}

	if len(p.PostParams) > 0 {
		for k, v := range p.PostParams {
			extraParams = append(extraParams, WkPostParam...)
			extraParams = append(extraParams, k, v)
		}
	}
	if len(p.Headers) > 0 {
		for k, v := range p.Headers {
			extraParams = append(extraParams, WkHeader...)
			extraParams = append(extraParams, k, v)
		}
	}

	return generateWKPDF(p.URL, extraParams)
}

func hookForAMQP(r *pdfRequest) {
	pdfResult := generateFromPDFRequest(r)
	writer, err := os.Create(path.Join(r.TargetFileDest, r.TargetFileName))
	if err != nil {
		log.Println(err)
	}
	writer.Write(pdfResult)
}
