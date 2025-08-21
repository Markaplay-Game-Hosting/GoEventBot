import { HttpClient } from '@angular/common/http';
import { Injectable } from '@angular/core';
import { Router } from '@angular/router';
import { Observable } from 'rxjs';
import { TagList } from '../../models/pages/tags.model';

@Injectable({
  providedIn: 'root'
})
export class TagService {
  constructor(private http: HttpClient, private router: Router) {}
  getTags(): Observable<TagList> {
    return this.http.get<TagList>('/v1/jobs');
  }
}
