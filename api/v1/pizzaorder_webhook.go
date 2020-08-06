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

package v1

import (
	apierrors "k8s.io/apimachinery/pkg/api/errors"
	"k8s.io/apimachinery/pkg/runtime"
	"k8s.io/apimachinery/pkg/runtime/schema"
	"k8s.io/apimachinery/pkg/util/validation/field"
	ctrl "sigs.k8s.io/controller-runtime"
	logf "sigs.k8s.io/controller-runtime/pkg/log"
	"sigs.k8s.io/controller-runtime/pkg/webhook"
)

// log is for logging in this package.
var pizzaorderlog = logf.Log.WithName("pizzaorder-resource")

func (r *PizzaOrder) SetupWebhookWithManager(mgr ctrl.Manager) error {
	return ctrl.NewWebhookManagedBy(mgr).
		For(r).
		Complete()
}

// EDIT THIS FILE!  THIS IS SCAFFOLDING FOR YOU TO OWN!

// +kubebuilder:webhook:path=/mutate-store-example-zoncoen-net-v1-pizzaorder,mutating=true,failurePolicy=fail,groups=store.example.zoncoen.net,resources=pizzaorders,verbs=create;update,versions=v1,name=mpizzaorder.kb.io

var _ webhook.Defaulter = &PizzaOrder{}

// Default implements webhook.Defaulter so a webhook will be registered for the type
func (r *PizzaOrder) Default() {
	pizzaorderlog.Info("default", "name", r.Name)

	for i, item := range r.Spec.Items {
		if item.Amount == nil {
			amount := 1
			r.Spec.Items[i].Amount = &amount
		}
	}

	if r.Status.Phase == "" {
		r.Status.Phase = OrderPhaseCreated
	}
}

// TODO(user): change verbs to "verbs=create;update;delete" if you want to enable deletion validation.
// +kubebuilder:webhook:verbs=create;update,path=/validate-store-example-zoncoen-net-v1-pizzaorder,mutating=false,failurePolicy=fail,groups=store.example.zoncoen.net,resources=pizzaorders,versions=v1,name=vpizzaorder.kb.io

var _ webhook.Validator = &PizzaOrder{}

// ValidateCreate implements webhook.Validator so a webhook will be registered for the type
func (r *PizzaOrder) ValidateCreate() error {
	pizzaorderlog.Info("validate create", "name", r.Name)

	var errs field.ErrorList
	if len(r.Spec.Items) == 0 {
		errs = append(errs, &field.Error{
			Type:     field.ErrorTypeRequired,
			Field:    "items",
			BadValue: r.Spec.Items,
			Detail:   "empty order",
		})
	}

	if len(errs) == 0 {
		return nil
	}
	return apierrors.NewInvalid(schema.GroupKind{Group: "store.example.zoncoen.net", Kind: "PizzaOrder"}, r.Name, errs)
}

// ValidateUpdate implements webhook.Validator so a webhook will be registered for the type
func (r *PizzaOrder) ValidateUpdate(old runtime.Object) error {
	pizzaorderlog.Info("validate update", "name", r.Name)

	// TODO(user): fill in your validation logic upon object update.
	return nil
}

// ValidateDelete implements webhook.Validator so a webhook will be registered for the type
func (r *PizzaOrder) ValidateDelete() error {
	pizzaorderlog.Info("validate delete", "name", r.Name)

	// TODO(user): fill in your validation logic upon object deletion.
	return nil
}
