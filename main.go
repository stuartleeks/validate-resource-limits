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
	log.Printf("validate-resource-limits\n")
	port := "8080"

	// TODO allow overriding these from args!
	appCertPath := "/app/certs/cert.pem"
	appKeyPath := "/app/certs/key.pem"

	server := http.Server{
		Addr: fmt.Sprintf(":%s", port),
	}

	mux := http.NewServeMux()
	mux.HandleFunc("/validate", validate)

	server.Handler = mux

	log.Printf("Starting server (listening on %s)\n", port)
	if err := server.ListenAndServeTLS(appCertPath, appKeyPath); err != nil {
		log.Printf("Failed to start server: %v", err)
	}
	log.Printf("Exiting\n")
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
