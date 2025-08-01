"use client";

import { useEffect, useRef, useCallback, useTransition } from "react";
import { useState } from "react";
import { cn } from "@/lib/utils";
import { SendIcon, LoaderIcon } from "lucide-react";
import { motion, AnimatePresence } from "framer-motion";
import * as React from "react";
import { Switch } from "@/components/ui/switch";
import NotiPopup from "../ui/noti-popup";
import {
  Select,
  SelectContent,
  SelectItem,
  SelectTrigger,
  SelectValue,
} from "@/components/ui/select";

// Helper functions for YouTube URL validation
const isYouTubeUrl = (url: string): boolean => {
  const youtubePatterns = [
    "https://youtube.com",
    "https://www.youtube.com",
    "https://youtu.be",
    "www.youtube.com",
    "youtu.be",
    "youtube.com",
  ];
  return youtubePatterns.some((pattern) => url.startsWith(pattern));
};

const hasPlaylist = (url: string): boolean => {
  return url.includes("list=");
};

interface UseAutoResizeTextareaProps {
  minHeight: number;
  maxHeight?: number;
}

function useAutoResizeTextarea({
  minHeight,
  maxHeight,
}: UseAutoResizeTextareaProps) {
  const textareaRef = useRef<HTMLTextAreaElement>(null);

  const adjustHeight = useCallback(
    (reset?: boolean) => {
      const textarea = textareaRef.current;
      if (!textarea) return;

      if (reset) {
        textarea.style.height = `${minHeight}px`;
        return;
      }

      textarea.style.height = `${minHeight}px`;
      const newHeight = Math.max(
        minHeight,
        Math.min(textarea.scrollHeight, maxHeight ?? Number.POSITIVE_INFINITY)
      );

      textarea.style.height = `${newHeight}px`;
    },
    [minHeight, maxHeight]
  );

  useEffect(() => {
    const textarea = textareaRef.current;
    if (textarea) {
      textarea.style.height = `${minHeight}px`;
    }
  }, [minHeight]);

  useEffect(() => {
    const handleResize = () => adjustHeight();
    window.addEventListener("resize", handleResize);
    return () => window.removeEventListener("resize", handleResize);
  }, [adjustHeight]);

  return { textareaRef, adjustHeight };
}

interface TextareaProps
  extends React.TextareaHTMLAttributes<HTMLTextAreaElement> {
  containerClassName?: string;
  showRing?: boolean;
}

const Textarea = React.forwardRef<HTMLTextAreaElement, TextareaProps>(
  ({ className, containerClassName, showRing = true, ...props }, ref) => {
    const [isFocused, setIsFocused] = React.useState(false);

    return (
      <div className={cn("relative", containerClassName)}>
        <textarea
          className={cn(
            "flex min-h-[80px] w-full rounded-md border border-input bg-background px-3 py-2 text-sm",
            "transition-all duration-200 ease-in-out",
            "placeholder:text-muted-foreground",
            "disabled:cursor-not-allowed disabled:opacity-50",
            showRing
              ? "focus-visible:outline-none focus-visible:ring-0 focus-visible:ring-offset-0"
              : "",
            className
          )}
          ref={ref}
          onFocus={() => setIsFocused(true)}
          onBlur={() => setIsFocused(false)}
          {...props}
        />

        {showRing && isFocused && (
          <motion.span
            className="absolute inset-0 rounded-md pointer-events-none ring-2 ring-offset-0 ring-violet-500/30"
            initial={{ opacity: 0 }}
            animate={{ opacity: 1 }}
            exit={{ opacity: 0 }}
            transition={{ duration: 0.2 }}
          />
        )}

        {props.onChange && (
          <div
            className="absolute bottom-2 right-2 opacity-0 w-2 h-2 bg-violet-500 rounded-full"
            style={{
              animation: "none",
            }}
            id="textarea-ripple"
          />
        )}
      </div>
    );
  }
);
Textarea.displayName = "Textarea";

const rippleKeyframes = `
@keyframes ripple {
  0% { transform: scale(0.5); opacity: 0.6; }
  100% { transform: scale(2); opacity: 0; }
}
`;

