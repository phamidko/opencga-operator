package statefulset

import (
	"github.com/phamidko/opencga-operator/pkg/kube/annotations"
	"github.com/phamidko/opencga-operator/pkg/util/merge"
	apiErrors "k8s.io/apimachinery/pkg/api/errors"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	"k8s.io/apimachinery/pkg/types"
	"sigs.k8s.io/controller-runtime/pkg/client"

	appsv1 "k8s.io/api/apps/v1"
	corev1 "k8s.io/api/core/v1"
)

const (
	notFound = -1
)

type Getter interface {
	GetDeployment(objectKey client.ObjectKey) (appsv1.Deployment, error)
}

type Updater interface {
	UpdateDeployment(sts appsv1.Deployment) (appsv1.Deployment, error)
}

type Creator interface {
	CreateDeployment(sts appsv1.Deployment) error
}

type Deleter interface {
	DeleteDeployment(objectKey client.ObjectKey) error
}

type GetUpdater interface {
	Getter
	Updater
}

type GetUpdateCreator interface {
	Getter
	Updater
	Creator
}

type GetUpdateCreateDeleter interface {
	Getter
	Updater
	Creator
	Deleter
}

// CreateOrUpdate creates the given Deployment if it doesn't exist,
// or updates it if it does.
func CreateOrUpdate(getUpdateCreator GetUpdateCreator, dep appsv1.Deployment) (appsv1.Deployment, error) {
	_, err := getUpdateCreator.GetDeployment(types.NamespacedName{Name: dep.Name, Namespace: dep.Namespace})
	if err != nil {
		if apiErrors.IsNotFound(err) {
			return appsv1.Deployment{}, getUpdateCreator.CreateDeployment(dep)
		}
		return appsv1.Deployment{}, err
	}
	return getUpdateCreator.UpdateDeployment(dep)
}

// GetAndUpdate applies the provided function to the most recent version of the object
func GetAndUpdate(getUpdater GetUpdater, nsName types.NamespacedName, updateFunc func(*appsv1.Deployment)) (appsv1.Deployment, error) {
	dep, err := getUpdater.GetDeployment(nsName)
	if err != nil {
		return appsv1.Deployment{}, err
	}
	// apply the function on the most recent version of the resource
	updateFunc(&dep)
	return getUpdater.UpdateDeployment(dep)
}

// VolumeMountData contains values required for the MountVolume function
type VolumeMountData struct {
	Name      string
	MountPath string
	Volume    corev1.Volume
	ReadOnly  bool
}

func CreateVolumeFromConfigMap(name, sourceName string, options ...func(v *corev1.Volume)) corev1.Volume {
	volume := &corev1.Volume{
		Name: name,
		VolumeSource: corev1.VolumeSource{
			ConfigMap: &corev1.ConfigMapVolumeSource{
				LocalObjectReference: corev1.LocalObjectReference{
					Name: sourceName,
				},
			},
		},
	}

	for _, option := range options {
		option(volume)
	}
	return *volume
}

func CreateVolumeFromSecret(name, sourceName string, options ...func(v *corev1.Volume)) corev1.Volume {
	permission := int32(416)
	volumeMount := &corev1.Volume{
		Name: name,
		VolumeSource: corev1.VolumeSource{
			Secret: &corev1.SecretVolumeSource{
				SecretName:  sourceName,
				DefaultMode: &permission,
			},
		},
	}
	for _, option := range options {
		option(volumeMount)
	}
	return *volumeMount

}

func CreateVolumeFromEmptyDir(name string) corev1.Volume {
	return corev1.Volume{
		Name: name,
		VolumeSource: corev1.VolumeSource{
			// No options EmptyDir means default storage medium and size.
			EmptyDir: &corev1.EmptyDirVolumeSource{},
		},
	}
}

// CreateVolumeMount returns a corev1.VolumeMount with options.
func CreateVolumeMount(name, path string, options ...func(*corev1.VolumeMount)) corev1.VolumeMount {
	volumeMount := &corev1.VolumeMount{
		Name:      name,
		MountPath: path,
	}
	for _, option := range options {
		option(volumeMount)
	}
	return *volumeMount
}

// NOOP is a valid Modification which applies no changes
func NOOP() Modification {
	return func(dep *appsv1.Deployment) {}
}

func WithSecretDefaultMode(mode *int32) func(*corev1.Volume) {
	return func(v *corev1.Volume) {
		if v.VolumeSource.Secret == nil {
			v.VolumeSource.Secret = &corev1.SecretVolumeSource{}
		}
		v.VolumeSource.Secret.DefaultMode = mode
	}
}

// WithSubPath sets the SubPath for this VolumeMount
func WithSubPath(subPath string) func(*corev1.VolumeMount) {
	return func(v *corev1.VolumeMount) {
		v.SubPath = subPath
	}
}

// WithReadOnly sets the ReadOnly attribute of this VolumeMount
func WithReadOnly(readonly bool) func(*corev1.VolumeMount) {
	return func(v *corev1.VolumeMount) {
		v.ReadOnly = readonly
	}
}

