package scrapper

import (
	"fmt"
	"io"
	"io/ioutil"
	"net/http"
	"os"
	"regexp"
)

type Link struct {
	Description string
	Value       string
}

func List(search *regexp.Regexp, url string) ([]Link, error) {
	links := make([]Link, 0)
	res, err := http.Get(url)
	if err != nil {
		return nil, fmt.Errorf("listaLinks: %w", err)
	}
	defer res.Body.Close()
	b, err := ioutil.ReadAll(res.Body)
	if err != nil {
		return nil, fmt.Errorf("listaLinks: %w", err)
	}
	doc := string(b)
	matches := search.FindAllStringSubmatch(doc, -1)
	for _, m := range matches {
		l := Link{}
		l.Description = m[0]
		l.Value = m[1]
		links = append(links, l)
	}
	return links, nil
}

func Download(path, file string) error {
	res, err := http.Get(path)
	if err != nil {
		return fmt.Errorf("download: %w", err)
	}
	defer res.Body.Close()
	f, err := os.Create(file)
	if err != nil {
		return fmt.Errorf("download: %w", err)
	}
	defer f.Close()
	if _, err := io.Copy(f, res.Body); err != nil {
		return fmt.Errorf("download: %w", err)
	}
	return nil
}

func MakeDir(path string) error {
	_, err := os.Stat(path)
	if err != nil && !os.IsExist(err) {
		if err := os.Mkdir(path, 0744); err != nil && err != os.ErrExist {
			return fmt.Errorf("makeDir: %w", err)
		}
	}
	return nil
}