export function BoxChat() {
  const [value, setValue] = useState("");
  const [isTyping, setIsTyping] = useState(false);
  const [inputFocused, setInputFocused] = useState(false);
  const [mousePosition, setMousePosition] = useState({ x: 0, y: 0 });
  const [selectedOption, setSelectedOption] = useState("auto");
  const [thumbnailEnabled, setThumbnailEnabled] = useState(false);
  const [isPending, startTransition] = useTransition();
  const { textareaRef: autoResizeTextareaRef, adjustHeight } =
    useAutoResizeTextarea({
      minHeight: 60,
      maxHeight: 200,
    });

  // Use the auto-resize textarea ref
  const textareaRefToUse = autoResizeTextareaRef;
  const [isChecked, setIsChecked] = useState(true);

  // Inject ripple styles on client side only
  useEffect(() => {
    if (typeof document !== "undefined") {
      // Check if the style already exists to prevent duplicates
      const existingStyle = document.querySelector(
        "style[data-ripple-keyframes]"
      );
      if (!existingStyle) {
        const style = document.createElement("style");
        style.setAttribute("data-ripple-keyframes", "true");
        style.innerHTML = rippleKeyframes;
        document.head.appendChild(style);
      }
    }

    // Cleanup function to remove the style when component unmounts
    return () => {
      if (typeof document !== "undefined") {
        const style = document.querySelector("style[data-ripple-keyframes]");
        if (style) {
          style.remove();
        }
      }
    };
  }, []);

  useEffect(() => {
    const handleMouseMove = (e: MouseEvent) => {
      setMousePosition({ x: e.clientX, y: e.clientY });
    };

    window.addEventListener("mousemove", handleMouseMove);
    return () => {
      window.removeEventListener("mousemove", handleMouseMove);
    };
  }, []);

  const handleSendMessage = () => {
    if (value.trim()) {
      startTransition(() => {
        setIsTyping(true);
        setTimeout(() => {
          setIsTyping(false);
          setValue("");
          adjustHeight(true);
        }, 2000);
      });
    }
  };

  return (
    <div className="flex flex-col w-full items-center justify-center bg-transparent text-white p-6 relative overflow-hidden select-none">
      <div className=" max-w-2xl mx-auto relative">
        <motion.div
          className="relative z-10 space-y-12"
          initial={{ opacity: 0, y: 20 }}
          animate={{ opacity: 1, y: 0 }}
          transition={{ duration: 0.6, ease: "easeOut" }}
        >
          <div className="text-center space-y-3">
            <motion.div
              initial={{ opacity: 0, y: 10 }}
              animate={{ opacity: 1, y: 0 }}
              transition={{ delay: 0.2, duration: 0.5 }}
              className="inline-block"
            >
              <h1 className="text-3xl font-medium tracking-tight bg-clip-text text-transparent bg-gradient-to-r from-white/90 to-white/40 pb-1">
                What do you wanna clip?
              </h1>
              <motion.div
                className="h-px bg-gradient-to-r from-transparent via-white/20 to-transparent"
                initial={{ width: 0, opacity: 0 }}
                animate={{ width: "100%", opacity: 1 }}
                transition={{ delay: 0.5, duration: 0.8 }}
              />
            </motion.div>
          </div>

          <motion.div
            className="relative backdrop-blur-2xl bg-white/[0.02] rounded-2xl border border-white/[0.05] shadow-2xl"
            initial={{ scale: 0.98 }}
            animate={{ scale: 1 }}
            transition={{ delay: 0.1 }}
          >
            <div className="p-4 flex pr-8">
              <Textarea
                ref={textareaRefToUse}
                value={value}
                onChange={(e) => {
                  setValue(e.target.value);
                  adjustHeight();
                }}
                onFocus={() => setInputFocused(true)}
                onBlur={() => setInputFocused(false)}
                placeholder="Paste video url here..."
                containerClassName="w-full"
                className={cn(
                  "w-full px-4 py-3",
                  "resize-none",
                  "bg-transparent",
                  "border-none",
                  "text-white/90 text-sm",
                  "focus:outline-none",
                  "placeholder:text-white/20",
                  "min-h-[60px]"
                )}
                style={{
                  overflow: "hidden",
                }}
                showRing={false}
              />
              <motion.button
                type="button"
                onClick={handleSendMessage}
                whileHover={{ scale: 1.01 }}
                whileTap={{ scale: 0.98 }}
                disabled={isTyping || !value.trim()}
                aria-label="Send message"
                aria-describedby="send-button-description"
                className={cn(
                  "px-4 py-2 rounded-lg text-sm font-medium transition-all h-10",
                  "flex items-center gap-2",
                  value.trim()
                    ? "bg-white text-[#0A0A0B] shadow-lg shadow-white/10 cursor-pointer"
                    : "bg-white/[0.05] text-white/40"
                )}
              >
                {isTyping ? (
                  <LoaderIcon className="w-4 h-4 animate-[spin_2s_linear_infinite]" />
                ) : (
                  <SendIcon className="w-3 h-3" />
                )}
              </motion.button>
              <div id="send-button-description" className="sr-only">
                {isTyping
                  ? "Processing request..."
                  : "Send the video URL for processing"}
              </div>
            </div>

            <div className="p-3 sm:p-4 border-t border-white/[0.05] flex flex-col sm:flex-row items-start sm:items-center justify-between gap-4 pl-6 sm:pl-8">
              <div className="flex flex-col gap-4 w-full">
                {/* Quality Select */}
                <div className="flex flex-col sm:flex-row items-start sm:items-center gap-4 sm:gap-6 justify-start">
                  <div className="flex flex-col gap-1 w-12/13 sm:max-w-[200px]">
                    <label
                      className="text-xs text-white/60 mb-1"
                      htmlFor="quality-select"
                    >
                      Option:
                    </label>
                    <Select
                      value={selectedOption}
                      onValueChange={(val) => {
                        setSelectedOption(val);
                        if (val === "audio-only" && thumbnailEnabled) {
                          setThumbnailEnabled(false);
                        }
                      }}
                      aria-label="Select video quality option"
                    >
                      <SelectTrigger
                        id="quality-select"
                        className="w-full bg-white/[0.05] border-none text-white/90 text-sm rounded-lg px-3 py-2 focus:ring-0 focus:outline-none cursor-pointer"
                        style={{
                          minWidth: 0,
                          width: "100%",
                          minHeight: 40,
                        }}
                      >
                        <SelectValue />
                      </SelectTrigger>
                      <SelectContent className="bg-[#18181b] border-none text-white/90">
                        <SelectItem value="auto" className="cursor-pointer">
                          Best quality (Up to 1080p)
                        </SelectItem>
                        <SelectItem
                          value="audio-only"
                          className="cursor-pointer"
                        >
                          Audio only
                        </SelectItem>
                        <SelectItem value="mute" className="cursor-pointer">
                          Muted video
                        </SelectItem>
                      </SelectContent>
                    </Select>
                  </div>
                </div>

                {/* Audio and Subtitles Switches */}
                <div className="flex flex-col sm:flex-row gap-4 sm:gap-6 justify-between">
                  {/* Audio Switch */}
                  <div className="flex flex-col gap-1 items-start min-w-[120px]">
                    <label
                      className="text-xs text-white/60 mb-1"
                      htmlFor="audio-switch"
                    >
                      SponsorBlock:
                    </label>
                    <div className="flex items-center gap-2">
                      <Switch
                        id="audio-switch"
                        className="cursor-pointer"
                        checked={isChecked}
                        onCheckedChange={setIsChecked}
                        aria-label="Auto block sponsor segments"
                        aria-describedby="sponsor-description"
                      />
                      <span
                        id="sponsor-description"
                        className="text-sm text-white/80"
                      >
                        Auto block sponsor
                      </span>
                    </div>
                  </div>

                  {/* Thumbnail Switch - moved under SponsorBlock for mobile */}
                  <div className="flex flex-col gap-1 items-start pr-4 sm:pr-4">
                    <label
                      className={cn(
                        "text-xs mb-1",
                        selectedOption === "audio-only"
                          ? "text-white/30"
                          : "text-white/60"
                      )}
                      htmlFor="thumbnail-switch"
                    >
                      Thumbnail:
                    </label>
                    <div className="flex items-center gap-2">
                      <Switch
                        id="thumbnail-switch"
                        className="cursor-pointer"
                        disabled={selectedOption === "audio-only"}
                        checked={thumbnailEnabled}
                        onCheckedChange={setThumbnailEnabled}
                        aria-label="Auto import subtitles"
                        aria-describedby="thumbnail-description"
                      />
                      <span
                        id="thumbnail-description"
                        className={cn(
                          "text-sm",
                          selectedOption === "audio-only"
                            ? "text-white/30"
                            : "text-white/80"
                        )}
                      >
                        Auto import thumbnail
                      </span>
                    </div>
                  </div>
                </div>
              </div>
            </div>
          </motion.div>

          <div className="flex flex-wrap items-center justify-center -mt-6">
            <motion.p
              className="text-sm text-white/40"
              initial={{ opacity: 0 }}
              animate={{ opacity: 1 }}
              transition={{ delay: 0.3 }}
            >
              Support youtube links from youtube.com and youtu.be.{" "}
              <a
                target="_blank"
                href="https://github.com/verse91/clippy-ytb/issues/new?template=bug-report.yml"
                className="underline cursor-pointer hover:text-white/60 transition-colors"
              >
                Report bugs?
              </a>
            </motion.p>
          </div>
        </motion.div>
      </div>

      {isTyping ? (
        <NotiPopup
          isVisible={isTyping}
          icon={<i className="bxl bx-youtube" style={{ color: "#ff0000" }}></i>}
          text={
            hasPlaylist(value) && isYouTubeUrl(value)
              ? "Playlist is not supported"
              : isYouTubeUrl(value)
              ? "Processing"
              : "This is not a Youtube video link"
          }
          showTypingDots={!hasPlaylist(value) && isYouTubeUrl(value)}
          widthClassName={
            hasPlaylist(value) && isYouTubeUrl(value)
              ? "w-[44vw] max-w-md sm:w-auto sm:max-w-none"
              : isYouTubeUrl(value)
              ? ""
              : "w-[53vw] max-w-md sm:w-auto sm:max-w-none"
          }
        />
      ) : null}

      {inputFocused && (
        <motion.div
          className="fixed w-[50rem] h-[50rem] rounded-full pointer-events-none z-0 opacity-[0.02] bg-gradient-to-r from-violet-500 via-fuchsia-500 to-indigo-500 blur-[96px]"
          animate={{
            x: mousePosition.x - 400,
            y: mousePosition.y - 400,
          }}
          transition={{
            type: "spring",
            damping: 25,
            stiffness: 150,
            mass: 0.5,
          }}
        />
      )}
    </div>
  );
}

