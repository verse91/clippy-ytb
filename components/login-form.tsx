"use client";

import { Button } from "@/components/ui/button";
import { useState } from "react";
// import { signIn } from "@/lib/auth-client";
import { cn } from "@/lib/utils";
import {
  Dialog,
  DialogContent,
  DialogTitle,
  DialogTrigger,
  DialogClose,
} from "@/components/ui/dialog";
import { Label } from "@/components/ui/label";
import { ScrollArea } from "@/components/ui/scroll-area";

interface SignInModalProps {
  trigger?: React.ReactNode;
}

import { TERMS_TEXT } from "@/lib/terms";

export default function SignInModal({ trigger }: SignInModalProps) {
  const [loading, setLoading] = useState(false);
  const [open, setOpen] = useState(false);
  const [accepted, setAccepted] = useState(false);

  //   const handleSignIn = async () => {
  //     await signIn.social(
  //       {
  //         provider: "google",
  //         callbackURL: "/",
  //       },
  //       {
  //         onRequest: () => setLoading(true),
  //         onResponse: () => {
  //           setLoading(false);
  //           setOpen(false);
  //         },
  //       }
  //     );
  //   };

  return (
    <Dialog open={open} onOpenChange={setOpen}>
      <DialogTrigger asChild>{trigger}</DialogTrigger>
      <DialogContent className="sm:max-w-md">
        <DialogClose asChild>
            <i
              className="bxr  bxs-backspace absolute top-7 right-6 transition-transform hover:scale-105 text-3xl cursor-pointer"
              style={{ color: "#ffffff" }}
              title="ESC"
            ></i>
        </DialogClose>
        <div className="flex flex-col gap-2">
          <DialogTitle className="text-2xl font-semibold">Sign In</DialogTitle>
          {/* Terms and Conditions scroll area */}
          <ScrollArea className="h-48 border rounded-md p-3 bg-muted/30 my-2">
            <div className="space-y-4 text-sm text-muted-foreground">
              {TERMS_TEXT.map((section, idx) => (
                <div key={idx}>
                  <div className="text-white/80 font-semibold mb-1">
                    {section.title}
                  </div>
                  {/* Highlight my email in terms */}
                  <div
                    dangerouslySetInnerHTML={{
                      __html: section.content
                        .replace(
                          "versedev.store@proton.me",
                          '<a href="mailto:versedev.store@proton.me"class="text-white underline underline-offset-4">versedev.store@proton.me</a>'
                        )
                        .replace(/\n/g, "<br/>")
                        .replace(/^- /gm, "• "),
                    }}
                  />
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
          <Button
            size="lg"
            className={cn("w-full gap-2 cursor-pointer")}
            disabled={loading || !accepted}
            // onClick={handleSignIn}
          >
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
            Sign in with Google
          </Button>
        </div>
      </DialogContent>
    </Dialog>
  );
}
