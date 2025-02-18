package controller

import (
	"encoding/json"
	"errors"
	"github.com/xpitr256/home-vision-challenge/model"
	"github.com/xpitr256/home-vision-challenge/service"
	"image"
	"log"
	"net/http"
	"strconv"
)

type CheckboxResponse struct {
	ImageName              string           `json:"image_name"`
	TotalDetections        int              `json:"total_detections"`
	CheckboxSizeInPixels   int              `json:"checkbox_size_in_pixels"`
	ImageWithCheckboxesUrl string           `json:"image_with_checkboxes_url"`
	Checkboxes             []model.Checkbox `json:"checkboxes"`
}

func GetCheckboxSize(r *http.Request) (int, error) {
	checkboxSizeInPixels := model.CheckboxDefaultSizeInPixels
	checkboxSizeInPixelsStr := r.URL.Query().Get("size")
	if checkboxSizeInPixelsStr != "" {
		converted, err := strconv.Atoi(checkboxSizeInPixelsStr)
		if err != nil {
			return 0, errors.New("invalid 'size' parameter, must be an integer between 1 and 200")
		}
		checkboxSizeInPixels = converted
	}
	if checkboxSizeInPixels <= 0 || checkboxSizeInPixels > 200 {
		return 0, errors.New("'size' parameter must be between 1 and 200")
	}
	return checkboxSizeInPixels, nil
}

// Acts as a Template Method
func processCheckboxRequest(w http.ResponseWriter, r *http.Request, loadImageFunc func(*http.Request) (image.Image, string, error)) {
	w.Header().Set("Content-Type", "application/json")
	// TODO: defer Add a latency metric here to measure the whole time taken
	checkboxSizeInPixels, err := GetCheckboxSize(r)
	if err != nil {
		log.Printf("Error getting checkbox size: %v", err)
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}

	// TODO: Add a latency metric here to measure the time taken to load the image
	formImage, imageName, err := loadImageFunc(r)
	if err != nil {
		log.Printf("Error loading image: %v", err)
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	// TODO: Add a latency metric here to measure the time taken to process checkboxes
	checkboxes, responseImageUrl, err := service.GetCheckboxes(checkboxSizeInPixels, formImage)
	if err != nil {
		log.Printf("Error processing checkboxes: %v", err)
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	w.WriteHeader(http.StatusOK)
	response := CheckboxResponse{
		ImageName:              imageName,
		TotalDetections:        len(checkboxes),
		Checkboxes:             checkboxes,
		CheckboxSizeInPixels:   checkboxSizeInPixels,
		ImageWithCheckboxesUrl: responseImageUrl,
	}
	json.NewEncoder(w).Encode(response)
}

func checkboxGetHandler(w http.ResponseWriter, r *http.Request) {
	// TODO: Add a request count metric here for GET requests t
	log.Printf("Entering checkboxGetHandler with method: %s, URL: %s", r.Method, r.URL)
	processCheckboxRequest(w, r, service.LoadTestImage)
}

func checkboxPostHandler(w http.ResponseWriter, r *http.Request) {
	// TODO: Add a request count metric here for POST requests
	log.Printf("Entering checkboxPostHandler with method: %s, Content-Type: %s, Content-Length: %d", r.Method, r.Header.Get("Content-Type"), r.ContentLength)
	processCheckboxRequest(w, r, service.LoadImageFromRequest)
}

func CheckboxHandler(w http.ResponseWriter, r *http.Request) {
	switch r.Method {
	case http.MethodPost:
		checkboxPostHandler(w, r)
	case http.MethodGet:
		checkboxGetHandler(w, r)
	default:
		http.Error(w, "Method not allowed", http.StatusMethodNotAllowed)
	}
}
