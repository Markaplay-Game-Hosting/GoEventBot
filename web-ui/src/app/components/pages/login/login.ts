import { Component, DestroyRef, inject } from '@angular/core';
import { Auth } from '../../../services/auth/auth';
import { takeUntilDestroyed } from '@angular/core/rxjs-interop';

@Component({
  selector: 'app-login',
  imports: [],
  templateUrl: './login.html',
  styleUrl: './login.css'
})
export class Login {
  destroyRef = inject(DestroyRef);
  constructor(private authService: Auth) {}

  ngOnInit() {
    // Initialization logic if needed
  }

  login() : void {
    window.location.href = 'http://localhost:8080/oauth/authenticate';
  }
}
