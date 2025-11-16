"use client";
import React, { type JSX } from "react";
import { motion } from "framer-motion";
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
};

export function TextShimmerWave({
    children,
    as: Component = "p",
    className,
    duration = 3,
}: TextShimmerWave) {
    const MotionComponent = motion.create(
        Component as keyof JSX.IntrinsicElements
    );

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
                "text-shimmer-wave", // Use CSS animation instead
                className
            )}
            style={{
                color: "var(--base-color)",
                animationDuration: `${duration}s`,
            }}
        >
            {prefersReducedMotion ? (
                children
            ) : (
                <span className="shimmer-text-content">{children}</span>
            )}
        </MotionComponent>
    );
}
