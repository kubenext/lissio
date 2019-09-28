// Copyright (c) 2019 VMware, Inc. All Rights Reserved.
// SPDX-License-Identifier: Apache-2.0
//

import { Component, Input, OnChanges, SimpleChanges } from '@angular/core';
import { ContainerDef, ContainersView } from 'src/app/models/content';

@Component({
  selector: 'app-view-containers',
  templateUrl: './containers.component.html',
  styleUrls: ['./containers.component.scss'],
})
export class ContainersComponent implements OnChanges {
  @Input() view: ContainersView;
  containers: ContainerDef[];

  constructor() {}

  ngOnChanges(changes: SimpleChanges): void {
    if (changes.view.currentValue) {
      const view = changes.view.currentValue as ContainersView;
      this.containers = view.config.containers;
    }
  }

  trackItem(index: number, item: ContainerDef): string {
    return item.name;
  }
}
