# Octant Roadmap

The Octant team is focused on improving and stabilizing Octant with the goal of reaching 1.0 stability.

Octant is a web based developer tool that offers introspection into a Kubernetes cluster on which their applications and workloads reside. It provides a dashboard that shows how the cluster is configured, its components, and its resources – helping a developer better understand what is going on inside of a cluster. The goals for this project are to provide enough functionality to be a graphical analog for kubectl, to be part of the must-install tooling all developers who use Kubernetes have on their workstations, and to make Kubernetes more accessible by lowering the barriers around learning.

## Working with your Kubernetes workloads

### View and Edit Kubernetes objects

Octant currently can visualize many Kubernetes types. In order to round out the offering, additional support will be added for Custom Resource Definitions, Nodes, Roles, and Role Bindings.

### Upgrading the "kubectl" experience

Octant aims to be a companion to `kubectl` rather than a replacement. There are powerful capabilities buried inside of `kubectl` that can not simply be replaced with a graphical view. With this in mind, Octant will attempt to embed the `kubectl` command line experience inside of Octant. 

### Providing additional information about objects

Since Octant has access to the OpenAPI schemas that Kubernetes is built upon, it is possible make additional help available when applicable. When viewing an object, such as a PodTemplate, Octant will provide help in context using the descriptions the Kubernetes authors have provided.

## Answering the question, "why isn't my workload working?"

In the Go client libraries for Kubernetes, there are multiple `Condition` objects. With the goal of unifying conditions and making them easier to work with, Octant will provide a _duck typed_ condition system that will allow describing the status of any object or group of objects in Kubernetes.

### YAML editor

At times, it is useful to have direct access to YAML source for a Kubernetes object. To satisfy this need, Octant will provide an object-aware YAML editor to ensure that users know what can be modified and how to do it in a safe manner.

### Graphical views of resources

Octant provides an object graph view that allows users to see how workloads are constructed using multiple objects. To provide more fidelity about a workload's health down to the pod level, Octant will include a heat map to help user's understand where problems are happening in pods.

## Extending Octant to work with any Kubernetes type

### Plugin API

Octant comes with a dashboard that allows users to visual objects in their cluster, but it isn't possible for the Octant team to provide detailed views for all Kubernetes types due to the extensibility of the Kubernetes API. With that in mind, Octant provides a plugin API to extend existing views and add additional pages and sections.

Using the plugin API, it will be possible to customize Octant to your Kubernetes needs.

Plugins can be created in any language that speaks GRPC. To help users getting started, more examples would need to be created.

### More components

Currently, Octant contains a library of components used to build out the detailed object information panes. To enable more rich and information packed views, additional components will be provided. Specifically, buttons, modals, spinners, alerts, drop down menus, additional navigation options, tooltips, and myriad graphing options are on the list.

## Providing better defaults

There are multiple ways to describe an application or workload. Octant aims to automatically group objects together that define a workload. This view will help users understand the scope of their workloads and quickly identify components that are causing issues. The project will use [common labels](https://kubernetes.io/docs/concepts/overview/working-with-objects/common-labels/) as a starting point.

Octant should also know about prevalent Kubernetes applications, such as Helm, and use their metadata to help group objects into a workload. 

## Being a better ecosystem participant

### Speeding up Octant

Collating the data that makes up Octant views is a complex task since the Kubernetes API is not specifically tailored to applications which want to query large amounts of specific data. Instead, Octant relies on Kubernetes' informers to provide real time data updates. Populating the caches the informers utilize can be further optimized which will in part speed up operations in Octant, like switching namespaces and viewing resource graphs. 

### Graduating to an application

Octant is a local web application that connects to Kubernetes clusters to provide detailed views about Kubernetes objects. If Octant could be free of of the browser, it could take advantage of items like system alerts, task and status bars provided by Linux, macOS and Windows. 

### Managing local clusters

It should be possible to allow Octant to build and manage local developer clusters.
