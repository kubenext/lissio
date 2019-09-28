/*
Copyright (c) 2019 VMware, Inc. All Rights Reserved.
SPDX-License-Identifier: Apache-2.0
*/

package printer

import (
	"context"
	"testing"
	"time"

	"github.com/golang/mock/gomock"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
	corev1 "k8s.io/api/core/v1"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	"k8s.io/apimachinery/pkg/apis/meta/v1/unstructured"

	"github.com/vmware/octant/pkg/store"
	storefake "github.com/vmware/octant/pkg/store/fake"
	"github.com/vmware/octant/pkg/view/component"
)

func Test_EventListHandler(t *testing.T) {
	controller := gomock.NewController(t)
	defer controller.Finish()

	tpo := newTestPrinterOptions(controller)
	printOptions := tpo.ToOptions()

	object := &corev1.EventList{
		Items: []corev1.Event{
			{
				ObjectMeta: metav1.ObjectMeta{
					Name:      "event-2",
					Namespace: "default",
				},
				TypeMeta: metav1.TypeMeta{
					APIVersion: "v1",
					Kind:       "Event",
				},
				InvolvedObject: corev1.ObjectReference{
					APIVersion: "apps/v1",
					Kind:       "Deployment",
					Name:       "d2",
					Namespace:  "default",
				},
				Count:          1234,
				Message:        "message",
				Reason:         "Reason",
				Type:           "Type",
				FirstTimestamp: metav1.Time{Time: time.Unix(1548424410, 0)},
				LastTimestamp:  metav1.Time{Time: time.Unix(1548424410, 0)},
			},
			{
				ObjectMeta: metav1.ObjectMeta{
					Name:      "event-1",
					Namespace: "default",
				},
				TypeMeta: metav1.TypeMeta{
					APIVersion: "v1",
					Kind:       "Event",
				},
				InvolvedObject: corev1.ObjectReference{
					APIVersion: "apps/v1",
					Kind:       "Deployment",
					Name:       "d1",
					Namespace:  "default",
				},
				Count:          1234,
				Message:        "message",
				Reason:         "Reason",
				Type:           "Type",
				FirstTimestamp: metav1.Time{Time: time.Unix(1548424410, 0)},
				LastTimestamp:  metav1.Time{Time: time.Unix(1548424420, 0)},
			},
		},
	}

	tpo.PathForObject(&object.Items[0], "message", "/event1")
	tpo.PathForObject(&object.Items[1], "message", "/event2")

	ctx := context.Background()
	got, err := EventListHandler(ctx, object, printOptions)
	require.NoError(t, err)

	cols := component.NewTableCols("Kind", "Message", "Reason", "Type",
		"First Seen", "Last Seen")
	expected := component.NewTableWithRows("Events", "We couldn't find any events!", cols, []component.TableRow{
		{
			"Kind":       component.NewLink("", "d1 (1234)", "/overview/namespace/default/workloads/deployments/d1"),
			"Message":    component.NewLink("", "message", "/event2"),
			"Reason":     component.NewText("Reason"),
			"Type":       component.NewText("Type"),
			"First Seen": component.NewTimestamp(time.Unix(1548424410, 0)),
			"Last Seen":  component.NewTimestamp(time.Unix(1548424420, 0)),
		},
		{
			"Kind":       component.NewLink("", "d2 (1234)", "/overview/namespace/default/workloads/deployments/d2"),
			"Message":    component.NewLink("", "message", "/event1"),
			"Reason":     component.NewText("Reason"),
			"Type":       component.NewText("Type"),
			"First Seen": component.NewTimestamp(time.Unix(1548424410, 0)),
			"Last Seen":  component.NewTimestamp(time.Unix(1548424410, 0)),
		},
	})

	component.AssertEqual(t, expected, got)
}