func IsReady(sts appsv1.StatefulSet, expectedReplicas int) bool {
	allUpdated := int32(expectedReplicas) == sts.Status.UpdatedReplicas
	allReady := int32(expectedReplicas) == sts.Status.ReadyReplicas
	atExpectedGeneration := sts.Generation == sts.Status.ObservedGeneration
	return allUpdated && allReady && atExpectedGeneration
}

type Modification func(*appsv1.Deployment)

func New(mods ...Modification) appsv1.Deployment {
	ocb := appsv1.Deployment{}
	for _, mod := range mods {
		mod(&ocb)
	}
	return ocb
}

func Apply(funcs ...Modification) func(*appsv1.Deployment) {
	return func(dep *appsv1.Deployment) {
		for _, f := range funcs {
			f(dep)
		}
	}
}

func WithName(name string) Modification {
	return func(dep *appsv1.Deployment) {
		dep.Name = name
	}
}

func WithNamespace(namespace string) Modification {
	return func(dep *appsv1.Deployment) {
		dep.Namespace = namespace
	}
}

func WithServiceName(svcName string) Modification {
	return func(sts *appsv1.Deployment) {
		// sts.Spec.ServiceName = svcName
	}
}

func WithLabels(labels map[string]string) Modification {
	return func(set *appsv1.Deployment) {
		set.Labels = copyMap(labels)
	}
}

func WithAnnotations(annotations map[string]string) Modification {
	return func(set *appsv1.Deployment) {
		set.Annotations = merge.StringToStringMap(set.Annotations, annotations)
	}
}

func WithMatchLabels(matchLabels map[string]string) Modification {
	return func(set *appsv1.Deployment) {
		if set.Spec.Selector == nil {
			set.Spec.Selector = &metav1.LabelSelector{}
		}
		set.Spec.Selector.MatchLabels = copyMap(matchLabels)
	}
}
func WithOwnerReference(ownerRefs []metav1.OwnerReference) Modification {
	ownerReference := make([]metav1.OwnerReference, len(ownerRefs))
	copy(ownerReference, ownerRefs)
	return func(set *appsv1.Deployment) {
		set.OwnerReferences = ownerReference
	}
}

func WithReplicas(replicas int) Modification {
	stsReplicas := int32(replicas)
	return func(sts *appsv1.Deployment) {
		sts.Spec.Replicas = &stsReplicas
	}
}

func WithRevisionHistoryLimit(revisionHistoryLimit int) Modification {
	rhl := int32(revisionHistoryLimit)
	return func(sts *appsv1.Deployment) {
		sts.Spec.RevisionHistoryLimit = &rhl
	}
}

func WithPodManagementPolicyType(policyType appsv1.PodManagementPolicyType) Modification {
	return func(set *appsv1.Deployment) {
		set.Spec.PodManagementPolicy = policyType
	}
}

func WithSelector(selector *metav1.LabelSelector) Modification {
	return func(set *appsv1.Deployment) {
		set.Spec.Selector = selector
	}
}

func WithUpdateStrategyType(strategyType appsv1.DeploymentStrategyType) Modification {
	return func(set *appsv1.Deployment) {
		set.Spec.Strategy = appsv1.DeploymentStrategy{
			Type: strategyType,
		}
	}
}

func WithPodSpecTemplate(templateFunc func(*corev1.PodTemplateSpec)) Modification {
	return func(set *appsv1.Deployment) {
		template := &set.Spec.Template
		templateFunc(template)
	}
}

func WithVolumeClaim(name string, f func(*corev1.PersistentVolumeClaim)) Modification {
	return func(set *appsv1.Deployment) {
		idx := findVolumeClaimIndexByName(name, set.Spec.VolumeClaimTemplates)
		if idx == notFound {
			set.Spec.VolumeClaimTemplates = append(set.Spec.VolumeClaimTemplates, corev1.PersistentVolumeClaim{})
			idx = len(set.Spec.VolumeClaimTemplates) - 1
		}
		pvc := &set.Spec.VolumeClaimTemplates[idx]
		f(pvc)
	}
}

func WithCustomSpecs(spec appsv1.Deployment) Modification {
	return func(set *appsv1.Deployment) {
		set.Spec = merge.StatefulSetSpecs(set.Spec, spec)
	}
}

func findVolumeClaimIndexByName(name string, pvcs []corev1.PersistentVolumeClaim) int {
	for idx, pvc := range pvcs {
		if pvc.Name == name {
			return idx
		}
	}
	return notFound
}

func VolumeMountWithNameExists(mounts []corev1.VolumeMount, volumeName string) bool {
	for _, mount := range mounts {
		if mount.Name == volumeName {
			return true
		}
	}
	return false
}

// ResetUpdateStrategy resets the statefulset update strategy to RollingUpdate.
// If a version change is in progress, it doesn't do anything.
func ResetUpdateStrategy(ocb annotations.Versioned, kubeClient GetUpdater) error {
	if !ocb.IsChangingVersion() {
		return nil
	}

	// if we changed the version, we need to reset the UpdatePolicy back to OnUpdate
	_, err := GetAndUpdate(kubeClient, ocb.NamespacedName(), func(dep *appsv1.Deployment) {
		//dep.Spec.UpdateStrategy.Type = appsv1.RollingUpdateDeploymentStrategyType
		dep.Spec.Strategy.Type = appsv1.RollingUpdateDeploymentStrategyType

	})
	return err
}
