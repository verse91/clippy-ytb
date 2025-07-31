"use client";
import React, { type JSX } from "react";
import { motion, Transition } from "framer-motion";
import { cn } from "@/lib/utils";

type TextShimmerWave = {
  children: string;
  as?: React.ElementType;
  className?: string;
  duration?: number;
  zDistance?: number;
  xDistance?: number;
  yDistance?: number;
  spread?: number;
  scaleDistance?: number;
  rotateYDistance?: number;
  transition?: Transition;
};

export function TextShimmerWave({
  children,
  as: Component = "p",
  className,
  duration = 1,
  zDistance = 10,
  xDistance = 2,
  yDistance = -2,
  spread = 1,
  scaleDistance = 1.1,
  rotateYDistance = 10,
  transition,
}: TextShimmerWave) {
  const MotionComponent = motion.create(
    Component as keyof JSX.IntrinsicElements
  );

  // Memoize the character array to prevent unnecessary re-renders
  const characters = React.useMemo(() => children.split(""), [children]);

  // Check for reduced motion preference
  const prefersReducedMotion = React.useMemo(() => {
    if (typeof window !== "undefined") {
      return window.matchMedia("(prefers-reduced-motion: reduce)").matches;
    }
    return false;
  }, []);

  return (
    <MotionComponent
      className={cn(
        "relative inline-block [perspective:500px]",
        "[--base-color:#a1a1aa] [--base-gradient-color:#000]",
        "dark:[--base-color:#71717a] dark:[--base-gradient-color:#ffffff]",
        className
      )}
      style={{ color: "var(--base-color)" }}
    >
      {characters.map((char, i) => {
        const delay = (i * duration * (1 / spread)) / characters.length;

        return (
          <motion.span
            key={`char-${i}-${char}`}
            className={cn(
              "inline-block whitespace-pre [transform-style:preserve-3d]"
            )}
            initial={{
              translateZ: 0,
              scale: 1,
              rotateY: 0,
              color: "var(--base-color)",
            }}
            animate={
              prefersReducedMotion
                ? {}
                : {
                    translateZ: [0, zDistance, 0],
                    translateX: [0, xDistance, 0],
                    translateY: [0, yDistance, 0],
                    scale: [1, scaleDistance, 1],
                    rotateY: [0, rotateYDistance, 0],
                    color: [
                      "var(--base-color)",
                      "var(--base-gradient-color)",
                      "var(--base-color)",
                    ],
                  }
            }
            transition={
              prefersReducedMotion
                ? {}
                : {
                    duration: duration,
                    repeat: Infinity,
                    repeatDelay: (characters.length * 0.05) / spread,
                    delay,
                    ease: "easeInOut",
                    ...transition,
                  }
            }
          >
            {char}
          </motion.span>
        );
      })}
    </MotionComponent>
  );
}
