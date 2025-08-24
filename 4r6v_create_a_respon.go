package main

import (
	"encoding/json"
	"fmt"
	"net/http"

	"github.com/gorilla/mux"
)

// Device represents an IoT device
type Device struct {
	ID          string `json:"id"`
	Name        string `json:"name"`
	Description string `json:"description"`
	Status      string `json:"status"`
}

// DeviceController is responsible for handling device operations
type DeviceController struct {
	devices map[string]Device
}

// NewDeviceController returns a new instance of DeviceController
func NewDeviceController() *DeviceController {
	return &DeviceController{
		devices: make(map[string]Device),
	}
}

// CreateDevice creates a new IoT device
func (dc *DeviceController) CreateDevice(w http.ResponseWriter, r *http.Request) {
	var device Device
	err := json.NewDecoder(r.Body).Decode(&device)
	if err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}
	dc.devices[device.ID] = device
	w.WriteHeader(http.StatusCreated)
}

// GetDevice returns a single IoT device by ID
func (dc *DeviceController) GetDevice(w http.ResponseWriter, r *http.Request) {
	params := mux.Vars(r)
	device, ok := dc.devices[params["id"]]
	if !ok {
		http.Error(w, "Device not found", http.StatusNotFound)
		return
	}
	json.NewEncoder(w).Encode(device)
}

// GetAllDevices returns a list of all IoT devices
func (dc *DeviceController) GetAllDevices(w http.ResponseWriter, r *http.Request) {
	json.NewEncoder(w).Encode(dc.devices)
}

// UpdateDevice updates an existing IoT device
func (dc *DeviceController) UpdateDevice(w http.ResponseWriter, r *http.Request) {
	params := mux.Vars(r)
	device, ok := dc.devices[params["id"]]
	if !ok {
		http.Error(w, "Device not found", http.StatusNotFound)
		return
	}
	err := json.NewDecoder(r.Body).Decode(&device)
	if err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}
	dc.devices[params["id"]] = device
	w.WriteHeader(http.StatusOK)
}

// DeleteDevice deletes an IoT device by ID
func (dc *DeviceController) DeleteDevice(w http.ResponseWriter, r *http.Request) {
	params := mux.Vars(r)
	delete(dc.devices, params["id"])
	w.WriteHeader(http.StatusNoContent)
}

func main() {
	dc := NewDeviceController()

	r := mux.NewRouter()
	r.HandleFunc("/devices", dc.CreateDevice).Methods("POST")
	r.HandleFunc("/devices/{id}", dc.GetDevice).Methods("GET")
	r.HandleFunc("/devices", dc.GetAllDevices).Methods("GET")
	r.HandleFunc("/devices/{id}", dc.UpdateDevice).Methods("PATCH")
	r.HandleFunc("/devices/{id}", dc.DeleteDevice).Methods("DELETE")

	fmt.Println("Server started on port 8080")
	http.ListenAndServe(":8080", r)
}