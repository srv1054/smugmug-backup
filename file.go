package main

import (
	"errors"
	"fmt"
	"io"
	"log"
	"os"
	"path"

	"golang.org/x/sys/unix"
)

func createFolder(path string) {
	_, err := os.Stat(path)

	if err != nil {
		log.Printf("Creating folder %s", path)
		os.MkdirAll(path, os.ModePerm)
	}
}

func checkFolderIsWritable(folderPath string) error {
	if !path.IsAbs(folderPath) {
		return errors.New("Destination path must be an absolute path")
	}

	info, err := os.Stat(folderPath)
	if err != nil {
		return errors.New("Destination path doesn't exist")
	}

	if !info.IsDir() {
		return errors.New("Destination path isn't a directory")
	}

	if unix.Access(folderPath, unix.W_OK) != nil {
		return errors.New("Destination path in not writeable")
	}

	return nil
}

func (c *smugMugConf) saveImages(images *[]albumImage, folder string) {
	for _, image := range *images {
		if image.IsVideo {
			c.saveVideo(&image, folder)
		} else {
			c.saveImage(&image, folder)
		}
	}
}

func (c *smugMugConf) saveImage(image *albumImage, folder string) {
	dest := fmt.Sprintf("%s/%s", folder, image.FileName)
	download(dest, image.ArchivedUri)
}

func (c *smugMugConf) saveVideo(image *albumImage, folder string) {
	dest := fmt.Sprintf("%s/%s", folder, image.FileName)

	var albumVideo albumVideo
	c.get(image.Uris.LargestVideo.Uri, &albumVideo)

	download(dest, albumVideo.Response.LargestVideo.Url)
}

func download(dest, downloadURL string) {
	if _, err := os.Stat(dest); err == nil {
		fmt.Printf("File exists: %s\n", downloadURL)
		return
	}
	fmt.Printf("Getting %s\n", downloadURL)

	response, err := makeAPICall(downloadURL)
	if err != nil {
		log.Fatal(err)
	}
	defer response.Body.Close()

	file, err := os.Create(dest)
	if err != nil {
		log.Fatal(err)
	}
	defer file.Close()

	_, err = io.Copy(file, response.Body)
	if err != nil {
		log.Fatal(err)
	}
	fmt.Printf("Saved %s\n\n", dest)
}
