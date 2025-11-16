import { useState, useEffect, useCallback, useRef } from "react";
interface UseUserCreditsReturn {
    userCredits: number;
    loading: boolean;
    error: string | null;
    refetch: () => void;
}
const API_BASE_URL =
    process.env.NEXT_PUBLIC_API_BASE_URL || "http://localhost:8080";

export const useUserCredits = (
    userId: string | undefined
): UseUserCreditsReturn => {
    const [userCredits, setUserCredits] = useState<number>(0);
    const [loading, setLoading] = useState(false);
    const [error, setError] = useState<string | null>(null);
    const abortControllerRef = useRef<AbortController | null>(null);

    const fetchUserCredits = useCallback(async () => {
        if (!userId) {
            setUserCredits(0);
            setError(null);
            return;
        }

        // Cancel previous request to prevent memory leaks
        if (abortControllerRef.current) {
            abortControllerRef.current.abort();
        }

        abortControllerRef.current = new AbortController();

        try {
            setLoading(true);
            setError(null);

            const response = await fetch(
                `${API_BASE_URL}/api/v1/user/${userId}/credits`,
                {
                    headers: {
                        "Content-Type": "application/json",
                    },
                    signal: abortControllerRef.current.signal,
                }
            );

            if (response.ok) {
                const data = await response.json();
                if (
                    data &&
                    typeof data === "object" &&
                    data.data &&
                    typeof data.data === "object" &&
                    typeof data.data.credits === "number"
                ) {
                    setUserCredits(data.data.credits);
                } else {
                    setUserCredits(0);
                    setError("Unexpected response structure when fetching user credits.");
                }
            } else {
                const errorMessage = `Failed to fetch user credits: ${response.status}`;
                console.error(errorMessage);
                setError(errorMessage);
                setUserCredits(0);
            }
        } catch (error) {
            // Ignore AbortError - request was cancelled intentionally
            if (error instanceof Error && error.name === 'AbortError') {
                return;
            }
            const errorMessage = `Error fetching user credits: ${error instanceof Error ? error.message : "Unknown error"}`;
            console.error(errorMessage);
            setError(errorMessage);
            setUserCredits(0);
        } finally {
            setLoading(false);
        }
    }, [userId]); // Removed API_BASE_URL from deps as it's a constant

    useEffect(() => {
        fetchUserCredits();
        return () => {
            // Cleanup: abort pending request on unmount
            if (abortControllerRef.current) {
                abortControllerRef.current.abort();
            }
        };
    }, [fetchUserCredits]);

    return {
        userCredits,
        loading,
        error,
        refetch: fetchUserCredits,
    };
};
