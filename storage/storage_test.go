package storage

import (
	"errors"
	"os"
	"path"
	"testing"

	"github.com/stretchr/testify/assert"
)

type MockArchive struct {
	extract func(string) ([]string, error)
}

func (m *MockArchive) Extract(target string) ([]string, error) {
	return m.extract(target)
}

func (m *MockArchive) Close() error {
	return nil
}

type MockStarter struct {
	err error
}

func (m *MockStarter) Start() error {
	return m.err
}

func TestNewStorage(t *testing.T) {
	storage := NewStorage(
		func(filename string) (Archive, error) {
			return nil, nil
		},
		func(name string, arg ...string) Starter {
			return nil
		},
		"/tmp",
	)
	assert.Equal(t, "/tmp", storage.storageDir)
}

func TestHasRelease(t *testing.T) {
	tmpDir := t.TempDir()

	release := "foo"
	storage := &Storage{storageDir: tmpDir}

	assert.False(t, storage.HasRelease(release))

	if os.Mkdir(path.Join(tmpDir, release), 0700) != nil {
		t.Fatal("Failed to create release dir")
	}
	assert.True(t, storage.HasRelease(release))
}

func TestExtractRelease(t *testing.T) {
	tmpDir := t.TempDir()

	var newArchiveErr error = nil
	archive := &MockArchive{}
	storage := &Storage{
		storageDir: tmpDir,
		newArchive: func(filename string) (Archive, error) {
			return archive, newArchiveErr
		},
	}
	filename := path.Join(tmpDir, "foo.7z")
	file, err := os.Create(filename)
	if err != nil {
		t.Fatalf("Failed to create dummy archive, %v", err)
	}
	file.Close()

	newArchiveErr = errors.New("open error")
	err = storage.ExtractRelease(filename)
	assert.NotNil(t, err)
	assert.EqualError(t, err, "open error")

	newArchiveErr = nil
	archive.extract = func(target string) ([]string, error) {
		return nil, errors.New("extract error")
	}
	err = storage.ExtractRelease(filename)
	assert.NotNil(t, err)
	assert.EqualError(t, err, "extract error")

	archive.extract = func(target string) ([]string, error) {
		assert.Equal(t, path.Join(tmpDir, "foo"), target)
		return nil, nil
	}
	assert.Nil(t, storage.ExtractRelease(filename))
	assert.NoFileExists(t, filename)
}

func TestCleanupReleases(t *testing.T) {
	tmpDir := t.TempDir()

	dirs := []string{"foo", "bar", "baz"}
	for _, dir := range dirs {
		if err := os.Mkdir(path.Join(tmpDir, dir), 0700); err != nil {
			t.Fatalf("Failed to create dummy release dir, %v", err)
		}
	}
	file, err := os.Create(path.Join(tmpDir, "bar.txt"))
	if err != nil {
		t.Fatalf("Failed to create random file, %v", err)
	}
	file.Close()

	storage := &Storage{storageDir: tmpDir}
	storage.CleanupReleases("bar")
	assert.DirExists(t, path.Join(tmpDir, "bar"))
	assert.NoDirExists(t, path.Join(tmpDir, "foo"))
	assert.NoDirExists(t, path.Join(tmpDir, "baz"))
	assert.FileExists(t, path.Join(tmpDir, "bar.txt"))
}

func TestRun(t *testing.T) {
	starter := &MockStarter{}
	storage := &Storage{
		storageDir: "/tmp",
		command: func(name string, arg ...string) Starter {
			assert.Equal(t, "/tmp/simc-905-01-win64-937b901/simc-905-01-win64/SimulationCraft.exe", name)
			return starter
		},
	}

	assert.Nil(t, storage.Run("simc-905-01-win64-937b901"))

	starter.err = errors.New("dummy error")
	err := storage.Run("simc-905-01-win64-937b901")
	assert.NotNil(t, err)
	assert.EqualError(t, err, "dummy error")
}
