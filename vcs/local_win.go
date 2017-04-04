package vcs

import (
	"bytes"
	"encoding/json"
	"io"
	"log"
	"os/exec"
	"path/filepath"
	"strings"
)

func init() {
	Register(newLocalWin, "local_win", "localFolderWindows")
}

type LocalWinDriver struct {
	Username string `json:"username"`
	Password string `json:"password"`
	Source string `source`
}

func newLocalWin(b []byte) (Driver, error) {
	var d LocalWinDriver

	if b != nil {
		if err := json.Unmarshal(b, &d); err != nil {
			return nil, err
		}
	}

	return &d, nil
}

func (g *LocalWinDriver) HeadRev(dir string) (string, error) {
	cmd := exec.Command(
		"cat",
		"version.txt")
	cmd.Dir = dir
	r, err := cmd.StdoutPipe()
	if err != nil {
		return "", err
	}
	defer r.Close()

	if err := cmd.Start(); err != nil {
		return "", err
	}

	var buf bytes.Buffer

	if _, err := io.Copy(&buf, r); err != nil {
		return "", err
	}

	return strings.TrimSpace(buf.String()), cmd.Wait()
}

func (g *LocalWinDriver) Pull(dir string) (string, error) {
	cmd := exec.Command(
	"cat",
	"source.txt")
	
	cmd.Dir = dir
	out, err := cmd.CombinedOutput()
	if err != nil {
		log.Printf("Failed to update %s, see output below\n%sContinuing...", dir, out)
		return "", err
	}

	source := string(out[:])
	par, _ := filepath.Split(dir)
	
	cmd2 := exec.Command(
		"robocopy",
		source,
		dir,
		"/MIR",
		"/mt:4")

	cmd2.Dir = par
	cmd2.CombinedOutput()
	return g.HeadRev(dir)
}

func (g *LocalWinDriver) Clone(dir, inurl string) (string, error) {
	par, _ := filepath.Split(dir)
	
	source:= strings.TrimPrefix(inurl, "file://")
	cmd := exec.Command(
		"robocopy",
		source,
		dir,
		"/MIR",
		"/mt:4")

	cmd.Dir = par
	cmd.CombinedOutput()

	return g.HeadRev(dir)
}

func (g *LocalWinDriver) SpecialFiles() []string {
	return []string{
		".dll", ".exe",
	}
}
