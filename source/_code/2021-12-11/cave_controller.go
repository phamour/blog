caveProtectionFinalizer := "abc/cave-protection"

// IN RECONCILATION

// Some reconcilation logic in before...

if cave.ObjectMeta.DeletionTimestamp.IsZero() {
	// Cave is not being deleted:
	// append finalizer to protect cave from being deleted before backend & hpo
	// otherwise backend & hpo might fail to unmount cave
	// which make kubelet hang
	if !containsFinalizer(cave.GetFinalizers(), caveProtectionFinalizer) {
		controllerutil.AddFinalizer(&cave, caveProtectionFinalizer)
		if err := r.Update(ctx, &cave); err != nil {
			reqLogger.Error(err, "unable to add finalizer on cave resource", "cave", cave.Name)
		}
		return ctrl.Result{Requeue: true}, nil
	}
} else {
	reqLogger.Info("got deletion request", "cave", cave.Name)
	// Cave is being deleted:
	// the deletion should be pending until backend pods are terminated
	if containsFinalizer(cave.GetFinalizers(), caveProtectionFinalizer) {
		// List backend pods
		var backendPods v1.PodList
		if err := r.List(ctx, &backendPods, client.InNamespace(req.Namespace), client.MatchingLabels{"app": "anylearn-backend"}); err != nil {
			reqLogger.Error(err, "unable to list backend pods")
			return ctrl.Result{}, err
		}
		// Count pods using current cave
		nbOccupants := 0
		for _, pod := range backendPods {
			for _, v := range pod.Spec.Volumes {
				if v.Name == cave.Name && v.NFS.Server == cave.Status.ServiceIP {
					nbOccupants++
					reqLogger.Info("Pod still using cave", pod.Name, cave.Name)
				}
			}
		}
		// Remove finalizer if backend pods are gone so that cave deletion can be truely proceeded
		if nbOccupants == 0 {
			controllerutil.RemoveFinalizer(&cave, caveProtectionFinalizer)
			if err := r.Update(ctx, &cave); err != nil {
				reqLogger.Error(err, "unable to remove protection finalizer from cave resource", "cave", cave.Name)
				return ctrl.Result{Requeue: true}, err
			}
			// End of reconcilation
			return ctrl.Result{}, nil
		} else {
			// Still got some pods using this cave, requeue the reconcilation
			return ctrl.Result{Requeue: true, RequeueAfter: 20 * time.Second}, nil
		}
	}
}

// Some reconcilation logic in after...
