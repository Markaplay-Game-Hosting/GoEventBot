import { Component } from '@angular/core';
import { TagService } from '../../../services/tags/tags';
import { Tag } from '../../../models/pages/tags.model';

@Component({
  selector: 'app-tags',
  imports: [],
  templateUrl: './tags.html',
  styleUrl: './tags.css'
})
export class Tags {
  data: Tag[] = [];
  
  constructor(private tagsService: TagService) {}

  ngOnInit() {
    this.tagsService.getTags().subscribe({
      next: (data) => {
        this.data = data.tags;
      },
      error: (err) => {
        console.error('Error fetching tags:', err);
      }
    });
  }
}
