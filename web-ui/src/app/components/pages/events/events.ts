import { Component } from '@angular/core';
import { FontAwesomeModule } from '@fortawesome/angular-fontawesome';
import { 
  faCircleCheck,
  faCircleXmark,
  faPencil
} from '@fortawesome/pro-regular-svg-icons'
import { EventServices } from '../../../services/events/event-services';
import { Event } from '../../../models/pages/events.model';

@Component({
  selector: 'app-events',
  imports: [FontAwesomeModule],
  templateUrl: './events.html',
  styleUrl: './events.css'
})
export class Events {
  data: Event[] = [];
  constructor(private eventsService: EventServices) {}
  faCircleCheck = faCircleCheck;
  faCircleXmark = faCircleXmark;
  faPencil = faPencil;

  ngOnInit() {
    this.eventsService.getEvents().subscribe({
      next: (data) => {
        this.data = data.events;
      },
      error: (err) => {
        console.error('Error fetching events:', err);
      }
    });
  }
}
