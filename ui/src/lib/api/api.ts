// src/lib/api.ts
import {apiBase} from "$lib/api/index";
import {goto} from "$app/navigation";

export async function apiFetch<T>(path: string, init?: RequestInit): Promise<T|null> {
    const url = `${apiBase}${path}`;
    const res = await fetch(url, {
        credentials: 'include',
        ...init
    });

    if (res.status === 401) {
        const redirectTo = encodeURIComponent(window.location.pathname + window.location.search);
        goto(`/login?redirectTo=${redirectTo}`, { replaceState: true });
        return null;
    }

    if (!res.ok) {
        const text = await res.text();
        throw new Error(`HTTP ${res.status}: ${text}`);
    }

    return (await res.json()) as T;
}