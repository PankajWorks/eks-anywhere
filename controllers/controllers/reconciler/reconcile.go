package reconciler

import (
	"context"

	"github.com/pkg/errors"
	"sigs.k8s.io/controller-runtime/pkg/client"
)

const fieldManager = "eks-a-controller"

func ReconcileYaml(ctx context.Context, c client.Client, yaml []byte) error {
	objs, err := YamlToClientObjects(yaml)
	if err != nil {
		return err
	}

	return ReconcileObjects(ctx, c, objs)
}

func ReconcileObjects(ctx context.Context, c client.Client, objs []client.Object) error {
	for _, o := range objs {
		if err := ReconcileObject(ctx, c, o); err != nil {
			return err
		}
	}

	return nil
}

func ReconcileObject(ctx context.Context, c client.Client, obj client.Object) error {
	// Server side apply
	err := c.Patch(ctx, obj, client.Apply, &client.PatchOptions{FieldManager: fieldManager})
	if err != nil {
		return errors.Wrapf(err, "failed to reconcile object %s, %s/%s", obj.GetObjectKind().GroupVersionKind(), obj.GetNamespace(), obj.GetName())
	}

	return nil
}
