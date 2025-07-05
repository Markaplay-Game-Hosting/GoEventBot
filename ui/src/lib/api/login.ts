import type { Login } from '../type'
import {apiBase} from "$lib/api/index";

export async function login(loginInfo: Login) {
    const res = await fetch(`${apiBase}/login`, {
        method: "POST",
        body: JSON.stringify(loginInfo),
        headers: {
            "Content-Type": "application/json"
        }
    })
    
    if (!res.ok) {
        throw new Error(res.statusText)
    }
    
    return await res.json()
}

export async function discordLogin() {
    const res = await fetch(`${apiBase}/oauth/authenticate`)
    if (!res.ok) {
        throw new Error(res.statusText)
    }
    
    return await res.json()
}