func Test_ReplicaSetEvents(t *testing.T) {
	controller := gomock.NewController(t)
	defer controller.Finish()

	tpo := newTestPrinterOptions(controller)
	printOptions := tpo.ToOptions()

	now := time.Unix(1547211430, 0)

	object := &corev1.EventList{
		Items: []corev1.Event{
			{
				ObjectMeta: metav1.ObjectMeta{
					Name: "frontend",
					CreationTimestamp: metav1.Time{
						Time: now,
					},
				},
				Count:  1,
				Type:   corev1.EventTypeNormal,
				Reason: "SuccessfulCreate",
				Source: corev1.EventSource{
					Component: "replicaset-controller",
				},
				Message:        "Created pod: frontend-97k6z",
				FirstTimestamp: metav1.Time{Time: time.Unix(1548424410, 0)},
				LastTimestamp:  metav1.Time{Time: time.Unix(1548424410, 0)},
			},
			{
				ObjectMeta: metav1.ObjectMeta{
					Name: "frontend",
					CreationTimestamp: metav1.Time{
						Time: now,
					},
				},
				Count:  1,
				Type:   corev1.EventTypeNormal,
				Reason: "SuccessfulCreate",
				Source: corev1.EventSource{
					Component: "replicaset-controller",
				},
				Message:        "Created pod: frontend-8n77p",
				FirstTimestamp: metav1.Time{Time: time.Unix(1548424410, 0)},
				LastTimestamp:  metav1.Time{Time: time.Unix(1548424410, 0)},
			},
			{
				ObjectMeta: metav1.ObjectMeta{
					Name: "frontend",
					CreationTimestamp: metav1.Time{
						Time: now,
					},
				},
				Count:  1,
				Type:   corev1.EventTypeNormal,
				Reason: "SuccessfulCreate",
				Source: corev1.EventSource{
					Component: "replicaset-controller",
				},
				Message:        "Created pod: frontend-b7fxf",
				FirstTimestamp: metav1.Time{Time: time.Unix(1548424410, 0)},
				LastTimestamp:  metav1.Time{Time: time.Unix(1548424410, 0)},
			},
		},
	}

	got, err := PrintEvents(object, printOptions)
	require.NoError(t, err)

	expected := component.NewTable("Events", "There are no events!", objectEventCols)

	expected.Add(component.TableRow{
		"Message":    component.NewText("Created pod: frontend-97k6z"),
		"Reason":     component.NewText("SuccessfulCreate"),
		"Type":       component.NewText("Normal"),
		"First Seen": component.NewTimestamp(time.Unix(1548424410, 0)),
		"Last Seen":  component.NewTimestamp(time.Unix(1548424410, 0)),
		"From":       component.NewText("replicaset-controller"),
		"Count":      component.NewText("1"),
	})

	expected.Add(component.TableRow{
		"Message":    component.NewText("Created pod: frontend-8n77p"),
		"Reason":     component.NewText("SuccessfulCreate"),
		"Type":       component.NewText("Normal"),
		"First Seen": component.NewTimestamp(time.Unix(1548424410, 0)),
		"Last Seen":  component.NewTimestamp(time.Unix(1548424410, 0)),
		"From":       component.NewText("replicaset-controller"),
		"Count":      component.NewText("1"),
	})

	expected.Add(component.TableRow{
		"Message":    component.NewText("Created pod: frontend-b7fxf"),
		"Reason":     component.NewText("SuccessfulCreate"),
		"Type":       component.NewText("Normal"),
		"First Seen": component.NewTimestamp(time.Unix(1548424410, 0)),
		"Last Seen":  component.NewTimestamp(time.Unix(1548424410, 0)),
		"From":       component.NewText("replicaset-controller"),
		"Count":      component.NewText("1"),
	})

	component.AssertEqual(t, expected, got)
}

