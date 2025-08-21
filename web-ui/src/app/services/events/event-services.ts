import { HttpClient } from '@angular/common/http';
import { Injectable } from '@angular/core';
import { Router } from '@angular/router';
import { Observable } from 'rxjs';
import { EventList } from '../../models/pages/events.model';

@Injectable({
  providedIn: 'root'
})
export class EventServices {
  constructor(private http: HttpClient, private router: Router) {}
  getEvents(): Observable<EventList> {
    return this.http.get<EventList>('/v1/events');
  }
}
