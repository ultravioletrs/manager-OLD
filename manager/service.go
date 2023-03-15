// Copyright (c) Mainflux
// SPDX-License-Identifier: Apache-2.0

package manager

import (
	"errors"

	"github.com/digitalocean/go-libvirt"
)

var (
	// ErrMalformedEntity indicates malformed entity specification (e.g.
	// invalid username or password).
	ErrMalformedEntity = errors.New("malformed entity specification")

	// ErrUnauthorizedAccess indicates missing or invalid credentials provided
	// when accessing a protected resource.
	ErrUnauthorizedAccess = errors.New("missing or invalid credentials provided")
)

// Service specifies an API that must be fullfiled by the domain service
// implementation, and all of its decorators (e.g. logging & metrics).
type Service interface {
	// Ping compares a given string with secret
	Ping(string) (string, error)
}

type managerService struct {
	secret  string
	libvirt *libvirt.Libvirt
}

var _ Service = (*managerService)(nil)

// New instantiates the manager service implementation.
func New(secret string, libvirtConn *libvirt.Libvirt) Service {
	return &managerService{
		secret:  secret,
		libvirt: libvirtConn,
	}
}

func (ks *managerService) Ping(secret string) (string, error) {
	dom, err := createDomain(ks.libvirt, poolXML, volXML, domXML)
	if err != nil {
		return "", ErrMalformedEntity
	}
	_ = dom

	if ks.secret != secret {
		return "", ErrUnauthorizedAccess
	}
	return "Hello World :)", nil
}
