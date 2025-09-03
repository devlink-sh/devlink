import { cn } from "@/lib/utils";
import {
  IconShield,
  IconBolt,
  IconUsers,
  IconLock,
  IconWorld,
  IconCode,
} from "@tabler/icons-react";

const FeaturesSection = () => {
  const features = [
    {
      title: "Secure by Design",
      description:
        "End-to-end encryption with zero-trust architecture for maximum security.",
      icon: <IconShield />,
    },
    {
      title: "Lightning Fast",
      description:
        "Peer-to-peer connections eliminate server bottlenecks for instant collaboration.",
      icon: <IconBolt />,
    },
    {
      title: "Team Collaboration",
      description:
        "Real-time collaboration with multiple developers on the same project.",
      icon: <IconUsers />,
    },
    {
      title: "Access Control",
      description: "Granular permissions and role-based access management.",
      icon: <IconLock />,
    },
    {
      title: "Global Reach",
      description:
        "Connect with developers worldwide without geographical limitations.",
      icon: <IconWorld />,
    },
    {
      title: "Developer First",
      description:
        "Built specifically for developers with familiar CLI and API interfaces.",
      icon: <IconCode />,
    },
  ];

  return (
    <section
      id="features"
      className="min-h-screen flex items-center justify-center px-6"
    >
      <div className="max-w-6xl mx-auto">
        <div className="text-center mb-16">
          <h2 className="text-4xl sm:text-5xl font-bold text-[#111827] dark:text-white mb-6">
            Powerful{" "}
            <span className="text-transparent bg-clip-text bg-gradient-to-r from-[#111827] dark:from-white to-[#EC4899]">
              Features
            </span>
          </h2>
          <p className="body-large max-w-3xl mx-auto text-[#6B7280]">
            Everything you need for secure, efficient development collaboration
            — end‑to‑end encryption, real‑time P2P sync, per‑environment access
            control, and instant onboarding.
            <br className="hidden sm:block" />
            Ship faster with isolated tunnels, private registries, and native
            pair programming built in.
          </p>
        </div>

        <div className="grid grid-cols-1 md:grid-cols-2 lg:grid-cols-3 relative z-10 py-4">
          {features.map((feature, index) => (
            <Feature key={feature.title} {...feature} index={index} />
          ))}
        </div>
      </div>
    </section>
  );
};

const Feature = ({
  title,
  description,
  icon,
  index,
}: {
  title: string;
  description: string;
  icon: React.ReactNode;
  index: number;
}) => {
  return (
    <div
      className={cn(
        "flex flex-col lg:border-r py-10 relative group/feature border-[#E5E7EB] dark:border-neutral-800",
        (index === 0 || index === 3) && "lg:border-l",
        index < 3 && "lg:border-b"
      )}
    >
      {index < 3 && (
        <div className="opacity-0 group-hover/feature:opacity-100 transition duration-200 absolute inset-0 h-full w-full bg-gradient-to-t from-[#F9FAFB] dark:from-neutral-800 to-transparent pointer-events-none" />
      )}
      {index >= 3 && (
        <div className="opacity-0 group-hover/feature:opacity-100 transition duration-200 absolute inset-0 h-full w-full bg-gradient-to-b from-[#F9FAFB] dark:from-neutral-800 to-transparent pointer-events-none" />
      )}
      <div className="mb-4 relative z-10 px-10 text-[#6B7280] dark:text-neutral-400">
        {icon}
      </div>
      <div className="text-lg font-bold mb-2 relative z-10 px-10">
        <div className="absolute left-0 inset-y-0 h-6 group-hover/feature:h-8 w-1 rounded-tr-full rounded-br-full bg-[#E5E7EB] dark:bg-neutral-700 group-hover/feature:bg-[#EC4899] transition-all duration-200 origin-center" />
        <span className="group-hover/feature:translate-x-2 transition duration-200 inline-block text-[#111827] dark:text-neutral-100">
          {title}
        </span>
      </div>
      <p className="text-sm text-[#6B7280] dark:text-neutral-300 max-w-xs relative z-10 px-10">
        {description}
      </p>
    </div>
  );
};

export default FeaturesSection;
