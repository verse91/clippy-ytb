import { useState, useEffect } from "react";

interface UseUserCreditsReturn {
  userCredits: number;
  loading: boolean;
  error: string | null;
  refetch: () => void;
}

export const useUserCredits = (
  userId: string | undefined
): UseUserCreditsReturn => {
  const [userCredits, setUserCredits] = useState<number>(0);
  const [loading, setLoading] = useState(false);
  const [error, setError] = useState<string | null>(null);

  const API_BASE_URL =
    process.env.NEXT_PUBLIC_API_BASE_URL || "http://localhost:8080";

  const fetchUserCredits = async () => {
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
            "X-User-ID": userId,
          },
        }
      );

      if (response.ok) {
        const data = await response.json();
        setUserCredits(data.data?.credits || 0);
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
  };

  useEffect(() => {
    fetchUserCredits();
  }, [userId]);

  return {
    userCredits,
    loading,
    error,
    refetch: fetchUserCredits,
  };
};
