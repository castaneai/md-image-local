package main

import (
	"fmt"
	"github.com/google/uuid"
	"io"
	"io/ioutil"
	"log"
	"net/http"
	"os"
	"regexp"
	"strings"
)

var pattern = regexp.MustCompile(`!\[([^]]*)]\((https?://[^)]+)\)`)

func main() {
	md, err := ioutil.ReadAll(os.Stdin)
	if err != nil {
		log.Fatalf("failed to read markdown source: %+v", err)
	}
	mds := pattern.ReplaceAllStringFunc(string(md), func(s string) string {
		img, err := extractImage(s)
		if err != nil {
			log.Fatalf("failed to extract image part: %+v", err)
		}
		localPath, err := downloadToLocal(img.url)
		if err != nil {
			log.Fatalf("failed to download to local: %+v", err)
		}
		return fmt.Sprintf(`![%s](%s)`, img.alt, localPath)
	})
	_, _ = os.Stdout.Write([]byte(mds))
}

type mdImage struct {
	alt string
	url string
}

func extractImage(mdPart string) (*mdImage, error) {
	ms := pattern.FindStringSubmatch(mdPart)
	if len(ms) < 3 {
		return nil, fmt.Errorf("failed to match markdown image (input: %s)", mdPart)
	}
	return &mdImage{
		alt: ms[1],
		url: ms[2],
	}, nil
}

func downloadToLocal(imageURL string) (string, error) {
	res, err := http.Get(imageURL)
	if err != nil {
		return "", err
	}
	defer res.Body.Close()
	ext := ".png"
	urlExt := strings.ToLower(imageURL)
	if strings.HasSuffix(urlExt, ".jpg") {
		ext = ".jpg"
	}
	if strings.HasSuffix(urlExt, ".gif") {
		ext = ".gif"
	}
	filename := uuid.Must(uuid.NewRandom()).String() + ext
	file, err := os.Create(filename)
	if err != nil {
		return "", err
	}
	defer file.Close()
	if _, err := io.Copy(file, res.Body); err != nil {
		return "", err
	}
	return filename, nil
}