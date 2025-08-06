"use client";

import { useEffect, useRef, useState } from "react";
import { useRouter } from "next/navigation";
import { supabase } from "@/lib/supabaseClients";

export default function AuthCallback() {
  const router = useRouter();
  const isMounted = useRef(true);
  const [errorMessage, setErrorMessage] = useState<string | null>(null);

  // Helper function to handle redirect or window close
  const handleRedirectOrClose = (isError = false) => {
    if (window.opener) {
      if (!isError) {
        // Only notify parent on success
        window.opener.postMessage(
          { type: "AUTH_SUCCESS" },
          window.location.origin
        );
      }
      window.close();
    } else {
      if (isMounted.current) {
        router.push("/");
      }
    }
  };

  useEffect(() => {
    const handleAuthCallback = async () => {
      try {
        const { data, error } = await supabase.auth.getSession();

        if (error) {
          setErrorMessage("Authentication failed. Please try again.");
          handleRedirectOrClose(true);
          return;
        }

        if (data.session) {
          // Successfully authenticated
          handleRedirectOrClose(false);
        } else {
          // No session found
          setErrorMessage("No active session found. Please sign in again.");
          handleRedirectOrClose(true);
        }
      } catch (error) {
        setErrorMessage("An unexpected error occurred during authentication.");
        handleRedirectOrClose(true);
      }
    };

    handleAuthCallback();

    // Cleanup function to prevent memory leaks
    return () => {
      isMounted.current = false;
    };
  }, [router]);

  return (
    <div className="min-h-screen flex items-center justify-center bg-black">
      <div className="text-center">
        <div className="w-8 h-8 border-2 border-white/20 border-t-white rounded-full animate-spin mx-auto mb-4"></div>
        <p className="text-white/70">
          {errorMessage ? (
            <span className="text-red-400">{errorMessage}</span>
          ) : (
            "Completing sign in..."
          )}
        </p>
      </div>
    </div>
  );
}
