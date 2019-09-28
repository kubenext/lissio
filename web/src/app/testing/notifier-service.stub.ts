// Copyright (c) 2019 VMware, Inc. All Rights Reserved.
// SPDX-License-Identifier: Apache-2.0
//

export function notifierServiceStubFactory() {
  return {
    notifierSessionStub: jasmine.createSpyObj([
      'removeAllSignals',
      'pushSignal',
    ]),
    createSession() {
      return this.notifierSessionStub;
    },
  };
}

export const notifierServiceStub = notifierServiceStubFactory();
