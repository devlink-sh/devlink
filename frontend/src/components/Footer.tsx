import { Github, Twitter, Linkedin, Mail, Globe, Download } from "lucide-react";

const Footer = () => {
  const currentYear = new Date().getFullYear();

  return (
    <footer className="relative bg-white dark:bg-black border-t border-[#E5E7EB] dark:border-neutral-800 mt-20">
      <div className="max-w-7xl mx-auto px-6 py-14">
        <div className="grid grid-cols-1 md:grid-cols-2 lg:grid-cols-4 gap-10 mb-10">
          {/* Brand */}
          <div className="lg:col-span-2">
            <div className="text-2xl font-bold tracking-tight text-[#EC4899] mb-3">
              DevLink CLI
            </div>
            <p className="text-[#6B7280] leading-relaxed max-w-md mb-6">
              Secure, unified, peer‑to‑peer collaboration for modern developers.
            </p>
            <a
              href="#download"
              className="inline-flex items-center gap-2 bg-[#EC4899] text-white px-4 py-2 rounded-lg border-2 border-[#EC4899] hover:bg-[#F472B6] hover:border-[#F472B6] transition-colors font-medium"
            >
              <Download size={16} />
              Download CLI
            </a>
          </div>

          {/* Quick Links */}
          <div>
            <h3 className="text-sm font-semibold mb-4 text-[#111827] dark:text-white">
              Quick Links
            </h3>
            <ul className="space-y-3">
              <li>
                <a href="#intro" className="text-[#6B7280] hover:text-[#EC4899] transition-colors">
                  Introduction
                </a>
              </li>
              <li>
                <a href="#features" className="text-[#6B7280] hover:text-[#EC4899] transition-colors">
                  Features
                </a>
              </li>
              <li>
                <a href="#showcase" className="text-[#6B7280] hover:text-[#EC4899] transition-colors">
                  Showcase
                </a>
              </li>
              <li>
                <a href="#tech-stack" className="text-[#6B7280] hover:text-[#EC4899] transition-colors">
                  Tech Stack
                </a>
              </li>
              <li>
                <a href="#impact" className="text-[#6B7280] hover:text-[#EC4899] transition-colors">
                  Impact
                </a>
              </li>
            </ul>
          </div>

          {/* Resources */}
          <div>
            <h3 className="text-sm font-semibold mb-4 text-[#111827] dark:text-white">
              Resources
            </h3>
            <ul className="space-y-3">
              <li>
                <a href="#download" className="text-[#6B7280] hover:text-[#EC4899] transition-colors">
                  Download
                </a>
              </li>
              <li>
                <a href="https://docs.devlink.dev" className="text-[#6B7280] hover:text-[#EC4899] transition-colors">
                  Documentation
                </a>
              </li>
              <li>
                <a href="https://github.com/devlink-cli" className="text-[#6B7280] hover:text-[#EC4899] transition-colors">
                  GitHub
                </a>
              </li>
              <li>
                <a href="mailto:hello@devlink.dev" className="text-[#6B7280] hover:text-[#EC4899] transition-colors">
                  Contact
                </a>
              </li>
            </ul>
          </div>
        </div>

        {/* Bottom */}
        <div className="pt-8 border-t border-[#E5E7EB] dark:border-neutral-800">
          <div className="flex flex-col md:flex-row items-center justify-between gap-4">
            <p className="text-[#6B7280]">© {currentYear} DevLink CLI. All rights reserved.</p>
            <div className="flex items-center gap-4">
              <a href="https://github.com/devlink-cli" className="text-[#6B7280] hover:text-[#EC4899] transition-colors" aria-label="GitHub">
                <Github size={20} />
              </a>
              <a href="https://twitter.com/devlinkcli" className="text-[#6B7280] hover:text-[#EC4899] transition-colors" aria-label="Twitter">
                <Twitter size={20} />
              </a>
              <a href="https://linkedin.com/company/devlink-cli" className="text-[#6B7280] hover:text-[#EC4899] transition-colors" aria-label="LinkedIn">
                <Linkedin size={20} />
              </a>
              <a href="mailto:hello@devlink.dev" className="text-[#6B7280] hover:text-[#EC4899] transition-colors" aria-label="Email">
                <Mail size={20} />
              </a>
              <a href="https://devlink.dev" className="text-[#6B7280] hover:text-[#EC4899] transition-colors" aria-label="Website">
                <Globe size={20} />
              </a>
            </div>
          </div>
        </div>
      </div>
    </footer>
  );
};

export default Footer;
