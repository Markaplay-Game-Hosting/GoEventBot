
export type Login = {
    email: string;
    password: string;
}

export type EventInstance = {
    event_id: string;
    title: string;
    description: string;
    duration: string;
    start_date: Date;
    end_date: Date;
}

export type Event = {
    id: string;
    title: string;
    description: string;
    duration: string;
    rrule: string;
    channel_id: string;
    guild_id: string;
    is_active: boolean;
}

export type EventInstancesResponse = {
    events: EventInstance;
}

export type EventResponse = {
    events: Event[];
}

export type DiscordChannel = {
    id: string;
    name: string;
}

export type DiscordChannelsResponse = {
    channels: DiscordChannel[];
}

export type DiscordGuildInfo = {
    id: string;
    name: string;
}

export type DiscordGuildInfoResponse = {
    guild: DiscordGuildInfo;
}

export type DiscordRole = {
    id: string;
    name: string;
}

export type DiscordRolesResponse = {
    roles: DiscordRole[];
}