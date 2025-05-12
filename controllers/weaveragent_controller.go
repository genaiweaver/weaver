package controllers

import (
	"context"
	"fmt"

	corev1 "k8s.io/api/core/v1"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	"k8s.io/apimachinery/pkg/runtime"
	"k8s.io/apimachinery/pkg/util/intstr"
	ctrl "sigs.k8s.io/controller-runtime"
	"sigs.k8s.io/controller-runtime/pkg/client"
	"sigs.k8s.io/controller-runtime/pkg/controller/controllerutil"
	"sigs.k8s.io/controller-runtime/pkg/log"

	eventingv1 "knative.dev/eventing/pkg/apis/eventing/v1"
	duckv1 "knative.dev/pkg/apis/duck/v1"
	servingv1 "knative.dev/serving/pkg/apis/serving/v1"

	weaverv1alpha1 "github.com/weaver/weaver/api/v1alpha1"
)

// WeaverAgentReconciler reconciles a WeaverAgent object
type WeaverAgentReconciler struct {
	client.Client
	Scheme *runtime.Scheme
	KnativeServingAvailable bool
	KnativeEventingAvailable bool
}

// RBAC permissions
// +kubebuilder:rbac:groups=weaver.io,resources=weaveragents,verbs=get;list;watch;create;update;patch;delete
// +kubebuilder:rbac:groups=weaver.io,resources=weaveragents/status,verbs=get;patch;update
// +kubebuilder:rbac:groups=weaver.io,resources=weaveragents/finalizers,verbs=update
// +kubebuilder:rbac:groups=serving.knative.dev,resources=services,verbs=get;list;watch;create;update;delete
// +kubebuilder:rbac:groups=eventing.knative.dev,resources=triggers,verbs=get;list;watch;create;update;delete

