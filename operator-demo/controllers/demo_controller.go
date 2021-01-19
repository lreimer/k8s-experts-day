/*
Copyright 2021 Mario-Leander Reimer.
*/

package controllers

import (
	"context"

	"github.com/go-logr/logr"
	"k8s.io/apimachinery/pkg/runtime"
	ctrl "sigs.k8s.io/controller-runtime"
	"sigs.k8s.io/controller-runtime/pkg/client"

	k8sv1 "github.com/lreimer/k8s-experts-day/operator-demo/api/v1"
)

// DemoReconciler reconciles a Demo object
type DemoReconciler struct {
	client.Client
	Log    logr.Logger
	Scheme *runtime.Scheme
}

// +kubebuilder:rbac:groups=k8s.qaware.de,resources=demoes,verbs=get;list;watch;create;update;patch;delete
// +kubebuilder:rbac:groups=k8s.qaware.de,resources=demoes/status,verbs=get;update;patch

func (r *DemoReconciler) Reconcile(req ctrl.Request) (ctrl.Result, error) {
	_ = context.Background()
	_ = r.Log.WithValues("demo", req.NamespacedName)

	// your logic here

	return ctrl.Result{}, nil
}

func (r *DemoReconciler) SetupWithManager(mgr ctrl.Manager) error {
	return ctrl.NewControllerManagedBy(mgr).
		For(&k8sv1.Demo{}).
		Complete(r)
}
