package domain

import "time"

type FailedLoginEvent struct {
	Source     string // ssh, ftp, etc
	Username   string
	IP         string
	Timestamp  time.Time
	RawMessage string
}
