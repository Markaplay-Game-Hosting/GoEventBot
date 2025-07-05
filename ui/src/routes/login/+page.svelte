<script lang="ts">
    import type { PageData } from './$types'
    import { discordLogin, login } from "$lib/api/login";
    import {PUBLIC_BASE_URL} from "$env/static/public";

    const { data }: { data: PageData } = $props();
    let showPassword = false
    
    async function submit(){
        window.location.href = `${PUBLIC_BASE_URL}/oauth/authenticate`;
    }
</script>
<div class="flex min-h-screen items-center justify-center bg-base-200">
    <div class="m-4 min-h-[50vh] w-full max-w-sm lg:max-w-4xl">
        <!-- logo -->
        <div class="flex items-center justify-center gap-2 p-8">
            <svg
                    width="24"
                    height="24"
                    viewBox="0 0 1024 1024"
                    xmlns="http://www.w3.org/2000/svg">
                <rect x="256" y="670.72" width="512" height="256" rx="128" class="fill-base-content" />
                <circle cx="512" cy="353.28" r="256" class="fill-base-content" />
                <circle
                        cx="512"
                        cy="353.28"
                        r="261"
                        stroke="black"
                        stroke-opacity="0.2"
                        stroke-width="10" />
                <circle cx="512" cy="353.28" r="114.688" class="fill-base-200" />
            </svg>
            <h1 class="text-lg font-bold">Login to your account</h1>
        </div>
        <!-- /logo -->
        <main class="grid bg-base-100 lg:aspect-[2/1] lg:grid-cols-2">
            <!-- image -->
            <figure class="pointer-events-none bg-base-300 object-cover max-lg:hidden">
                <img src="https://picsum.photos/id/283/1200/1200" alt="Login" class="h-full" />
            </figure>
            <!-- /image -->
            <form class="flex flex-col justify-center gap-4 px-10 py-10 lg:px-16">
                <!-- email -->
                <div class="form-control">
                    <label class="label" for="input1"><span class="label-text">Email</span></label>
                    <input
                            type="email"
                            placeholder="email"
                            class="input input-bordered [&:user-invalid]:input-warning [&:user-valid]:input-success"
                            required
                            id="input1" />
                </div>
                <!-- /email -->
                <!-- password -->
                <div class="form-control">
                    <label class="label" for="input2"><span class="label-text">Password</span></label>
                    <input
                            type="password"
                            placeholder="password"
                            class="input input-bordered [&:user-invalid]:input-warning [&:user-valid]:input-success"
                            required
                            minlength="6"
                            for="input2" />
                </div>
                <!-- password -->
                <!-- remember me and recovery link -->
                <div class="flex items-center justify-between gap-3">
                    <label class="flex cursor-pointer gap-3 text-xs">
                        <input name="remember-me" type="checkbox" class="toggle toggle-xs" />
                        Remember me
                    </label>
                    <div class="label">
                        <a class="link-hover link label-text-alt" href="recovery.html">Forgot password?</a>
                    </div>
                </div>
                <!-- /remember me and recovery link -->
                <!-- submit -->
                <button class="btn btn-neutral" type="submit">Login</button>
                <!-- /submit -->
                <!-- login with discord -->
                <button type="button" class="btn" onclick={submit}>
                    <svg id="Discord-Logo"
                         xmlns="http://www.w3.org/2000/svg"
                         viewBox="0 0 126.644 96"
                        width="25"
                        height="25">
                        <defs>
                            <style>.cls-1{fill:#5865f2;}</style>
                        </defs>
                        <path id="Discord-Symbol-Blurple"
                              class="cls-1"
                              d="M81.15,0c-1.2376,2.1973-2.3489,4.4704-3.3591,6.794-9.5975-1.4396-19.3718-1.4396-28.9945,0-.985-2.3236-2.1216-4.5967-3.3591-6.794-9.0166,1.5407-17.8059,4.2431-26.1405,8.0568C2.779,32.5304-1.6914,56.3725.5312,79.8863c9.6732,7.1476,20.5083,12.603,32.0505,16.0884,2.6014-3.4854,4.8998-7.1981,6.8698-11.0623-3.738-1.3891-7.3497-3.1318-10.8098-5.1523.9092-.6567,1.7932-1.3386,2.6519-1.9953,20.281,9.547,43.7696,9.547,64.0758,0,.8587.7072,1.7427,1.3891,2.6519,1.9953-3.4601,2.0457-7.0718,3.7632-10.835,5.1776,1.97,3.8642,4.2683,7.5769,6.8698,11.0623,11.5419-3.4854,22.3769-8.9156,32.0509-16.0631,2.626-27.2771-4.496-50.9172-18.817-71.8548C98.9811,4.2684,90.1918,1.5659,81.1752.0505l-.0252-.0505ZM42.2802,65.4144c-6.2383,0-11.4159-5.6575-11.4159-12.6535s4.9755-12.6788,11.3907-12.6788,11.5169,5.708,11.4159,12.6788c-.101,6.9708-5.026,12.6535-11.3907,12.6535ZM84.3576,65.4144c-6.2637,0-11.3907-5.6575-11.3907-12.6535s4.9755-12.6788,11.3907-12.6788,11.4917,5.708,11.3906,12.6788c-.101,6.9708-5.026,12.6535-11.3906,12.6535Z"/>
                    </svg>
                    Login with Discord
                </button>
                <!-- /login with google -->
                <!-- signup -->
                <div class="label justify-end">
                    <a class="link-hover link label-text-alt" href="signup.html">Create new account</a>
                </div>
                <!-- /signup -->
            </form>
        </main>
    </div>
</div>