func (r *WeaverAgentReconciler) Reconcile(ctx context.Context, req ctrl.Request) (ctrl.Result, error) {
	logger := log.FromContext(ctx)

	// 1. Load the WeaverAgent resource
	var agent weaverv1alpha1.WeaverAgent
	if err := r.Get(ctx, req.NamespacedName, &agent); err != nil {
		return ctrl.Result{}, client.IgnoreNotFound(err)
	}

	// Check if Knative Serving is available
	if !r.KnativeServingAvailable {
		logger.Info("Knative Serving is not available, skipping service reconciliation")
		agent.Status.ObservedGeneration = agent.Generation
		agent.Status.Conditions = append(agent.Status.Conditions, metav1.Condition{
			Type:    "KnativeServingAvailable",
			Status:  metav1.ConditionFalse,
			Reason:  "NotInstalled",
			Message: "Knative Serving is not installed in the cluster",
		})
		if err := r.Status().Update(ctx, &agent); err != nil {
			return ctrl.Result{}, err
		}
		return ctrl.Result{}, nil
	}

	// 2. Reconcile each Node as a Knative Service
	for _, node := range agent.Spec.Nodes {
		svc := &servingv1.Service{
			ObjectMeta: metav1.ObjectMeta{
				Name:      node.ServiceName,
				Namespace: agent.Namespace,
			},
			Spec: servingv1.ServiceSpec{
				ConfigurationSpec: servingv1.ConfigurationSpec{
					Template: servingv1.RevisionTemplateSpec{
						Spec: servingv1.RevisionSpec{
							PodSpec: corev1.PodSpec{
								Containers: []corev1.Container{{
									Image: fmt.Sprintf("registry.local/%s:latest", node.ServiceName),
									ReadinessProbe: &corev1.Probe{
										ProbeHandler: corev1.ProbeHandler{
											HTTPGet: &corev1.HTTPGetAction{
												Path: "/health",
												Port: intstr.FromInt(8080),
											},
										},
									},
								}},
							},
						},
					},
				},
			},
		}
		// Owner-reference so Service is garbage-collected
		if err := controllerutil.SetControllerReference(&agent, svc, r.Scheme); err != nil {
			return ctrl.Result{}, err
		}
		if _, err := controllerutil.CreateOrUpdate(ctx, r.Client, svc, func() error {
			// This is the mutate function that should update the object
			svc.Spec = servingv1.ServiceSpec{
				ConfigurationSpec: servingv1.ConfigurationSpec{
					Template: servingv1.RevisionTemplateSpec{
						Spec: servingv1.RevisionSpec{
							PodSpec: corev1.PodSpec{
								Containers: []corev1.Container{{
									Image: fmt.Sprintf("registry.local/%s:latest", node.ServiceName),
									ReadinessProbe: &corev1.Probe{
										ProbeHandler: corev1.ProbeHandler{
											HTTPGet: &corev1.HTTPGetAction{
												Path: "/health",
												Port: intstr.FromInt(8080),
											},
										},
									},
								}},
							},
						},
					},
				},
			}
			return nil
		}); err != nil {
			logger.Error(err, "unable to create or update Knative Service", "service", node.ServiceName)
			return ctrl.Result{}, err
		}
	}

	// Check if Knative Eventing is available
	if !r.KnativeEventingAvailable {
		logger.Info("Knative Eventing is not available, skipping trigger reconciliation")
		agent.Status.ObservedGeneration = agent.Generation
		agent.Status.Conditions = append(agent.Status.Conditions, metav1.Condition{
			Type:    "KnativeEventingAvailable",
			Status:  metav1.ConditionFalse,
			Reason:  "NotInstalled",
			Message: "Knative Eventing is not installed in the cluster",
		})
		if err := r.Status().Update(ctx, &agent); err != nil {
			return ctrl.Result{}, err
		}
		return ctrl.Result{}, nil
	}

	// 3. Reconcile each Edge as a Knative Trigger
	for i, edge := range agent.Spec.Edges {
		trig := &eventingv1.Trigger{
			ObjectMeta: metav1.ObjectMeta{
				Name:      fmt.Sprintf("%s-edge-%d", agent.Name, i),
				Namespace: agent.Namespace,
			},
			Spec: eventingv1.TriggerSpec{
				Broker: agent.Spec.Broker,
				Filter: &eventingv1.TriggerFilter{
					Attributes: map[string]string{"type": edge.EventType},
				},
				Subscriber: duckv1.Destination{
					Ref: &duckv1.KReference{
						Kind:       "Service",
						APIVersion: "serving.knative.dev/v1",
						Name:       edge.To,
					},
				},
			},
		}
		if err := controllerutil.SetControllerReference(&agent, trig, r.Scheme); err != nil {
			return ctrl.Result{}, err
		}
		if _, err := controllerutil.CreateOrUpdate(ctx, r.Client, trig, func() error {
			// This is the mutate function that should update the object
			trig.Spec = eventingv1.TriggerSpec{
				Broker: agent.Spec.Broker,
				Filter: &eventingv1.TriggerFilter{
					Attributes: map[string]string{"type": edge.EventType},
				},
				Subscriber: duckv1.Destination{
					Ref: &duckv1.KReference{
						Kind:       "Service",
						APIVersion: "serving.knative.dev/v1",
						Name:       edge.To,
					},
				},
			}
			return nil
		}); err != nil {
			logger.Error(err, "unable to create or update Trigger", "trigger", trig.Name)
			return ctrl.Result{}, err
		}
	}

	// 4. Update status to reflect successful reconciliation
	agent.Status.ObservedGeneration = agent.Generation
	if err := r.Status().Update(ctx, &agent); err != nil {
		return ctrl.Result{}, err
	}

	return ctrl.Result{}, nil
}

// SetupWithManager registers the controller with the manager
func (r *WeaverAgentReconciler) SetupWithManager(mgr ctrl.Manager) error {
	logger := log.Log.WithName("weaveragent-controller")

	// Only try to add Knative schemes if the respective features are available
	if r.KnativeServingAvailable {
		logger.Info("Knative Serving is available, registering Serving types")
		if err := servingv1.AddToScheme(mgr.GetScheme()); err != nil {
			logger.Error(err, "Failed to add Knative Serving scheme")
			r.KnativeServingAvailable = false
		}
	} else {
		logger.Info("Knative Serving is not available. To enable Knative Serving features, install Knative Serving in your cluster")
	}

	if r.KnativeEventingAvailable {
		logger.Info("Knative Eventing is available, registering Eventing types")
		if err := eventingv1.AddToScheme(mgr.GetScheme()); err != nil {
			logger.Error(err, "Failed to add Knative Eventing scheme")
			r.KnativeEventingAvailable = false
		}
	} else {
		logger.Info("Knative Eventing is not available. To enable Knative Eventing features, install Knative Eventing in your cluster")
	}

	// Create a controller builder that only watches WeaverAgent resources
	builder := ctrl.NewControllerManagedBy(mgr).
		For(&weaverv1alpha1.WeaverAgent{})

	// Only watch Knative resources if they're available and successfully registered
	if r.KnativeServingAvailable {
		logger.Info("Setting up watches for Knative Serving resources")
		builder = builder.Owns(&servingv1.Service{})
	}

	if r.KnativeEventingAvailable {
		logger.Info("Setting up watches for Knative Eventing resources")
		builder = builder.Owns(&eventingv1.Trigger{})
	}

	return builder.Complete(r)
}
