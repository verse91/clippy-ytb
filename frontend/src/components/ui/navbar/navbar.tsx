"use client";
import {
  Navbar,
  NavBody,
  NavItems,
  MobileNav,
  NavbarLogo,
  NavbarButton,
  MobileNavHeader,
  MobileNavToggle,
  MobileNavMenu,
} from "@/components/ui/navbar/resizable-navbar";
import { BoxChat } from "@/components/homepage/box-chat";
import { useState } from "react";

export function NavbarMain() {
  const navItems = [
    {
      name: "Features",
      link: "#features",
    },
    {
      name: "Pricing",
      link: "#pricing",
    },
    {
      name: "Contact",
      link: "#contact",
    },
  ];

  const [isMobileMenuOpen, setIsMobileMenuOpen] = useState(false);

  return (
    <div className="relative w-full">
      <Navbar>
        {/* Desktop Navigation */}
        <NavBody>
          <NavbarLogo />
          <NavItems items={navItems} />
          <div className="flex items-center gap-4">
            <NavbarButton
              variant="secondary"
            >
              Login
            </NavbarButton>
            <NavbarButton
              target="_blank"
              rel="noopener noreferrer"
              title="Star it on GitHub ⭐"
              href="https://github.com/verse91/clippy-ytb"
              onClick={() => setIsMobileMenuOpen(false)}
              variant="primary"
              className="w-full flex items-center gap-2"
            >
              <i className="bxl bx-github text-2xl text-black transition-all group-hover:text-gray-300 group-hover:scale-110"></i>
              Star on GitHub
            </NavbarButton>
          </div>
        </NavBody>

        {/* Mobile Navigation */}
        <MobileNav>
          <MobileNavHeader>
            <NavbarLogo />
            <MobileNavToggle
              isOpen={isMobileMenuOpen}
              onClick={() => setIsMobileMenuOpen(!isMobileMenuOpen)}
            />
          </MobileNavHeader>

          <MobileNavMenu
            isOpen={isMobileMenuOpen}
            onClose={() => setIsMobileMenuOpen(false)}
          >
            {navItems.map((item, idx) => (
              <a
                key={`mobile-link-${idx}`}
                href={item.link}
                onClick={() => setIsMobileMenuOpen(false)}
                className="relative text-neutral-600 dark:text-neutral-300"
              >
                <span className="block">{item.name}</span>
              </a>
            ))}
            <div className="flex w-full flex-col gap-4">
              <NavbarButton
                onClick={() => setIsMobileMenuOpen(false)}
                variant="primary"
                className="w-full"
              >
                Login
              </NavbarButton>
              <NavbarButton
                target="_blank"
                rel="noopener noreferrer"
                title="Star it on GitHub ⭐"
                href="https://github.com/verse91/clippy-ytb"
                onClick={() => setIsMobileMenuOpen(false)}
                variant="primary"
                className="w-full"
              >
                Star on GitHub
              </NavbarButton>
            </div>
          </MobileNavMenu>
        </MobileNav>
      </Navbar>
        {/* <BoxChat /> */}
    </div>
  );
}
