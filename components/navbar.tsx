"use client";
import React from "react";
import SignInModal from "@/components/login-form";
import { Button } from "./ui/button";

export default function Navbar() {
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
        <div className="flex items-center gap-4">
          <a
            href="https://github.com/verse91/clippy-ytb"
            target="_blank"
            rel="noopener noreferrer"
            title="Star it on GitHub ⭐"
            className="transition-colors group"
          >
            <i
              className="bxl bx-github text-4xl text-white transition-all group-hover:text-gray-300 group-hover:scale-110"
            ></i>
          </a>
          <SignInModal
            trigger={
              <button className="p-3 cursor-pointer">
                <i className="bx bxs-user text-3xl text-white hover:text-gray-300 transition-all hover:scale-110"></i>
              </button>
            }
          />
        </div>
      </nav>
    </>
  );
}
