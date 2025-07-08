import React from "react";

export default function Navbar() {
  return (
    <nav className="absolute top-0 left-0 w-full flex items-center justify-between px-8 py-4 z-10">
      <span className="text-3xl font-extrabold text-red-500 select-none" style={{ fontFamily: 'FTV' }}>
        CLIP<span className="text-white">PY</span>
      </span>
      <button className="p-3">
        <i className="bxr bxs-user text-3xl" style={{ color: "#ffffff" }}></i>
      </button>
    </nav>
  );
}
