import TerminalDemoSection from "./TerminalDemoSection";

const SnapshotsSection = () => {
  return (
    <section id="showcase" className="py-20 px-6">
      <div className="max-w-6xl mx-auto">
        <div className="text-center mb-12">
          <h2 className="text-4xl sm:text-5xl font-bold text-[#111827] dark:text-white mb-10">
            Snapshots &{" "}
            <span className="text-transparent bg-clip-text bg-gradient-to-r from-[#111827] dark:from-white to-[#EC4899]">
              Demo
            </span>
          </h2>
          <p className="body-large text-[#6B7280] max-w-3xl mx-auto">
            Explore real command snapshots alongside an interactive terminal
            demo. Switch between DevLink commands, compare before/after outputs,
            and see how workflows come together in real time.
          </p>
        </div>

        {/* Existing Snapshots content would go here if present */}

        <div className="mt-12">
          <TerminalDemoSection />
        </div>
      </div>
    </section>
  );
};

export default SnapshotsSection;
