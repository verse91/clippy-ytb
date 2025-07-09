"use client";

import { useEffect, useRef, useCallback, useTransition } from "react";
import { useState } from "react";
import { cn } from "@/lib/utils";
import {
  ImageIcon,
  FileUp,
  Figma,
  MonitorIcon,
  CircleUserRound,
  ArrowUpIcon,
  Paperclip,
  PlusIcon,
  SendIcon,
  XIcon,
  LoaderIcon,
  Sparkles,
  Command,
} from "lucide-react";
import { motion, AnimatePresence } from "framer-motion";
import * as React from "react";
import { Switch } from "@/components/ui/switch";
import NotiPopup from "./ui/noti-popup";
import {
  Select,
  SelectContent,
  SelectItem,
  SelectTrigger,
  SelectValue,
} from "@/components/ui/select";

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

interface CommandSuggestion {
  icon: React.ReactNode;
  label: string;
  description: string;
  prefix: string;
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

export function BoxChat() {
  const [value, setValue] = useState("");
  const [attachments, setAttachments] = useState<string[]>([]);
  const [isTyping, setIsTyping] = useState(false);
  const [isPending, startTransition] = useTransition();
  const [activeSuggestion, setActiveSuggestion] = useState<number>(-1);
  const [showCommandPalette, setShowCommandPalette] = useState(false);
  const [recentCommand, setRecentCommand] = useState<string | null>(null);
  const [mousePosition, setMousePosition] = useState({ x: 0, y: 0 });
  const { textareaRef, adjustHeight } = useAutoResizeTextarea({
    minHeight: 60,
    maxHeight: 200,
  });
  const [inputFocused, setInputFocused] = useState(false);
  const commandPaletteRef = useRef<HTMLDivElement>(null);

  const commandSuggestions: CommandSuggestion[] = [
    {
      icon: <ImageIcon className="w-4 h-4" />,
      label: "Clone UI",
      description: "Generate a UI from a screenshot",
      prefix: "/clone",
    },
    {
      icon: <Figma className="w-4 h-4" />,
      label: "Import Figma",
      description: "Import a design from Figma",
      prefix: "/figma",
    },
    {
      icon: <MonitorIcon className="w-4 h-4" />,
      label: "Create Page",
      description: "Generate a new web page",
      prefix: "/page",
    },
    {
      icon: <Sparkles className="w-4 h-4" />,
      label: "Improve",
      description: "Improve existing UI design",
      prefix: "/improve",
    },
  ];

  useEffect(() => {
    if (value.startsWith("/") && !value.includes(" ")) {
      setShowCommandPalette(true);

      const matchingSuggestionIndex = commandSuggestions.findIndex((cmd) =>
        cmd.prefix.startsWith(value)
      );

      if (matchingSuggestionIndex >= 0) {
        setActiveSuggestion(matchingSuggestionIndex);
      } else {
        setActiveSuggestion(-1);
      }
    } else {
      setShowCommandPalette(false);
    }
  }, [value]);

  useEffect(() => {
    const handleMouseMove = (e: MouseEvent) => {
      setMousePosition({ x: e.clientX, y: e.clientY });
    };

    window.addEventListener("mousemove", handleMouseMove);
    return () => {
      window.removeEventListener("mousemove", handleMouseMove);
    };
  }, []);

  useEffect(() => {
    const handleClickOutside = (event: MouseEvent) => {
      const target = event.target as Node;
      const commandButton = document.querySelector("[data-command-button]");

      if (
        commandPaletteRef.current &&
        !commandPaletteRef.current.contains(target) &&
        !commandButton?.contains(target)
      ) {
        setShowCommandPalette(false);
      }
    };

    document.addEventListener("mousedown", handleClickOutside);
    return () => {
      document.removeEventListener("mousedown", handleClickOutside);
    };
  }, []);

  const handleKeyDown = (e: React.KeyboardEvent<HTMLTextAreaElement>) => {
    if (showCommandPalette) {
      if (e.key === "ArrowDown") {
        e.preventDefault();
        setActiveSuggestion((prev) =>
          prev < commandSuggestions.length - 1 ? prev + 1 : 0
        );
      } else if (e.key === "ArrowUp") {
        e.preventDefault();
        setActiveSuggestion((prev) =>
          prev > 0 ? prev - 1 : commandSuggestions.length - 1
        );
      } else if (e.key === "Tab" || e.key === "Enter") {
        e.preventDefault();
        if (activeSuggestion >= 0) {
          const selectedCommand = commandSuggestions[activeSuggestion];
          setValue(selectedCommand.prefix + " ");
          setShowCommandPalette(false);

          setRecentCommand(selectedCommand.label);
          setTimeout(() => setRecentCommand(null), 3500);
        }
      } else if (e.key === "Escape") {
        e.preventDefault();
        setShowCommandPalette(false);
      }
    } else if (e.key === "Enter" && !e.shiftKey) {
      e.preventDefault();
      if (value.trim()) {
        handleSendMessage();
      }
    }
  };

  const handleSendMessage = () => {
    if (value.trim()) {
      startTransition(() => {
        setIsTyping(true);
        setTimeout(() => {
          setIsTyping(false);
          setValue("");
          adjustHeight(true);
        }, 3000);
      });
    }
  };

  const handleAttachFile = () => {
    const mockFileName = `file-${Math.floor(Math.random() * 1000)}.pdf`;
    setAttachments((prev) => [...prev, mockFileName]);
  };

  const removeAttachment = (index: number) => {
    setAttachments((prev) => prev.filter((_, i) => i !== index));
  };

  const selectCommandSuggestion = (index: number) => {
    const selectedCommand = commandSuggestions[index];
    setValue(selectedCommand.prefix + " ");
    setShowCommandPalette(false);

    setRecentCommand(selectedCommand.label);
    setTimeout(() => setRecentCommand(null), 2000);
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
                ref={textareaRef}
                value={value}
                onChange={(e) => {
                  setValue(e.target.value);
                  adjustHeight();
                }}
                onKeyDown={handleKeyDown}
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
            </div>

            <div className="p-4 border-t border-white/[0.05] flex items-center justify-between gap-4 pl-8">
              <div className="flex flex-col gap-4 w-full">
                {/* Quality Select */}
                <div className="flex flex-col gap-1 max-w-[220px]">
                  <label
                    className="text-xs text-white/60 mb-1"
                    htmlFor="quality-select"
                  >
                    Quality:
                  </label>
                  <Select>
                    <SelectTrigger className="w-full bg-white/[0.05] border-none text-white/90 text-sm rounded-lg px-3 py-2 focus:ring-0 focus:outline-none cursor-pointer">
                      <SelectValue placeholder="Select quality" />
                    </SelectTrigger>
                    <SelectContent className="bg-[#18181b] border-none text-white/90">
                      <SelectItem value="480p" className="cursor-pointer">
                        720p - HD (FREE)
                      </SelectItem>
                      <SelectItem value="1080p" className="cursor-pointer">
                        1080p - FHD (FREE)
                      </SelectItem>
                      <SelectItem
                        value="720p"
                        className="cursor-pointer text-amber-400"
                        disabled
                      >
                        2k - QHD (PRO)
                      </SelectItem>
                      <SelectItem
                        value="720p"
                        className="cursor-pointer text-amber-400"
                        disabled
                      >
                        4k - UHD (PRO)
                      </SelectItem>
                    </SelectContent>
                  </Select>
                </div>

                {/* Audio and Subtitles Switches */}
                <div className="flex items-center gap-6 justify-between">
                  {/* Audio Switch */}
                  <div className="flex flex-col gap-1 items-start min-w-[120px]">
                    <label
                      className="text-xs text-white/60 mb-1"
                      htmlFor="audio-switch"
                    >
                      Audio:
                    </label>
                    <div className="flex items-center gap-2">
                      <Switch id="audio-switch" className="cursor-pointer" />
                      <span className="text-sm text-white/80">
                        Export audio (.mp3)
                      </span>
                    </div>
                  </div>

                  {/* Subtitles Switch */}
                  <div className="flex flex-col gap-1 items-start pr-4">
                    <label
                      className="text-xs text-white/60 mb-1"
                      htmlFor="subtitles-switch"
                    >
                      Subtitles:
                    </label>
                    <div className="flex items-center gap-2">
                      <Switch
                        id="subtitles-switch"
                        className="cursor-pointer"
                      />
                      <span className="text-sm text-white/80">
                        Export subtitles (.srt)
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
              <a href="" className="underline cursor-pointer hover:text-white/60 transition-colors">
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
            value.startsWith("https://youtube.com") ||
            value.startsWith("https://www.youtube.com") ||
            value.startsWith("https://youtu.be") ||
            value.startsWith("www.youtube.com") ||
            value.startsWith("youtu.be") ||
            value.startsWith("youtube.com")
              ? "Processing"
              : "This is not a Youtube video link"
          }
          showTypingDots={
            value.startsWith("https://youtube.com") ||
            value.startsWith("https://www.youtube.com") ||
            value.startsWith("https://youtu.be") ||
            value.startsWith("www.youtube.com") ||
            value.startsWith("youtu.be") ||
            value.startsWith("youtube.com")
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

const rippleKeyframes = `
@keyframes ripple {
  0% { transform: scale(0.5); opacity: 0.6; }
  100% { transform: scale(2); opacity: 0; }
}
`;

if (typeof document !== "undefined") {
  const style = document.createElement("style");
  style.innerHTML = rippleKeyframes;
  document.head.appendChild(style);
}
