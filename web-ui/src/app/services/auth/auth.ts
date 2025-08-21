import { HttpClient, HttpHeaders } from '@angular/common/http';
import { Injectable } from '@angular/core';
import { Router } from '@angular/router';
import { Observable } from 'rxjs';

@Injectable({
  providedIn: 'root'
})
export class Auth {
  constructor(private http: HttpClient, private router: Router) {}
  login(): Observable<Object> {
    return this.http.get('/oauth/authenticate');
  }

  isAuthenticated(): Observable<boolean> {
    const headers = new HttpHeaders({
      'Content-Type': 'application/json'
    });

    return this.http.get<boolean>('/v1/events', {
      headers,
      withCredentials: true
    });
  }
}
