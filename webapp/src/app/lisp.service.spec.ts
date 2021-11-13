import { TestBed } from '@angular/core/testing';

import { LispService } from './lisp.service';

describe('LispService', () => {
  let service: LispService;

  beforeEach(() => {
    TestBed.configureTestingModule({});
    service = TestBed.inject(LispService);
  });

  it('should be created', () => {
    expect(service).toBeTruthy();
  });
});
