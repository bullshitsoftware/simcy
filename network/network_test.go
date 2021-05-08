package network

import (
	"bytes"
	"errors"
	"io/ioutil"
	"net/http"
	"os"
	"strings"
	"testing"

	"github.com/stretchr/testify/assert"
)

const downloadsPageFixture = "downloads_page_fixture.html"

type MockHttpClient struct {
	get func(url string) (*http.Response, error)
}

func (c *MockHttpClient) Get(url string) (*http.Response, error) {
	return c.get(url)
}

func TestNewNetwork(t *testing.T) {
	const downloadsUrl = "http://example.com"
	client := &MockHttpClient{}
	network := NewNetwork(client, downloadsUrl)

	assert.Equal(t, client, network.client)
	assert.Equal(t, downloadsUrl, downloadsUrl)
}

func TestGetLastReleaseName(t *testing.T) {
	const downloadsUrl = "http://example.com"

	client := &MockHttpClient{}
	network := &Network{client, downloadsUrl}

	{
		client.get = func(url string) (*http.Response, error) {
			assert.Equal(t, downloadsUrl, url)
			body, _ := os.Open(downloadsPageFixture)
			return &http.Response{Body: body}, nil
		}
		release, err := network.GetLastReleaseName()
		assert.NotEqual(t, "", release)
		assert.Nil(t, err)
		assert.Equal(t, "simc-905-01-win64-937b901", release)
	}

	{
		client.get = func(url string) (*http.Response, error) {
			assert.Equal(t, downloadsUrl, url)
			return &http.Response{Body: ioutil.NopCloser(strings.NewReader(""))}, nil
		}
		release, err := network.GetLastReleaseName()
		assert.Equal(t, "", release)
		assert.NotNil(t, err)
		assert.EqualError(t, err, "failed to find any release data")
	}
}

func TestDownloadRelease(t *testing.T) {
	const downloadsUrl = "http://example.com"
	const release = "foo"

	client := &MockHttpClient{}
	network := &Network{client, downloadsUrl}
	tmpDir := t.TempDir()

	{
		data := []byte{1, 2, 3}
		client.get = func(url string) (*http.Response, error) {
			assert.Equal(t, downloadsUrl+"/foo.7z", url)
			body := ioutil.NopCloser(bytes.NewReader(data))
			return &http.Response{Body: body}, nil
		}
		filename, err := network.DownloadRelease(release, tmpDir)
		assert.NotEqual(t, "", filename)
		assert.Nil(t, err)
		downloadedBytes, _ := ioutil.ReadFile(filename)
		assert.True(t, bytes.Equal(data, downloadedBytes))
	}

	{
		client.get = func(url string) (*http.Response, error) {
			assert.Equal(t, downloadsUrl+"/foo.7z", url)
			return &http.Response{}, errors.New("random error")
		}
		filename, err := network.DownloadRelease(release, tmpDir)
		assert.Equal(t, "", filename)
		assert.NotNil(t, err)
		assert.EqualError(t, err, "random error")
	}
}
