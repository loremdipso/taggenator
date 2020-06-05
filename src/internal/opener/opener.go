package opener

import (
	"errors"
	"fmt"
	"internal/data"
	"os"
	"os/exec"
	"path/filepath"
	"strings"

	"github.com/kballard/go-shellquote"
)

type Opener struct {
	Config data.OpenerConfig

	lastExtension string
	lastPath      string
	lastConfig    *data.OpenerFileTypeConfig
	lastPID       int
}

func New(config data.OpenerConfig) *Opener {
	return &Opener{Config: config}
}

func (self *Opener) Close() {
	// TODO: anything?
	fmt.Println("Closing")
	if self.lastConfig != nil && self.lastConfig.Close != "" {
		executeConfig(self.lastConfig.Close, self.lastPath, true)
	}
}

const (
	default_key = "default"
	video_key   = "video"
	comic_key   = "comic"
	music_key   = "music"
)

func (self *Opener) Open(path string) bool {
	updateTempFile(path)

	extension := strings.ToLower(filepath.Ext(path))
	return self.
		tryHandleByExtension(path, extension) || self.
		tryHandleVideo(path, extension) || self.
		tryHandleComic(path, extension) || self.
		tryHandleMusic(path, extension) || self.
		tryHandleDefault(path, extension)
}

// TODO: make more generic
func updateTempFile(path string) {
	f, err := os.Create("/tmp/trim_filename.txt")
	defer f.Close()

	if err != nil {
		fmt.Println(err)
		return
	}

	abspath, _ := filepath.Abs(path)
	fmt.Fprintln(f, abspath)
}

func (self *Opener) getConfig(extension string) *data.OpenerFileTypeConfig {
	config, ok := self.Config[extension]
	if ok {
		return config
	}
	return nil
}

func (self *Opener) tryHandleByExtension(path string, extension string) bool {
	return self.handleConfig(self.getConfig(extension), path)
}

func (self *Opener) tryHandleVideo(path string, extension string) bool {
	switch extension {
	// TODO: is this exhaustive
	case ".flv", ".mp4", ".m4a", ".wmv":
		return self.tryHandleByExtension(path, video_key)
	}
	return false
}

func (self *Opener) tryHandleComic(path string, extension string) bool {
	switch extension {
	// TODO: is this exhaustive
	case ".cbz":
		return self.tryHandleByExtension(path, comic_key)
	}
	return false
}

func (self *Opener) tryHandleMusic(path string, extension string) bool {
	switch extension {
	// TODO: is this exhaustive
	case ".mp3", ".wav":
		return self.tryHandleByExtension(path, music_key)
	}
	return false
}

func (self *Opener) tryHandleDefault(path string, extension string) bool {
	return self.tryHandleByExtension(path, default_key)
}

func (self *Opener) handleConfig(config *data.OpenerFileTypeConfig, path string) bool {
	if config == nil {
		return false
	}

	if self.lastConfig != nil {
		if self.lastConfig.Update != "" && self.lastConfig == config {
			self.lastConfig = config
			self.lastPID, _ = executeConfig(config.Update, path, false)
			self.lastPath = path
			return true
		} else if self.lastConfig.Close != "" {
			executeConfig(config.Close, path, true)
		}
	}

	if config.Open != "" {
		// NOTE: lastPID might be -1 if there was an error
		self.lastConfig = config
		self.lastPID, _ = executeConfig(config.Open, path, false)
		self.lastPath = path
	} else {
		self.lastConfig = nil
	}

	return true
}

func executeConfig(estr string, path string, doWait bool) (int, error) {
	// fmt.Printf(estr, path)
	// cmd := exec.Command("ssh", "-i", keyFile, "-o", "ExitOnForwardFailure yes", "-fqnNTL", fmt.Sprintf("%d:127.0.0.1:%d", port, port), fmt.Sprintf("%s@%s", serverUser, serverIP))
	// path = "./" + path
	// wd, _ := os.Getwd()
	// path = filepath.Join(wd, path)
	cstr, _ := truncatingSprintf(estr, path) // TODO: is this unsafe?
	words, err := shellquote.Split(cstr)

	// TODO: handle err
	if err != nil {
		return -1, err
	}

	cmd := exec.Command(words[0], words[1:]...)

	cmd.Stdout = nil
	cmd.Stderr = nil

	// cmd.D
	err = cmd.Start()
	if err != nil {
		fmt.Println(err)
	}
	pid := cmd.Process.Pid
	if doWait {
		cmd.Wait()
	}
	return pid, nil
}

func truncatingSprintf(str string, args ...interface{}) (string, error) {
	n := strings.Count(str, "%s")
	if n > len(args) {
		return "", errors.New("Unexpected string:" + str)
	}
	return fmt.Sprintf(str, args[:n]...), nil
}
