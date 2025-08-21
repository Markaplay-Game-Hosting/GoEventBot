import { Component, EventEmitter } from '@angular/core';
import { RouterOutlet } from '@angular/router';
import { FontAwesomeModule } from '@fortawesome/angular-fontawesome';
import { 
  faCalendar,
  faGear, 
  faHourglassClock,
  faHouse,
  faTags,
  faTruckRampBox,
  faUser
} from '@fortawesome/pro-regular-svg-icons'

@Component({
  selector: 'app-main-layout',
  imports: [RouterOutlet,FontAwesomeModule],
  templateUrl: './main-layout.html',
  styleUrl: './main-layout.css'
})
export class MainLayout {
  faCalendar = faCalendar;
  faGear = faGear;
  faHourglassClock = faHourglassClock;
  faHouse = faHouse;
  faTags = faTags;
  faTruckRampBox = faTruckRampBox;
  faUser = faUser;

}
