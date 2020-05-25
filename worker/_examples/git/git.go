package main

import (
	"log"
	"os"
	"path/filepath"

	"github.com/Kenza-AI/worker/initialization"
)

func main() {
	vcsClient := initialization.NewVCS()

	repo := "https://github.com/Kenza-AI/sagify.git"
	branch := "refs/heads/master"
	commitID, gitHubAccessToken := "", ""
	path := filepath.Join(os.TempDir(), repo)

	commitID, err := vcsClient.Checkout(repo, branch, commitID, path, gitHubAccessToken)
	if err != nil {
		log.Println(err)
	}

	log.Println("Commit ID:", commitID)
}
