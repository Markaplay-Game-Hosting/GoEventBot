import {apiBase} from "$lib/api/index";
import type {DiscordChannelsResponse, DiscordGuildInfoResponse, DiscordRolesResponse} from "$lib/type";


export async function getDiscordChannels(guildId: string) : Promise<DiscordChannelsResponse> {
    const response = await fetch(`${apiBase}/v1/discord/${guildId}/channels`, {
        credentials: "include",
    });

    if (!response.ok) {
        throw new Error('Failed to fetch events');
    }
    
    return response.json();
}

export async function getDiscordRolesByGuildId(guildId: string) : Promise<DiscordRolesResponse> {
    const response = await fetch(`${apiBase}/v1/discord/${guildId}/roles`, {
        credentials: "include",
    });

    if (!response.ok) {
        throw new Error('Failed to fetch events');
    }

    return response.json();
}

export async function getDiscordGuildById(guildId: string) : Promise<DiscordGuildInfoResponse> {
    const response = await fetch(`${apiBase}/v1/discord/${guildId}`, {
        credentials: "include",
    });

    if (!response.ok) {
        throw new Error('Failed to fetch events');
    }

    return response.json();
}