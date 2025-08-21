import { TestBed } from '@angular/core/testing';

import { Discord } from './discord';

describe('Discord', () => {
  let service: Discord;

  beforeEach(() => {
    TestBed.configureTestingModule({});
    service = TestBed.inject(Discord);
  });

  it('should be created', () => {
    expect(service).toBeTruthy();
  });
});
