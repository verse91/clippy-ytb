"use client"
import React, { useState } from "react";
// Import your LoginForm component
import { LoginForm }     from "@/components/login-form"; // Adjust the path as needed
import { Button } from "./ui/button";

export default function Navbar() {
  const [showLogin, setShowLogin] = useState(false);

  const handleOpenLogin = () => setShowLogin(true);
  const handleCloseLogin = () => setShowLogin(false);

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
          >
            <i className="bxl  bx-github text-4xl" style={{color:"#fff"}}></i>
          </a>
          <button className="p-3 cursor-pointer" onClick={handleOpenLogin}>
            <i
              className="bxr bxs-user text-3xl"
              style={{ color: "#ffffff" }}
            ></i>
          </button>
        </div>
      </nav>
      {showLogin && (
        <div className="fixed inset-0 z-50 flex items-center justify-center bg-black/50">
          <div>
            <button
              className="absolute top-2 right-2 text-gray-500 hover:text-gray-700"
              onClick={handleCloseLogin}
              aria-label="Close"
            >
              &times;
            </button>
            <LoginForm />
          </div>
        </div>
      )}
    </>
  );
}
