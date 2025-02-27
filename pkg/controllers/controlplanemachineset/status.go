/*
Copyright 2022 Red Hat, Inc.

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

package controlplanemachineset

import (
	"context"

	"github.com/go-logr/logr"
	machinev1 "github.com/openshift/api/machine/v1"
	"github.com/openshift/cluster-control-plane-machine-set-operator/pkg/machineproviders"
	"sigs.k8s.io/controller-runtime/pkg/client"
)

const (
	// updatingStatus is a log message used to inform users that the ControlPlaneMachineSet status is being updated.
	updatingStatus = "Updating control plane machine set status"

	// notUpdatingStatus is a log message used to inform users that the ControlPlaneMachineSet status is not being updated.
	notUpdatingStatus = "No update to control plane machine set status required"
)

// updateControlPlaneMachineSetStatus ensures that the status of the ControlPlaneMachineSet is up to date after
// the resource has been reconciled.
func (r *ControlPlaneMachineSetReconciler) updateControlPlaneMachineSetStatus(ctx context.Context, logger logr.Logger, cpms *machinev1.ControlPlaneMachineSet, patchBase client.Patch) error {
	// TODO: Observe the data from the patchBase.Data() to determine if we actually need to update the status or not.
	// Then log and update as per the test cases.
	return nil
}

// reconcileStatusWithMachineInfo takes the information gathered in the machineInfos and reconciles the status of the
// ControlPlaneMachineSet to match the data gathered.
// In particular, it will update the ObservedGeneration, Replicas, ReadyReplicas, UnreadyReplicas and UpdatedReplicas
// fields based on the information gathered, and then set any relevant conditions if applicable.
// Eg. If there are any unowned Control Plane Nodes, or failed ControlPlaneMachines, it should add conditions to
// to denote the state of the Machines.
// It observes the following rules for setting the status:
// - Replicas is the number of Machines present (a missing index or unmanaged node does not count towards this).
// - ReadyReplicas is the number of the above Replicas which are reporting as Ready
// - UpdatedReplicas is the number of Ready Replicas that do not need an update (this should be at most 1 per index).
// - UnavailableReplicas is the number of Machines required to satisfy the requirement of at least 1 Ready Replica per
//   index. Eg. if one index has no ready replicas, this is 1, if an index has 2 ready replicas, this does not count as
//   2 available replicas.
func reconcileStatusWithMachineInfo(logger logr.Logger, cpms *machinev1.ControlPlaneMachineSet, machineInfos []machineproviders.MachineInfo) error {
	return nil
}
