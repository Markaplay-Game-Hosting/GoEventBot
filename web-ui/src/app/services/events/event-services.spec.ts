import { TestBed } from '@angular/core/testing';

import { EventServices } from './event-services';

describe('EventServices', () => {
  let service: EventServices;

  beforeEach(() => {
    TestBed.configureTestingModule({});
    service = TestBed.inject(EventServices);
  });

  it('should be created', () => {
    expect(service).toBeTruthy();
  });
});
