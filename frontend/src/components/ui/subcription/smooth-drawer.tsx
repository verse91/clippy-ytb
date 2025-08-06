"use client";
import * as React from "react";

import { Button } from "@/components/ui/button";
import {
  Drawer,
  DrawerClose,
  DrawerContent,
  DrawerFooter,
  DrawerHeader,
  DrawerTitle,
  DrawerTrigger,
} from "@/components/ui/subcription/drawer";
import Image from "next/image";
import Link from "next/link";
import { motion } from "motion/react";
import { Fingerprint, Check } from "lucide-react";

interface CreditOption {
  id: string;
  credits: number;
  price: number;
  originalPrice: number;
  discount: number;
  popular?: boolean;
}

const creditOptions: CreditOption[] = [
  {
    id: "basic",
    credits: 60,
    price: 5.99,
    originalPrice: 5.99,
    discount: 0,
  },
  {
    id: "popular",
    credits: 250,
    price: 14.99,
    originalPrice: 24.95,
    discount: 40,
    popular: true,
  },
  {
    id: "premium",
    credits: 600,
    price: 29.99,
    originalPrice: 59.99,
    discount: 50,
  },
];

const features = [
  "High quality 1080p exports",
  "Credits stack with existing balance",
  "Credits never expire",
  "1 credit = 1 minute of 1080p processing",
];

interface PriceTagProps {
  option: CreditOption;
  isSelected: boolean;
  onSelect: (option: CreditOption) => void;
}

function PriceTag({ option, isSelected, onSelect }: PriceTagProps) {
  return (
    <div
      className={`relative p-3 rounded-lg border-2 cursor-pointer transition-all duration-200 ${
        isSelected
          ? "border-rose-500 bg-rose-50 dark:bg-rose-950/20"
          : "border-zinc-200 dark:border-zinc-800 hover:border-zinc-300 dark:hover:border-zinc-700"
      }`}
      onClick={() => onSelect(option)}
    >
      {option.popular && (
        <div className="absolute -top-1.5 left-1/2 transform -translate-x-1/2">
          <span className="bg-rose-500 text-white text-xs font-semibold px-2.5 py-0.5 rounded-full font-['SF-Pro-Display']">
            Popular
          </span>
        </div>
      )}

      <div className="flex items-center justify-between mb-1.5">
        <div className="flex items-baseline gap-1.5">
          <span className="text-xl font-bold bg-gradient-to-br from-zinc-900 to-zinc-700 dark:from-white dark:to-zinc-300 bg-clip-text text-transparent font-['SF-Pro-Display']">
            ${option.price}
          </span>
          {option.discount > 0 && (
            <span className="text-xs line-through text-zinc-400 dark:text-zinc-500 font-['SF-Pro-Display']">
              ${option.originalPrice}
            </span>
          )}
        </div>
        {isSelected && (
          <div className="w-4 h-4 bg-rose-500 rounded-full flex items-center justify-center">
            <Check className="w-2.5 h-2.5 text-white" />
          </div>
        )}
      </div>

      <div className="flex items-center justify-between">
        <span className="text-sm font-semibold text-zinc-900 dark:text-zinc-100 font-['SF-Pro-Display']">
          {option.credits} Credits
        </span>
        {option.discount > 0 && (
          <span className="text-xs text-rose-600 dark:text-rose-400 font-['SF-Pro-Display'] font-semibold">
            {option.discount}% OFF
          </span>
        )}
      </div>
    </div>
  );
}

interface DrawerDemoProps extends React.HTMLAttributes<HTMLDivElement> {
  title?: string;
  description?: string;
  primaryButtonText?: string;
  secondaryButtonText?: string;
  onPrimaryAction?: () => void;
  onSecondaryAction?: () => void;
  isUserLoggedIn?: boolean;
  userCredits?: number;
  trigger?: React.ReactNode;
}

const drawerVariants = {
  hidden: {
    y: "100%",
    opacity: 0,
    rotateX: 5,
    transition: {
      type: "spring",
      stiffness: 300,
      damping: 30,
    },
  },
  visible: {
    y: 0,
    opacity: 1,
    rotateX: 0,
    transition: {
      type: "spring",
      stiffness: 300,
      damping: 30,
      mass: 0.8,
      staggerChildren: 0.07,
      delayChildren: 0.2,
    },
  },
};

const itemVariants = {
  hidden: {
    y: 20,
    opacity: 0,
    transition: {
      type: "spring",
      stiffness: 300,
      damping: 30,
    },
  },
  visible: {
    y: 0,
    opacity: 1,
    transition: {
      type: "spring",
      stiffness: 300,
      damping: 30,
      mass: 0.8,
    },
  },
};

