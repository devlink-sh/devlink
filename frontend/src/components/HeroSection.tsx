import { ArrowDown, Copy, Check } from "lucide-react";
import { cn } from "@/lib/utils";
import { useState } from "react";
import { FlipWords } from "@/components/ui/flip-words";
import { Dialog, DialogContent, DialogHeader, DialogTitle } from "@/components/ui/dialog";
import { Button } from "@/components/ui/button";
import { Input } from "@/components/ui/input";
import { useToast } from "@/hooks/use-toast";

const HeroSection = () => {
  const [isDownloadDialogOpen, setIsDownloadDialogOpen] = useState(false);
  const [copied, setCopied] = useState(false);
  const { toast } = useToast();
  const words = ["Secure", "Unified", "Private", "Real‑time"];

  const command = "curl -fsSL http://localhost:8080/downloads/install.sh | bash";

  const copyToClipboard = async () => {
    try {
      await navigator.clipboard.writeText(command);
      setCopied(true);
      toast({
        title: "Command copied!",
        description: "The CLI installation command has been copied to your clipboard.",
      });
      setTimeout(() => setCopied(false), 2000);
    } catch (err) {
      toast({
        title: "Failed to copy",
        description: "Please copy the command manually.",
        variant: "destructive",
      });
    }
  };

  return (
    <section
      id="intro"
      className="relative min-h-screen flex items-center justify-center px-6 overflow-hidden"
    >
      {/* Dot Background */}
      <div
        className={cn(
          "absolute inset-0 pointer-events-none",
          "[background-size:20px_20px]",
          "[background-image:radial-gradient(#d4d4d4_1px,transparent_1px)]",
          "dark:[background-image:radial-gradient(#404040_1px,transparent_1px)]"
        )}
      />

      {/* Radial gradient overlay for faded look */}
      <div className="pointer-events-none absolute inset-0 flex items-center justify-center bg-white [mask-image:radial-gradient(ellipse_at_center,transparent_20%,black)] dark:bg-black"></div>

      {/* Content */}
      <div className="relative z-20 max-w-6xl mx-auto text-center">
        {/* Title */}
        <h1 className="mb-6 tracking-tight text-6xl sm:text-7xl font-semibold text-[#111827] dark:text-white">
          DevLink CLI
        </h1>

        {/* Subheading with FlipWords */}
        <p className="mb-8 text-3xl sm:text-4xl font-medium text-[#6B7280]">
          <span className="text-transparent bg-clip-text bg-gradient-to-r from-[#111827] dark:from-white to-[#EC4899] inline-block align-middle">
            <FlipWords words={words} />
          </span>
          <span className="align-middle">, Peer‑to‑Peer</span>
          <br />
          <span className="align-middle">Developer </span>
          <span className="text-transparent bg-clip-text bg-gradient-to-r from-[#111827] dark:from-white to-[#EC4899] align-middle">
            Workflow
          </span>
        </p>

        {/* Supporting copy */}
        <div className="max-w-3xl mx-auto mb-10">
          <p className="text-xl sm:text-2xl text-[#6B7280]">
            Fast development shouldn't mean slow and insecure collaboration.
            DevLink transforms your workflow with secure peer‑to‑peer
            connections, instant setup, and zero‑config sharing.
          </p>
        </div>

        {/* Highlights */}
        <div className="flex flex-wrap justify-center gap-3 mb-12">
          <span className="px-3 py-1 rounded-full border border-[#E5E7EB] text-sm text-[#6B7280] bg-white/70 dark:bg-transparent dark:border-neutral-800">
            Encrypted env & tunnels
          </span>
          <span className="px-3 py-1 rounded-full border border-[#E5E7EB] text-sm text-[#6B7280] bg-white/70 dark:bg-transparent dark:border-neutral-800">
            Real‑time git sync
          </span>
          <span className="px-3 py-1 rounded-full border border-[#E5E7EB] text-sm text-[#6B7280] bg-white/70 dark:bg-transparent dark:border-neutral-800">
            Native pair programming
          </span>
        </div>

        {/* CTAs */}
        <div className="mb-12 flex items-center justify-center gap-3">
          <button
            onClick={() => setIsDownloadDialogOpen(true)}
            className="inline-block px-8 py-3.5 bg-[#EC4899] hover:bg-[#F472B6] text-white font-medium rounded-lg border-2 border-[#EC4899] hover:border-[#F472B6] transition-all duration-200 shadow-lg hover:shadow-xl"
          >
            Download CLI
          </button>
          <a
            href="#features"
            className="inline-block px-8 py-3.5 rounded-lg border-2 border-[#E5E7EB] text-[#111827] dark:text-white hover:border-[#EC4899] hover:text-[#EC4899] transition-colors"
          >
            Explore Features
          </a>
        </div>

        <div className="animate-bounce">
          <ArrowDown size={32} className="mx-auto text-[#EC4899]" />
        </div>
      </div>

      {/* Download CLI Dialog */}
      <Dialog open={isDownloadDialogOpen} onOpenChange={setIsDownloadDialogOpen}>
        <DialogContent className="sm:max-w-md">
          <DialogHeader>
            <DialogTitle className="text-center">Download DevLink CLI</DialogTitle>
          </DialogHeader>
          <div className="space-y-4">
            <p className="text-sm text-muted-foreground text-center">
              Run this command in your terminal to install DevLink CLI:
            </p>
            <div className="flex items-center space-x-2">
              <Input
                value={command}
                readOnly
                className="font-mono text-sm"
              />
              <Button
                size="sm"
                variant="outline"
                onClick={copyToClipboard}
                className="shrink-0"
              >
                {copied ? (
                  <Check className="h-4 w-4" />
                ) : (
                  <Copy className="h-4 w-4" />
                )}
              </Button>
            </div>
            <div className="text-xs text-muted-foreground text-center">
              <p>Make sure you have curl installed on your system.</p>
              <p className="mt-1">This will download and run the installation script.</p>
            </div>
          </div>
        </DialogContent>
      </Dialog>
    </section>
  );
};

export default HeroSection;
