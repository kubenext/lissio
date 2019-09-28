// Copyright (c) 2019 VMware, Inc. All Rights Reserved.
// SPDX-License-Identifier: Apache-2.0
//

import { Component, Input, OnChanges, SimpleChanges } from '@angular/core';
import { ErrorView } from 'src/app/models/content';

@Component({
  selector: 'app-view-error',
  templateUrl: './error.component.html',
  styleUrls: ['./error.component.scss'],
})
export class ErrorComponent implements OnChanges {
  @Input() view: ErrorView;

  source: string;

  constructor() {}

  ngOnChanges(changes: SimpleChanges): void {
    if (changes.view.currentValue) {
      const view = changes.view.currentValue as ErrorView;
      this.source = view.config.data;
    }
  }
}
