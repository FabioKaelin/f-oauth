package image

import (
	"bytes"
	"fmt"
	"image"
	"image/png"
	"io"
	"mime/multipart"
	"net/http"

	"github.com/fabiokaelin/f-oauth/config"
)

func SaveImage(url string, userid string) error {
	// Create new file name
	newFileName := "profileimage-" + userid + ".png"
	fmt.Println("newFileName", newFileName)

	// Get the data
	resp, err := http.Get(url)
	if err != nil {
		return err
	}
	defer resp.Body.Close()

	// Decode the image
	imageFile, _, err := image.Decode(resp.Body)
	if err != nil {
		return err
	}

	buf := new(bytes.Buffer)
	if err := png.Encode(buf, imageFile); err != nil {
		return err
	}

	// push to INTERNAL_IMAGE_SERVICE as formFile with name "image"
	var b bytes.Buffer
	w := multipart.NewWriter(&b)

	// Create a new form file
	fw, err := w.CreateFormFile("image", newFileName)
	if err != nil {
		return err
	}

	// Write the image data to the form file
	if _, err = io.Copy(fw, buf); err != nil {
		return err
	}

	// Close the multipart writer
	if err = w.Close(); err != nil {
		return err
	}

	// Create a new HTTP request
	req, err := http.NewRequest("POST", config.InternalImageService+"/api/users/"+userid, &b)
	if err != nil {
		return err
	}

	// Set the content type, this is very important
	req.Header.Set("Content-Type", w.FormDataContentType())

	// Do the request
	client := &http.Client{}
	res, err := client.Do(req)
	if err != nil {
		return err
	}
	if res.StatusCode != http.StatusOK {
		return fmt.Errorf("bad status: %s", res.Status)
	}

	return nil
}
