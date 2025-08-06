"use client";

import { Button } from "@/components/ui/button";
import { useState, useEffect, useRef } from "react";
import { cn } from "@/lib/utils";
import { supabase } from "@/lib/supabaseClients";
import { useAuth } from "@/lib/auth-context";
import {
  Dialog,
  DialogContent,
  DialogTitle,
  DialogTrigger,
} from "@/components/ui/dialog";
import { Label } from "@/components/ui/label";
import { ScrollArea } from "@/components/ui/scroll-area";
import { TERMS_TEXT } from "@/lib/terms";
import { LoaderIcon } from "lucide-react";

interface SignInModalProps {
  trigger?: React.ReactNode;
}

export default function SignInModal({ trigger }: SignInModalProps) {
  const { user } = useAuth();
  const [loading, setLoading] = useState(false);
  const [open, setOpen] = useState(false);
  const [accepted, setAccepted] = useState(false);
  const [error, setError] = useState<string | null>(null);
  const popupRef = useRef<Window | null>(null);
  const checkClosedIntervalRef = useRef<NodeJS.Timeout | null>(null);
  const messageListenerRef = useRef<((event: MessageEvent) => void) | null>(
    null
  );

  const handleSignIn = async () => {
    setLoading(true);
    setError(null);

    // Clean up any existing popup and listeners
    if (popupRef.current && !popupRef.current.closed) {
      popupRef.current.close();
    }
    if (checkClosedIntervalRef.current) {
      clearInterval(checkClosedIntervalRef.current);
    }
    if (messageListenerRef.current) {
      window.removeEventListener("message", messageListenerRef.current);
    }

    try {
      // Check if it's a mobile device
      const isMobile =
        /Android|webOS|iPhone|iPad|iPod|BlackBerry|IEMobile|Opera Mini/i.test(
          navigator.userAgent
        );

      if (isMobile) {
        // For mobile, use a simpler approach without redirect
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
          console.error("Login error:", error.message);
          setError(error.message);
          setLoading(false);
        } else {
          // For mobile, we'll let the auth context handle the state
          console.log("Mobile sign in initiated");
          // Keep loading state to prevent modal from closing too quickly
          // The useEffect will handle closing the modal once user is authenticated
        }
      } else {
        // Use popup for desktop devices
        const width = 400;
        const height = 500;
        const left = (window.screen.width - width) / 2;
        const top = (window.screen.height - height) / 2;

        const popup = window.open(
          `${window.location.origin}/auth/signin`,
          "auth-popup",
          `width=${width},height=${height},left=${left},top=${top},scrollbars=yes,resizable=yes`
        );

        if (!popup) {
          setError(
            "Please allow popups to sign in. You can enable popups in your browser settings."
          );
          setLoading(false);
          return;
        }

        popupRef.current = popup;

        // Listen for messages from the popup
        const handleMessage = (event: MessageEvent) => {
          if (event.origin !== window.location.origin) return;

          if (event.data.type === "AUTH_SUCCESS") {
            // Authentication successful
            setLoading(false);
            setOpen(false);
            cleanupPopup();
          } else if (event.data.type === "AUTH_ERROR") {
            // Authentication failed
            setError(event.data.error || "Authentication failed");
            setLoading(false);
            cleanupPopup();
          }
        };

        messageListenerRef.current = handleMessage;
        window.addEventListener("message", handleMessage);

        // Check if popup is closed
        const checkClosed = setInterval(() => {
          if (popup.closed) {
            clearInterval(checkClosed);
            setLoading(false);
            cleanupPopup();
          }
        }, 1000);

        checkClosedIntervalRef.current = checkClosed;
      }
    } catch (err) {
      console.error("Login error:", err);
      setError("An unexpected error occurred. Please try again.");
      setLoading(false);
    }
  };

  const cleanupPopup = () => {
    if (popupRef.current && !popupRef.current.closed) {
      popupRef.current.close();
      popupRef.current = null;
    }
    if (checkClosedIntervalRef.current) {
      clearInterval(checkClosedIntervalRef.current);
      checkClosedIntervalRef.current = null;
    }
    if (messageListenerRef.current) {
      window.removeEventListener("message", messageListenerRef.current);
      messageListenerRef.current = null;
    }
  };

  // Cleanup on component unmount
  useEffect(() => {
    return () => {
      cleanupPopup();
    };
  }, []);

  // Close modal when user is successfully authenticated
  useEffect(() => {
    if (user && open) {
      // Add a small delay to prevent immediate closure on mobile
      const timer = setTimeout(() => {
        setLoading(false);
        setOpen(false);
      }, 1000);

      return () => clearTimeout(timer);
    }
  }, [user, open]);

  // Reset error and loading when modal closes, but preserve terms acceptance for better UX
  useEffect(() => {
    if (!open) {
      setError(null);
      setLoading(false);
      cleanupPopup();
      // Don't reset accepted state to improve mobile UX
      // setAccepted(false);
    }
  }, [open]);

  // Reset terms acceptance only when component unmounts or user changes
  useEffect(() => {
    if (!user) {
      setAccepted(false);
    }
  }, [user]);

  return (
    <Dialog open={open} onOpenChange={setOpen}>
      <DialogTrigger asChild>{trigger}</DialogTrigger>
      <DialogContent className="sm:max-w-lg max-h-[80vh]">
        <div className="flex flex-col gap-2">
          <DialogTitle className="text-2xl font-semibold">Sign In</DialogTitle>
          {/* Terms and Conditions scroll area */}
          <ScrollArea className="h-80 border rounded-md p-3 bg-muted/30 my-2">
            <div className="space-y-4 text-sm text-muted-foreground">
              {TERMS_TEXT.map((section, idx) => (
                <div key={idx}>
                  <div className="text-white/80 font-semibold mb-1">
                    {section.title}
                  </div>
                  {/* Highlight my email in terms */}
                  <div className="text-sm text-muted-foreground">
                    {section.content
                      .split("\n")
                      .map((line, lineIndex) => {
                        // Replace email with mailto link
                        const emailRegex = /versedev\.store@proton\.me/g;
                        const parts = line.split(emailRegex);

                        if (parts.length === 1) {
                          // No email found, render as plain text
                          if (line === "") {
                            return <br key={lineIndex} />;
                          }
                          // Handle bullet points
                          if (line.startsWith("- ")) {
                            return (
                              <div
                                key={lineIndex}
                                className="flex items-start gap-2"
                              >
                                <span className="text-white/60">•</span>
                                <span>{line.substring(2)}</span>
                              </div>
                            );
                          }
                          return <span key={lineIndex}>{line}</span>;
                        }

                        // Email found, render with link
                        const elements: React.ReactElement[] = [];
                        parts.forEach((part, partIndex) => {
                          if (partIndex > 0) {
                            // Add email link
                            elements.push(
                              <a
                                key={`email-${lineIndex}-${partIndex}`}
                                href="mailto:versedev.store@proton.me"
                                className="text-white underline underline-offset-4"
                              >
                                versedev.store@proton.me
                              </a>
                            );
                          }
                          if (part) {
                            elements.push(
                              <span key={`part-${lineIndex}-${partIndex}`}>
                                {part}
                              </span>
                            );
                          }
                        });

                        // Handle bullet points for lines with email
                        if (line.startsWith("- ")) {
                          return (
                            <div
                              key={lineIndex}
                              className="flex items-start gap-2"
                            >
                              <span className="text-white/60">•</span>
                              <div>{elements}</div>
                            </div>
                          );
                        }
                        return <div key={lineIndex}>{elements}</div>;
                      })
                      .map((element, index) => {
                        // Handle bullet points - element is now a React element, not a string
                        // The bullet point handling should be done in the first map
                        return element;
                      })}
                  </div>
                </div>
              ))}
            </div>
          </ScrollArea>
          {/* Acceptance Checkbox */}
          <div className="flex items-center space-x-2 pt-2">
            <input
              type="checkbox"
              id="accept-terms"
              checked={accepted}
              onChange={(e) => setAccepted(e.target.checked)}
              className="rounded border-gray-300 text-primary focus:ring-primary"
            />
            <Label htmlFor="accept-terms" className="text-sm">
              I have read and agree to the Terms & Conditions.
            </Label>
          </div>
        </div>
        <div className="flex w-full items-center justify-center pt-4">
          {error && (
            <div className="w-full mb-4 p-3 bg-red-500/10 border border-red-500/20 rounded-md">
              <p className="text-red-400 text-sm">{error}</p>
            </div>
          )}
          <Button
            size="lg"
            className={cn("w-full gap-2 cursor-pointer")}
            disabled={loading || !accepted}
            onClick={handleSignIn}
          >
            {loading ? (
              <LoaderIcon className="w-4 h-4 animate-spin" />
            ) : (
              <svg
                xmlns="http://www.w3.org/2000/svg"
                width="0.98em"
                height="1em"
                viewBox="0 0 256 262"
              >
                <path
                  fill="#4285F4"
                  d="M255.878 133.451c0-10.734-.871-18.567-2.756-26.69H130.55v48.448h71.947c-1.45 12.04-9.283 30.172-26.69 42.356l-.244 1.622l38.755 30.023l2.685.268c24.659-22.774 38.875-56.282 38.875-96.027"
                ></path>
                <path
                  fill="#34A853"
                  d="M130.55 261.1c35.248 0 64.839-11.605 86.453-31.622l-41.196-31.913c-11.024 7.688-25.82 13.055-45.257 13.055c-34.523 0-63.824-22.773-74.269-54.25l-1.531.13l-40.298 31.187l-.527 1.465C35.393 231.798 79.49 261.1 130.55 261.1"
                ></path>
                <path
                  fill="#FBBC05"
                  d="M56.281 156.37c-2.756-8.123-4.351-16.827-4.351-25.82c0-8.994 1.595-17.697 4.206-25.82l-.073-1.73L15.26 71.312l-1.335.635C5.077 89.644 0 109.517 0 130.55s5.077 40.905 13.925 58.602z"
                ></path>
                <path
                  fill="#EB4335"
                  d="M130.55 50.479c24.514 0 41.05 10.589 50.479 19.438l36.844-35.974C195.245 12.91 165.798 0 130.55 0C79.49 0 35.393 29.301 13.925 71.947l42.211 32.783c10.59-31.477 39.891-54.251 74.414-54.251"
                ></path>
              </svg>
            )}
            {loading ? "Signing in..." : "Sign in with Google"}
          </Button>
        </div>
      </DialogContent>
    </Dialog>
  );
}
