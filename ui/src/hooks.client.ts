// src/hooks.client.ts
import { goto } from '$app/navigation';
import {redirect} from "@sveltejs/kit";
import {apiBase} from "$lib/api";

/**
 * This runs for *every* fetch() call in the browser.
 * If the response is 401, we redirect to /login.
 */
export async function handleFetch({ request, fetch }) {
    console.log(`Fetching ${request.url}`);
    if (!request.url.startsWith(apiBase)) {
        return;
    }

    const response = await fetch(request);

    if (response.status === 401) {
        // preserve the current URL
        const params = encodeURIComponent(window.location.pathname + window.location.search);
        // this forces a full-page reload → picks up your server-hook next
        redirect(303,`/login?redirectTo=${params}`);
        // return a stub so the caller doesn’t try to parse JSON
        //return new Response(null, { status: 401 });
    }

    return response;
}
