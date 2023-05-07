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
	v1 "k8s.io/api/core/v1"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	"k8s.io/apimachinery/pkg/runtime"
	"k8s.io/apimachinery/pkg/types"
	"k8s.io/apimachinery/pkg/util/intstr"
	ctrl "sigs.k8s.io/controller-runtime"
	"sigs.k8s.io/controller-runtime/pkg/client"
	"sigs.k8s.io/controller-runtime/pkg/log"
	"time"
)

// DOReconciler reconciles a DO object
type DOReconciler struct {
	client.Client
	Scheme *runtime.Scheme
}

//+kubebuilder:rbac:groups=operators.adityajoshi.online,resources=does,verbs=get;list;watch;create;update;patch;delete
//+kubebuilder:rbac:groups=operators.adityajoshi.online,resources=does/status,verbs=get;update;patch
//+kubebuilder:rbac:groups=operators.adityajoshi.online,resources=does/finalizers,verbs=update

// Reconcile is part of the main kubernetes reconciliation loop which aims to
// move the current state of the cluster closer to the desired state.
// TODO(user): Modify the Reconcile function to compare the state specified by
// the DO object against the actual cluster state, and then
// perform operations to make the cluster state reflect the state specified by
// the user.
//
// For more details, check Reconcile and its Result here:
// - https://pkg.go.dev/sigs.k8s.io/controller-runtime@v0.12.2/pkg/reconcile
func (r *DOReconciler) Reconcile(ctx context.Context, req ctrl.Request) (ctrl.Result, error) {

	logger := log.FromContext(ctx)

	// TODO(user): your logic here
	logger.Info("reconcile called")

	var doCR operatorsv1alpha1.DO

	err := r.Get(ctx, req.NamespacedName, &doCR)
	if err != nil {
		logger.Error(err, "unable to get the CR from the cluster")
	}
	fmt.Println("CR fetched from cluster", doCR)

	podObj := getPodObject(doCR.Name, doCR.Namespace, doCR.Spec.Image)

	key := types.NamespacedName{
		Namespace: doCR.Namespace,
		Name:      doCR.Name,
	}
	err = r.Get(ctx, key, &podObj)
	if err == nil {
		logger.Info("got the pod, skipping creating")
		return ctrl.Result{}, nil
	}

	// POD creation
	err = r.Create(ctx, &podObj, &client.CreateOptions{})
	if err != nil {
		logger.Error(err, "Failed to create pod")
		return ctrl.Result{}, err
	}

	// SVC creation
	svc := getSvcObject(doCR.Name, doCR.Namespace)
	err = r.Create(ctx, &svc, &client.CreateOptions{})
	if err != nil {
		logger.Error(err, "Failed to create pod")
		return ctrl.Result{}, err
	}

	return ctrl.Result{RequeueAfter: time.Duration(time.Second) * 30}, nil
}

func getPodObject(name, namespace, image string) v1.Pod {
	pod := v1.Pod{
		TypeMeta: metav1.TypeMeta{},
		ObjectMeta: metav1.ObjectMeta{
			Name:      name,
			Namespace: namespace,
			Labels: map[string]string{
				"app": name,
			},
		},
		Spec: v1.PodSpec{
			Containers: []v1.Container{
				{
					Name:  name,
					Image: image,
				},
			},
		},
	}
	return pod

}

func getSvcObject(name, namespace string) v1.Service {
	return v1.Service{
		TypeMeta: metav1.TypeMeta{},
		ObjectMeta: metav1.ObjectMeta{
			Name:      name,
			Namespace: namespace,
		},
		Spec: v1.ServiceSpec{
			Selector: map[string]string{
				"app": name,
			},
			Ports: []v1.ServicePort{
				{
					Name:       "http",
					Protocol:   "TCP",
					Port:       80,
					TargetPort: intstr.IntOrString{},
				},
			},
		},
	}
}

// SetupWithManager sets up the controller with the Manager.
func (r *DOReconciler) SetupWithManager(mgr ctrl.Manager) error {
	return ctrl.NewControllerManagedBy(mgr).
		For(&operatorsv1alpha1.DO{}).
		Complete(r)
}
