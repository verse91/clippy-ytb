"use client";
import React from "react";
import SignInModal from "@/components/login-form";
import { Button } from "./ui/button";
import { useAuth } from "@/lib/auth-context";
import Image from "next/image";
import {
  DropdownMenu,
  DropdownMenuContent,
  DropdownMenuItem,
  DropdownMenuLabel,
  DropdownMenuSeparator,
  DropdownMenuTrigger,
} from "@/components/ui/dropdown-menu";



export default function Navbar() {
  const { user, signOut, loading } = useAuth();

  const handleSignOut = async () => {
    try {
      await signOut();
    } catch (error) {
      console.error("Error signing out:", error);
    }
  };

  return (
    <>
      <nav className="absolute top-0 left-0 w-full flex items-center justify-between px-8 py-4 z-10">
        {/* <a href="/">
          <span
            className="text-3xl font-extrabold text-red-500 select-none"
            style={{ fontFamily: "FTV" }}
          >
            CLIP<span className="text-white">PY</span>
          </span>
        </a> */}
        <a href="/">
          <img
            src="/assets/icons/logo-no-bg.png"
            alt="Clippy Logo"
            className="h-12 w-auto select-none"
          />
        </a>
        <div className="flex items-center gap-6">
          <a
            href="https://github.com/verse91/clippy-ytb"
            target="_blank"
            rel="noopener noreferrer"
            title="Star it on GitHub ⭐"
            className="transition-colors group"
          >
            <i className="bxl bx-github text-4xl text-white transition-all group-hover:text-gray-300 group-hover:scale-110"></i>
          </a>
          {loading ? (
            <div className="p-3">
              <div className="w-6 h-6 border-2 border-white/20 border-t-white rounded-full animate-spin"></div>
            </div>
          ) : user ? (
            <div className="flex items-center gap-3">
              <DropdownMenu>
                <DropdownMenuTrigger asChild>
                  <Image
                    src={user.user_metadata.picture}
                    alt="User Avatar"
                    width={35}
                    height={35}
                    className="rounded-full border shadow cursor-pointer mb-1 hover:scale-105"
                  />
                </DropdownMenuTrigger>
                <DropdownMenuContent className="w-56" align="start">
                  <DropdownMenuLabel>
                    {user.user_metadata.name}
                  </DropdownMenuLabel>
                  <DropdownMenuLabel className="text-xs text-muted-foreground -mt-3 mb-3">
                    {user.email}
                  </DropdownMenuLabel>
                  <DropdownMenuItem className="cursor-pointer">
                    <i className="bx bxs-credit-card-alt text-sm text-white"></i>
                    Buy credits
                  </DropdownMenuItem>
                  <DropdownMenuSeparator />
                  <DropdownMenuItem
                    className="cursor-pointer"
                    onClick={handleSignOut}
                  >
                    <i className="bx bxs-arrow-out-right-square-half text-sm text-white"></i>
                    Sign out
                  </DropdownMenuItem>
                </DropdownMenuContent>
              </DropdownMenu>
            </div>
          ) : (
            <SignInModal
              trigger={
                <Button
                  className="hidden sm:inline-flex hover:scale-105 cursor-pointer rounded-xl max-w-24"
                  variant="default"
                  title="Sign In"
                >
                  <i className="bx bxs-arrow-in-right-square-half text-3xl text-black"></i>
                  <p
                    className="font-bold mr-1"
                    style={{ fontFamily: "SF-Pro-Display" }}
                  >
                    Sign in
                  </p>
                </Button>
              }
            />
          )}
        </div>
      </nav>
    </>
  );
}
