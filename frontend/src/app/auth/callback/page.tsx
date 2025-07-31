"use client";

import { useEffect } from "react";
import { useRouter } from "next/navigation";
import { supabase } from "@/lib/supabaseClients";

export default function AuthCallback() {
  const router = useRouter();

  useEffect(() => {
    const handleAuthCallback = async () => {
      try {
        const { data, error } = await supabase.auth.getSession();

        if (error) {
          console.error("Auth callback error:", error);
          window.close();
          return;
        }

        if (data.session) {
          // Successfully authenticated
          console.log("Authentication successful");

          // Close the popup and notify the parent window
          if (window.opener) {
            window.opener.postMessage(
              { type: "AUTH_SUCCESS" },
              window.location.origin
            );
          }

          window.close();
        } else {
          // No session found, close popup
          window.close();
        }
      } catch (error) {
        console.error("Auth callback error:", error);
        window.close();
      }
    };

    handleAuthCallback();
  }, []);

  return (
    <div className="min-h-screen flex items-center justify-center bg-black">
      <div className="text-center">
        <div className="w-8 h-8 border-2 border-white/20 border-t-white rounded-full animate-spin mx-auto mb-4"></div>
        <p className="text-white/70">Completing sign in...</p>
      </div>
    </div>
  );
}
