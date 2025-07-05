import type {Event, EventResponse} from "$lib/type";
import {apiBase} from '$lib/api/index'


export async function listEvents(): Promise<EventResponse> {
    const response = await fetch(`${apiBase}/v1/events`, {
        credentials: "include"
    });
    if (!response.ok) {
        throw new Error('Failed to fetch events');
    }
    return await response.json() as EventResponse;
}