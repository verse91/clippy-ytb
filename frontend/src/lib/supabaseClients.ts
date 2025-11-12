"use client";

import "./localStorage-polyfill";

import { createClient } from "@supabase/supabase-js";

const supabaseUrl = process.env.NEXT_PUBLIC_SUPABASE_URL;
const supabaseAnonKey = process.env.NEXT_PUBLIC_SUPABASE_ANON_KEY;

if (!supabaseUrl || !supabaseAnonKey || supabaseUrl === "YOUR_URL" || supabaseAnonKey === "YOUR_ANON_KEY") {
    throw new Error(
        "Missing Supabase environment variables. Please check that NEXT_PUBLIC_SUPABASE_URL and NEXT_PUBLIC_SUPABASE_ANON_KEY are defined in your .env file with actual values (not placeholders)."
    );
}

const createStorageAdapter = () => {
    const memoryStorage: Record<string, string> = {};

    return {
        getItem: (key: string): string | null => {
            if (typeof window === "undefined") {
                // SSR: use in-memory storage
                return memoryStorage[key] || null;
            }
            try {
                return window.localStorage.getItem(key);
            } catch (error) {
                console.warn("localStorage.getItem failed:", error);
                return memoryStorage[key] || null;
            }
        },
        setItem: (key: string, value: string): void => {
            if (typeof window === "undefined") {
                memoryStorage[key] = value;
                return;
            }
            try {
                window.localStorage.setItem(key, value);
            } catch (error) {
                console.warn("localStorage.setItem failed:", error);
                memoryStorage[key] = value;
            }
        },
        removeItem: (key: string): void => {
            if (typeof window === "undefined") {
                delete memoryStorage[key];
                return;
            }
            try {
                window.localStorage.removeItem(key);
            } catch (error) {
                console.warn("localStorage.removeItem failed:", error);
                delete memoryStorage[key];
            }
        },
    };
};

let supabaseInstance: ReturnType<typeof createClient> | null = null;

function getSupabaseClient(): ReturnType<typeof createClient> {
    if (supabaseInstance) {
        return supabaseInstance;
    }

    if (typeof window === "undefined" && typeof (global as any).localStorage === "undefined") {
        (global as any).localStorage = createStorageAdapter();
    }

    if (typeof window === "undefined") {
        supabaseInstance = createClient(supabaseUrl!, supabaseAnonKey!, {
            auth: {
                storage: createStorageAdapter(),
                autoRefreshToken: false,
                persistSession: false,
                detectSessionInUrl: false,
            },
        });
    } else {
        supabaseInstance = createClient(supabaseUrl!, supabaseAnonKey!, {
            auth: {
                storage: createStorageAdapter(),
                autoRefreshToken: true,
                persistSession: true,
                detectSessionInUrl: true,
            },
        });
    }

    return supabaseInstance;
}

export const supabase = new Proxy({} as ReturnType<typeof createClient>, {
    get(_target, prop) {
        const client = getSupabaseClient();
        const value = client[prop as keyof typeof client];
        if (typeof value === "function") {
            return value.bind(client);
        }
        return value;
    },
});
