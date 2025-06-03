// Copyright 2020-2023 Project Capsule Authors.
// SPDX-License-Identifier: Apache-2.0

package controllers

import (
	"context"
	"fmt"

	"github.com/pkg/errors"
	k8serrors "k8s.io/apimachinery/pkg/api/errors"
	"k8s.io/apimachinery/pkg/types"
	"k8s.io/apimachinery/pkg/util/sets"
	"sigs.k8s.io/controller-runtime/pkg/client"
	"sigs.k8s.io/controller-runtime/pkg/reconcile"

	capsulev1beta2 "github.com/projectcapsule/capsule/api/v1beta2"
)

type CapsuleConfiguration struct {
	Client                      client.Client
	CapsuleConfigurationName    string
	DeprecatedCapsuleUserGroups []string
}

//nolint:gochecknoglobals
var CapsuleUserGroups sets.Set[string]

func (c *CapsuleConfiguration) Start(ctx context.Context) error {
	if len(c.DeprecatedCapsuleUserGroups) > 0 {
		CapsuleUserGroups = sets.New[string](c.DeprecatedCapsuleUserGroups...)

		return nil
	}

	capsuleConfig := &capsulev1beta2.CapsuleConfiguration{}
	if err := c.Client.Get(ctx, types.NamespacedName{Name: c.CapsuleConfigurationName}, capsuleConfig); err != nil {
		if k8serrors.IsNotFound(err) {
			return fmt.Errorf("CapsuleConfiguration %s does not exist", c.CapsuleConfigurationName)
		}

		return errors.Wrap(err, "unable to retrieve CapsuleConfiguration")
	}

	CapsuleUserGroups = sets.New(capsuleConfig.Spec.UserGroups...)

	return nil
}

// Reconcile is kept for legacy or possible future controller-based use cases,
// but should not be wired into a controller manager if running standalone.
func (c *CapsuleConfiguration) Reconcile(ctx context.Context, request reconcile.Request) (reconcile.Result, error) {
	capsuleConfig := &capsulev1beta2.CapsuleConfiguration{}
	if err := c.Client.Get(ctx, types.NamespacedName{Name: request.Name}, capsuleConfig); err != nil {
		panic(err)
	}

	CapsuleUserGroups = sets.New(capsuleConfig.Spec.UserGroups...)

	return reconcile.Result{}, nil
}
