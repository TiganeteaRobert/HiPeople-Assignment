package image

import (
	"errors"
	"fmt"
	"io"
	"net/http"
	"os"
	"regexp"
	"strings"
)

var (
	allowedExtensions = []string{`png`, `jpg`, `svg`}
)

// addImage adds an image to the document store
func addImage(w http.ResponseWriter, r *http.Request) {
	// handle logging
	InfoLogger.Printf(`started addImage`)
	defer InfoLogger.Printf(`finished addImage flow`)

	// parse the request form, allow maximum 32 MB file size for the image
	err := r.ParseMultipartForm(32 << 20)
	if err != nil {
		InfoLogger.Printf(`received file is over the size limit. error: %v`, err)
		w.WriteHeader(http.StatusBadRequest)
		w.Write([]byte(`file size is over the 32MB limit`))
		return
	}
	// extract the image from the form
	originalImage, imgHeader, err := r.FormFile("image")
	if err != nil {
		InfoLogger.Printf(`could not extract image from the form. error: %v`, err)
		w.WriteHeader(http.StatusBadRequest)
		w.Write([]byte(`could not extract image from the form, image key not found`))
		return
	}

	InfoLogger.Printf(`extracted file, filename: %s`, imgHeader.Filename)

	// split the image's filename into name and extension
	splitFilename := strings.Split(imgHeader.Filename, ".")
	var extensionAllowed bool
	name := splitFilename[0]
	var extension string
	// check if the file's extension is allowed
	// also make sure files without an extension are also forbidden
	if len(splitFilename) > 1 {
		extension = splitFilename[1]
		for _, allowedExtension := range allowedExtensions {
			if extension == allowedExtension {
				extensionAllowed = true
			}
		}
	}
	if !extensionAllowed {
		InfoLogger.Printf(`file extension is not allowed`)
		w.WriteHeader(http.StatusBadRequest)
		w.Write([]byte(`file extension is not allowed`))
		return
	}

	// get all the document store image names
	images, err := os.ReadDir(`document_store`)
	if err != nil {
		w.WriteHeader(http.StatusInternalServerError)
		return
	}

	// prepare regex for finding images with the same name
	re, err := regexp.Compile(name + `\(.+\)\..+`)
	if err != nil {
		ErrorLogger.Printf(`compiling filename regex. error: %v`, err)
		w.WriteHeader(http.StatusInternalServerError)
		return
	}

	// count the images that have the same name in the document_store
	count := 0
	for _, img := range images {
		if match := re.FindString(img.Name()); match != `` {
			count += 1
		}
	}

	InfoLogger.Printf(`found %d images with the same name`, count)

	// create a new filename, adding the index of the image
	// e.g. first upload of an image example.png will be stored as example(0).png
	// subsequent upload will be stored as example(1).png...
	imgHeader.Filename = name + fmt.Sprintf(`(%d).`, count) + extension

	// create a new file with the image's filename in the document store
	newImage, err := os.Create("document_store/" + imgHeader.Filename)
	if err != nil {
		ErrorLogger.Printf(`creating new file in document store. error: %v`, err)
		w.WriteHeader(http.StatusInternalServerError)
		return
	}
	defer newImage.Close()

	// copy the image to the newly created file
	_, err = io.Copy(newImage, originalImage)
	if err != nil {
		ErrorLogger.Printf(`copying original image to the file in document store. error: %v`, err)
		w.WriteHeader(http.StatusInternalServerError)
		return
	}

	InfoLogger.Printf(`created new image, filename: %s`, imgHeader.Filename)

	// return success status and image filename
	w.WriteHeader(200)
	w.Write([]byte(imgHeader.Filename))
}

// getImage retrieves an image from the document store
func getImage(w http.ResponseWriter, r *http.Request) {
	imageName, ok := r.URL.Query()["image"]
	if !ok {
		InfoLogger.Printf(`could not extract image ID from the image URL parameter`)
		w.WriteHeader(http.StatusBadRequest)
		w.Write([]byte(`could not extract image ID from the image URL parameter`))
		return
	}

	InfoLogger.Printf(`extracted image ID from URL parameter: %s`, imageName[0])

	img, err := os.ReadFile(`document_store/` + imageName[0])
	if err != nil {
		if errors.Is(err, os.ErrNotExist) {
			InfoLogger.Printf(`no image found for the requested name`)
			w.WriteHeader(http.StatusNotFound)
			return
		} else if errors.Is(err, os.ErrPermission) {
			InfoLogger.Printf(`permission denied`)
			w.WriteHeader(http.StatusForbidden)
			return
		} else {
			ErrorLogger.Printf(`reading image from document store. error: %v`, err)
			w.WriteHeader(http.StatusInternalServerError)
			return
		}
	}

	InfoLogger.Printf(`successfully retreived image from document store`)

	w.WriteHeader(http.StatusOK)
	w.Header().Set("Content-Type", "application/octet-stream")
	w.Write(img)
}
