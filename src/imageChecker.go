package main

import (
	"bufio"
	"fmt"
	"os"
	"strings"
)

func imageMatchesAllowed(image string, allowed []string) bool {
	imageParts := strings.Split(image, ":")
	imageName := imageParts[0]
	for _, allowedImage := range allowed {
		imageParts := strings.Split(imageName, "/")
		allowedParts := strings.Split(allowedImage, "/")
		match := true
		for i, part := range allowedParts {
			if i >= len(imageParts) || part != imageParts[i] {
				match = false
				break
			}
		}
		if match {
			return true
		}
	}
	return false
}

func main() {
	defaultDockerfilePath := "Dockerfile"
	allowedImagesList, exists := os.LookupEnv("ALLOWED_IMAGES")
	if !exists {
		fmt.Println("The environment variable ALLOWED_IMAGES is not set. Exiting.")
		os.Exit(1)
	}
	allowedImages := strings.Split(allowedImagesList, ",")
	dockerfilePath := defaultDockerfilePath
	if len(os.Args) > 1 {
		dockerfilePath = os.Args[1]
	}
	file, err := os.Open(dockerfilePath)
	if err != nil {
		fmt.Printf("Error opening Dockerfile at '%s': %s\n", dockerfilePath, err)
		os.Exit(1)
	}
	defer file.Close()
	anyImageNotAllowed := false
	scanner := bufio.NewScanner(file)
	for scanner.Scan() {
		line := scanner.Text()
		if strings.HasPrefix(line, "FROM ") {
			parts := strings.Fields(line)
			if len(parts) > 1 {
				fromImage := parts[1]
				if !imageMatchesAllowed(fromImage, allowedImages) {
					anyImageNotAllowed = true
					fmt.Printf("The FROM image '%s' does not match any allowed configuration.\n", fromImage)
				}
			}
		}
	}
	if err := scanner.Err(); err != nil {
		fmt.Printf("Error reading Dockerfile: %s\n", err)
		os.Exit(1)
	}
	if anyImageNotAllowed {
		os.Exit(1)
	} else {
		fmt.Println("All FROM images match an allowed configuration.")
	}
}
