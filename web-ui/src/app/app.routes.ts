import { Routes } from '@angular/router';
import { MainLayout } from './components/common/main-layout/main-layout';

export const routes: Routes = [
    {
        path: '',
        component: MainLayout,
        children: [
            {
                path: 'events',
                loadComponent: () => import('./components/pages/events/events').then((m) => m.Events),
            },
            {
                path: 'jobs',
                loadComponent: () => import('./components/pages/jobs/jobs').then((m) => m.Jobs),
            },
            {
                path: 'tags',
                loadComponent: () => import('./components/pages/tags/tags').then((m) => m.Tags),
            },
            {
                path: 'settings',
                loadComponent: () => import('./components/pages/settings/settings').then((m) => m.Settings),
                data: { hideOutlet: true },
            },
        ]
    },
    {
        path: 'login',
        loadComponent: () =>
        import('./components/pages/login/login').then((m) => m.Login),
    },
];
