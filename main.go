package main

import (
	"flag"
	"fmt"
	"log"
	"net/http"
	"net/url"
	"os"
	"os/exec"
	"strings"
)

const queueSize int = 10

var (
	configFile = flag.String("c", "repos.json", "Repo configuration")
)

type Service struct {
	config   *Config
	incoming chan *Payload
	queue    map[string]chan *Payload
}

func (s *Service) Watch() {
	var repos []string
	for p := range s.incoming {
		repos = append(repos, p.Repository.Name)

		for _, repo := range repos {
			if !s.config.hasRepo(repo) {
				continue
			}

			repositories := s.config.getReposWithName(repo)

			for _, r := range repositories {
				channel := make(chan *Payload, queueSize)
				s.queue[repo+"#"+r.Branch] = channel

				go s.Process(channel, r)
			}
		}
	}
}

func (s *Service) ServeHTTP(res http.ResponseWriter, req *http.Request) {
	log.Printf("New request from %s\n", req.RemoteAddr)
	if req.Method != "POST" {
		log.Printf("Not a POST but a %s from %s\n", req.Method, req.RemoteAddr)
		res.WriteHeader(http.StatusBadRequest)
		return
	}

	GET, err := url.ParseQuery(req.URL.RawQuery)
	if err != nil {
		log.Printf("Invalid query from %s\n", req.RemoteAddr)
		res.WriteHeader(http.StatusForbidden)
		return
	}

	secret, ok := GET["secret"]
	if !ok || secret[0] != s.config.Secret {
		if !ok {
			log.Printf("No secret from %s\n", req.RemoteAddr)
		} else {
			log.Printf("Bad secret from %s [secret: %s]\n", req.RemoteAddr, secret[0])
		}
		res.WriteHeader(http.StatusForbidden)
		return
	}

	err = req.ParseForm()
	if err != nil {
		log.Printf("Bad form from %s\n", req.RemoteAddr)
		res.WriteHeader(http.StatusBadRequest)
		return
	}

	content := req.FormValue("payload")
	payload, err := parsePayload(content)

	if err != nil {
		log.Printf("Invalid payload from %s\n", req.RemoteAddr)
		res.WriteHeader(http.StatusBadRequest)
		return
	}

	s.incoming <- payload
}

func (s *Service) Process(ch chan *Payload, repo *RepoConfig) {
	log.Printf("Run `%s` on '%s'\n", repo.Command, repo.Name)

	for payload := range ch {
		log.Printf("New changeset:\n - %s\n", strings.Join(payload.GetCommitMessages(), "\n - "))
		log.Printf("Run command on %s: %s", payload.Repository.Name, repo.Command)

		command := strings.Split(repo.Command, " ")
		cmd := exec.Command(command[0], command[1:]...)
		cmd.Stdout = os.Stdout
		cmd.Stderr = os.Stderr

		if err := cmd.Run(); err != nil {
			log.Println(err)
		}
	}
}

func main() {
	flag.Parse()
	config, err := ParseConfigurationFile(*configFile)

	if err != nil {
		log.Fatal(err)
	}

	service := &Service{
		config:   config,
		incoming: make(chan *Payload, 5),
		queue:    make(map[string]chan *Payload),
	}

	go service.Watch()

	address := fmt.Sprintf(":%d", config.Port)

	log.Printf("Listening on %s\n", address)
	http.ListenAndServe(address, service)
}
