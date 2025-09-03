import React, { useState } from "react";
import Terminal, { ColorMode, TerminalOutput } from "react-terminal-ui";

const TrafficLights = ({ label }: { label: string }) => (
  <div className="flex items-center justify-between px-4 py-2 bg-[#0f0f10] border-b border-[#1f1f22]">
    <div className="flex items-center gap-2">
      <span className="w-3 h-3 rounded-full bg-[#ff5f56]" />
      <span className="w-3 h-3 rounded-full bg-[#ffbd2e]" />
      <span className="w-3 h-3 rounded-full bg-[#27c93f]" />
      <span className="ml-3 text-[11px] font-medium tracking-wide text-[#9CA3AF]">{label}</span>
    </div>
  </div>
);

const TerminalDemoSection = () => {
  const [activeFeature, setActiveFeature] = useState(0);

  const features = [
    {
      name: "Environment Sharing",
      command: "devlink env",
      oldWay: [
        <TerminalOutput># Complex Environment Setup</TerminalOutput>,
        <TerminalOutput>$ docker-compose up -d</TerminalOutput>,
        <TerminalOutput>$ npm install -g nodemon</TerminalOutput>,
        <TerminalOutput>$ export NODE_ENV=development</TerminalOutput>,
        <TerminalOutput>$ export DATABASE_URL=postgresql://...</TerminalOutput>,
        <TerminalOutput>$ export REDIS_URL=redis://...</TerminalOutput>,
        <TerminalOutput># 10+ environment variables</TerminalOutput>,
        <TerminalOutput># Multiple configuration files</TerminalOutput>,
        <TerminalOutput># Team members struggle to sync</TerminalOutput>,
      ],
      newWay: [
        <TerminalOutput># Instant Environment Sharing</TerminalOutput>,
        <TerminalOutput>$ devlink env share</TerminalOutput>,
        <TerminalOutput>✓ Environment shared securely</TerminalOutput>,
        <TerminalOutput>$ devlink env join team-project</TerminalOutput>,
        <TerminalOutput>✓ All configs synchronized</TerminalOutput>,
        <TerminalOutput>✓ Dependencies installed</TerminalOutput>,
        <TerminalOutput>✓ Ready to develop!</TerminalOutput>,
        <TerminalOutput># Total time: 30 seconds</TerminalOutput>,
      ],
    },
    {
      name: "Database Tunneling",
      command: "devlink db",
      oldWay: [
        <TerminalOutput># Shared Staging Database</TerminalOutput>,
        <TerminalOutput>
          $ ssh -L 5432:db.staging.com:5432 user@jump
        </TerminalOutput>,
        <TerminalOutput>$ psql -h localhost -U dev -d app</TerminalOutput>,
        <TerminalOutput># Multiple developers sharing</TerminalOutput>,
        <TerminalOutput># Data conflicts and corruption</TerminalOutput>,
        <TerminalOutput># Expensive cloud resources</TerminalOutput>,
        <TerminalOutput># Security concerns</TerminalOutput>,
        <TerminalOutput># Slow connections</TerminalOutput>,
      ],
      newWay: [
        <TerminalOutput># Secure Database Tunneling</TerminalOutput>,
        <TerminalOutput>$ devlink db tunnel</TerminalOutput>,
        <TerminalOutput>✓ Secure tunnel established</TerminalOutput>,
        <TerminalOutput>$ devlink db share schema</TerminalOutput>,
        <TerminalOutput>✓ Schema synchronized</TerminalOutput>,
        <TerminalOutput>✓ Zero data conflicts</TerminalOutput>,
        <TerminalOutput>✓ End-to-end encrypted</TerminalOutput>,
        <TerminalOutput># Your own secure instance</TerminalOutput>,
      ],
    },
    {
      name: "Git Collaboration",
      command: "devlink git",
      oldWay: [
        <TerminalOutput># Complex Git Workflow</TerminalOutput>,
        <TerminalOutput>$ git checkout -b feature/new-api</TerminalOutput>,
        <TerminalOutput>$ git push origin feature/new-api</TerminalOutput>,
        <TerminalOutput>$ git checkout main</TerminalOutput>,
        <TerminalOutput>$ git pull origin main</TerminalOutput>,
        <TerminalOutput>$ git merge feature/new-api</TerminalOutput>,
        <TerminalOutput># Merge conflicts...</TerminalOutput>,
        <TerminalOutput># Slow code review cycles</TerminalOutput>,
        <TerminalOutput># Branch management overhead</TerminalOutput>,
      ],
      newWay: [
        <TerminalOutput># Real-time Git Collaboration</TerminalOutput>,
        <TerminalOutput>$ devlink git sync</TerminalOutput>,
        <TerminalOutput>✓ Branch synchronized</TerminalOutput>,
        <TerminalOutput>$ devlink git pair</TerminalOutput>,
        <TerminalOutput>✓ Pair programming session</TerminalOutput>,
        <TerminalOutput>✓ Real-time code sync</TerminalOutput>,
        <TerminalOutput>✓ Zero merge conflicts</TerminalOutput>,
        <TerminalOutput># Instant collaboration</TerminalOutput>,
      ],
    },
    {
      name: "Pair Programming",
      command: "devlink pair",
      oldWay: [
        <TerminalOutput># Traditional Pair Programming</TerminalOutput>,
        <TerminalOutput>$ screen -S pair-session</TerminalOutput>,
        <TerminalOutput>$ tmux new-session -d -s pair</TerminalOutput>,
        <TerminalOutput># Screen sharing tools</TerminalOutput>,
        <TerminalOutput># Network latency issues</TerminalOutput>,
        <TerminalOutput># Security concerns</TerminalOutput>,
        <TerminalOutput># Limited collaboration</TerminalOutput>,
        <TerminalOutput># Poor user experience</TerminalOutput>,
      ],
      newWay: [
        <TerminalOutput># Native Pair Programming</TerminalOutput>,
        <TerminalOutput>$ devlink pair start</TerminalOutput>,
        <TerminalOutput>✓ Pair session initiated</TerminalOutput>,
        <TerminalOutput>
          $ devlink pair invite alice@devlink.dev
        </TerminalOutput>,
        <TerminalOutput>✓ Alice joined the session</TerminalOutput>,
        <TerminalOutput>✓ Real-time code sync</TerminalOutput>,
        <TerminalOutput>✓ Encrypted connection</TerminalOutput>,
        <TerminalOutput># Seamless collaboration</TerminalOutput>,
      ],
    },
    {
      name: "Package Registry",
      command: "devlink registry",
      oldWay: [
        <TerminalOutput># Public Package Registry</TerminalOutput>,
        <TerminalOutput>$ npm install lodash</TerminalOutput>,
        <TerminalOutput>$ npm install express</TerminalOutput>,
        <TerminalOutput># Version conflicts</TerminalOutput>,
        <TerminalOutput># Security vulnerabilities</TerminalOutput>,
        <TerminalOutput># Public exposure</TerminalOutput>,
        <TerminalOutput># Dependency hell</TerminalOutput>,
        <TerminalOutput># Slow downloads</TerminalOutput>,
      ],
      newWay: [
        <TerminalOutput># Private Package Registry</TerminalOutput>,
        <TerminalOutput>$ devlink registry init</TerminalOutput>,
        <TerminalOutput>✓ Private registry ready</TerminalOutput>,
        <TerminalOutput>$ devlink registry publish</TerminalOutput>,
        <TerminalOutput>✓ Package published securely</TerminalOutput>,
        <TerminalOutput>✓ Zero vulnerabilities</TerminalOutput>,
        <TerminalOutput>✓ Peer-to-peer distribution</TerminalOutput>,
        <TerminalOutput># Fast, secure, private</TerminalOutput>,
      ],
    },
  ];

  const currentFeature = features[activeFeature];

  return (
    <section id="terminal-demo" className="py-12 px-6">
      <div className="max-w-6xl mx-auto">
        {/* Feature Selector */}
        <div className="flex flex-wrap justify-center items-center gap-3 mb-10">
          {features.map((feature, index) => (
            <button
              key={index}
              onClick={() => setActiveFeature(index)}
              className={`flex items-center space-x-2 px-5 py-2.5 rounded-full border transition-all duration-200 backdrop-blur ${
                activeFeature === index
                  ? "border-[#EC4899] bg-[#F472B6]/20 shadow-sm"
                  : "border-[#E5E7EB] hover:border-[#EC4899]/50 hover:bg-[#F472B6]/10"
              }`}
            >
              <span className="font-medium text-[#111827]">{feature.command}</span>
            </button>
          ))}
        </div>

        {/* Active Feature Display */}
        <div className="max-w-5xl mx-auto">
          <div className="grid md:grid-cols-2 gap-6">
            {/* Before Panel */}
            <div className="rounded-2xl border border-[#E5E7EB] dark:border-neutral-800 shadow-lg overflow-hidden bg-[#0b0b0c]">
              <TrafficLights label="Before" />
              <div className="bg-[#0b0b0c]">
                <Terminal
                  name="Before"
                  colorMode={ColorMode.Dark}
                  height="360px"
                  TopButtonsPanel={() => null}
                >
                  {currentFeature.oldWay}
                </Terminal>
              </div>
            </div>

            {/* DevLink Panel */}
            <div className="rounded-2xl border border-[#EC4899] shadow-lg overflow-hidden bg-[#0b0b0c]">
              <TrafficLights label="DevLink" />
              <div className="bg-[#0b0b0c]">
                <Terminal
                  name="DevLink"
                  colorMode={ColorMode.Dark}
                  height="360px"
                  TopButtonsPanel={() => null}
                >
                  {currentFeature.newWay}
                </Terminal>
              </div>
            </div>
          </div>
        </div>
      </div>
    </section>
  );
};

export default TerminalDemoSection;
