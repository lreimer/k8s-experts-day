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
	"k8s.io/apimachinery/pkg/util/intstr"
)

//MutatingServerHandler listen to admission requests and serve responses
type MutatingServerHandler struct {
}

func (vs *MutatingServerHandler) serve(w http.ResponseWriter, r *http.Request) {
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

	if r.URL.Path != "/mutate" {
		log.Fatalf("No /mutate URL called")
		http.Error(w, "No /mutate URL called", http.StatusBadRequest)
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

	// add the liveness probe to the pod
	if pod.Spec.Containers[0].LivenessProbe == nil {
		log.Println("Added Liveness Probe to pod.")
		pod.Spec.Containers[0].LivenessProbe = &v1.Probe{
			InitialDelaySeconds: 5,
			TimeoutSeconds:      1,
			PeriodSeconds:       10,
			Handler: v1.Handler{
				HTTPGet: &v1.HTTPGetAction{
					Port: intstr.FromInt(80),
					Path: "/",
				},
			},
		}
	}

	// add the readoness probe to the pod
	if pod.Spec.Containers[0].ReadinessProbe == nil {
		log.Println("Added Readiness Probe to pod.")
		pod.Spec.Containers[0].ReadinessProbe = &v1.Probe{
			InitialDelaySeconds: 1,
			TimeoutSeconds:      1,
			PeriodSeconds:       5,
			Handler: v1.Handler{
				HTTPGet: &v1.HTTPGetAction{
					Port: intstr.FromInt(80),
					Path: "/",
				},
			},
		}
	}

	containersBytes, err := json.Marshal(&pod.Spec.Containers)
	if err != nil {
		http.Error(w, fmt.Sprintf("marshall containers: %v", err), http.StatusInternalServerError)
		return
	}

	// build json patch
	patch := []JSONPatchEntry{
		JSONPatchEntry{
			OP:    "add",
			Path:  "/metadata/labels/probes-added",
			Value: []byte(`"OK"`),
		},
		JSONPatchEntry{
			OP:    "replace",
			Path:  "/spec/containers",
			Value: containersBytes,
		},
	}

	patchBytes, err := json.Marshal(&patch)
	if err != nil {
		http.Error(w, fmt.Sprintf("marshall jsonpatch: %v", err), http.StatusInternalServerError)
		return
	}

	patchType := admissionv1.PatchTypeJSONPatch

	// build admission response
	admissionResponse := &admissionv1.AdmissionResponse{
		UID:       arRequest.Request.UID,
		Allowed:   true,
		Patch:     patchBytes,
		PatchType: &patchType,
	}

	arResponse := &admissionv1.AdmissionReview{
		TypeMeta: metav1.TypeMeta{
			Kind:       "AdmissionReview",
			APIVersion: "admission.k8s.io/v1",
		},
		Response: admissionResponse,
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

type JSONPatchEntry struct {
	OP    string          `json:"op"`
	Path  string          `json:"path"`
	Value json.RawMessage `json:"value,omitempty"`
}
