package controllers

import (
	"context"
	"errors"
	"os"
	"reflect"

	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	"k8s.io/client-go/tools/record"
	"sigs.k8s.io/controller-runtime/pkg/event"
	"sigs.k8s.io/controller-runtime/pkg/handler"
	"sigs.k8s.io/controller-runtime/pkg/predicate"
	"sigs.k8s.io/controller-runtime/pkg/source"

	"github.com/go-logr/logr"
	"k8s.io/apimachinery/pkg/runtime"
	ctrl "sigs.k8s.io/controller-runtime"
	"sigs.k8s.io/controller-runtime/pkg/client"

	serverv1alpha1 "github.com/naveensrinivasan/api/v1alpha1"
	corev1 "k8s.io/api/core/v1"
	apierrors "k8s.io/apimachinery/pkg/api/errors"
)

const TargetSecretName = "dockerpullsecret"

// SecretReconcile reconciles a Secret object
type SecretReconcile struct {
	client.Client
	Log           logr.Logger
	Scheme        *runtime.Scheme
	EventRecorder record.EventRecorder
}

var namespaceLabels map[string]string

func init() {
	namespaceLabels = map[string]string{"inject-docker-secret": "true"}
}

const dockersecret = "dockersecret"

// +kubebuilder:rbac:groups=server.naveensrinivasan.dev,resources=secrets,verbs=get;list;watch;create;update;patch;delete
// +kubebuilder:rbac:groups=server.naveensrinivasan.dev,resources=secrets/status,verbs=get;update;patch
// +kubebuilder:rbac:groups=core,resources=secret,verbs=list;watch;get;patch;create
// +kubebuilder:rbac:groups=apps,resources=namespace,verbs=list;watch;get;update

func (r *SecretReconcile) Reconcile(req ctrl.Request) (ctrl.Result, error) {
	// This is MY_POD_NAMESPACE should have been set by the downward API to identify the namespace which this controller is running from
	MyPodNamespace := GetEnvDefault("MY_POD_NAMESPACE", "")
	if len(MyPodNamespace) == 0 {
		return ctrl.Result{}, apierrors.NewInternalError(errors.New("the env variable MY_POD_NAMESPACE hasn't been set"))
	}

	ctx := context.Background()
	log := r.Log.WithValues("secret", req.NamespacedName)
	log.Info("In the reconcile")

	var sec serverv1alpha1.Secret
	if err := r.Get(ctx, req.NamespacedName, &sec); err != nil {
		if apierrors.IsNotFound(err) {
			// This is the case where the namespace is just created and watch event kicks in.
			if req.Namespace == "" {
				log.Info("Scheduling for namespace", "namespace", req.Name)
				o := metav1.ObjectMeta{Namespace: req.Name, Name: dockersecret}
				err := r.Create(ctx, &serverv1alpha1.Secret{ObjectMeta: o})
				if err != nil {
					log.Error(err, "Error in scheduling for the namespace")
					return ctrl.Result{}, err
				}
				log.Info("scheduled for namespace", "namespace", req.Name)
				return ctrl.Result{}, nil
			}
			return ctrl.Result{}, client.IgnoreNotFound(err)
		}
		log.Info("error", "error", err)
		return ctrl.Result{}, err
	}

	storeSecret := corev1.Secret{}
	// This is the secret that would be created in every namespace
	log.Info("Looking for the original docker secret", "namespace", MyPodNamespace)
	if err := r.Get(ctx, client.ObjectKey{Name: dockersecret, Namespace: MyPodNamespace}, &storeSecret); err != nil {
		if apierrors.IsNotFound(err) {
			return ctrl.Result{}, errors.New("docker secret not found")
		}
	}

	if err := r.Get(ctx, client.ObjectKey{Name: TargetSecretName, Namespace: req.Namespace}, &storeSecret); err != nil {
		if apierrors.IsNotFound(err) {
			result, err2 := r.createDockerSecret(log, sec, ctx, storeSecret)
			if err2 != nil {
				return result, err2
			}
			log.Info("Created Secret", "namespace", req.Namespace)
		}
	}

	log.Info("found the secret", "secret", sec)
	return ctrl.Result{}, nil
}

// createDockerSecret creates the docker pull secret for the new namespace that was created based on label
func (r *SecretReconcile) createDockerSecret(log logr.Logger, sec serverv1alpha1.Secret, ctx context.Context, storeSecret corev1.Secret) (ctrl.Result, error) {
	log.Info("Creating Secret", "namespace", sec.Namespace, "TargetSecretName", sec.Name)
	newSecret := corev1.Secret{
		ObjectMeta: metav1.ObjectMeta{
			Name:      TargetSecretName,
			Namespace: sec.Namespace,
		},
		Type: corev1.DockerConfigKey,
		Data: storeSecret.Data,
	}
	storeSecret.Type = corev1.DockerConfigKey
	err := r.Create(ctx, &newSecret)
	if err != nil {
		log.Error(err, "unable to create Secret")
		return ctrl.Result{}, err
	}
	if err = ctrl.SetControllerReference(&sec, &newSecret, r.Scheme); err != nil {
		log.Error(err, "Error in setting the controller reference")
		return ctrl.Result{}, err
	}

	log.Info("Created secret", "namespace", newSecret.Namespace, "TargetSecretName", newSecret.Name)
	return ctrl.Result{}, nil
}

func (r *SecretReconcile) SetupWithManager(mgr ctrl.Manager) error {
	return ctrl.NewControllerManagedBy(mgr).
		For(&serverv1alpha1.Secret{}).
		Owns(&corev1.Secret{}).
		Watches(&source.Kind{Type: &corev1.Namespace{}},
			&handler.EnqueueRequestForObject{}).
		WithEventFilter(predicate.Funcs{CreateFunc: func(event event.CreateEvent) bool {
			// Look for only namespace with Kubeflow
			_, ok := event.Object.(*serverv1alpha1.Secret)
			if ok {
				return true
			}

			// Only for namespaces that have specific labels
			if ns, e := event.Object.(*corev1.Namespace); e {
				//compares the actual namespace labels
				return reflect.DeepEqual(ns.Labels, namespaceLabels)
			}
			return false
		}}).Complete(r)
}

// GetEnvDefault returns the value of the given environment variable or a
// default value if the given environment variable is not set.
func GetEnvDefault(variable string, defaultVal string) string {
	envVar, exists := os.LookupEnv(variable)
	if !exists {
		return defaultVal
	}
	return envVar
}
