<script lang="ts">
    import FullCalendar, { type CalendarOptions } from 'svelte-fullcalendar';
    import dayGridPlugin from '@fullcalendar/daygrid';
    import { datetime, RRule, RRuleSet, rrulestr } from 'rrule'
    
    rrulestr()
    
    let options: CalendarOptions = {
        initialView: 'dayGridMonth',
        plugins: [dayGridPlugin],
        windowResize: ({ view }) => {true},
        events: [
            { title: 'Event 1', date: '2025-07-01' },
            { title: 'Event 2', date: '2025-07-05' }
        ]
    };
</script>
<svelte:head>
    <title>Dashboard</title>
</svelte:head>
<FullCalendar {options} />