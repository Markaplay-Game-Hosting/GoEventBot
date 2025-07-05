<script lang="ts">
	import '../app.css';
    import { page } from '$app/state';
    import { FontAwesomeIcon } from '@fortawesome/svelte-fontawesome';
    import '$lib/icons';
    
	let { children } = $props();
    const isLoginPage = $derived(() => page.url.pathname === '/login');
</script>


{#if isLoginPage()}
    {@render children()}
{:else}
    <!-- Layout with sidebar -->
    <div class="drawer min-h-screen bg-base-200 lg:drawer-open">
        <input id="my-drawer" type="checkbox" class="drawer-toggle" />
        <main class="drawer-content">
            <div class="grid grid-cols-12 gap-y-12 p-4 lg:gap-x-12 lg:p-10">
                <header class="col-span-12 flex items-center gap-2 lg:gap-4">
                    <label for="my-drawer" class="btn btn-square btn-ghost drawer-button lg:hidden">
                        <svg data-src="https://unpkg.com/heroicons/20/solid/bars-3.svg" class="h-5 w-5"></svg>
                    </label>
                    <div class="grow">
                        <h1 class="lg:text-2xl lg:font-light">{page.data.title ?? "Go Event Dispatcher"}</h1>
                    </div>
                    <!-- dropdown -->
                    <div class="dropdown-end dropdown z-10">
                        <div class="avatar btn btn-circle btn-ghost">
                            <div class="w-10 rounded-full">
                                <FontAwesomeIcon icon="fa-regular fa-user" />
                            </div>
                        </div>
                        <ul class="menu dropdown-content mt-3 w-52 rounded-box bg-base-100 p-2 shadow-2xl">
                            <li><a href="/logout">Logout</a></li>
                        </ul>
                    </div>
                    <!-- /dropdown -->
                </header>
                <section class="stats stats-vertical col-span-12 w-full shadow-sm xl:stats-horizontal">
                    {@render children()}
                </section>
            </div>
        </main>
        <aside class="drawer-side z-10">
            <label for="my-drawer" class="drawer-overlay"></label>
            <nav class="flex min-h-screen w-72 flex-col gap-2 overflow-y-auto bg-base-100 px-6 py-10">
                <!-- your sidebar menu here -->
                <div class="mx-4 flex items-center gap-2 font-black">
                    GoEventManager
                </div>
                <ul class="menu">
                    <li class="active">
                        <a href="/">
                            <FontAwesomeIcon icon="fa-regular fa-house" />
                            Dashboard
                        </a>
                    </li>
                    <li>
                        <a href="/events">
                            <FontAwesomeIcon icon="fa-regular fa-calendar" />
                            Events
                        </a>
                    </li>
                    <li>
                        <a href="/jobs">
                            <FontAwesomeIcon icon="fa-regular fa-hourglass-clock" />
                            Jobs
                        </a>
                    </li>
                    <li>
                        <a href="/tags">
                            <FontAwesomeIcon icon="fa-regular fa-tags" />
                            Tags
                        </a>
                    </li>
                    <li>
                        <a href="/Settings">
                            <FontAwesomeIcon icon="fa-regular fa-gear" />
                            Settings
                        </a>
                    </li>
                </ul>
            </nav>
        </aside>
    </div>
{/if}