package main

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

type MockNetwork struct {
	getLastReleaseName func() (string, error)
	downloadRelease    func(release, destDir string) (string, error)
}

func (m *MockNetwork) GetLastReleaseName() (string, error) {
	return m.getLastReleaseName()
}

func (m *MockNetwork) DownloadRelease(release, destDir string) (string, error) {
	return m.downloadRelease(release, destDir)
}

type MockStorage struct {
	hasRelease      func(release string) bool
	extractRelease  func(filename string) error
	cleanupReleases func(keepRelease string) error
	run             func(release string) error
}

func (m *MockStorage) HasRelease(release string) bool {
	return m.hasRelease(release)
}

func (m *MockStorage) ExtractRelease(filename string) error {
	return m.extractRelease(filename)
}

func (m *MockStorage) CleanupReleases(keepRelease string) error {
	return m.cleanupReleases(keepRelease)
}

func (m *MockStorage) Run(release string) error {
	return m.run(release)
}

func TestRun(t *testing.T) {
	tmpDir := "/tmp"
	network := &MockNetwork{
		getLastReleaseName: func() (string, error) {
			return "foo", nil
		},
		downloadRelease: func(release, destDir string) (string, error) {
			assert.Equal(t, "foo", release)
			assert.Equal(t, tmpDir, destDir)

			return "/tmp/foo.7z", nil
		},
	}
	storage := &MockStorage{
		hasRelease: func(release string) bool {
			assert.Equal(t, "foo", release)
			return true
		},
		extractRelease: func(filename string) error {
			assert.Equal(t, "/tmp/foo.7z", filename)
			return nil
		},
		cleanupReleases: func(keepRelease string) error {
			assert.Equal(t, "foo", keepRelease)
			return nil
		},
		run: func(release string) error {
			assert.Equal(t, "foo", release)
			return nil
		},
	}
	assert.Equal(t, 0, run(network, storage))
}
