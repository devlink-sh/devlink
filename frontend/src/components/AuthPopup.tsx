import React, { useState } from "react";
import { X, Mail, Lock, User, ArrowRight, Apple, Command } from "lucide-react";
import { useNavigate } from "react-router-dom";
import { cn } from "@/lib/utils";

interface AuthPopupProps {
  isOpen: boolean;
  onClose: () => void;
}

const AuthPopup = ({ isOpen, onClose }: AuthPopupProps) => {
  const [isSignUp, setIsSignUp] = useState(false);
  const [email, setEmail] = useState("");
  const [password, setPassword] = useState("");
  const [name, setName] = useState("");
  const navigate = useNavigate();

  const handleSubmit = (e: React.FormEvent) => {
    e.preventDefault();
    // Close popup and navigate to flow page
    onClose();
    navigate("/flow");
  };

  if (!isOpen) return null;

  return (
    <div className="fixed inset-0 z-50 flex items-center justify-center p-4">
      {/* Backdrop with blur */}
      <div
        className="absolute inset-0 bg-black/20 backdrop-blur-md"
        onClick={onClose}
      />

      {/* Main Popup */}
      <div className="relative w-full max-w-md transform transition-all duration-300 ease-out">
        {/* Glass morphism container */}
        <div className="relative bg-white/90 dark:bg-black/90 backdrop-blur-xl border border-white/20 dark:border-white/10 rounded-3xl shadow-2xl overflow-hidden">
          {/* Subtle gradient overlay */}
          <div className="absolute inset-0 bg-gradient-to-br from-white/50 to-transparent dark:from-black/50"></div>

          {/* Content */}
          <div className="relative z-10">
            {/* Header */}
            <div className="flex items-center justify-between p-8 pb-6">
              <div className="flex items-center space-x-3">
                <div className="w-10 h-10 bg-gradient-to-br from-[#EC4899] to-[#F472B6] rounded-2xl flex items-center justify-center">
                  <Command size={20} className="text-white" />
                </div>
                <div>
                  <h2 className="text-2xl font-semibold text-gray-900 dark:text-white">
                    {isSignUp ? "Create Account" : "Welcome Back"}
                  </h2>
                  <p className="text-sm text-gray-500 dark:text-gray-400">
                    {isSignUp ? "Join DevLink CLI" : "Sign in to continue"}
                  </p>
                </div>
              </div>
              <button
                onClick={onClose}
                className="p-2 hover:bg-gray-100 dark:hover:bg-gray-800 rounded-full transition-all duration-200 group"
              >
                <X
                  size={20}
                  className="text-gray-400 group-hover:text-gray-600 dark:group-hover:text-gray-300"
                />
              </button>
            </div>

            {/* Form */}
            <form onSubmit={handleSubmit} className="px-8 pb-8 space-y-6">
              {isSignUp && (
                <div className="space-y-2">
                  <label className="text-sm font-medium text-gray-700 dark:text-gray-300">
                    Full Name
                  </label>
                  <div className="relative">
                    <input
                      type="text"
                      value={name}
                      onChange={(e) => setName(e.target.value)}
                      className="w-full pl-4 pr-4 py-4 bg-gray-50 dark:bg-gray-800/50 border border-gray-200 dark:border-gray-700 rounded-2xl focus:outline-none focus:ring-2 focus:ring-[#EC4899]/50 focus:border-[#EC4899] transition-all duration-200 text-gray-900 dark:text-white placeholder-gray-400"
                      placeholder="Enter your full name"
                      required
                    />
                  </div>
                </div>
              )}

              <div className="space-y-2">
                <label className="text-sm font-medium text-gray-700 dark:text-gray-300">
                  Email Address
                </label>
                <div className="relative">
                  <input
                    type="email"
                    value={email}
                    onChange={(e) => setEmail(e.target.value)}
                    className="w-full pl-4 pr-4 py-4 bg-gray-50 dark:bg-gray-800/50 border border-gray-200 dark:border-gray-700 rounded-2xl focus:outline-none focus:ring-2 focus:ring-[#EC4899]/50 focus:border-[#EC4899] transition-all duration-200 text-gray-900 dark:text-white placeholder-gray-400"
                    placeholder="Enter your email"
                    required
                  />
                </div>
              </div>

              <div className="space-y-2">
                <label className="text-sm font-medium text-gray-700 dark:text-gray-300">
                  Password
                </label>
                <div className="relative">
                  <input
                    type="password"
                    value={password}
                    onChange={(e) => setPassword(e.target.value)}
                    className="w-full pl-4 pr-4 py-4 bg-gray-50 dark:bg-gray-800/50 border border-gray-200 dark:border-gray-700 rounded-2xl focus:outline-none focus:ring-2 focus:ring-[#EC4899]/50 focus:border-[#EC4899] transition-all duration-200 text-gray-900 dark:text-white placeholder-gray-400"
                    placeholder="Enter your password"
                    required
                  />
                </div>
              </div>

              {/* Submit Button */}
              <button
                type="submit"
                className="w-full bg-gradient-to-r from-[#EC4899] to-[#F472B6] hover:from-[#F472B6] hover:to-[#EC4899] text-white py-4 px-6 rounded-2xl font-semibold transition-all duration-300 transform hover:scale-[1.02] active:scale-[0.98] shadow-lg hover:shadow-xl"
              >
                <span className="flex items-center justify-center gap-2">
                  {isSignUp ? "Create Account" : "Sign In"}
                  <ArrowRight
                    size={18}
                    className="transition-transform group-hover:translate-x-1"
                  />
                </span>
              </button>

              {/* Divider */}
              <div className="relative">
                <div className="absolute inset-0 flex items-center">
                  <div className="w-full border-t border-gray-200 dark:border-gray-700"></div>
                </div>
                <div className="relative flex justify-center text-sm">
                  <span className="px-2 bg-white dark:bg-black text-gray-500">
                    or
                  </span>
                </div>
              </div>

              {/* Apple Sign In Button */}
              <button
                type="button"
                className="w-full bg-black dark:bg-white text-white dark:text-black py-4 px-6 rounded-2xl font-semibold transition-all duration-200 hover:bg-gray-800 dark:hover:bg-gray-100 flex items-center justify-center gap-3"
              >
                <Apple size={20} />
                Continue with Apple
              </button>
            </form>

            {/* Toggle */}
            <div className="px-8 pb-8 text-center">
              <p className="text-sm text-gray-500 dark:text-gray-400">
                {isSignUp
                  ? "Already have an account?"
                  : "Don't have an account?"}
                <button
                  onClick={() => setIsSignUp(!isSignUp)}
                  className="ml-1 text-[#EC4899] hover:text-[#F472B6] font-semibold transition-colors duration-200"
                >
                  {isSignUp ? "Sign In" : "Sign Up"}
                </button>
              </p>
            </div>
          </div>
        </div>
      </div>
    </div>
  );
};

export default AuthPopup;
