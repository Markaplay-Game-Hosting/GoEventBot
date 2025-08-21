import { TestBed } from '@angular/core/testing';

import { Tags } from './tags';

describe('Tags', () => {
  let service: Tags;

  beforeEach(() => {
    TestBed.configureTestingModule({});
    service = TestBed.inject(Tags);
  });

  it('should be created', () => {
    expect(service).toBeTruthy();
  });
});
