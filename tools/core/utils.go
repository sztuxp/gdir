package core

import (
	"fmt"
	"io"
	"io/ioutil"
	"os"
	"path/filepath"
	"regexp"
)

func PromptYesNo(question string) bool {
	var line string
	for {
		fmt.Printf("%s (y/n) ", question)
		fmt.Scanln(&line)
		if regexp.MustCompile(`(?i)^\s*y\s*$`).MatchString(line) {
			return true
		} else if regexp.MustCompile(`(?i)^\s*n\s*$`).MatchString(line) {
			return false
		}
	}
}

func PromptYesNoWithDefault(question string, defaultYes bool) bool {
	var line string
	for {
		fmt.Printf("%s (", question)
		if defaultYes {
			fmt.Printf("Y/n) ")
		} else {
			fmt.Printf("y/N) ")
		}
		fmt.Scanln(&line)
		if regexp.MustCompile(`(?i)^\s*$`).MatchString(line) {
			return defaultYes
		} else if regexp.MustCompile(`(?i)^\s*y\s*$`).MatchString(line) {
			return true
		} else if regexp.MustCompile(`(?i)^\s*n\s*$`).MatchString(line) {
			return false
		}
	}
}

func CopyFile(src, dst string) (err error) {
	var srcfd *os.File
	var dstfd *os.File
	var srcinfo os.FileInfo

	if srcfd, err = os.Open(src); err != nil {
		return
	}
	defer srcfd.Close()

	if dstfd, err = os.Create(dst); err != nil {
		return
	}
	defer dstfd.Close()

	if _, err = io.Copy(dstfd, srcfd); err != nil {
		return
	}
	if srcinfo, err = os.Stat(src); err != nil {
		return
	}
	return os.Chmod(dst, srcinfo.Mode())
}

func CopyDir(src string, dst string) (err error) {
	var fds []os.FileInfo
	var srcinfo os.FileInfo

	if srcinfo, err = os.Stat(src); err != nil {
		return
	}

	if err = os.MkdirAll(dst, srcinfo.Mode()); err != nil {
		return
	}

	if fds, err = ioutil.ReadDir(src); err != nil {
		return
	}
	for _, fd := range fds {
		srcfp := filepath.Join(src, fd.Name())
		dstfp := filepath.Join(dst, fd.Name())

		if fd.IsDir() {
			if err = CopyDir(srcfp, dstfp); err != nil {
				return
			}
		} else {
			if err = CopyFile(srcfp, dstfp); err != nil {
				return
			}
		}
	}
	return nil
}
