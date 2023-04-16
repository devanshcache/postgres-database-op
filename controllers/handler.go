package controllers

import (
	"context"
	databasev1alpha1 "github.com/devansh/database-op/api/v1alpha1"
	v12 "k8s.io/api/apps/v1"
	v1 "k8s.io/api/core/v1"
	"k8s.io/apimachinery/pkg/api/errors"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	"k8s.io/apimachinery/pkg/types"
	ctrl "sigs.k8s.io/controller-runtime"
	"time"
)

func (r *PostgresReconciler) handleInstanceChange(ctx context.Context, instance *databasev1alpha1.Postgres) (ctrl.Result, error) {
	log.Info("handleInstanceChange() starting... Event arrived handle it")

	// Get Namespace if not present create
	ns := &v1.NamespaceList{}
	err := r.List(ctx, ns)
	if err != nil {
		log.Error("error getting namespace list:" + err.Error())
		return RequeueAfter10Sec, err
	}

	nsFound := false
	for _, n := range ns.Items {
		if n.Name == instance.Spec.Image.Namespace {
			nsFound = true
		}
	}
	if !nsFound {
		n := &v1.Namespace{
			ObjectMeta: metav1.ObjectMeta{
				Name: instance.Spec.Image.Namespace,
			},
		}
		err := r.Create(ctx, n)
		if err != nil {
			log.Error("error creating namespace list:" + err.Error())
			return RequeueAfter10Sec, err
		}
	}

	return r.handleCreateDeployment(ctx, instance)
}

func (r *PostgresReconciler) handleCreateDeployment(ctx context.Context, instance *databasev1alpha1.Postgres) (ctrl.Result, error) {
	d := &v12.Deployment{}
	err := r.Get(ctx, types.NamespacedName{Namespace: instance.Spec.Image.Namespace,
		Name: instance.Spec.Image.Name + "-deploy"}, d)
	if err != nil {
		log.Info("Error getting deployment instance: " + err.Error())
		if errors.IsNotFound(err) {
			log.Info("Database.Postgres Deployment was not found")

			newDeploy := &v12.Deployment{
				ObjectMeta: metav1.ObjectMeta{
					Name:      instance.Spec.Image.Name + "-deploy",
					Namespace: instance.Spec.Image.Namespace,
				},
				Spec: v12.DeploymentSpec{
					Replicas: getReplicas(instance),
					Selector: &metav1.LabelSelector{
						MatchLabels: map[string]string{
							"app": "postgres",
						},
					},
					Template: v1.PodTemplateSpec{
						ObjectMeta: metav1.ObjectMeta{
							Labels: map[string]string{
								"app": "postgres", // Matching labels for Pods in the deployment
							},
						},
						Spec: v1.PodSpec{
							Containers: []v1.Container{
								{
									Name:  "postgres-container", // Container name
									Image: instance.Spec.Image.Name + ":" + instance.Spec.Image.Tag,
									Ports: []v1.ContainerPort{
										{
											ContainerPort: 80, // Container port (exposed through Service later)
										},
									},
									Env: []v1.EnvVar{
										{
											Name:  "POSTGRES_DB",
											Value: "postgres",
										},
										{
											Name:  "POSTGRES_USER",
											Value: "postgres_user",
										},
										{
											Name:  "POSTGRES_PASSWORD",
											Value: "password",
										},
									},
								},
							},
						},
					},
				},
			}
			err := r.Create(ctx, newDeploy)
			if err != nil {
				log.Error("error creating postgres deployment:" + err.Error())
				return RequeueAfter10Sec, err
			}

			return RequeueAfter10Sec, nil
		}
		// Error reading the object
		log.Error("Error reading the getting deployment instance: " + err.Error())
		return RequeueAfter10Sec, nil
	}

	if *getReplicas(instance) != *d.Spec.Replicas {
		*d.Spec.Replicas = *getReplicas(instance)
		err := r.Update(ctx, d)
		if err != nil {
			log.Error("Error updating the deployment: " + err.Error())
			return ctrl.Result{}, err
		}
	}

	log.Info("handleCreateDeployment() finishes...")
	return RequeueAfter30Sec, nil
}

func getReplicas(instance *databasev1alpha1.Postgres) *int32 {
	var replicas int32 = 1
	currHour := time.Now().UTC().Hour()
	if currHour >= instance.Spec.ScaleAt.StartHour && currHour <= instance.Spec.ScaleAt.EndHour {
		replicas = int32(instance.Spec.ScaleAt.Replicas)
	}

	return &replicas
}
