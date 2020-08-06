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
	"path/filepath"
	"testing"

	. "github.com/onsi/ginkgo"
	. "github.com/onsi/gomega"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	"k8s.io/apimachinery/pkg/types"
	"k8s.io/client-go/kubernetes/scheme"
	"k8s.io/client-go/rest"
	ctrl "sigs.k8s.io/controller-runtime"
	"sigs.k8s.io/controller-runtime/pkg/client"
	"sigs.k8s.io/controller-runtime/pkg/envtest"
	"sigs.k8s.io/controller-runtime/pkg/envtest/printer"
	logf "sigs.k8s.io/controller-runtime/pkg/log"
	"sigs.k8s.io/controller-runtime/pkg/log/zap"

	storev1 "github.com/zoncoen-sample/k8s-pizza-order-controller/api/v1"
	// +kubebuilder:scaffold:imports
)

// These tests use Ginkgo (BDD-style Go testing framework). Refer to
// http://onsi.github.io/ginkgo/ to learn more about Ginkgo.

var cfg *rest.Config
var k8sClient client.Client
var testEnv *envtest.Environment

func TestAPIs(t *testing.T) {
	RegisterFailHandler(Fail)

	RunSpecsWithDefaultAndCustomReporters(t,
		"Controller Suite",
		[]Reporter{printer.NewlineReporter{}})
}

var _ = BeforeSuite(func(done Done) {
	logf.SetLogger(zap.LoggerTo(GinkgoWriter, true))

	By("bootstrapping test environment")
	testEnv = &envtest.Environment{
		CRDDirectoryPaths: []string{filepath.Join("..", "config", "crd", "bases")},
	}

	var err error
	cfg, err = testEnv.Start()
	Expect(err).ToNot(HaveOccurred())
	Expect(cfg).ToNot(BeNil())

	err = storev1.AddToScheme(scheme.Scheme)
	Expect(err).NotTo(HaveOccurred())

	// +kubebuilder:scaffold:scheme

	k8sClient, err = client.New(cfg, client.Options{Scheme: scheme.Scheme})
	Expect(err).ToNot(HaveOccurred())
	Expect(k8sClient).ToNot(BeNil())

	ctx := context.Background()
	amount := 1
	order := storev1.PizzaOrder{
		ObjectMeta: metav1.ObjectMeta{
			Name:      "1",
			Namespace: "pizzaorder",
		},
		Spec: storev1.PizzaOrderSpec{
			Items: []storev1.Item{
				{
					Name:   "margherita",
					Amount: &amount,
				},
			},
		},
		Status: storev1.PizzaOrderStatus{
			Phase: storev1.OrderPhaseCreated,
		},
	}
	err = k8sClient.Create(ctx, &order)
	Expect(err).ToNot(HaveOccurred())

	reconciler := &PizzaOrderReconciler{
		Client: k8sClient,
		Log:    logf.Log,
	}
	req := ctrl.Request{
		NamespacedName: types.NamespacedName{
			Name:      order.ObjectMeta.Name,
			Namespace: order.ObjectMeta.Namespace,
		},
	}
	_, err = reconciler.Reconcile(req)
	Expect(err).ToNot(HaveOccurred())

	var result storev1.PizzaOrder
	err = k8sClient.Get(ctx, req.NamespacedName, &result)
	Expect(err).ToNot(HaveOccurred())
	Expect(result.Status.Phase).To(Equal(storev1.OrderPhaseAccepted))
	Expect(result.Status.OrderedAt).NotTo(BeNil())

	close(done)
}, 60)

var _ = AfterSuite(func() {
	By("tearing down the test environment")
	err := testEnv.Stop()
	Expect(err).ToNot(HaveOccurred())
})
