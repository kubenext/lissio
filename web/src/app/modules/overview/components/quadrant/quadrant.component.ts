// Copyright (c) 2019 VMware, Inc. All Rights Reserved.
// SPDX-License-Identifier: Apache-2.0
//

import { Component, Input, OnChanges, SimpleChanges } from '@angular/core';
import { QuadrantValue, QuadrantView } from 'src/app/models/content';
import { ViewService } from '../../services/view/view.service';

const emptyQuadrantValue = { value: '', label: '' };

@Component({
  selector: 'app-view-quadrant',
  templateUrl: './quadrant.component.html',
  styleUrls: ['./quadrant.component.scss'],
})
export class QuadrantComponent implements OnChanges {
  @Input() view: QuadrantView;

  title: string;
  nw: QuadrantValue = emptyQuadrantValue;
  ne: QuadrantValue = emptyQuadrantValue;
  sw: QuadrantValue = emptyQuadrantValue;
  se: QuadrantValue = emptyQuadrantValue;

  constructor(private viewService: ViewService) {}

  ngOnChanges(changes: SimpleChanges): void {
    if (changes.view.currentValue) {
      const view = changes.view.currentValue as QuadrantView;
      this.title = this.viewService.viewTitleAsText(view);
      this.nw = view.config.nw;
      this.ne = view.config.ne;
      this.sw = view.config.sw;
      this.se = view.config.se;
    }
  }
}
