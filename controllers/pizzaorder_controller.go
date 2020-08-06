/*


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
	"time"

	"github.com/go-logr/logr"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	"k8s.io/apimachinery/pkg/runtime"
	ctrl "sigs.k8s.io/controller-runtime"
	"sigs.k8s.io/controller-runtime/pkg/client"

	storev1 "github.com/zoncoen-sample/k8s-pizza-order-controller/api/v1"
)

// PizzaOrderReconciler reconciles a PizzaOrder object
type PizzaOrderReconciler struct {
	client.Client
	Log    logr.Logger
	Scheme *runtime.Scheme
}

// +kubebuilder:rbac:groups=store.example.zoncoen.net,resources=pizzaorders,verbs=get;list;watch;create;update;patch;delete
// +kubebuilder:rbac:groups=store.example.zoncoen.net,resources=pizzaorders/status,verbs=get;update;patch

func (r *PizzaOrderReconciler) Reconcile(req ctrl.Request) (ctrl.Result, error) {
	ctx := context.Background()
	log := r.Log.WithValues("pizzaorder", req.NamespacedName)

	var order storev1.PizzaOrder
	if err := r.Get(ctx, req.NamespacedName, &order); err != nil {
		log.Error(err, "unable to fetch PizzaOrder")
		return ctrl.Result{}, client.IgnoreNotFound(err)
	}

	if order.Status.Phase != storev1.OrderPhaseCreated {
		return ctrl.Result{}, nil
	}

	// TODO: send order request
	log.Info("send order request")

	order.Status.Phase = storev1.OrderPhaseAccepted
	order.Status.OrderedAt = &metav1.Time{Time: time.Now()}
	if err := r.Update(ctx, &order); err != nil {
		log.Error(err, "unable to update PizzaOrder")
		return ctrl.Result{}, err
	}

	return ctrl.Result{}, nil
}

func (r *PizzaOrderReconciler) SetupWithManager(mgr ctrl.Manager) error {
	return ctrl.NewControllerManagedBy(mgr).
		For(&storev1.PizzaOrder{}).
		Complete(r)
}
