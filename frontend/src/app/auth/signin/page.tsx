"use client";

import { useEffect } from "react";
import { supabase } from "@/lib/supabaseClients";

export default function SignInPage() {
  useEffect(() => {
    const handleSignIn = async () => {
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
          window.close();
        }
      } catch (error) {
        console.error("Sign in error:", error);
        window.close();
      }
    };

    // Start the sign-in process immediately
    handleSignIn();
  }, []);

  return (
    <div className="min-h-screen flex items-center justify-center bg-black">
      <div className="text-center">
        <div className="w-8 h-8 border-2 border-white/20 border-t-white rounded-full animate-spin mx-auto mb-4"></div>
        <p className="text-white/70">Opening Google sign in...</p>
      </div>
    </div>
  );
}
