package controllers

import (
    "context"
    "time"

    metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
    "k8s.io/apimachinery/pkg/runtime"
    ctrl "sigs.k8s.io/controller-runtime"
    "sigs.k8s.io/controller-runtime/pkg/client"
    "sigs.k8s.io/controller-runtime/pkg/log"

    weaverv1alpha1 "github.com/weaver/weaver/api/v1alpha1"
)

// WeaverNodeReconciler reconciles a WeaverNode object
type WeaverNodeReconciler struct {
    client.Client
    Scheme *runtime.Scheme
}

// +kubebuilder:rbac:groups=weaver.io,resources=weavernodes,verbs=get;list;watch;create;update;patch;delete
// +kubebuilder:rbac:groups=weaver.io,resources=weavernodes/status,verbs=get;update;patch

func (r *WeaverNodeReconciler) Reconcile(ctx context.Context, req ctrl.Request) (ctrl.Result, error) {
    logger := log.FromContext(ctx)

    var node weaverv1alpha1.WeaverNode
    if err := r.Get(ctx, req.NamespacedName, &node); err != nil {
        return ctrl.Result{}, client.IgnoreNotFound(err)
    }

    // Example: update heartbeat timestamp
    now := metav1.Now()
    node.Status.LastHeartbeat = now
    node.Status.Healthy = true // you could add real health-check logic
    if err := r.Status().Update(ctx, &node); err != nil {
        logger.Error(err, "unable to update WeaverNode status")
        return ctrl.Result{}, err
    }

    // If Redis caching is enabled, write to Redis here (omitted for brevity)

    // Requeue after 30s to refresh health
    return ctrl.Result{RequeueAfter: 30 * time.Second}, nil
}

func (r *WeaverNodeReconciler) SetupWithManager(mgr ctrl.Manager) error {
    return ctrl.NewControllerManagedBy(mgr).
        For(&weaverv1alpha1.WeaverNode{}).
        Complete(r)
}
