// src/routes/events/+page.ts
import type {PageLoad} from './$types';
import {apiFetch} from "$lib/api/api";
import type {DiscordRole, EventResponse, Event} from "$lib/type";
import {getDiscordChannels, getDiscordGuildById, getDiscordRolesByGuildId} from "$lib/api/discord";

export const load: PageLoad = async ({ fetch, url }) => {
    const eventsResponse = await apiFetch<EventResponse>( '/v1/events');
    const events = eventsResponse.events

    const guildIds = Array.from(new Set(events.map(e => e.guild_id)));

    // Fetch guild info for each guild ID
    const guilds = await Promise.all(
        guildIds.map(async (id) => ({
            id,
            info: await getDiscordGuildById(id)
        }))
    );

    // Create a map from guild ID to guild name
    const guildIdToName: Record<string, string> = {};
    guilds.forEach(g => {
        guildIdToName[g.id] = g.info.guild.name;
    });

    // Create a map from channel ID to channel name
    const channelIdToName: Record<string, string> = {};
    const discordRolesByGuildId : Record<string, Record<string, string>> = {};
    for (const guildId of guildIds) {
        // initialize the inner map
        discordRolesByGuildId[guildId] = {};

        // fetch channels
        const { channels } = await getDiscordChannels(guildId);
        channels.forEach(c => {
            channelIdToName[c.id] = c.name;
        });

        // fetch roles
        const { roles } = await getDiscordRolesByGuildId(guildId);
        roles.forEach(role => {
            discordRolesByGuildId[guildId][role.id] = role.name;
        });
    }

    return { events, channelIdToName, guildIdToName, discordRolesByGuildId  };
};
