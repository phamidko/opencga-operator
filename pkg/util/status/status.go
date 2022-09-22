package status

import (
	"context"

	opencgav1 "github.com/phamidko/opencga-operator/api/v1"

	"sigs.k8s.io/controller-runtime/pkg/client"
	"sigs.k8s.io/controller-runtime/pkg/reconcile"
)

type Option interface {
	ApplyOption(ocb *opencgav1.OpenCGACommunity)
	GetResult() (reconcile.Result, error)
}

type OptionBuilder interface {
	GetOptions() []Option
}

// Update takes the options provided by the given option builder, applies them all and then updates the resource
func Update(statusWriter client.StatusWriter, ocb *opencgav1.OpenCGACommunity, optionBuilder OptionBuilder) (reconcile.Result, error) {
	options := optionBuilder.GetOptions()
	for _, opt := range options {
		opt.ApplyOption(ocb)
	}

	if err := statusWriter.Update(context.TODO(), ocb); err != nil {
		return reconcile.Result{}, err
	}

	return determineReconciliationResult(options)
}

func determineReconciliationResult(options []Option) (reconcile.Result, error) {
	// if there are any errors in any of our options, we return those first
	for _, opt := range options {
		res, err := opt.GetResult()
		if err != nil {
			return res, err
		}
	}
	// otherwise we might need to re-queue
	for _, opt := range options {
		res, _ := opt.GetResult()
		if res.Requeue || res.RequeueAfter > 0 {
			return res, nil
		}
	}
	// it was a successful reconciliation, nothing to do
	return reconcile.Result{}, nil
}
