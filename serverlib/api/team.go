package api

import (
	"time"
)

// Team represents the domain for our team entity
type Team struct {
	ID           int64          `json:"id"`
	Name         string         `json:"name"`
	CreatedAt    time.Time      `json:"created_at"`
	Version      int            `json:"version"`
	NumMembers   int            `json:"num_members"`
	Meta         meta           `json:"meta"`
	UserAccounts []*UserAccount `json:"user_accounts"`
}

type meta struct {
	GitURL    string `json:"git_url"`
	ServerURL string `json:"server_url"`
	Version   int64  `json:"version"`
	ID        int64  `json:"id"`
}