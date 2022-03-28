package serviceimport

import (
	"context"

	meshv1beta1 "bitbucket.org/realtimeai/kubeslice-operator/api/v1beta1"
	"bitbucket.org/realtimeai/kubeslice-operator/controllers"
	"bitbucket.org/realtimeai/kubeslice-operator/internal/logger"
	corev1 "k8s.io/api/core/v1"
	"k8s.io/apimachinery/pkg/api/errors"
	"k8s.io/apimachinery/pkg/types"
	ctrl "sigs.k8s.io/controller-runtime"
)

// reconcile istio resources like virtualService and serviceentry to facilitate service to service connectivity
func (r *Reconciler) reconcileIstio(ctx context.Context, serviceimport *meshv1beta1.ServiceImport) (ctrl.Result, error, bool) {
	log := logger.FromContext(ctx).WithValues("type", "Istio")
	debugLog := log.V(1)

	slice, err := controllers.GetSlice(ctx, r.Client, serviceimport.Spec.Slice)
	if err != nil {
		log.Error(err, "Unable to fetch slice for serviceexport")
		return ctrl.Result{}, err, true
	}

	if slice.Status.SliceConfig.ExternalGatewayConfig == nil ||
		slice.Status.SliceConfig.ExternalGatewayConfig.GatewayType != "istio" {
		debugLog.Info("istio not enabled for slice, skipping reconcilation")
		return ctrl.Result{}, nil, false
	}

	debugLog.Info("reconciling istio")

	// Create k8s service for app pods to connect
	res, err, requeue := r.ReconcileService(ctx, serviceimport)
	if requeue {
		return res, err, requeue
	}

	egressGatewayConfig := slice.Status.SliceConfig.ExternalGatewayConfig.Egress

	// Add service entries and virtualServices for app pods in the slice to be able to route traffic through slice.
	// Endpoints are load balanced at equal weights.
	if egressGatewayConfig != nil && egressGatewayConfig.Enabled {
		res, err, requeue = r.ReconcileServiceEntries(ctx, serviceimport, controllers.ControlPlaneNamespace)
		if requeue {
			return res, err, requeue
		}
		res, err, requeue = r.ReconcileVirtualServiceEgress(ctx, serviceimport)
		if requeue {
			return res, err, requeue
		}
	} else {
		res, err, requeue = r.ReconcileServiceEntries(ctx, serviceimport, serviceimport.Namespace)
		if requeue {
			return res, err, requeue
		}
		res, err, requeue = r.ReconcileVirtualServiceNonEgress(ctx, serviceimport)
		if requeue {
			return res, err, requeue
		}
	}

	return ctrl.Result{}, nil, false
}

func (r *Reconciler) ReconcileService(ctx context.Context, serviceimport *meshv1beta1.ServiceImport) (ctrl.Result, error, bool) {
	log := logger.FromContext(ctx).WithValues("type", "Istio Service")

	svc := &corev1.Service{}
	err := r.Get(ctx, types.NamespacedName{
		Name:      serviceimport.Name,
		Namespace: serviceimport.Namespace,
	}, svc)
	if err != nil {
		if errors.IsNotFound(err) {
			// Define a new service
			svc := r.serviceForServiceImport(serviceimport)
			log.Info("Creating a new Service", "Namespace", svc.Namespace, "Name", svc.Name)
			err = r.Create(ctx, svc)
			if err != nil {
				log.Error(err, "Failed to create new Service", "Namespace", svc.Namespace, "Name", svc.Name)
				return ctrl.Result{}, err, true
			}
			return ctrl.Result{Requeue: true}, nil, true
		}
		log.Error(err, "Failed to get Service")
		return ctrl.Result{}, err, true
	}

	// TODO handle change of ports

	return ctrl.Result{}, nil, false
}
