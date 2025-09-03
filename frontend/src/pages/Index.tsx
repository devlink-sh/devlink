import Navbar from "@/components/Navbar";
import HeroSection from "@/components/HeroSection";
import FeaturesSection from "@/components/FeaturesSection";
import SnapshotsSection from "@/components/SnapshotsSection";
import TechStackSection from "@/components/TechStackSection";
import ImpactSection from "@/components/ImpactSection";
import Footer from "@/components/Footer";

const Index = () => {
  return (
    <div className="min-h-screen bg-background flex flex-col">
      <Navbar />
      <main className="flex-1">
        <HeroSection />
        <FeaturesSection />
        <SnapshotsSection />
        <TechStackSection />
        <ImpactSection />
      </main>
      <Footer />
    </div>
  );
};

export default Index;
