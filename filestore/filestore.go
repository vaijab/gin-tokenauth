package filestore

import (
	"io/ioutil"
	"os"
	"path/filepath"
	"sync"

	"github.com/fsnotify/fsnotify"
	"gopkg.in/yaml.v2"

	"github.com/orendo/gin-tokenauth"
)

// Store is a file based token store.
type Store struct {
	file string

	mu     sync.Mutex
	Tokens []tokenauth.Token `yaml:"tokens"`
}

// IsTokenValid returns a bool to indicate validity of the token.
func (s *Store) IsTokenValid(token string) bool {
	for _, t := range s.Tokens {
		if token == t.Token && !t.IsDisabled {
			return true
		}
	}
	return false
}

// loadTokens loads tokens from a yaml file.
func (s *Store) loadTokens() error {
	data, err := ioutil.ReadFile(s.file)
	if err != nil {
		// non-existing file is not an error, the file can be created after main starts.
		if os.IsNotExist(err) {
			return nil
		}
		return err
	}

	s.mu.Lock()
	defer s.mu.Unlock()
	return yaml.Unmarshal(data, &s)
}

func (s *Store) clearTokens() {
	s.mu.Lock()
	s.Tokens = nil
	s.mu.Unlock()
}

func (s *Store) watchTokensFile() error {
	absPath, err := filepath.Abs(s.file)
	if err != nil {
		return err
	}

	watcher, err := fsnotify.NewWatcher()
	if err != nil {
		return err
	}

	go func() {
		for {
			select {
			case event := <-watcher.Events:
				if event.Name == absPath {
					if event.Op&fsnotify.Write == fsnotify.Write {
						if err := s.loadTokens(); err != nil {
							// TODO log errors
						}
					}
					if event.Op&fsnotify.Remove == fsnotify.Remove {
						s.clearTokens()
					}
				}
			case err = <-watcher.Errors:
				// TODO log errors
			}
		}
	}()

	// We want to watch file dir instead of a file, because watching
	// non-existing file is not possible. Also, we want to pick up on token
	// file deletion and creation.
	return watcher.Add(filepath.Dir(absPath))
}

// New loads tokens and initializes a file watcher on tokens file and returns a TokenStore.
func New(f string) (*Store, error) {
	store := &Store{
		file: f,
	}

	// Attempt to load tokens file.
	if err := store.loadTokens(); err != nil {
		// TODO log errors - initial load is allowed to fail
	}

	// Create a tokens file watcher.
	if err := store.watchTokensFile(); err != nil {
		return store, err
	}

	return store, nil
}