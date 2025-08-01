"use client";

import { motion, AnimatePresence } from "framer-motion";
import { ReactNode } from "react";

interface TypingPopupProps {
  isVisible: boolean;
  icon?: ReactNode;
  text?: string;
  className?: string;
  showTypingDots?: boolean;
  widthClassName?: string;
}

// Typing dots animation component
function TypingDots() {
  return (
    <div className="flex space-x-1">
      {[0, 1, 2].map((i) => (
        <motion.div
          key={i}
          className="w-1.5 h-1.5 bg-white/70 rounded-full"
          animate={{
            scale: [1, 1.2, 1],
            opacity: [0.5, 1, 0.5],
          }}
          transition={{
            duration: 1,
            repeat: Infinity,
            delay: i * 0.1,
          }}
        />
      ))}
    </div>
  );
}

export default function NotiPopup({
  isVisible,
  icon,
  text = "Processing...",
  className = "",
  showTypingDots = true,
  widthClassName = "",
}: TypingPopupProps) {
  return (
    <AnimatePresence>
      {isVisible && (
        <motion.div
          role="alert"
          aria-live="polite"
          className={`fixed left-1/2 bottom-8 -translate-x-1/2 backdrop-blur-2xl bg-white/[0.02] rounded-full px-4 py-2 shadow-lg border border-white/[0.05] z-50 ${widthClassName} ${className}`}
          initial={{ opacity: 0, y: 20 }}
          animate={{ opacity: 1, y: 0 }}
          exit={{ opacity: 0, y: 20 }}
        >
          <div className="flex items-center gap-3">
            {icon && (
              <div className="rounded-full bg-white/[0.05] flex items-center justify-center text-center">
                {icon}
              </div>
            )}
            <div className="flex items-center gap-2 text-sm text-white/70">
              <span>{text}</span>
              {showTypingDots && <TypingDots />}
            </div>
          </div>
        </motion.div>
      )}
    </AnimatePresence>
  );
}
