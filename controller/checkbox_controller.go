package controller

import (
	"encoding/json"
	"errors"
	"github.com/xpitr256/home-vision-challenge/model"
	"github.com/xpitr256/home-vision-challenge/service"
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

func CheckboxHandler(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json")
	checkboxSizeInPixels, err := GetCheckboxSize(r)
	if err != nil {
		log.Printf("Error getting checkbox size: %v", err)
		w.WriteHeader(http.StatusBadRequest)
		json.NewEncoder(w).Encode(map[string]interface{}{
			"error":  err.Error(),
			"status": http.StatusBadRequest,
		})
	}
	checkboxes, imageName, responseImageUrl, err := service.GetCheckboxes(checkboxSizeInPixels)
	if err != nil {
		log.Printf("Error processing checkboxes: %v", err)
		w.WriteHeader(http.StatusInternalServerError)
		json.NewEncoder(w).Encode(map[string]interface{}{
			"error":  err.Error(),
			"status": http.StatusInternalServerError,
		})
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
