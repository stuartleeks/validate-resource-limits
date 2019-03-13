package main

import (
	"encoding/json"
	"fmt"
	"io/ioutil"
	"log"
	"net/http"

	"k8s.io/api/admission/v1beta1"
	"k8s.io/api/core/v1"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
)

func main() {
	server := http.Server{
		Addr: ":8080",
	}

	mux := http.NewServeMux()
	mux.HandleFunc("/validate", validate)

	server.Handler = mux
}

func validate(w http.ResponseWriter, r *http.Request) {
	var body []byte
	if r.Body == nil {
		log.Printf("Missing body\n")
		http.Error(w, "Missing body!", http.StatusBadRequest)
		return
	}
	data, err := ioutil.ReadAll(r.Body)
	if err != nil {
		log.Printf("Failed to read body\n")
		http.Error(w, "Failed to read body!", http.StatusBadRequest)
		return
	}
	body = data

	admissionReview := v1beta1.AdmissionReview{}
	if err := json.Unmarshal(body, &admissionReview); err != nil {
		log.Printf("Unable to marshal context\n")
		http.Error(w, "Unable to Unmarshal content", http.StatusBadRequest)
		return
	}

	pod := v1.Pod{}
	if err := json.Unmarshal(admissionReview.Request.Object.Raw, &pod); err != nil {
		http.Error(w, "Unable to Unmarshal content", http.StatusBadRequest)
		return
	}
	for _, container := range pod.Spec.Containers {
		cpuLimit := container.Resources.Limits[v1.ResourceCPU]
		if cpuLimit.IsZero() {
			log.Printf("Missing CPU limit for container '%s'\n", container.Name)
			admissionReviewResponse := v1beta1.AdmissionReview{
				Response: &v1beta1.AdmissionResponse{
					Allowed: false,
					Result: &metav1.Status{
						Message: fmt.Sprintf("Missing CPU limit for container '%s'", container.Name),
					},
				},
			}
			response, err := json.Marshal(admissionReviewResponse)
			if err != nil {
				log.Printf("Failed to marshal response\n")
				http.Error(w, "Failed to marshal response", http.StatusInternalServerError)
				return
			}
			_, err = w.Write(response)
			if err != nil {
				log.Printf("Failed to write response\n")
				http.Error(w, "Failed to write response", http.StatusInternalServerError)
				return
			}
		}
	}
}
