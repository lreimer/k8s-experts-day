package main

import (
	"encoding/json"
	"fmt"
	"io/ioutil"
	"log"
	"net/http"

	admissionv1 "k8s.io/api/admission/v1"
	v1 "k8s.io/api/core/v1"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
)

//ValidatingServerHandler listen to admission requests and serve responses
type ValidatingServerHandler struct {
}

func (vs *ValidatingServerHandler) serve(w http.ResponseWriter, r *http.Request) {
	var body []byte
	if r.Body != nil {
		if data, err := ioutil.ReadAll(r.Body); err == nil {
			body = data
		}
	}
	if len(body) == 0 {
		log.Fatalf("Empty body")
		http.Error(w, "Empty body", http.StatusBadRequest)
		return
	}
	log.Println("Received request")

	if r.URL.Path != "/validate" {
		log.Fatalf("No /validate URL called")
		http.Error(w, "No /validate URL called", http.StatusBadRequest)
		return
	}

	arRequest := admissionv1.AdmissionReview{}
	if err := json.Unmarshal(body, &arRequest); err != nil {
		log.Fatalf("Incorrect body")
		http.Error(w, "Incorrect body", http.StatusBadRequest)
	}

	raw := arRequest.Request.Object.Raw
	pod := v1.Pod{}
	if err := json.Unmarshal(raw, &pod); err != nil {
		log.Fatalf("Error deserializing pod")
		return
	}

	// check if there are Liveness and Readiness probes defined
	arResponse := admissionv1.AdmissionReview{}

	if pod.Spec.Containers[0].LivenessProbe != nil && pod.Spec.Containers[0].ReadinessProbe != nil {
		log.Printf("Pod container spec valid.")
		arResponse = admissionv1.AdmissionReview{
			TypeMeta: metav1.TypeMeta{
				Kind:       "AdmissionReview",
				APIVersion: "admission.k8s.io/v1",
			},
			Response: &admissionv1.AdmissionResponse{
				UID:     arRequest.Request.UID,
				Allowed: true,
			},
		}
	} else {
		log.Printf("Pod container spec invalid.")
		arResponse = admissionv1.AdmissionReview{
			TypeMeta: metav1.TypeMeta{
				Kind:       "AdmissionReview",
				APIVersion: "admission.k8s.io/v1",
			},
			Response: &admissionv1.AdmissionResponse{
				UID:     arRequest.Request.UID,
				Allowed: false,
				Result: &metav1.Status{
					Message: "Readiness and Liveness probes are required.",
				},
			},
		}
	}

	resp, err := json.Marshal(arResponse)
	if err != nil {
		log.Fatalf("Can't encode response: %v", err)
		http.Error(w, fmt.Sprintf("could not encode response: %v", err), http.StatusInternalServerError)
	}
	if _, err := w.Write(resp); err != nil {
		log.Fatalf("Can't write response: %v", err)
		http.Error(w, fmt.Sprintf("could not write response: %v", err), http.StatusInternalServerError)
	}
}
