// Copyright (c) 2017 Masashi Fujita
//
// This software is released under the MIT License.
// https://opensource.org/licenses/MIT
package main

import (
	"flag"
	"fmt"
	"io/ioutil"
	"os"
	"path/filepath"
	"regexp"

	"github.com/pkg/errors"
)

var (
	option struct {
		verbose bool
		dryRun  bool
	}

	// ProgramName holds the path to the executable (if possible)
	ProgramName = getProgramName()
)

func main() {
	flag.Usage = func() {
		fmt.Fprintf(os.Stderr, "Usage: %s [options] [<groovy-home>]\n", ProgramName)
		flag.PrintDefaults()
		os.Exit(1)
	}
	flag.BoolVar(&option.verbose, "v", false, "be verbose")
	flag.BoolVar(&option.dryRun, "N", false, "Don't modify anything")
	flag.Parse()

	var groovyHome string
	if 0 < flag.NArg() {
		groovyHome = filepath.ToSlash(filepath.Clean(flag.Arg(0)))
	} else {
		var err error
		groovyHome, err = getGroovyDirectory()
		if err != nil {
			fmt.Fprintf(os.Stderr, "%s: Failed to obtain Groovy directory (%v)\n", ProgramName, err)
			os.Exit(1)
		}
	}
	if err := indynize(groovyHome); err != nil {
		fmt.Fprintf(os.Stderr, "%s: Failed to indynize (%v)\n", ProgramName, err)
		os.Exit(1)
	}
	os.Exit(0)
}

func indynize(groovyHome string) error {
	libDir := filepath.Join(groovyHome, "lib")
	origLibDir := filepath.Join(groovyHome, "lib.orig")

	if !Exists(origLibDir) {
		if !Exists(libDir) {
			return errors.Errorf("missing \"%s\"", libDir)
		}
		if err := doRename(libDir, origLibDir); err != nil {
			return errors.Wrapf(err, "failed to back up the original to \"%s\"", origLibDir)
		}
	}
	if Exists(libDir) {
		if err := doRemoveAll(libDir); err != nil {
			return errors.Wrapf(err, "failed to remove \"%s\"", libDir)
		}
	}
	if err := doMakeDirectory(libDir); err != nil {
		return errors.Wrapf(err, "failed to create new library directory \"%s\"", libDir)
	}
	verbose("%s: Linking `indy` enabled jars...\n", ProgramName)
	{
		indyDir := filepath.Join(groovyHome, "indy")
		indy, err := ioutil.ReadDir(indyDir)
		if err != nil {
			return errors.Wrapf(err, "failed to obtain `indy` enabled file list")
		}
		rxIndyJar := regexp.MustCompile(`^(?P<stem>.*)-indy.jar$`)
		for _, f := range indy {
			match := rxIndyJar.FindStringSubmatch(f.Name())
			if match == nil {
				continue
			}
			stem := match[1]
			src := filepath.Join(indyDir, f.Name())
			dst := filepath.Join(libDir, stem+".jar")
			if err := doLink(src, dst); err != nil {
				return err
			}
		}
	}
	verbose("%s: Copying originals...\n.", ProgramName)
	{
		items, err := ioutil.ReadDir(origLibDir)
		if err != nil {
			return errors.Wrapf(err, "failed to read \"%s\"", origLibDir)
		}
		for _, f := range items {
			src := filepath.Join(origLibDir, f.Name())
			dst := filepath.Join(libDir, f.Name())
			if !Exists(dst) {
				doLink(src, dst)
			}
		}
	}
	return nil
}

func doLink(src string, dst string) error {
	if option.dryRun || option.verbose {
		fmt.Fprintf(os.Stderr, "%s: Link \"%s\" to \"%s\"\n", ProgramName, src, dst)
	}
	if option.dryRun {
		return nil
	}
	return os.Link(src, dst)
}

func doMakeDirectory(dir string) error {
	if option.dryRun || option.verbose {
		fmt.Fprintf(os.Stderr, "%s: Create directory \"%s\"\n", ProgramName, dir)
	}
	if option.dryRun {
		return nil
	}
	return os.Mkdir(dir, os.ModeDir|0755)
}

func doRename(oldPath string, newPath string) error {
	if option.dryRun || option.verbose {
		fmt.Fprintf(os.Stderr, "%s: Rename \"%s\" to \"%s\"\n", ProgramName, oldPath, newPath)
	}
	if option.dryRun {
		return nil
	}
	return os.Rename(oldPath, newPath)
}

func doRemoveAll(path string) error {
	if option.dryRun || option.verbose {
		fmt.Fprintf(os.Stderr, "%s: Remove \"%s\" and its children\n", ProgramName, path)
	}
	if option.dryRun {
		return nil
	}
	return os.RemoveAll(path)
}

func getProgramName() string {
	p, err := os.Executable()
	if err == nil {
		return filepath.ToSlash(filepath.Clean(p))
	}
	return "indynize"
}

// Exists checks `path` existance (on FS).
func Exists(path string) bool {
	_, err := os.Stat(path)
	return err == nil
}

func getGroovyDirectory() (string, error) {
	d, ok := os.LookupEnv("GROOVY_HOME")
	if ok {
		return filepath.ToSlash(filepath.Clean(d)), nil
	}
	home, ok := os.LookupEnv("HOME")
	if ok {
		return filepath.ToSlash(filepath.Clean(filepath.Join(home, ".gvm/groovy/current"))), nil
	}
	return "", errors.Errorf("failed to obtain Groovy directory")
}

func verbose(format string, args ...interface{}) {
	if option.verbose {
		fmt.Fprintf(os.Stderr, format, args...)
	}
}
