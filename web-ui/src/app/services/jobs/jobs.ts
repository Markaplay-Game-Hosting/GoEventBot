import { HttpClient } from '@angular/common/http';
import { Injectable } from '@angular/core';
import { Router } from '@angular/router';
import { JobList } from '../../models/pages/jobs.model';
import { Observable } from 'rxjs';

@Injectable({
  providedIn: 'root'
})
export class JobServices {
  constructor(private http: HttpClient, private router: Router) {}
  getJobs(): Observable<JobList> {
    return this.http.get<JobList>('/v1/jobs');
  }
}
