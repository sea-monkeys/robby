package main

import (
	_ "embed"
	"encoding/json"
	"github.com/extism/go-pdk"
)

//go:embed rsrc.calculator.md
var rsrcCalculator string

//go:embed rsrc.addition.md
var rsrcAddition string

//go:embed rsrc.subtraction.md
var rsrcSubtraction string

//go:embed rsrc.multiplication.md
var rsrcMultiplication string

//go:embed rsrc.division.md
var rsrcDivision string

// -------------------------------------------------
//	Resources
// -------------------------------------------------
type Resource struct {
	URI         string `json:"uri"`
	Name        string `json:"name"`
	Description string `json:"description,omitempty"`
	MimeType    string `json:"mimeType,omitempty"`
	Text        string `json:"text,omitempty"`
	Blob        string `json:"blob,omitempty"`
}

//go:export resources_information
func ResourcesInformation() {
	// Define the resources information
	resources := []Resource{
		{
			URI:         "info:///calculator",
			Name:        "calc_rsrc_about",
			Description: "Information about the calculator plugin",
			MimeType:    "text/markdown",
		},
		{
			URI:         "info:///add",
			Name:        "calc_rsrc_add",
			Description: "Help information for the addition tool",
			MimeType:    "text/markdown",
		},
		{
			URI:         "info:///subtract",
			Name:        "calc_rsrc_subtract",
			Description: "Help information for the subtraction tool",
			MimeType:    "text/markdown",
		},
		{
			URI:         "info:///multiply",
			Name:        "calc_rsrc_multiply",
			Description: "Help information for the multiplication tool",
			MimeType:    "text/markdown",
		},
		{
			URI:         "info:///divide",
			Name:        "calc_rsrc_divide",
			Description: "Help information for the division tool",
			MimeType:    "text/markdown",
		},
	}
	jsonData, _ := json.Marshal(resources)
	pdk.OutputString(string(jsonData))
}

//go:export calc_rsrc_about
func CalcRsrcAbout() {
	aboutInfo := Resource{
		URI:         "info:///calculator",
		Name:        "calc_rsrc_about",
		Description: "Information about this calculator plugin",
		MimeType:    "text/markdown",
		Text:        rsrcCalculator,
	}
	jsonData, _ := json.Marshal(aboutInfo)
	pdk.OutputString(string(jsonData))
}

//go:export calc_rsrc_add
func CalcRsrcAdd() {
	additionInfo := Resource{
		URI:         "info:///add",
		Name:        "calc_rsrc_add",
		Description: "Help information for the addition tool",
		MimeType:    "text/markdown",
		Text:        rsrcAddition,
	}
	jsonData, _ := json.Marshal(additionInfo)
	pdk.OutputString(string(jsonData))
}
//go:export calc_rsrc_subtract
func CalcRsrcSubtract() {
	subtractionInfo := Resource{
		URI:         "info:///subtract",
		Name:        "calc_rsrc_subtract",
		Description: "Help information for the subtraction tool",
		MimeType:    "text/markdown",
		Text:        rsrcSubtraction,
	}
	jsonData, _ := json.Marshal(subtractionInfo)
	pdk.OutputString(string(jsonData))
}
//go:export calc_rsrc_multiply
func CalcRsrcMultiply() {
	multiplicationInfo := Resource{
		URI:         "info:///multiply",
		Name:        "calc_rsrc_multiply",
		Description: "Help information for the multiplication tool",
		MimeType:    "text/markdown",
		Text:        rsrcMultiplication,
	}
	jsonData, _ := json.Marshal(multiplicationInfo)
	pdk.OutputString(string(jsonData))
}
//go:export calc_rsrc_divide
func CalcRsrcDivide() {
	divisionInfo := Resource{
		URI:         "info:///divide",
		Name:        "calc_rsrc_divide",
		Description: "Help information for the division tool",
		MimeType:    "text/markdown",
		Text:        rsrcDivision,
	}
	jsonData, _ := json.Marshal(divisionInfo)
	pdk.OutputString(string(jsonData))
}