func Test_EventHandler(t *testing.T) {
	controller := gomock.NewController(t)
	defer controller.Finish()

	tpo := newTestPrinterOptions(controller)
	printOptions := tpo.ToOptions()

	event := &corev1.Event{
		ObjectMeta: metav1.ObjectMeta{
			Name: "event-12345",
		},
		TypeMeta: metav1.TypeMeta{
			APIVersion: "v1",
			Kind:       "Event",
		},
		InvolvedObject: corev1.ObjectReference{
			APIVersion: "apps/v1",
			Kind:       "Deployment",
			Name:       "d1",
		},
		Count:          1234,
		Message:        "message",
		Reason:         "Reason",
		Type:           corev1.EventTypeNormal,
		FirstTimestamp: metav1.Time{Time: time.Unix(1548424410, 0)},
		LastTimestamp:  metav1.Time{Time: time.Unix(1548424410, 0)},
		Source: corev1.EventSource{
			Component: "component",
			Host:      "host",
		},
	}

	ctx := context.Background()
	got, err := EventHandler(ctx, event, printOptions)
	require.NoError(t, err)

	eventDetailSections := []component.SummarySection{
		{
			Header:  "Last Seen",
			Content: component.NewTimestamp(time.Unix(1548424410, 0)),
		},
		{
			Header:  "First Seen",
			Content: component.NewTimestamp(time.Unix(1548424410, 0)),
		},
		{
			Header:  "Count",
			Content: component.NewText("1234"),
		},
		{
			Header:  "Message",
			Content: component.NewText("message"),
		},
		{
			Header:  "Kind",
			Content: component.NewText("Deployment"),
		},
		{
			Header:  "Involved Object",
			Content: component.NewLink("", "d1", "/overview/workloads/deployments/d1"),
		},
		{
			Header:  "Type",
			Content: component.NewText("Normal"),
		},
		{
			Header:  "Reason",
			Content: component.NewText("Reason"),
		},
		{
			Header:  "Source",
			Content: component.NewText("component on host"),
		},
	}
	eventDetailView := component.NewSummary("Event Detail", eventDetailSections...)

	expected := component.NewFlexLayout("Event")
	expected.AddSections([]component.FlexLayoutSection{
		{
			{Width: component.WidthFull, View: eventDetailView},
		},
	}...)

	component.AssertEqual(t, expected, got)
}

func Test_eventsForObject(t *testing.T) {
	object := &corev1.Pod{
		TypeMeta: metav1.TypeMeta{
			APIVersion: "v1",
			Kind:       "Pod",
		},
		ObjectMeta: metav1.ObjectMeta{
			Name:      "pod",
			Namespace: "default",
		},
	}

	controller := gomock.NewController(t)
	defer controller.Finish()

	o := storefake.NewMockStore(controller)
	key := store.Key{
		Namespace:  "default",
		APIVersion: "v1",
		Kind:       "Event",
	}

	events := &unstructured.UnstructuredList{}
	events.Items = append(events.Items, []unstructured.Unstructured{
		{
			Object: map[string]interface{}{
				"involvedObject": map[string]interface{}{
					"namespace":  "default",
					"apiVersion": "v1",
					"kind":       "Pod",
					"name":       "pod",
				},
				"message": "pod",
			},
		},
		{
			Object: map[string]interface{}{
				"involvedObject": map[string]interface{}{
					"namespace":  "default",
					"apiVersion": "v1",
					"kind":       "Pod",
					"name":       "pod2",
				},
				"message": "pod2",
			},
		},
	}...)

	o.EXPECT().List(gomock.Any(), gomock.Eq(key)).Return(events, false, nil)

	ctx := context.Background()
	got, err := eventsForObject(ctx, object, o)
	require.NoError(t, err)

	expected := &corev1.EventList{
		Items: []corev1.Event{
			{
				InvolvedObject: corev1.ObjectReference{
					Namespace:  "default",
					APIVersion: "v1",
					Kind:       "Pod",
					Name:       "pod",
				},
				Message: "pod",
			},
		},
	}

	assert.Equal(t, expected.Items, got.Items)
}
