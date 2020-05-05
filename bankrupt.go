package main

import (
	"bytes"
	"fmt"
	"net/http"
	"strings"

	"github.com/360EntSecGroup-Skylar/excelize"
)

type Cell struct {
	bin                  int
	rnn                  int
	taxpayerOrganization int
	taxpayerName         int
	ownerName            int
	ownerIin             int
	ownerRnn             int
	courtDecision        int
}

type Bankrupt struct {
	bin                  string
	rnn                  string
	taxpayerOrganization string
	taxpayerName         string
	ownerName            string
	ownerIin             string
	ownerRnn             string
	courtDecision        string
}

func (p Bankrupt) toString() string {
	var id string

	if p.bin != "" {
		id = "\"_id\": \"" + p.bin + "\""
	}
	return "{ \"index\": {" + id + "}} \n" +
		"{ \"bin\":\"" + p.bin + "\"" +
		", \"rnn\":\"" + p.rnn + "\"" +
		", \"taxpayer_organization\":\"" + p.taxpayerOrganization + "\"" +
		", \"taxpayer_name\":\"" + p.taxpayerName + "\"" +
		", \"owner_name\":\"" + p.ownerName + "\"" +
		", \"owner_iin\":\"" + p.ownerIin + "\"" +
		", \"owner_rnn\":\"" + p.ownerRnn + "\"" +
		", \"court_decision\":\"" + p.courtDecision + "\"" +
		"}\n"
}

func parseAndSendToES(TaxInfoDescription string, f *excelize.File) error {
	cell := Cell{1, 2, 3, 4, 5,
		6, 7, 8}

	replacer := strings.NewReplacer(
		"\"", "'",
		"\\", "/",
		"\n", "",
		"\n\n", "",
		"\r", "")

	for _, name := range f.GetSheetMap() {
		// Get all the rows in the name
		rows := f.GetRows(name)
		var input strings.Builder
		for i, row := range rows {
			if i < 3 {
				continue
			}
			bankrupt := new(Bankrupt)
			for j, colCell := range row {
				switch j {
				case cell.bin:
					bankrupt.bin = replacer.Replace(colCell)
				case cell.rnn:
					bankrupt.rnn = replacer.Replace(colCell)
				case cell.taxpayerOrganization:
					bankrupt.taxpayerOrganization = replacer.Replace(colCell)
				case cell.taxpayerName:
					bankrupt.taxpayerName = replacer.Replace(colCell)
				case cell.ownerName:
					bankrupt.ownerName = replacer.Replace(colCell)
				case cell.ownerIin:
					bankrupt.ownerIin = replacer.Replace(colCell)
				case cell.ownerRnn:
					bankrupt.ownerRnn = replacer.Replace(colCell)
				case cell.courtDecision:
					bankrupt.courtDecision = replacer.Replace(colCell)
				}
			}
			// if bankrupt.bin != "" {
			input.WriteString(bankrupt.toString())
			// }
			if i%20000 == 0 {
				if errorT := sendPost(TaxInfoDescription, input.String()); errorT != nil {
					return errorT
				}
				input.Reset()
			}
		}
		if input.Len() != 0 {
			if errorT := sendPost(TaxInfoDescription, input.String()); errorT != nil {
				return errorT
			}
		}
	}
	return nil
}

func sendPost(TaxInfoDescription string, query string) error {
	data := []byte(query)
	r := bytes.NewReader(data)
	resp, err := http.Post("http://localhost:9200/bankrupt/companies/_bulk", "application/json", r)
	if err != nil {
		fmt.Println("Could not send the data to elastic search " + TaxInfoDescription)
		fmt.Println(err)
		return err
	}
	fmt.Println(TaxInfoDescription + " " + resp.Status)
	return nil
}
