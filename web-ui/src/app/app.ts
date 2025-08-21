import { Component } from '@angular/core';
import { RouterOutlet } from '@angular/router';
import { FontAwesomeModule } from '@fortawesome/angular-fontawesome';
import { 
  faCalendar,
  faCircleCheck,
  faCircleXmark,
  faGear, 
  faHourglassClock,
  faHouse,
  faTags,
  faTruckRampBox,
  faUser
} from '@fortawesome/pro-regular-svg-icons'

@Component({
  selector: 'app-root',
  imports: [
    RouterOutlet,
    FontAwesomeModule,
  ],
  templateUrl: './app.html',
  styleUrl: './app.css'
})
export class App {
  protected title = 'web-ui';
  faCalendar = faCalendar;
  faCircleCheck = faCircleCheck;
  faCircleXmark = faCircleXmark;
  faGear = faGear;
  faHourglassClock = faHourglassClock;
  faHouse = faHouse;
  faTags = faTags;
  faTruckRampBox = faTruckRampBox;
  faUser = faUser;
}
