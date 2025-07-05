<svelte:head>
    <title>Events</title>
</svelte:head>
<script lang="ts">
    import {FontAwesomeIcon} from "@fortawesome/svelte-fontawesome";
    import type {Event} from "$lib/type";

    export let data: {events: Event[]; channelIdToName: Record<string, string>, guildIdToName: Record<string, string>, discordRolesByGuildId: Record<string, Record<string, string>>};

    // just split out id + raw text
    function splitSegments(message: string) {
        const re = /<@&(\d+)>/g;
        const segs: Array<
            { type: "text"; text: string } |
            { type: "mention"; text: string; id: string }
        > = [];
        let last = 0, m: RegExpExecArray;

        while ((m = re.exec(message))) {
            if (m.index > last) {
                segs.push({ type: "text", text: message.slice(last, m.index) });
            }
            segs.push({ type: "mention", text: m[0], id: m[1] });
            last = m.index + m[0].length;
        }
        if (last < message.length) {
            segs.push({ type: "text", text: message.slice(last) });
        }
        return segs;
    }
</script>

<div class="overflow-x-auto rounded-box border border-base-content/5 bg-base-100">
    <table class="table">
        <thead>
            <tr>
                <th>Name</th>
                <th>Description</th>
                <th>Channel</th>
                <th>Guild</th>
                <th>Status</th>
            </tr>
        </thead>
        <tbody>
        {#each data.events as partyEvent}
            <tr>
                <td>{partyEvent.title}</td>
                <td>
                    {#each splitSegments(partyEvent.description) as seg}
                        {#if seg.type === 'text'}
                            {@html seg.text}
                        {:else}
                          <span class="inline-flex items-center space-x-2">
                            <!-- raw mention stays the same -->
                            
                                                  <!-- badge pulled directly from your nested map -->
                            <span class="badge badge-outline badge-primary">
                              @{
                                data.discordRolesByGuildId[partyEvent.guild_id]?.[seg.id]
                                ?? seg.id
                            }
                            </span>
                          </span>
                        {/if}
                    {/each}
                </td>
                <td>{data.channelIdToName[partyEvent.channel_id] ?? partyEvent.channel_id}</td>
                <td>{data.guildIdToName[partyEvent.guild_id] ?? partyEvent.guild_id}</td>
                <td>
                    {#if partyEvent.is_active === true}
                        <FontAwesomeIcon icon="fa-regular fa-check-circle" style="color: green" size="2x"  />
                    {:else}
                        <FontAwesomeIcon icon="fa-regular fa-circle-xmark" style="color: red" size="2x" />
                    {/if}
                </td>
            </tr>
        {/each}
        </tbody>
    </table>
</div>