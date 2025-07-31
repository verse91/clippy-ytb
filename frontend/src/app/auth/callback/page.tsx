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
          // For mobile, redirect back to home
          if (!window.opener) {
            router.push("/");
          } else {
            window.close();
          }
          return;
        }

        if (data.session) {
          // Successfully authenticated
          console.log("Authentication successful");

          if (window.opener) {
            // Popup flow - notify parent and close
            window.opener.postMessage(
              { type: "AUTH_SUCCESS" },
              window.location.origin
            );
            window.close();
          } else {
            // Mobile redirect flow - redirect back to home
            router.push("/");
          }
        } else {
          // No session found
          if (window.opener) {
            window.close();
          } else {
            router.push("/");
          }
        }
      } catch (error) {
        console.error("Auth callback error:", error);
        if (window.opener) {
          window.close();
        } else {
          router.push("/");
        }
      }
    };

    handleAuthCallback();
  }, [router]);

  return (
    <div className="min-h-screen flex items-center justify-center bg-black">
      <div className="text-center">
        <div className="w-8 h-8 border-2 border-white/20 border-t-white rounded-full animate-spin mx-auto mb-4"></div>
        <p className="text-white/70">Completing sign in...</p>
      </div>
    </div>
  );
}
