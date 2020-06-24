package probe

import (
	"context"

	operatorsv1alpha1 "github.com/awgreene/status-probe-operator/pkg/apis/operators/v1alpha1"
)

func (r *ReconcileProbe) updateProbeStatus(probe *operatorsv1alpha1.Probe) error {
	probe.Status.ProbeResources = make([]operatorsv1alpha1.ProbeResourceStatus, len(probe.Spec.ProbeResources))
	for i, resource := range probe.Spec.ProbeResources {
		probe.Status.ProbeResources[i] = operatorsv1alpha1.ProbeResourceStatus{
			Name: resource.Name,
		}
	}
	return r.client.Status().Update(context.TODO(), probe)
}
