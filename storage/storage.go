package storage

import (
	"os"
	"path"
	"path/filepath"
	"strings"
)

type NewArchive func(string) (Archive, error)
type Archive interface {
	Extract(string) ([]string, error)
	Close() error
}

type Command func(name string, arg ...string) Starter
type Starter interface {
	Start() error
}

type Storage struct {
	newArchive NewArchive
	command    Command
	storageDir string
}

func NewStorage(newArchive NewArchive, command Command, storageDir string) *Storage {
	return &Storage{
		newArchive,
		command,
		storageDir,
	}
}

func (s *Storage) HasRelease(release string) bool {
	entries, _ := os.ReadDir(s.storageDir)
	for _, entry := range entries {
		if entry.IsDir() && release == entry.Name() {
			return true
		}
	}
	return false
}

func (s *Storage) ExtractRelease(filename string) error {
	archive, err := s.newArchive(filename)
	if err != nil {
		return err
	}
	_, err = archive.Extract(strings.TrimSuffix(filename, filepath.Ext(filename)))
	if err != nil {
		return err
	}
	archive.Close()

	return os.Remove(filename)
}

func (s *Storage) CleanupReleases(keepRelease string) error {
	entries, err := os.ReadDir(s.storageDir)
	if err != nil {
		return err
	}

	for _, entry := range entries {
		if !entry.IsDir() {
			continue
		}

		if entry.Name() != keepRelease {
			os.RemoveAll(path.Join(s.storageDir, entry.Name()))
		}
	}

	return nil
}

func (s *Storage) Run(release string) error {
	executable := path.Join(s.storageDir, release, release[:len(release)-8], "SimulationCraft.exe")
	cmd := s.command(executable)

	return cmd.Start()
}
