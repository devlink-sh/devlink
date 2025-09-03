import { useState } from "react";
import { Menu, X } from "lucide-react";
import { cn } from "@/lib/utils";

const Navbar = () => {
  const [isOpen, setIsOpen] = useState(false);

  const navItems = [
    { href: "#intro", label: "Intro" },
    { href: "#features", label: "Features" },
    { href: "#showcase", label: "Showcase" },
    { href: "#tech-stack", label: "Tech Stack" },
    { href: "#impact", label: "Impact" },
  ];

  return (
    <nav className="fixed top-0 left-0 right-0 z-50 bg-white dark:bg-black shadow-sm">
      {/* Content */}
      <div className="max-w-7xl mx-auto px-6 py-4">
        <div className="flex items-center justify-between">
          {/* Logo/Brand - Left side */}
          <div className="flex-shrink-0">
            <div className="text-xl font-medium tracking-tight text-[#EC4899]">
              DevLink CLI
            </div>
          </div>

          {/* Desktop Navigation - Center aligned */}
          <div className="hidden md:flex items-center justify-center flex-1">
            <div className="flex items-center space-x-8">
              {navItems.map((item) => (
                <a
                  key={item.href}
                  href={item.href}
                  className="nav-link text-sm font-medium px-2 py-1 rounded-md hover:bg-[#F472B6]/20 hover:text-[#EC4899] text-gray-700 dark:text-gray-300 transition-colors duration-200"
                >
                  {item.label}
                </a>
              ))}
            </div>
          </div>

          {/* Mobile Menu Button - Right side */}
          <div className="flex-shrink-0">
            <button
              className="md:hidden p-2 rounded-md hover:bg-[#F472B6]/20 hover:text-[#EC4899] text-gray-700 dark:text-gray-300 transition-colors duration-200"
              onClick={() => setIsOpen(!isOpen)}
              aria-label="Toggle menu"
            >
              {isOpen ? <X size={24} /> : <Menu size={24} />}
            </button>
          </div>
        </div>

        {/* Mobile Navigation */}
        {isOpen && (
          <div className="md:hidden mt-4 pb-4">
            <div className="flex flex-col space-y-2 pt-4">
              {navItems.map((item) => (
                <a
                  key={item.href}
                  href={item.href}
                  className="nav-link text-sm font-medium px-3 py-2 rounded-md hover:bg-[#F472B6]/20 hover:text-[#EC4899] text-gray-700 dark:text-gray-300 transition-colors duration-200"
                  onClick={() => setIsOpen(false)}
                >
                  {item.label}
                </a>
              ))}
            </div>
          </div>
        )}
      </div>
    </nav>
  );
};

export default Navbar;
