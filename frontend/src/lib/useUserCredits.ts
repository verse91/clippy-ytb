import { useState, useEffect, useCallback } from "react";
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

  // Memoize fetchUserCredits with useCallback and correct dependencies
  const fetchUserCredits = useCallback(async () => {
    if (!userId) {
      setUserCredits(0);
      setError(null);
      return;
    }

    try {
      setLoading(true);
      setError(null);

      const response = await fetch(
        `${API_BASE_URL}/api/v1/user/${userId}/credits`,
        {
          headers: {
            "Content-Type": "application/json",
          },
        }
      );

      if (response.ok) {
        const data = await response.json();
        // Validate the expected data structure before accessing credits
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
      const errorMessage = `Error fetching user credits: ${
        error instanceof Error ? error.message : "Unknown error"
      }`;
      console.error(errorMessage);
      setError(errorMessage);
      setUserCredits(0);
    } finally {
      setLoading(false);
    }
  }, [userId, API_BASE_URL]);

  // Add fetchUserCredits to the dependency array
  useEffect(() => {
    fetchUserCredits();
  }, [fetchUserCredits]);

  return {
    userCredits,
    loading,
    error,
    refetch: fetchUserCredits,
  };
};
