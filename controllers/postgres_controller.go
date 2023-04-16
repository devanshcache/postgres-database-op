/*
Copyright 2024.

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
	"go.uber.org/zap"
	"k8s.io/apimachinery/pkg/api/errors"
	"time"

	"k8s.io/apimachinery/pkg/runtime"
	ctrl "sigs.k8s.io/controller-runtime"
	"sigs.k8s.io/controller-runtime/pkg/client"

	databasev1alpha1 "github.com/devansh/database-op/api/v1alpha1"
)

// PostgresReconciler reconciles a Postgres object
type PostgresReconciler struct {
	client.Client
	Scheme *runtime.Scheme
}

var (
	RequeueAfter10Sec = ctrl.Result{RequeueAfter: time.Second * 10}
	RequeueAfter30Sec = ctrl.Result{RequeueAfter: time.Second * 30}
)

//+kubebuilder:rbac:groups=database.devansh.com,resources=postgres,verbs=get;list;watch;create;update;patch;delete
//+kubebuilder:rbac:groups=database.devansh.com,resources=postgres/status,verbs=get;update;patch
//+kubebuilder:rbac:groups=database.devansh.com,resources=postgres/finalizers,verbs=update

// Reconcile is part of the main kubernetes reconciliation loop which aims to
// move the current state of the cluster closer to the desired state.
// Modify the Reconcile function to compare the state specified by
// the Postgres object against the actual cluster state, and then
// perform operations to make the cluster state reflect the state specified by
// the user.
//
// For more details, check Reconcile and its Result here:
// - https://pkg.go.dev/sigs.k8s.io/controller-runtime@v0.14.1/pkg/reconcile
func (r *PostgresReconciler) Reconcile(ctx context.Context, req ctrl.Request) (ctrl.Result, error) {
	log.Info("Reconciling", zap.String("namespace", req.Namespace), zap.String("Request.Name", req.Name))

	// add a finalizer here

	// Fetch the Postgres CR
	instance := &databasev1alpha1.Postgres{}
	err := r.Get(ctx, req.NamespacedName, instance)
	if err != nil {
		log.Info("Error getting Database.Postgres instance: " + err.Error())
		if errors.IsNotFound(err) {
			log.Info("Database.Postgres CR was not found")
			return RequeueAfter10Sec, nil
		}
		// Error reading the object
		log.Error("Error reading the Database.Postgres instance: " + err.Error())
		return RequeueAfter10Sec, nil
	}

	return r.handleInstanceChange(ctx, instance)
}

// SetupWithManager sets up the controller with the Manager.
func (r *PostgresReconciler) SetupWithManager(mgr ctrl.Manager) error {
	return ctrl.NewControllerManagedBy(mgr).
		For(&databasev1alpha1.Postgres{}).
		Complete(r)
}
