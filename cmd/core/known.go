package main

import (
	"github.com/romeq/logfront/internal/consumers"
	"github.com/romeq/logfront/internal/consumers/ntfy_sh"
	"github.com/romeq/logfront/internal/sources"
	"github.com/romeq/logfront/internal/sources/ftp"
	"github.com/romeq/logfront/internal/sources/ssh"
)

// registerInit handles
func initRegisters() {
	// register known sources
	sources.Register(ssh.ConfigName, ssh.NewSource)
	sources.Register(ftp.ConfigName, ftp.NewSource)

	// register known consumers
	consumers.Register(ntfy_sh.ConfigName, ntfy_sh.NewConsumer)
}
