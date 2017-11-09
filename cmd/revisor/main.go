// Package main contains the cli app for static asset revisioning.
package main

import (
	"flag"
	"log"
	"os"
	"path/filepath"
	"strings"

	"github.com/denisbrodbeck/revisor/rev"
)

var (
	successExitCode    = 0
	errorExitCode      = 1
	errorParseExitCode = 2
)

func main() {
	wd, err := os.Getwd()
	if err != nil {
		log.Println("failed to get working directory:", err)
		os.Exit(errorExitCode)
	}
	var webroot string
	var baseURL string
	var help bool
	flag.StringVar(&webroot, "root", "./public/", "")
	flag.StringVar(&baseURL, "base", "https://www.dukud.com/", "")
	flag.BoolVar(&help, "h", false, "")
	flag.BoolVar(&help, "help", false, "")

	usage := func() {
		log.Println(usageStr)
		os.Exit(errorParseExitCode)
	}

	log.SetFlags(0)
	flag.Usage = usage
	flag.Parse()

	if help || flag.NArg() == 1 {
		arg := strings.ToLower(flag.Arg(0))
		if help || arg == "help" {
			usage()
		}
	}

	if filepath.IsAbs(webroot) == false {
		webroot = filepath.Join(wd, webroot)
	}

	atomTypes := []string{".png", ".gif", ".jpg", ".jpeg", ".ico", ".woff", ".woff2", ".ttf", ".otf", ".eot"}
	moleculeTypes := []string{".css", ".js"}
	pageTypes := []string{".html", ".xml"}

	if err := rev.Revision(webroot, baseURL, atomTypes, moleculeTypes, pageTypes); err != nil {
		log.Println(err)
		os.Exit(errorExitCode)
	}
}

const usageStr = `revisor is a tool for static asset revisioning by appending content hashes to filenames

Usage: revisor [-root <path>] [-base <url>]

Flags:
  -root  <path>  webroot directory containing all ready to deploy files (default: ./public/)
  -base  <url>   base URL for revisioned assets (default: ./)

Try:
revisor -root ./public/ --base https://www.domain.com/`