export default function SmoothDrawer({
  title = "Clippy - Pro",
  description = "High quality 1080p exports Credits stack with existing balance • Credits never expire • 1 credit = 1 minute of 1080p processing",
  primaryButtonText = "Buy",
  secondaryButtonText = "Maybe Later",
  onSecondaryAction,
  isUserLoggedIn = false,
  userCredits = 0,
  trigger,
}: DrawerDemoProps) {
  const [selectedOption, setSelectedOption] = React.useState<CreditOption>(
    () => {
      const popularOption = creditOptions.find((option) => option.popular);
      return popularOption || creditOptions[0] || creditOptions[1];
    }
  );

  const handleSecondaryClick = () => {
    onSecondaryAction?.();
  };

  return (
    <Drawer>
      <DrawerTrigger asChild>
        {trigger || <Button variant="outline">Open Drawer</Button>}
      </DrawerTrigger>
      {/* Make the DrawerContent a bit wider, but keep the color unchanged */}
      <DrawerContent className="max-w-[400px] mx-auto p-6 rounded-2xl shadow-xl">
        <motion.div
          variants={drawerVariants as any}
          initial="hidden"
          animate="visible"
          className="mx-auto w-full max-w-[400px] space-y-5"
        >
          <motion.div variants={itemVariants as any}>
            <DrawerHeader className="px-0 space-y-2">
              <DrawerTitle className="text-2xl font-semibold flex items-center gap-2.5 tracking-tighter font-['SF-Pro-Display']">
                <motion.div variants={itemVariants as any}>
                  <div className="p-1.5 rounded-xl bg-gradient-to-br from-zinc-100 to-zinc-200 dark:from-zinc-800 dark:to-zinc-900 shadow-inner">
                    <Image
                      src="/assets/icons/logo-no-bg.png"
                      alt="Logo"
                      width={32}
                      height={32}
                      className="hidden dark:block"
                    />
                    <Image
                      src="/assets/icons/logo-no-bg.png"
                      alt="Logo"
                      width={32}
                      height={32}
                      className="block dark:hidden"
                    />
                  </div>
                </motion.div>
                <motion.span variants={itemVariants as any}>
                  {title}
                </motion.span>
              </DrawerTitle>
              {isUserLoggedIn && (
                <motion.div variants={itemVariants as any}></motion.div>
              )}
              <motion.div variants={itemVariants as any}>
                <ul className="space-y-1.5">
                  {features.map((feature, index) => (
                    <li key={index} className="flex items-center gap-2">
                      <Check className="w-4 h-4 text-green-500" />
                      <span className="text-xs text-zinc-600 dark:text-zinc-400 font-['SF-Pro-Display']">
                        {feature}
                      </span>
                    </li>
                  ))}
                </ul>
              </motion.div>
            </DrawerHeader>
          </motion.div>

          <motion.div variants={itemVariants as any}>
            <div className="space-y-2.5">
              {creditOptions.map((option) => (
                <PriceTag
                  key={option.id}
                  option={option}
                  isSelected={selectedOption.id === option.id}
                  onSelect={setSelectedOption}
                />
              ))}
            </div>
          </motion.div>

          <motion.div variants={itemVariants as any}>
            <DrawerFooter className="flex flex-col gap-2.5 px-0">
              <div className="w-full">
                <Link
                  href="https://stripe.com"
                  target="_blank"
                  rel="noopener noreferrer"
                  className="group w-full relative overflow-hidden inline-flex items-center justify-center h-10 rounded-lg bg-gradient-to-r from-rose-500 to-pink-500 dark:from-rose-600 dark:to-pink-600 text-white text-xs font-semibold tracking-wide shadow-lg shadow-rose-500/20 transition-all duration-500 hover:shadow-xl hover:shadow-rose-500/30 hover:from-rose-600 hover:to-pink-600 dark:hover:from-rose-500 dark:hover:to-pink-500 font-['SF-Pro-Display']"
                  aria-label={`Purchase ${selectedOption.credits} credits`}
                >
                  <motion.span
                    className="absolute inset-0 bg-gradient-to-r from-transparent via-white/20 to-transparent translate-x-[-200%]"
                    whileHover={{
                      x: ["-200%", "200%"],
                    }}
                    transition={{
                      duration: 1.5,
                      ease: "easeInOut",
                      repeat: 0,
                    }}
                  />
                  <motion.div
                    initial={{ opacity: 0 }}
                    animate={{ opacity: 1 }}
                    transition={{ duration: 0.3 }}
                    className="relative flex items-center gap-2 tracking-tighter"
                  >
                    {primaryButtonText} {selectedOption.credits} credits
                    <motion.div
                      animate={{
                        rotate: [0, 15, -15, 0],
                        y: [0, -2, 2, 0],
                      }}
                      transition={{
                        duration: 2,
                        ease: "easeInOut",
                        repeat: Number.POSITIVE_INFINITY,
                        repeatDelay: 1,
                      }}
                    >
                      <Fingerprint className="w-3.5 h-3.5" />
                    </motion.div>
                  </motion.div>
                </Link>
              </div>
              <DrawerClose asChild>
                <Button
                  variant="outline"
                  onClick={handleSecondaryClick}
                  className="w-full h-10 rounded-lg border-zinc-200 dark:border-zinc-800 hover:bg-zinc-100 dark:hover:bg-zinc-800/80 text-xs font-semibold transition-colors tracking-tighter font-['SF-Pro-Display'] cursor-pointer"
                >
                  {secondaryButtonText}
                </Button>
              </DrawerClose>
            </DrawerFooter>
          </motion.div>
        </motion.div>
      </DrawerContent>
    </Drawer>
  );
}
