package network

import (
	"errors"
	"io"
	"net/http"
	"net/url"
	"os"
	"path"
	"path/filepath"
	"strings"
	"time"

	"github.com/PuerkitoBio/goquery"
)

const (
	platform   = "win64"
	archiveExt = ".7z"
)

type HttpClient interface {
	Get(string) (*http.Response, error)
}

type Network struct {
	client       HttpClient
	downloadsUrl string
}

func NewNetwork(client HttpClient, downloadsUrl string) *Network {
	return &Network{
		client,
		downloadsUrl,
	}
}

func (n *Network) GetLastReleaseName() (string, error) {
	response, err := n.client.Get(n.downloadsUrl)
	if err != nil {
		return "", err
	}

	defer response.Body.Close()
	doc, err := goquery.NewDocumentFromReader(response.Body)
	if err != nil {
		return "", err
	}

	var releaseName = ""
	var releaseDate = time.Time{}
	doc.Find("tr").Each(func(i int, selection *goquery.Selection) {
		tds := selection.Find("td")
		nameNode := tds.Eq(1).Find("a")
		href, exist := nameNode.Attr("href")
		if !exist || !strings.Contains(href, platform) || !strings.HasSuffix(href, archiveExt) {
			return
		}

		dateNode := tds.Eq(2)
		date, err := time.Parse("2006-01-02 15:04", strings.Trim(dateNode.Text(), " "))
		if err != nil {
			return
		}

		if releaseName == "" || releaseDate.Before(date) {
			releaseName = strings.TrimSuffix(href, filepath.Ext(href))
			releaseDate = date
		}
	})

	if releaseName == "" {
		return "", errors.New("failed to find any release data")
	}

	return releaseName, nil
}

func (n *Network) DownloadRelease(release, destDir string) (string, error) {
	releaseUrl, err := url.Parse(n.downloadsUrl)
	if err != nil {
		return "", err
	}
	releaseUrl.Path = path.Join(releaseUrl.Path, release+archiveExt)

	response, err := n.client.Get(releaseUrl.String())
	if err != nil {
		return "", err
	}
	defer response.Body.Close()

	var destFile = path.Join(destDir, release+archiveExt)
	localFile, err := os.Create(destFile)
	if err != nil {
		return "", err
	}
	defer localFile.Close()

	_, err = io.Copy(localFile, response.Body)
	if err != nil {
		os.Remove(destFile)
		return "", err
	}

	return destFile, nil
}
