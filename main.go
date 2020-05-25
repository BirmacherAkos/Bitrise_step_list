package main

import (
	"encoding/json"
	"fmt"
	"io/ioutil"
	"net/http"
	"os"

	"github.com/bitrise-io/go-utils/log"
)

// BasicData ...
type BasicData struct {
	FormatVersion        string `json:"format_version"`
	GeneratedAtTimestamp int    `json:"generated_at_timestamp"`
	SteplibSource        string `json:"steplib_source"`
	DownloadLocations    []struct {
		Type string `json:"type"`
		Src  string `json:"src"`
	} `json:"download_locations"`
	AssetsDownloadBaseURI string          `json:"assets_download_base_uri"`
	Steps                 map[string]Step `json:"steps"`
}

// Step ...
type Step struct {
	Info struct {
		DeprecateNotes string `json:"deprecate_notes"`
		AssetUrls      struct {
			IconSvg string `json:"icon.svg"`
		} `json:"asset_urls"`
	} `json:"info"`
	LatestVersionNumber string                            `json:"latest_version_number"`
	Versions            map[string]map[string]interface{} `json:"versions"`
	StepID              string
}

func fetchSteps() (BasicData, error) {
	response, err := http.Get("https://bitrise-steplib-collection.s3.amazonaws.com/spec.json")
	if err != nil {
		return BasicData{}, err
	}

	var data BasicData

	b, err := ioutil.ReadAll(response.Body)
	if err != nil {
		return BasicData{}, err
	}

	if err := json.Unmarshal(b, &data); err != nil {
		return BasicData{}, err
	}

	return data, nil
}

func logPretty(v interface{}) string {
	b, err := json.MarshalIndent(v, "", "  ")
	if err != nil {
		fmt.Println("error:", err)
	}

	return fmt.Sprintf("%+v\n", string(b))
}

func failf(format string, v ...interface{}) {
	log.Errorf(format, v...)
	os.Exit(1)
}

func contains(arr []string, str string) bool {
	for _, a := range arr {
		if a == str {
			return true
		}
	}
	return false
}

func main() {
	log.Infof("Fetching step list")
	d, err := fetchSteps()
	if err != nil {
		failf("Failed to fetch the step list from the server, error: %s", err)
	}
	log.Successf("Done")
	fmt.Println()

	log.Infof("List deprecated steps")
	var deprecatedSteps []string
	for stepID, step := range d.Steps {
		step.StepID = stepID
		if step.Info.DeprecateNotes != "" {
			// log.Printf("%s - note: : ", stepID, step.Info.DeprecateNotes)
			deprecatedSteps = append(deprecatedSteps, stepID)
		}

	}
	log.Printf(logPretty(deprecatedSteps))
	fmt.Println()

}
