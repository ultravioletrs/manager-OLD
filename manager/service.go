// Copyright (c) Mainflux
// SPDX-License-Identifier: Apache-2.0

package manager

import (
	"errors"
	"os"
	"strings"

	"github.com/digitalocean/go-libvirt"
	"github.com/golang/protobuf/ptypes/timestamp"
	"github.com/mainflux/mainflux"
)

const (
	poolXML = "xml/pool.xml"
	volXML  = "xml/vol.xml"
	domXML  = "xml/dom.xml"
)

var (
	// ErrMalformedEntity indicates malformed entity specification (e.g.
	// invalid username or password).
	ErrMalformedEntity = errors.New("malformed entity specification")

	// ErrUnauthorizedAccess indicates missing or invalid credentials provided
	// when accessing a protected resource.
	ErrUnauthorizedAccess = errors.New("missing or invalid credentials provided")

	// ErrNotFound indicates a non-existent entity request.
	ErrNotFound = errors.New("entity not found")
)

// Service specifies an API that must be fullfiled by the domain service
// implementation, and all of its decorators (e.g. logging & metrics).
type Service interface {
	CreateDomain(pool, volume, domain string) (string, error)
	Run(comp Computation) (string, error)
}

type managerService struct {
	secret     string
	libvirt    *libvirt.Libvirt
	idProvider mainflux.IDProvider
}

var _ Service = (*managerService)(nil)

// New instantiates the manager service implementation.
func New(secret string, libvirtConn *libvirt.Libvirt, idp mainflux.IDProvider) Service {
	return &managerService{
		secret:     secret,
		libvirt:    libvirtConn,
		idProvider: idp,
	}
}

func (ms *managerService) CreateDomain(poolXML, volXML, domXML string) (string, error) {
	wd, _ := os.Getwd()

	poolStr, err := readXMLFile(poolXML, "pool.xml")
	if err != nil {
		return "", err
	}
	poolStr = replaceSubstring(poolStr, "./", wd+"/")

	volStr, err := readXMLFile(volXML, "vol.xml")
	if err != nil {
		return "", err
	}
	volStr = replaceSubstring(volStr, "./", wd+"/")

	domStr, err := readXMLFile(domXML, "dom.xml")
	if err != nil {
		return "", err
	}
	domStr = replaceSubstring(domStr, "./", wd+"/")

	dom, err := createDomain(ms.libvirt, poolStr, volStr, domStr)
	if err != nil {
		return "", ErrMalformedEntity
	}

	return dom.Name, nil
}

func (ms *managerService) Run(comp Computation) (string, error) {
	// Generate a unique ID for the computation
	runID, err := ms.idProvider.ID()
	if err != nil {
		return "", err
	}

	// Initialize the Computation object
	comp.ID = runID
	comp.Status = ""
	comp.StartTime = &timestamp.Timestamp{}
	comp.EndTime = &timestamp.Timestamp{}

	// // Save the Computation object to the database
	// if err := ms.db.SaveComputation(comp); err != nil {
	// 	return "", err
	// }

	// // Start the computation process
	// go ms.processComputation(comp)

	return runID, nil
}

func readXMLFile(filename string, defaultFilename string) (string, error) {
	if filename == "" {
		filename = "./xml/" + defaultFilename
	}

	xmlBytes, err := os.ReadFile(filename)
	if err != nil {
		return "", ErrNotFound
	}

	return string(xmlBytes), nil
}

func replaceSubstring(xml, substring, replacement string) string {
	// Split the file text into lines
	lines := strings.Split(xml, "\n")

	// Create a variable to hold the resulting string
	var result strings.Builder

	// Iterate over each line
	for _, line := range lines {
		// Replace the substring with the replacement
		newLine := strings.ReplaceAll(line, substring, replacement)

		// Append the modified line to the resulting string
		result.WriteString(newLine)
		result.WriteString("\n")
	}

	return result.String()
}
