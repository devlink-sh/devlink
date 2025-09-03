const TechStackSection = () => {
  const technologies = [
    {
      name: "Go",
      description: "High‑performance, concurrent backend",
    },
    {
      name: "Cobra",
      description: "Powerful CLI framework",
    },
    {
      name: "OpenZiti",
      description: "Zero‑trust network security",
    },
    {
      name: "Docker",
      description: "Containerized environments",
    },
    {
      name: "React",
      description: "Modern web interfaces",
    },
  ];

  return (
    <section id="tech-stack" className="py-20 px-6">
      <div className="max-w-6xl mx-auto w-full">
        <div className="text-center mb-16">
          <h2 className="text-4xl sm:text-5xl font-bold text-[#111827] dark:text-white mb-4">
            Built on{" "}
            <span className="text-transparent bg-clip-text bg-gradient-to-r from-[#111827] dark:from-white to-[#EC4899]">
              Modern Tech
            </span>
          </h2>
          <p className="body-large text-[#6B7280] max-w-4xl mx-auto">
            DevLink combines proven tooling with innovative peer‑to‑peer
            networking to deliver a secure, scalable, and developer‑first
            collaboration platform. From encrypted tunnels and containerized
            environments to real‑time git sync and private registries — every
            layer is optimized for speed, security, and simplicity.
          </p>

          <div className="mt-6 flex flex-wrap justify-center gap-3">
            <span className="px-3 py-1 rounded-full border border-[#E5E7EB] text-sm text-[#6B7280] bg-white/70 dark:bg-transparent dark:border-neutral-800">
              Zero‑trust networking
            </span>
            <span className="px-3 py-1 rounded-full border border-[#E5E7EB] text-sm text-[#6B7280] bg-white/70 dark:bg-transparent dark:border-neutral-800">
              Containerized workflows
            </span>
            <span className="px-3 py-1 rounded-full border border-[#E5E7EB] text-sm text-[#6B7280] bg-white/70 dark:bg-transparent dark:border-neutral-800">
              Real‑time sync
            </span>
            <span className="px-3 py-1 rounded-full border border-[#E5E7EB] text-sm text-[#6B7280] bg-white/70 dark:bg-transparent dark:border-neutral-800">
              Private registries
            </span>
          </div>
        </div>

        <div className="grid grid-cols-2 sm:grid-cols-3 md:grid-cols-5 gap-6">
          {technologies.map((tech, index) => (
            <div
              key={index}
              className="text-center rounded-2xl border border-[#E5E7EB] dark:border-neutral-800 bg-white dark:bg-gray-900 p-6 shadow-sm hover:shadow-md transition-shadow"
            >
              <div className="w-16 h-16 rounded-full bg-[#F472B6]/20 flex items-center justify-center mx-auto mb-4 text-2xl font-bold text-[#EC4899]">
                {tech.name[0]}
              </div>
              <h3 className="text-lg font-semibold text-[#111827] dark:text-white mb-1">
                {tech.name}
              </h3>
              <p className="text-sm text-[#6B7280] max-w-[160px] mx-auto">
                {tech.description}
              </p>
            </div>
          ))}
        </div>
      </div>
    </section>
  );
};

export default TechStackSection;
