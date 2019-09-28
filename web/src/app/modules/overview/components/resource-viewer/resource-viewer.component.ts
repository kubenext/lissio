// Copyright (c) 2019 VMware, Inc. All Rights Reserved.
// SPDX-License-Identifier: Apache-2.0
//

import {
  AfterViewInit,
  Component,
  Input,
  OnChanges,
  SimpleChanges,
  ViewEncapsulation,
} from '@angular/core';
import _ from 'lodash';
import { Node, ResourceViewerView } from 'src/app/models/content';
import { ElementsDefinition, Stylesheet } from 'cytoscape';

const statusColorCodes = {
  ok: '#60b515',
  warning: '#f57600',
  error: '#e12200',
};

const edgeColorCode = '#003d79';

@Component({
  selector: 'app-view-resource-viewer',
  templateUrl: './resource-viewer.component.html',
  styleUrls: ['./resource-viewer.component.scss'],
  encapsulation: ViewEncapsulation.None,
})
export class ResourceViewerComponent implements OnChanges, AfterViewInit {
  @Input() view: ResourceViewerView;

  currentView: ResourceViewerView;
  selected: string;
  selectedNode: Node;

  layout = {
    name: 'dagre',
    padding: 30,
    rankDir: 'TB',
    directed: true,
    animate: false,
  };

  zoom = {
    min: 0.5,
    max: 2.0,
  };

  style: Stylesheet[] = [
    {
      selector: 'node',
      css: {
        shape: 'rectangle',
        width: 'label',
        height: 'label',
        content: 'data(name)',
        'background-color': 'data(colorCode)',
        color: '#fff',
        'font-size': 12,
        'text-wrap': 'wrap',
        'text-valign': 'center',
        'padding-left': '10px',
        'padding-right': '10px',
        'padding-top': '10px',
        'padding-bottom': '10px',
      },
    },

    {
      selector: 'node:selected',
      css: {
        'curve-style': 'bezier',
        'line-color': 'data(colorCode)',
        'source-arrow-color': 'data(colorCode)',
        'target-arrow-color': 'data(colorCode)',
        'border-width': 1,
        'border-color': '#313131',
        'border-style': 'solid',
      },
    },

    {
      selector: 'edge',
      css: {
        'curve-style': 'bezier',
        opacity: 0.666,
        width: 'mapData(strength, 70, 100, 1, 3)',
        'line-color': 'data(colorCode)',
        'source-arrow-color': 'data(colorCode)',
        'target-arrow-color': 'data(colorCode)',
      },
    },
  ];

  graphData: ElementsDefinition;

  private afterFirstChange: boolean;

  constructor() {}

  ngAfterViewInit(): void {
    if (this.afterFirstChange) {
      this.graphData = this.generateGraphData();
    }
  }

  ngOnChanges(changes: SimpleChanges): void {
    const isEqual = _.isEqual(
      changes.view.currentValue,
      changes.view.previousValue
    );

    if (changes.view.isFirstChange() || !isEqual) {
      this.currentView = changes.view.currentValue as ResourceViewerView;
      this.select(this.currentView.config.selected);
      this.graphData = this.generateGraphData();
      this.afterFirstChange = true;
    }
  }

  nodeChange(event) {
    this.select(event.id);
  }

  generateGraphData() {
    return {
      nodes: this.nodes(),
      edges: this.edges(),
    };
  }

  private nodes() {
    if (!this.currentView || !this.currentView.config.nodes) {
      return [];
    }

    const nodes = Object.entries(this.currentView.config.nodes).map(
      ([name, details]) => {
        const colorCode =
          statusColorCodes[details.status] || statusColorCodes.error;

        return {
          data: {
            id: name,
            name: `${details.name}\n${details.apiVersion} ${details.kind}`,
            weight: 100,
            colorCode,
          },
        };
      }
    );

    return Array.prototype.concat(...nodes);
  }

  private edges() {
    if (!this.currentView || !this.currentView.config.edges) {
      return [];
    }

    const edges = Object.entries(this.currentView.config.edges).map(
      ([parent, maps]) => {
        return maps.map(edge => {
          return {
            data: {
              source: parent,
              target: edge.node,
              colorCode: edgeColorCode,
              strength: 10,
            },
          };
        });
      }
    );

    return Array.prototype.concat(...edges);
  }

  private select(id: string) {
    this.selected = id;

    const nodes = this.currentView.config.nodes;

    if (nodes && nodes[id]) {
      this.selectedNode = nodes[id];
    }
  }
}
