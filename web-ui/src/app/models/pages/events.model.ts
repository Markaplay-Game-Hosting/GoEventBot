export interface Event {
    "event_id": string,
    "title": string,
    "description": string,
    "duration": string,
    "start_date": Date,
    "end_date": Date
    "channel_id": string,
    "guild_id": string,
    "is_active": boolean
}

export interface EventList {
    "events": Event[]
}