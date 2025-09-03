import { TrendingUp, Users, Clock, Shield, Zap, Globe } from "lucide-react";

const ImpactSection = () => {
  const impacts = [
    {
      icon: TrendingUp,
      metric: "10x",
      title: "Faster Development",
      description: "Eliminate server bottlenecks and configuration delays",
    },
    {
      icon: Users,
      metric: "100%",
      title: "Team Productivity",
      description: "Real-time collaboration without geographical limitations",
    },
    {
      icon: Clock,
      metric: "90%",
      title: "Time Saved",
      description: "Reduce setup and deployment time significantly",
    },
    {
      icon: Shield,
      metric: "Zero",
      title: "Security Breaches",
      description: "End-to-end encryption with zero-trust architecture",
    },
    {
      icon: Zap,
      metric: "Instant",
      title: "Connections",
      description: "Peer-to-peer networking for immediate collaboration",
    },
    {
      icon: Globe,
      metric: "Global",
      title: "Access",
      description: "Connect with developers worldwide seamlessly",
    },
  ];

  return (
    <section id="impact" className="py-20 px-6">
      <div className="max-w-6xl mx-auto">
        <div className="text-center mb-16">
          <h2 className="text-4xl sm:text-5xl font-bold text-[#111827] dark:text-white mb-4">
            DevLink {" "}
            <span className="text-transparent bg-clip-text bg-gradient-to-r from-[#111827] dark:from-white to-[#EC4899]">
              Impact
            </span>
          </h2>
          <p className="body-large text-[#6B7280] max-w-3xl mx-auto">
            Transform your development workflow with measurable improvements — faster onboarding,
            fewer incidents, and real‑time collaboration that scales with your team and stack.
          </p>
        </div>

        <div className="grid grid-cols-1 md:grid-cols-2 lg:grid-cols-3 gap-8">
          {impacts.map((impact, index) => (
            <div
              key={index}
              className="text-center rounded-2xl border border-[#E5E7EB] dark:border-neutral-800 bg-white dark:bg-gray-900 p-8 shadow-sm hover:shadow-md transition-shadow"
            >
              <div className="inline-flex items-center justify-center w-16 h-16 bg-[#F472B6]/20 rounded-full mb-6">
                <impact.icon size={32} className="text-[#EC4899]" />
              </div>
              <div className="text-4xl font-bold text-[#EC4899] mb-2">
                {impact.metric}
              </div>
              <h3 className="text-xl font-semibold text-[#111827] dark:text-white mb-3">
                {impact.title}
              </h3>
              <p className="text-[#6B7280] leading-relaxed">
                {impact.description}
              </p>
            </div>
          ))}
        </div>
      </div>
    </section>
  );
};

export default ImpactSection;
