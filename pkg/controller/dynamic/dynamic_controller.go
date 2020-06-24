package dynamic

import (
	"context"
	"fmt"

	operatorsv1alpha1 "github.com/awgreene/status-probe-operator/pkg/apis/operators/v1alpha1"
	"github.com/yalp/jsonpath"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	"k8s.io/apimachinery/pkg/apis/meta/v1/unstructured"
	"k8s.io/apimachinery/pkg/runtime"
	"k8s.io/apimachinery/pkg/runtime/schema"
	"k8s.io/client-go/dynamic"
	"k8s.io/klog"
	"sigs.k8s.io/controller-runtime/pkg/controller"
	"sigs.k8s.io/controller-runtime/pkg/handler"
	logf "sigs.k8s.io/controller-runtime/pkg/log"
	"sigs.k8s.io/controller-runtime/pkg/manager"
	"sigs.k8s.io/controller-runtime/pkg/reconcile"

	"sigs.k8s.io/controller-runtime/pkg/source"
)

var log = logf.Log.WithName("controller_dynamic")

/**
* USER ACTION REQUIRED: This is a scaffold file intended for the user to modify with their own Controller
* business logic.  Delete these comments after modifying this file.*
 */

// Add creates a new Dynamic Controller and adds it to the Manager. The Manager will set fields on the Controller
// and Start it when the Manager is Started.
func Add(mgr manager.Manager) error {
	return add(mgr, newReconciler(mgr))
}

// newReconciler returns a new reconcile.Reconciler
func newReconciler(mgr manager.Manager) reconcile.Reconciler {
	dynClient, errClient := dynamic.NewForConfig(mgr.GetConfig())
	if errClient != nil {
		klog.Fatalf("Error received creating client %v", errClient)
	}
	return &ReconcileDynamic{client: dynClient, scheme: mgr.GetScheme()}
}

// add adds a new Controller to mgr with r as the reconcile.Reconciler
func add(mgr manager.Manager, r reconcile.Reconciler) error {
	// Create a new controller
	c, err := controller.New("dynamic-controller", mgr, controller.Options{Reconciler: r})
	if err != nil {
		return err
	}

	obj := unstructured.Unstructured{}
	obj.GetObjectKind().SetGroupVersionKind(schema.GroupVersionKind{
		Group:   "example.com",
		Version: "v1alpha1",
		Kind:    "Foo",
	})

	err = c.Watch(&source.Kind{Type: &obj}, &handler.EnqueueRequestForObject{})
	if err != nil {
		return err
	}

	return nil
}

// blank assignment to verify that ReconcileDynamic implements reconcile.Reconciler
var _ reconcile.Reconciler = &ReconcileDynamic{}

// ReconcileDynamic reconciles a Dynamic object
type ReconcileDynamic struct {
	// This client, initialized using mgr.Client() above, is a split client
	// that reads objects from the cache and writes to the apiserver
	client dynamic.Interface
	scheme *runtime.Scheme
}

const (
	componentConditionsJSONPath = "$.status.conditions"
)

// Reconcile reads that state of the cluster for a Dynamic object and makes changes based on the state read
// and what is in the Dynamic.Spec
// TODO(user): Modify this Reconcile function to implement your Controller logic.  This example creates
// a Pod as an example
// Note:
// The Controller will requeue the Request to be processed again if the returned error is non-nil or
// Result.Requeue is true, otherwise upon completion it will remove the work from the queue.
func (r *ReconcileDynamic) Reconcile(request reconcile.Request) (reconcile.Result, error) {
	reqLogger := log.WithValues("Request.Namespace", request.Namespace, "Request.Name", request.Name)
	reqLogger.Info("Reconciling Dynamic")

	instance, err := r.client.Resource(schema.GroupVersionResource{
		Group:    "example.com",
		Version:  "v1alpha1",
		Resource: "foos", // Plural
	}).Namespace(request.Namespace).Get(context.TODO(), request.Name, metav1.GetOptions{})
	if err != nil {
		return reconcile.Result{}, err
	}

	conditionsBlob, err := jsonpath.Read(instance.UnstructuredContent(), componentConditionsJSONPath)
	if err != nil {
		return reconcile.Result{Requeue: false}, err
	}

	conditionsArray, ok := conditionsBlob.([]interface{})
	if !ok {
		return reconcile.Result{Requeue: false}, fmt.Errorf("Unable to get conditions array: %v", err)
	}

	operatorProbeConditions := []operatorsv1alpha1.ProbeCondition{}
	for _, c := range conditionsArray {
		conditionMap, ok := c.(map[string]interface{})
		if !ok {
			return reconcile.Result{Requeue: false}, fmt.Errorf("Unable to convert object in condition array to a map: %v", err)
		}
		conditionType, ok := conditionMap["type"]
		if !ok {
			return reconcile.Result{Requeue: false}, fmt.Errorf("Unable to find condition type: %v", err)
		}

		// Abstract code that identifies if type maps to Operator Status Condition.
		typeString, ok := conditionType.(string)
		if !ok {
			return reconcile.Result{Requeue: false}, fmt.Errorf("Condition Type should be a string")
		}

		if IsOperatorProbeCondition(conditionType.(string)) {
			klog.Infof("Found target Type")
			probeCondition := operatorsv1alpha1.ProbeCondition{Type: MapCondition(typeString)}

			if reason, ok := conditionMap["reason"].(string); ok {
				probeCondition.Reason = reason
			}
			if message, ok := conditionMap["message"].(string); ok {
				probeCondition.Message = message
			}
			if lastTransitionTime, ok := conditionMap["lastTransitionTime"].(string); ok {
				probeCondition.LastTransitionTime = lastTransitionTime
			}
			if status, ok := conditionMap["status"].(string); ok {
				probeCondition.Status = status
			}
			operatorProbeConditions = append(operatorProbeConditions, probeCondition)
		}
	}

	klog.Infof("Conditions %v", operatorProbeConditions)

	/*if err = r.updateProbeStatus(instance); err != nil {
		reqLogger.Info("Error Updating probe status: %v\n", err)
	}*/

	return reconcile.Result{}, nil
}

func IsOperatorProbeCondition(s string) bool {
	if s == "test" || s == "critical" {
		return true
	}
	return false
}

func MapCondition(s string) string {
	if s == "test" {
		return "Upgradeable"
	}

	if s == "critical" {
		return "Important"
	}
	return ""
}