interface ActionButtonProps {
  icon: React.ReactNode;
  label: string;
}

function ActionButton({ icon, label }: ActionButtonProps) {
  const [isHovered, setIsHovered] = useState(false);

  return (
    <motion.button
      type="button"
      whileHover={{ scale: 1.05, y: -2 }}
      whileTap={{ scale: 0.97 }}
      onHoverStart={() => setIsHovered(true)}
      onHoverEnd={() => setIsHovered(false)}
      className="flex items-center gap-2 px-4 py-2 bg-neutral-900 hover:bg-neutral-800 rounded-full border border-neutral-800 text-neutral-400 hover:text-white transition-all relative overflow-hidden group"
    >
      <div className="relative z-10 flex items-center gap-2">
        {icon}
        <span className="text-xs relative z-10">{label}</span>
      </div>

      <AnimatePresence>
        {isHovered && (
          <motion.div
            className="absolute inset-0 bg-gradient-to-r from-violet-500/10 to-indigo-500/10"
            initial={{ opacity: 0 }}
            animate={{ opacity: 1 }}
            exit={{ opacity: 0 }}
            transition={{ duration: 0.2 }}
          />
        )}
      </AnimatePresence>

      <motion.span
        className="absolute bottom-0 left-0 w-full h-0.5 bg-gradient-to-r from-violet-500 to-indigo-500"
        initial={{ width: 0 }}
        whileHover={{ width: "100%" }}
        transition={{ duration: 0.3 }}
      />
    </motion.button>
  );
}
