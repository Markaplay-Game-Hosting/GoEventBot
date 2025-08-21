import { faCircleCheck, faCircleXmark, faGalaxy, faHourglassClock, IconDefinition } from '@fortawesome/pro-regular-svg-icons';

export interface Job {
  id: string;
  event_id: string;
  execution_date: Date;
  status: JobStatus;
}

export interface JobList {
  jobs: Job[];
}

export enum JobStatus {
  Unknown,
  Pending,
  Running,
  Completed,
  Failed,
  Cancel,
}

export const JobStatusIcons: Record<JobStatus, IconDefinition> = {
  [JobStatus.Unknown]: faGalaxy,
  [JobStatus.Pending]: faHourglassClock,
  [JobStatus.Running]: faCircleCheck,
  [JobStatus.Completed]: faCircleCheck,
  [JobStatus.Failed]: faCircleXmark,
  [JobStatus.Cancel]: faCircleXmark,
};
