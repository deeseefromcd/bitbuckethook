package main

import (
	"encoding/json"
)

type Payload struct {
	CanonURL   string     `json:"canon_url"`
	Commits    []Commit   `json:"commits"`
	Repository Repository `json:"repository"`
	User       string     `json:"user"`
}

type Commit struct {
	Author       string   `json:"author"`
	Branch       string   `json:"branch"`
	Files        []File   `json:"files"`
	Message      string   `json:"message"`
	Node         string   `json:"node"`
	Parents      []string `json:"parents"`
	RawAuthor    string   `json:"raw_author"`
	RawNode      string   `json:"raw_node"`
	Timestamp    string   `json:"timestamp"`
	UTCTimestamp string   `json:"utc_timestamp"`
}

type File struct {
	File string `json:"file"`
	Type string `json:"type"`
}

type Repository struct {
	AbsoluteUrl string `json:"absolute_url"`
	IsFork      bool   `json:"fork"`
	IsPrivate   bool   `json:"is_private"`
	Name        string `json:"name"`
	Owner       string `json:"owner"`
	SCM         string `json:"scm"`
	Slug        string `json:"slug"`
	Website     string `json:"website"`
}

func parsePayload(content string) (*Payload, error) {
	payload := new(Payload)
	err := json.Unmarshal([]byte(content), payload)

	if err != nil {
		return nil, err
	}

	return payload, nil
}

func (p *Payload) GetCommitMessages() []string {
	var messages []string

	for _, commit := range p.Commits {
		messages = append(messages, commit.Message)
	}

	return messages
}
