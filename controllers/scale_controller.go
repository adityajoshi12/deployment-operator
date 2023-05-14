/*
Copyright 2023.

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
	operatorsv1alpha1 "github.com/adityajoshi12/deployment-operator/api/v1alpha1"
	appsv1 "k8s.io/api/apps/v1"
	"k8s.io/apimachinery/pkg/runtime"
	"k8s.io/apimachinery/pkg/types"
	ctrl "sigs.k8s.io/controller-runtime"
	"sigs.k8s.io/controller-runtime/pkg/client"
	"sigs.k8s.io/controller-runtime/pkg/log"
	"time"
)

// ScaleReconciler reconciles a Scale object
type ScaleReconciler struct {
	client.Client
	Scheme *runtime.Scheme
}

//+kubebuilder:rbac:groups=operators.adityajoshi.online,resources=scales,verbs=get;list;watch;create;update;patch;delete
//+kubebuilder:rbac:groups=operators.adityajoshi.online,resources=scales/status,verbs=get;update;patch
//+kubebuilder:rbac:groups=operators.adityajoshi.online,resources=scales/finalizers,verbs=update

// Reconcile is part of the main kubernetes reconciliation loop which aims to
// move the current state of the cluster closer to the desired state.
// TODO(user): Modify the Reconcile function to compare the state specified by
// the Scale object against the actual cluster state, and then
// perform operations to make the cluster state reflect the state specified by
// the user.
//
// For more details, check Reconcile and its Result here:
// - https://pkg.go.dev/sigs.k8s.io/controller-runtime@v0.12.2/pkg/reconcile
func (r *ScaleReconciler) Reconcile(ctx context.Context, req ctrl.Request) (ctrl.Result, error) {
	logger := log.FromContext(ctx)

	logger.Info("reconcile called")

	var scaleCR operatorsv1alpha1.Scale

	err := r.Get(ctx, req.NamespacedName, &scaleCR)
	if err != nil {
		logger.Error(err, "unable to get the CR from the cluster")
	}
	fmt.Println("CR fetched from cluster", scaleCR)
	// handle the timing
	var deployment appsv1.Deployment

	key := types.NamespacedName{
		Name:      scaleCR.Spec.DeploymentName,
		Namespace: scaleCR.Spec.DeploymentNamespace,
	}

	currentHour := time.Now().UTC().Hour()
	startTime := scaleCR.Spec.StartTime
	endTime := scaleCR.Spec.EndTime
	fmt.Println("currentHour :", currentHour)

	if startTime <= currentHour && currentHour <= endTime {
		//TODO scale down the instance
		err = r.Get(ctx, key, &deployment)
		if err != nil {
			return ctrl.Result{RequeueAfter: time.Second * 30}, err
		}
		replicas := int32(scaleCR.Spec.MinReplicas)
		deployment.Spec.Replicas = &replicas
		err := r.Update(ctx, &deployment)
		if err != nil {
			return ctrl.Result{RequeueAfter: time.Second * 30}, err
		}
	} else {
		err = r.Get(ctx, key, &deployment)
		if err != nil {
			return ctrl.Result{RequeueAfter: time.Duration(time.Second) * 30}, err
		}
		replicas := int32(scaleCR.Spec.MaxReplicas)
		deployment.Spec.Replicas = &replicas
		err = r.Update(ctx, &deployment)
		if err != nil {
			return ctrl.Result{RequeueAfter: time.Second * 30}, err
		}
	}

	// TODO find if the current time is in the range of startTime and endTime

	// Get the deployment and after that we will update the replicas in the deployment and then we will apply that deployment

	return ctrl.Result{}, nil
}

// SetupWithManager sets up the controller with the Manager.
func (r *ScaleReconciler) SetupWithManager(mgr ctrl.Manager) error {
	return ctrl.NewControllerManagedBy(mgr).
		For(&operatorsv1alpha1.Scale{}).
		Complete(r)
}
