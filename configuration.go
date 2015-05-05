package main

import (
	"encoding/json"
	"io/ioutil"
)

type RepoConfig struct {
	Name    string
	Branch  string
	Root    string
	Command string
}

type Config struct {
	Repos  []RepoConfig `json:"repos"`
	Port   int          `json:"port"`
	Secret string       `json:"secret"`
}

func (c *Config) hasRepo(name string) bool {
	for _, repo := range c.Repos {
		if repo.Name == name {
			return true
		}
	}

	return false
}

func (c *Config) getReposWithName(name string) []*RepoConfig {
	var repos []*RepoConfig

	for _, repo := range c.Repos {
		if repo.Name == name {
			repos = append(repos, &repo)
		}
	}

	return repos
}

func ParseConfigurationFile(fileName string) (*Config, error) {
	conf := Config{
		Port: 8765,
	}
	jsonContent, err := ioutil.ReadFile(fileName)

	if err != nil {
		return nil, err
	}

	err = json.Unmarshal(jsonContent, &conf)

	if err != nil {
		return nil, err
	}

	return &conf, nil
}
