package vcs

import (
	"bytes"
	"encoding/json"
	"io"
	"log"
	"os/exec"
	"path/filepath"
	"strings"
	"net/url"
	"io/ioutil"
	"os"
)

func init() {
	Register(newLocal, "local", "localfolder")
}

type LocalDriver struct {
	Username string `json:"username"`
	Password string `json:"password"`
	Source string `source`
}

func newLocal(b []byte) (Driver, error) {
	var d LocalDriver

	if b != nil {
		if err := json.Unmarshal(b, &d); err != nil {
			return nil, err
		}
	}

	return &d, nil
}

func (g *LocalDriver) HeadRev(dir string) (string, error) {
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

func (g *LocalDriver) Pull(dir string) (string, error) {
	cmd := exec.Command(
	"cat",
	"source.txt")
	cmd.Dir = dir
	out, err := cmd.CombinedOutput()
	if err != nil {
		log.Printf("Failed to SVN update %s, see output below\n%sContinuing...", dir, out)
		return "", err
	}

	source := string(out[:])
	par, _ := filepath.Split(dir)
	cmd2 := exec.Command(
		"rsync",
		"-zvr",
		source,
		dir)

	cmd2.Dir = par
	out1, err1 := cmd2.CombinedOutput()
	if err1 != nil {
		log.Printf("Failed to sync %s, see output below\n%sContinuing...", source, out1)
		return "", err1
	}

	return g.HeadRev(dir)
}

func check(e error) {
    if e != nil {
            panic(e)
	        }
		}

func WriteVersion(dir string) {
	cmd2 := exec.Command(
		"date",
		"--iso-8601")

	cmd2.Dir = dir
	out1, err1 := cmd2.CombinedOutput()
	if err1 != nil {
		log.Printf("Failed to get date %s...", out1)
		panic(err1)
	}

	var versionBuffer bytes.Buffer
	versionBuffer.WriteString(dir)
	versionBuffer.WriteString("/version.txt")

     err := ioutil.WriteFile(versionBuffer.String(), out1, 0644)
     check(err)
}

func WriteStringToFile(filepath, s string) error {
     fo, err := os.Create(filepath)
     	 if err != nil {
	    	return err
		       }
			defer fo.Close()

			_, err = io.Copy(fo, strings.NewReader(s))
			   if err != nil {
			      	  return err
				  	 }

					 return nil
					 }

func (g *LocalDriver) Clone(dir, inurl string) (string, error) {
	par, _ := filepath.Split(dir)
	 u, err := url.Parse(inurl)
	     if err != nil {
	             panic(err)
		         }

	cmd := exec.Command(
		"rsync",
		"-zvr",
		u.Path,
		dir)

	cmd.Dir = par
	out, err := cmd.CombinedOutput()
	if err != nil {
		log.Printf("Failed to checkout %s, see output below\n%sContinuing...", inurl, out)
		return "", err
	}

	var sourceBuffer bytes.Buffer
	sourceBuffer.WriteString(dir)
	sourceBuffer.WriteString("/source.txt")

	WriteStringToFile(sourceBuffer.String(), u.Path)

	WriteVersion(dir)
	return g.HeadRev(dir)
}

func (g *LocalDriver) SpecialFiles() []string {
	return []string{
		".dll", ".exe",
	}
}
