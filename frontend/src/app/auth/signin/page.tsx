"use client";

import { useEffect, useRef, useState } from "react";
import { supabase } from "@/lib/supabaseClients";

export default function SignInPage() {
  const [error, setError] = useState<string | null>(null);
  const hasStarted = useRef(false);

  useEffect(() => {
    const handleSignIn = async () => {
      // Prevent multiple sign-in attempts
      if (hasStarted.current) {
        return;
      }
      hasStarted.current = true;

      try {
        const { data, error } = await supabase.auth.signInWithOAuth({
          provider: "google",
          options: {
            redirectTo: `${window.location.origin}/auth/callback`,
            queryParams: {
              access_type: "offline",
              prompt: "consent",
            },
          },
        });

        if (error) {
          console.error("Sign in error:", error);
          setError(error.message || "Failed to start sign in process");
          // Don't close immediately, let user see the error
          setTimeout(() => {
            if (window.opener) {
              window.close();
            }
          }, 3000);
        }
      } catch (error) {
        console.error("Sign in error:", error);
        setError("An unexpected error occurred. Please try again.");
        setTimeout(() => {
          if (window.opener) {
            window.close();
          }
        }, 3000);
      }
    };

    // Start the sign-in process immediately
    handleSignIn();
  }, []);

  return (
    <div className="min-h-screen flex items-center justify-center bg-black">
      <div className="text-center">
        {error ? (
          <>
            <div className="w-8 h-8 border-2 border-red-500/20 border-t-red-500 rounded-full animate-spin mx-auto mb-4"></div>
            <p className="text-red-400/70 mb-2">Sign in failed</p>
            <p className="text-red-300/50 text-sm">{error}</p>
          </>
        ) : (
          <>
            <div className="w-8 h-8 border-2 border-white/20 border-t-white rounded-full animate-spin mx-auto mb-4"></div>
            <p className="text-white/70">Opening Google sign in...</p>
          </>
        )}
      </div>
    </div>
  );
}
