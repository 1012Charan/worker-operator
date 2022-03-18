package manifest

import (
	"context"

	appsv1 "k8s.io/api/apps/v1"
	corev1 "k8s.io/api/core/v1"
	rbacv1 "k8s.io/api/rbac/v1"
	"k8s.io/apimachinery/pkg/api/errors"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	"sigs.k8s.io/controller-runtime/pkg/client"
)

// Install istio ingress gw resources on the given cluster in a slice
// Resources:
//  deployment (adds annotations to add the ingress pod to the slice)
//  serviceaccount
//  role
//  rolebinding
//  service (type clusterip)
func InstallIngress(ctx context.Context, c client.Client, slice string) error {
	// TODO make the objects unique per slice by adding slice name

	deploy := &appsv1.Deployment{}
	err := NewManifest("../../files/ingress/ingress-deploy.json", slice).Parse(deploy)
	if err != nil {
		return err
	}

	// Add the deployment to slice
	deploy.Labels["slice"] = slice
	deploy.Annotations = map[string]string{
		"avesha.io/slice": slice,
	}

	svc := &corev1.Service{}
	err = NewManifest("../../files/ingress/ingress-svc.json", slice).Parse(svc)
	if err != nil {
		return err
	}

	role := &rbacv1.Role{}
	err = NewManifest("../../files/ingress/ingress-role.json", slice).Parse(role)
	if err != nil {
		return err
	}

	sa := &corev1.ServiceAccount{}
	err = NewManifest("../../files/ingress/ingress-sa.json", slice).Parse(sa)
	if err != nil {
		return err
	}

	rb := &rbacv1.RoleBinding{}
	err = NewManifest("../../files/ingress/ingress-rolebinding.json", slice).Parse(rb)
	if err != nil {
		return err
	}

	objects := []client.Object{
		deploy,
		svc,
		role,
		sa,
		rb,
	}

	for _, o := range objects {
		if err := c.Create(ctx, o); err != nil {
			return err
		}
	}

	return nil
}

// Uninstall istio ingress (EW) resources fo a slice from a given cluster
// Resources:
//  deployment
//  serviceaccount
//  role
//  rolebinding
//  service
func UninstallIngress(ctx context.Context, c client.Client, slice string) error {
	// TODO objects should be unique to slice

	log.Info("deleting EW ingress gw for the slice", "slice", slice)

	objects := []client.Object{
		&appsv1.Deployment{
			ObjectMeta: metav1.ObjectMeta{
				Name:      "istio-ingressgateway",
				Namespace: "kubeslice-system",
			},
		},
		&corev1.Service{
			ObjectMeta: metav1.ObjectMeta{
				Name:      "istio-ingressgateway",
				Namespace: "kubeslice-system",
			},
		},
		&rbacv1.Role{
			ObjectMeta: metav1.ObjectMeta{
				Name:      "istio-ingressgateway-sds",
				Namespace: "kubeslice-system",
			},
		},
		&corev1.ServiceAccount{
			ObjectMeta: metav1.ObjectMeta{
				Name:      "istio-ingressgateway-service-account",
				Namespace: "kubeslice-system",
			},
		},
		&rbacv1.RoleBinding{
			ObjectMeta: metav1.ObjectMeta{
				Name:      "istio-ingressgateway-sds",
				Namespace: "kubeslice-system",
			},
		},
	}

	for _, o := range objects {
		if err := c.Delete(ctx, o); err != nil {
			// Ignore the error if the resource is already deleted
			if !errors.IsNotFound(err) {
				return err
			}
		}
	}

	return nil
}
