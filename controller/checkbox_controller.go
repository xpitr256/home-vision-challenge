package controller

import (
	"encoding/json"
	"errors"
	"github.com/xpitr256/home-vision-challenge/model"
	"github.com/xpitr256/home-vision-challenge/service"
	"net/http"
	"strconv"
)

type CheckboxResponse struct {
	ImageName       string           `json:"image_name"`
	TotalDetections int              `json:"total_detections"`
	Checkboxes      []model.Checkbox `json:"checkboxes"`
	SizeInPixels    int              `json:"size_in_pixels"`
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
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}
	checkboxes, imageName, err := service.GetCheckboxes(checkboxSizeInPixels)
	if err != nil {
		w.WriteHeader(http.StatusInternalServerError)
		json.NewEncoder(w).Encode(map[string]string{"error": "Processing checkboxes"})
		return
	}
	w.WriteHeader(http.StatusOK)
	response := CheckboxResponse{
		ImageName:       imageName,
		TotalDetections: len(checkboxes),
		Checkboxes:      checkboxes,
		SizeInPixels:    checkboxSizeInPixels,
	}
	json.NewEncoder(w).Encode(response)
}
