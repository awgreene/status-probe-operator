package dynamic

import (
	operatorsv1alpha1 "github.com/awgreene/status-probe-operator/pkg/apis/operators/v1alpha1"
)

func (r *ReconcileDynamic) updateProbeStatus(probe *operatorsv1alpha1.Probe) error {
	probe.Status.ProbeResources = make([]operatorsv1alpha1.ProbeResourceStatus, len(probe.Spec.ProbeResources))
	for i, resource := range probe.Spec.ProbeResources {
		probe.Status.ProbeResources[i] = operatorsv1alpha1.ProbeResourceStatus{
			Name: resource.Name,
			Resources: []operatorsv1alpha1.ProbeResourceCondition{
				{
					UID:        probe.GetUID(),
					Conditions: []operatorsv1alpha1.ProbeCondition{},
				},
			},
		}
	}
	return nil //r.client.Resource().Namespace(probe.GetNamespace()).Update(context.TODO(), probe.(runtime.Object), metav1.UpdateOptions{})(context.TODO(), probe)
}
