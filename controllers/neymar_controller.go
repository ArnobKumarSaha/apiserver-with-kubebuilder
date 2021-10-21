/*
Copyright 2021.

Licensed under the Apache License, Version 2.0 (the "License");
you may not use this file except in compliance with the License.
You may obtain a copy of the License at

    http://www.apache.org/licenses/LICENSE-2.0

Unless required by applicable law or agreed to in writing, software
distributed under the License is distributed on an "AS IS" BASIS,
WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
See the License for the specific language governing permissions and
limitations under the License.
*/

package controllers

import (
	"context"
	"fmt"
	"github.com/go-logr/logr"
	corev1 "k8s.io/api/core/v1"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"

	appsv1 "k8s.io/api/apps/v1"

	"k8s.io/apimachinery/pkg/runtime"
	ctrl "sigs.k8s.io/controller-runtime"
	"sigs.k8s.io/controller-runtime/pkg/client"

	webappv1 "saha.com/mycrd/api/v1"
)

// NeymarReconciler reconciles a Neymar object
type NeymarReconciler struct {
	client.Client
	Log    logr.Logger
	Scheme *runtime.Scheme
}

//+kubebuilder:rbac:groups=webapp.saha.com,resources=neymars,verbs=get;list;watch;create;update;patch;delete
//+kubebuilder:rbac:groups=webapp.saha.com,resources=neymars/status,verbs=get;update;patch
//+kubebuilder:rbac:groups=webapp.saha.com,resources=neymars/finalizers,verbs=update

// Reconcile is part of the main kubernetes reconciliation loop which aims to
// move the current state of the cluster closer to the desired state.
// TODO(user): Modify the Reconcile function to compare the state specified by
// the Neymar object against the actual cluster state, and then
// perform operations to make the cluster state reflect the state specified by
// the user.
//
// For more details, check Reconcile and its Result here:
// - https://pkg.go.dev/sigs.k8s.io/controller-runtime@v0.8.3/pkg/reconcile
func (r *NeymarReconciler) Reconcile(ctx context.Context, req ctrl.Request) (ctrl.Result, error) {
	fmt.Println("req.name = ", req.Name, req.Namespace)
	// _ = log.FromContext(ctx)
	log := r.Log.WithValues("neymar", req.NamespacedName)

	// your logic here

	// 1. Load the Neymar by name
	var jr webappv1.Neymar
	if err := r.Get(ctx, req.NamespacedName, &jr); err != nil {
		log.Error(err, "unable to fetch neymar")
		// we'll ignore not-found errors, since they can't be fixed by an immediate
		// requeue (we'll need to wait for a new notification), and we can get them
		// on deleted requests.
		return ctrl.Result{}, client.IgnoreNotFound(err)
	}

	// 2. List all active deployments, and update the status
	var childDeps appsv1.DeploymentList
	if err := r.List(ctx, &childDeps, client.InNamespace(req.Namespace), client.MatchingFields{depOwnerKey: req.Name}); err != nil {
		log.Error(err, "unable to list child deployments")
		return ctrl.Result{}, err
	}

	// if no childDeployment found , create it on the cluster
	if len(childDeps.Items) == 0 {
		deploy := newDeployment(&jr)
		if err := r.Create(ctx, deploy); err != nil {
			log.Error(err, "unable to create deployment for Neymar", "Deploy", deploy)
			return ctrl.Result{}, err
		} else {
			fmt.Println("created childDeps = ", len(childDeps.Items))
		}
	}

	// Same for service
	var childSvcs corev1.ServiceList
	if err := r.List(ctx, &childSvcs, client.InNamespace(req.Namespace)); err != nil {
		log.Error(err, "unable to list child services")
		// ERROR	controllers.Neymar	unable to list child services
		// {"neymar": "default/neymar-sample", "error": "Index with name field:.metadata.controller does not exist"}
		return ctrl.Result{}, err
	}


	// if no childService found or the default 'kubernetes' service found, create one on the cluster
	if len(childSvcs.Items) == 0 || (len(childSvcs.Items) == 1 && childSvcs.Items[0].Name == "kubernetes") {
		svcObj := newService(&jr)
		if err := r.Create(ctx, svcObj); err != nil {
			log.Error(err, "unable to create service for Neymar", "service", svcObj)
			return ctrl.Result{}, err
		} else {
			fmt.Println("created childSvc = ", len(childSvcs.Items))
		}
	}
	// your logic here
	fmt.Println("Reconcilier function has been called")

	return ctrl.Result{}, nil
}

var (
	depOwnerKey = ".metadata.controller"
	svcOwnerKey = ".metadata.controller"
	apiGVStr    = webappv1.GroupVersion.String()
)

// SetupWithManager sets up the controller with the Manager.
func (r *NeymarReconciler) SetupWithManager(mgr ctrl.Manager) error {

	if err := mgr.GetFieldIndexer().IndexField(context.Background(), &appsv1.Deployment{}, depOwnerKey, func(rawObj client.Object) []string {
		// grab the deploy object, extract the owner...
		deploy := rawObj.(*appsv1.Deployment)
		owner := metav1.GetControllerOf(deploy)
		if owner == nil {
			return nil
		}
		// ...make sure it's a Nicedeploy...
		if owner.APIVersion != apiGVStr || owner.Kind != "Neymar" {
			return nil
		}

		// ...and if so, return it
		return []string{owner.Name}
	}); err != nil {
		return err
	}

	fmt.Println("SetupWithManager successful. ")
	return ctrl.NewControllerManagedBy(mgr).
		For(&webappv1.Neymar{}).
		Owns(&appsv1.Deployment{}).
		Owns(&corev1.Service{}).
		Complete(r)
}
