---
weight: 15
---

# Terminology

The dashboard plugin documentation contains a few terms that may be uncommon. Familiarity with terminology used in Kubernetes will be helpful. This page is used to define terms needed to get started with writing a plugin.

 * View components are reusable building blocks to serve content in various formats on the dashboard. Available components can be found in `pkg/view/component`. Additional components can be added through a plugin.

 * Links are a type of component used to reference a Kubernetes object through text.

 * GVK stands for Group, Version, Kind. GVK is used to interact with the Kubernetes API server in order to identify an object type.

 * Resource viewer is a graph containing nodes repesenting interactions between objects in a Kubernetes cluster. Nodes provide additional information about the status as well as additional components used to describe the object.

 * Flexlayout is a layout that builds on top of [flexbox](https://github.com/angular/flex-layout). Each layout contains sections that can contain one more view components. Contents within a flex layout will shrink or expand depending on the space available.

 * Workloads are rules used to describe how a pod should be run. Example of these are Deployments, StatefulSets, or DaemonSets.
