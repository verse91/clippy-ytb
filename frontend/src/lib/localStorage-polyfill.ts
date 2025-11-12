// localStorage polyfill for SSR
// This must be imported before any code that uses localStorage
// Import this file at the very top of files that need localStorage support

(function setupLocalStoragePolyfill() {
    if (typeof window !== "undefined") {
        return;
    }

    const memoryStorage: Record<string, string> = {};

    const localStoragePolyfill = {
        getItem: function(key: string): string | null {
            return memoryStorage[key] || null;
        },
        setItem: function(key: string, value: string): void {
            try {
                memoryStorage[key] = String(value);
            } catch (e) {
            }
        },
        removeItem: function(key: string): void {
            delete memoryStorage[key];
        },
        clear: function(): void {
            Object.keys(memoryStorage).forEach((key) => delete memoryStorage[key]);
        },
        key: function(index: number): string | null {
            const keys = Object.keys(memoryStorage);
            return keys[index] || null;
        },
        get length(): number {
            return Object.keys(memoryStorage).length;
        },
    };

    try {
        (global as any).localStorage = localStoragePolyfill;
    } catch (e) {
    }

    try {
        (globalThis as any).localStorage = localStoragePolyfill;
    } catch (e) {
    }

    try {
        if (typeof (global as any) !== "undefined") {
            Object.defineProperty(global as any, "localStorage", {
                value: localStoragePolyfill,
                writable: true,
                configurable: true,
            });
        }
    } catch (e) {
    }
})();
