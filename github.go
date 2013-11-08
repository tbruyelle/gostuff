package main

import "bufio"
import "bytes"
import "fmt"
import "github.com/google/go-github/github"
import "code.google.com/p/goauth2/oauth"
import "io/ioutil"
import "io"
import "os"

const (
	owner, repo = "github-username", "github-repo"
)

func ReadFile(path string) (io.Reader, error) {
	b, err := ioutil.ReadFile(path)
	if err != nil {
		return nil, err
	}
	return bytes.NewReader(b), nil
}

func ReadFileBuffer(path string) (io.Reader, error) {
	fi, err := os.Open(path)
	if err != nil {
		return nil, err
	}
	return bufio.NewReader(fi), nil
}

func main() {
	t := &oauth.Transport{Token: &oauth.Token{AccessToken: "your-personal-github-token"}}
	c := t.Client()
	client := github.NewClient(c)

	releases, resp, err := client.Repositories.ListReleases(owner, repo)
	if err != nil {
		fmt.Println("error Releases.List", err, resp)
		return
	}
	if len(releases) == 0 {
		fmt.Println("No release found for repo, create a release to make it work", repo)
		return
	}

	fmt.Printf("Release : %+v\n\n", releases[0])

	opt := &github.UploadOptions{Name: "testupload"}

	b, err := os.Open("test.txt") // pass directly the file -> ERROR 400 Bad Content length
	//b, err := ReadFileBuffer("test.txt") // pass a bufio.Reader -> ERROR 400 Bad Content length
	//b, err := ReadFile("test.txt") // Pass all file content -> OK
	if err != nil {
		fmt.Println("Unable to open file")
		return
	}
	asset, resp, err := client.Repositories.UploadReleaseAsset(owner, repo, *releases[0].ID, opt, b, "text/plain")
	if err != nil {
		fmt.Println("Error during asset upload", err)
		return
	}
	fmt.Println("New created asset", asset)
}
