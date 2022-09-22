package predicates

import (
	"reflect"

	opencgav1 "github.com/phamidko/opencga-operator/api/v1"
	"sigs.k8s.io/controller-runtime/pkg/event"
	"sigs.k8s.io/controller-runtime/pkg/predicate"
)

// OnlyOnSpecChange returns a set of predicates indicating
// that reconciliations should only happen on changes to the Spec of the resource.
// any other changes won't trigger a reconciliation. This allows us to freely update the annotations
// of the resource without triggering unintentional reconciliations.
func OnlyOnSpecChange() predicate.Funcs {
	return predicate.Funcs{
		UpdateFunc: func(e event.UpdateEvent) bool {
			oldResource := e.ObjectOld.(*opencgav1.OpenCGACommunity)
			newResource := e.ObjectNew.(*opencgav1.OpenCGACommunity)
			specChanged := !reflect.DeepEqual(oldResource.Spec, newResource.Spec)
			return specChanged
		},
	}
}
