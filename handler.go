package main

import (
	"encoding/csv"
	"encoding/json"
	"net/http"
	"os"
	"strings"

	"gopkg.in/yaml.v3"
)

type Redirect struct {
	Path string `yaml:"path" json:"path"`
	URL  string `yaml:"url" json:"url"`
}

func RedirectHandler(mode string, data string, raw bool, fallback http.Handler) http.HandlerFunc {
	paths, err := detectModeAndUnmarshall(mode, data, raw)
	if err != nil {
		panic(err)
	}
	pathsToUrls := mapFromRedirectSlice(paths)
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		requestedUrl := pathsToUrls[r.URL.Path]
		if requestedUrl != "" {
			http.Redirect(w, r, requestedUrl, http.StatusFound)
			return
		}
		fallback.ServeHTTP(w, r)

	})
}

func readRedirectData(path string) (*os.File, error) {
	var bytes []byte
	file, err := os.Open(path)
	if err != nil {
		return nil, err
	}
	file.Read(bytes)
	return file, nil
}

func detectModeAndUnmarshall(mode string, data string, raw bool) (*[]Redirect, error) {

	switch mode {
	case "yaml":
		return UnmarshallFromYAML(data, raw)

	case "csv":
		return UnmarshalFromCSV(data, raw)

	default:
		return UnmarshalFromJSON(data, raw)
	}

}

func UnmarshallFromYAML(redirectData string, raw bool) (*[]Redirect, error) {

	var urls []Redirect
	if raw {
		err := yaml.Unmarshal([]byte(redirectData), &urls)
		return &urls, err

	}
	data, err := readRedirectData(redirectData)

	if err != nil {
		panic(err)
	}

	err = yaml.NewDecoder(data).Decode(&urls)

	return &urls, err

}

func UnmarshalFromJSON(redirectData string, raw bool) (*[]Redirect, error) {
	var urls []Redirect
	if raw {
		err := json.Unmarshal([]byte(redirectData), &urls)
		return &urls, err

	}
	data, err := readRedirectData(redirectData)

	if err != nil {
		panic(err)
	}
	err = json.NewDecoder(data).Decode(&urls)

	return &urls, err

}

// TODO improve
func UnmarshalFromCSV(data string, raw bool) (*[]Redirect, error) {
	var urls []Redirect

	if raw {
		lines := strings.Split(data, "\n")

		for _, line := range lines {
			columns := strings.Split(line, ",")
			urls = append(urls, Redirect{
				Path: strings.TrimSpace(columns[0]),
				URL:  strings.Trim(strings.TrimSpace(columns[1]), "\""),
			})
		}

		return &urls, nil
	} else {
		file, err := os.Open(data)

		if err != nil {
			return nil, err
		}
		csvReader := csv.NewReader(file)

		slice, err := csvReader.ReadAll()
		if err != nil {
			return nil, err
		}
		for _, line := range slice {
			urls = append(urls, Redirect{
				Path: strings.TrimSpace(line[0]),
				URL:  strings.Trim(line[1], "\""),
			})
		}
		return &urls, nil

	}

}

func mapFromRedirectSlice(urls *[]Redirect) map[string]string {
	out := make(map[string]string)
	for _, v := range *urls {
		out[v.Path] = v.URL
	}
	return out
}
