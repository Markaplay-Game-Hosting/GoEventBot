import { Component } from '@angular/core';
import { FontAwesomeModule } from '@fortawesome/angular-fontawesome';
import { Job, JobStatus, JobStatusIcons } from '../../../models/pages/jobs.model';
import { JobServices } from '../../../services/jobs/jobs';
import { faCircleCheck, faCircleXmark, faGalaxy, faHourglassClock, IconDefinition } from '@fortawesome/pro-regular-svg-icons';

@Component({
  selector: 'app-jobs',
  imports: [FontAwesomeModule],
  templateUrl: './jobs.html',
  styleUrl: './jobs.css'
})
export class Jobs {
  data: Job[] = [];
  constructor(private jobServices: JobServices) {}
  faCircleCheck = faCircleCheck;
  faCircleXmark = faCircleXmark;
  faGalaxy = faGalaxy;
  faHourglassClock = faHourglassClock;

  ngOnInit() {
    this.jobServices.getJobs().subscribe({
      next: (data) => {
        this.data = data.jobs;
      },
      error: (err) => {
        console.error('Error fetching events:', err);
      }
    });
  }

  getIconName(status: JobStatus): IconDefinition {
    const iconName = JobStatusIcons[status];
    return iconName;
  }

  isIconStatus(status: JobStatus): boolean {
    const iconStatuses = [
      JobStatus.Running,
      JobStatus.Pending
    ];
    return !iconStatuses.includes(status);
  }

}
