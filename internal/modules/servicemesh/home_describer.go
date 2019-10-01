/*
 * Copyright (c) 2019 Kubenext, Inc. All Rights Reserved.
 * SPDX-License-Identifier: Apache-2.0
 */

package servicemesh

import (
	"context"
	"github.com/kubenext/lissio/internal/describer"
	"github.com/kubenext/lissio/pkg/view/component"
	"github.com/pkg/errors"
)

// HomeDescriberOption is an option for configuring HomeDescriber.
type HomeDescriberOption func(d *HomeDescriber)

// HomeDescriber describes content for servicemesh.
type HomeDescriber struct {
	summarizer Summarizer
}

var _ describer.Describer = (*HomeDescriber)(nil)

func NewHomeDescriber(options ...HomeDescriberOption) *HomeDescriber {
	d := &HomeDescriber{}

	for _, option := range options {
		option(d)
	}

	if d.summarizer == nil {
		d.summarizer = &summarizer{}
	}

	return d
}

func (h *HomeDescriber) Describe(ctx context.Context, namespace string, options describer.Options) (component.ContentResponse, error) {
	table, err := h.summarizer.Summarize(ctx, namespace, options)
	if err != nil {
		return component.EmptyContentResponse, errors.Wrap(err, "sumarize servicemesh")
	}

	contentResponse := component.ContentResponse{
		Title:      component.TitleFromString("Service Mesh"),
		Components: []component.Component{table},
		IconName:   "",
		IconSource: "",
	}

	return contentResponse, nil
}

func (h *HomeDescriber) PathFilters() []describer.PathFilter {
	return []describer.PathFilter{
		*describer.NewPathFilter("/", h),
	}
}

func (h *HomeDescriber) Reset(ctx context.Context) error {
	return nil
}
