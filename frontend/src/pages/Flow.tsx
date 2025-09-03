import React, { useState, useCallback, useEffect } from "react";
import {
  ArrowLeft,
  Command,
  Terminal,
  Database,
  GitBranch,
  Users,
  Package,
  X,
  Clock,
  User,
  Server,
  Activity,
  Globe,
  Shield,
  Database as DatabaseIcon,
  GitBranch as GitIcon,
  Users as UsersIcon,
  Package as PackageIcon,
} from "lucide-react";
import {
  ReactFlow,
  Background,
  Controls,
  Node,
  Edge,
  NodeChange,
  EdgeChange,
  useReactFlow,
} from "@xyflow/react";
import "@xyflow/react/dist/style.css";

// Types for your API response
interface APINode {
  id: string;
  label: string;
  type?: string;
  status?: string;
  command?: string;
  users?: string[];
  packages?: number;
  metrics?: string;
  position?: { x: number; y: number };
}

interface APIEdge {
  id: string;
  source: string;
  target: string;
  label?: string;
  animated?: boolean;
}

interface APIResponse {
  nodes: APINode[];
  edges: APIEdge[];
  lastUpdated: string;
  systemStatus: string;
}

// Detailed metrics interface for each node
interface NodeMetrics {
  id: string;
  label: string;
  type: string;
  status: string;
  details: {
    [key: string]: string | number | string[];
  };
}

const Flow = () => {
  const [nodes, setNodes] = useState<Node[]>([]);
  const [edges, setEdges] = useState<Edge[]>([]);
  const [isLoading, setIsLoading] = useState(true);
  const [lastUpdated, setLastUpdated] = useState<string>("");
  const [systemStatus, setSystemStatus] = useState<string>("offline");
  const [selectedNode, setSelectedNode] = useState<NodeMetrics | null>(null);
  const [showMetricsModal, setShowMetricsModal] = useState(false);

  // Demo metrics data for each node type
  const getNodeMetrics = (
    nodeId: string,
    nodeType: string,
    nodeLabel: string,
    nodeStatus: string
  ): NodeMetrics => {
    const baseMetrics = {
      id: nodeId,
      label: nodeLabel,
      type: nodeType,
      status: nodeStatus,
    };

    switch (nodeType) {
      case "input":
        return {
          ...baseMetrics,
          details: {
            Email: "kshitijnk08@gmail.com",
            "Account Token": "••••••••••••••••••••",
            "Account Type": "Developer",
            Created: "31/08/2025, 20:53:30",
            "Last Login": "31/08/2025, 20:53:30",
            Status: "Active",
            Plan: "Pro",
            Projects: 3,
          },
        };

      case "terminal":
        return {
          ...baseMetrics,
          details: {
            Address: "49.207.63.224",
            Created: "31/08/2025, 20:53:30",
            Description: "knk@KSHITIJ",
            Host: "knk; KSHITIJ; linux; ubuntu; debian; 24.04; 6.6.87.2-microsoft-standard-WSL2; x86_64",
            Updated: "31/08/2025, 20:53:30",
            Status: "Connected",
            Port: "22",
            Protocol: "SSH",
            Environment: "Development",
            Variables: 12,
          },
        };

      case "database":
        return {
          ...baseMetrics,
          details: {
            Host: "db.devlink.internal",
            Port: "5432",
            Database: "devlink_prod",
            Username: "devlink_user",
            Status: "Connected",
            Created: "31/08/2025, 20:53:30",
            "Last Sync": "31/08/2025, 20:53:30",
            Tables: 24,
            Size: "2.4 GB",
            Connections: 8,
            SSL: "Enabled",
          },
        };

      case "git":
        return {
          ...baseMetrics,
          details: {
            Repository: "devlink-flow",
            Branch: "main",
            "Last Commit": "feat: add real-time flow chart",
            Author: "kshitijnk08@gmail.com",
            Status: "Syncing",
            Created: "31/08/2025, 20:53:30",
            Updated: "31/08/2025, 20:53:30",
            Commits: 47,
            Branches: 3,
            Collaborators: 2,
            "Pull Requests": 1,
          },
        };

      case "users":
        return {
          ...baseMetrics,
          details: {
            "Session ID": "pair_20250831_205330",
            Host: "knk@KSHITIJ",
            Status: "Available",
            Created: "31/08/2025, 20:53:30",
            Updated: "31/08/2025, 20:53:30",
            Participants: ["Alice", "Bob"],
            "Max Participants": 4,
            "Room Code": "DEV-ABC-123",
            Duration: "2h 15m",
            "Shared Files": 8,
            "Chat Messages": 23,
          },
        };

      case "package":
        return {
          ...baseMetrics,
          details: {
            Registry: "devlink.internal",
            Packages: 15,
            Status: "Ready",
            Created: "31/08/2025, 20:53:30",
            Updated: "31/08/2025, 20:53:30",
            "Total Size": "156 MB",
            Downloads: 234,
            "Latest Version": "1.2.3",
            Dependencies: 8,
            Security: "Scanned",
            License: "MIT",
          },
        };

      case "output":
        return {
          ...baseMetrics,
          details: {
            "Deployment ID": "deploy_20250831_205330",
            Environment: "Production",
            Status: "Pending",
            Created: "31/08/2025, 20:53:30",
            Updated: "31/08/2025, 20:53:30",
            Target: "AWS ECS",
            Region: "us-east-1",
            Resources: "2 vCPU, 4GB RAM",
            "Health Check": "Passing",
            Rollback: "Available",
            Monitoring: "Enabled",
          },
        };

      default:
        return {
          ...baseMetrics,
          details: {
            Status: nodeStatus,
            Created: "31/08/2025, 20:53:30",
            Updated: "31/08/2025, 20:53:30",
            Type: nodeType,
          },
        };
    }
  };

  // Handle node click
  const onNodeClick = useCallback((event: React.MouseEvent, node: Node) => {
    const nodeId = typeof node.id === "string" ? node.id : String(node.id);
    const data: any = node.data ?? {};
    const nodeType = (data.type as string) || "default";
    const nodeLabel = (data.label as string) || "Unknown";
    const nodeStatus = (data.status as string) || "unknown";

    const metrics = getNodeMetrics(nodeId, nodeType, nodeLabel, nodeStatus);
    setSelectedNode(metrics);
    setShowMetricsModal(true);
  }, []);

  // Transform API data to React Flow format
  const transformAPIDataToFlow = useCallback((apiData: APIResponse) => {
    const transformedNodes: Node[] = apiData.nodes.map((apiNode, index) => {
      // Default positions if not provided by API
      const defaultPosition = apiNode.position || {
        x: 200 + index * 250,
        y: 100 + Math.floor(index / 3) * 200,
      };

      // Get icon based on node type
      const getIcon = (nodeType: string) => {
        switch (nodeType) {
          case "terminal":
            return <Terminal className="w-6 h-6 text-[#EC4899]" />;
          case "database":
            return <Database className="w-6 h-6 text-[#EC4899]" />;
          case "git":
            return <GitBranch className="w-6 h-6 text-[#EC4899]" />;
          case "users":
            return <Users className="w-6 h-6 text-[#EC4899]" />;
          case "package":
            return <Package className="w-6 h-6 text-[#EC4899]" />;
          default:
            return <Command className="w-6 h-6 text-[#EC4899]" />;
        }
      };

      // Get status color
      const getStatusColor = (status: string) => {
        switch (status) {
          case "active":
          case "connected":
          case "online":
            return "#10B981"; // green
          case "syncing":
          case "pending":
            return "#F59E0B"; // yellow
          case "error":
          case "disconnected":
            return "#EF4444"; // red
          default:
            return "#6B7280"; // gray
        }
      };

      return {
        id: apiNode.id,
        position: defaultPosition,
        data: {
          label: apiNode.label,
          icon: getIcon(apiNode.type || "default"),
          command: apiNode.command,
          status: apiNode.status,
          users: apiNode.users,
          packages: apiNode.packages,
          metrics: apiNode.metrics,
          type: apiNode.type || "default",
        },
        type:
          apiNode.type === "input"
            ? "input"
            : apiNode.type === "output"
            ? "output"
            : "default",
        style: {
          background: "white",
          border: `2px solid ${getStatusColor(apiNode.status || "default")}`,
          borderRadius: "12px",
          padding: "16px",
          minWidth: "200px",
          boxShadow: `0 4px 15px -3px ${getStatusColor(
            apiNode.status || "default"
          )}20`,
          cursor: "pointer",
        },
      };
    });

    const transformedEdges: Edge[] = apiData.edges.map((apiEdge) => ({
      id: apiEdge.id,
      source: apiEdge.source,
      target: apiEdge.target,
      label: apiEdge.label,
      type: "step",
      animated: apiEdge.animated,
      style: {
        stroke: "#EC4899",
        strokeWidth: apiEdge.animated ? 3 : 2,
      },
    }));

    return { nodes: transformedNodes, edges: transformedEdges };
  }, []);

  // Fetch data from your API endpoints
  const fetchFlowData = useCallback(async () => {
    try {
      setIsLoading(true);

      // This is where you'll call your actual API endpoints
      // For now, using mock data that simulates your API response
      const mockAPIResponse: APIResponse = {
        nodes: [
          {
            id: "start",
            label: "DevLink CLI Project",
            type: "input",
            status: "active",
            position: { x: 400, y: 100 },
          },
          {
            id: "env",
            label: "Environment Setup",
            type: "terminal",
            status: "active",
            command: "devlink env share",
            position: { x: 150, y: 250 },
          },
          {
            id: "db",
            label: "Database Connection",
            type: "database",
            status: "connected",
            command: "devlink db tunnel",
            position: { x: 400, y: 250 },
          },
          {
            id: "git",
            label: "Git Sync",
            type: "git",
            status: "syncing",
            command: "devlink git sync",
            position: { x: 650, y: 250 },
          },
          {
            id: "pair",
            label: "Pair Programming",
            type: "users",
            status: "available",
            command: "devlink pair start",
            users: ["Alice", "Bob"],
            position: { x: 200, y: 400 },
          },
          {
            id: "registry",
            label: "Package Registry",
            type: "package",
            status: "ready",
            command: "devlink registry init",
            packages: 15,
            position: { x: 450, y: 400 },
          },
          {
            id: "monitoring",
            label: "Live Monitoring",
            type: "default",
            status: "online",
            metrics: "CPU: 45%, RAM: 2.1GB",
            position: { x: 700, y: 400 },
          },
          {
            id: "deploy",
            label: "Deploy & Share",
            type: "output",
            status: "pending",
            position: { x: 400, y: 550 },
          },
        ],
        edges: [
          {
            id: "start-env",
            source: "start",
            target: "env",
            label: "Setup Environment",
            animated: true,
          },
          {
            id: "start-db",
            source: "start",
            target: "db",
            label: "Connect Database",
            animated: true,
          },
          {
            id: "start-git",
            source: "start",
            target: "git",
            label: "Sync Code",
            animated: true,
          },
          {
            id: "env-pair",
            source: "env",
            target: "pair",
            label: "Collaborate",
          },
          {
            id: "db-registry",
            source: "db",
            target: "registry",
            label: "Manage Packages",
          },
          {
            id: "git-pair",
            source: "git",
            target: "pair",
            label: "Real-time Sync",
          },
          {
            id: "git-monitoring",
            source: "git",
            target: "monitoring",
            label: "Track Changes",
          },
          {
            id: "pair-deploy",
            source: "pair",
            target: "deploy",
            label: "Deploy Together",
          },
          {
            id: "registry-deploy",
            source: "registry",
            target: "deploy",
            label: "Package & Deploy",
          },
          {
            id: "monitoring-deploy",
            source: "monitoring",
            target: "deploy",
            label: "Health Check",
          },
        ],
        lastUpdated: new Date().toISOString(),
        systemStatus: "online",
      };

      // Simulate API call delay
      await new Promise((resolve) => setTimeout(resolve, 500));

      // Transform API data to React Flow format
      const { nodes: flowNodes, edges: flowEdges } =
        transformAPIDataToFlow(mockAPIResponse);

      setNodes(flowNodes);
      setEdges(flowEdges);
      setLastUpdated(mockAPIResponse.lastUpdated);
      setSystemStatus(mockAPIResponse.systemStatus);
    } catch (error) {
      console.error("Failed to fetch flow data:", error);
      setSystemStatus("error");
    } finally {
      setIsLoading(false);
    }
  }, [transformAPIDataToFlow]);

  // Initial data fetch
  useEffect(() => {
    fetchFlowData();
  }, [fetchFlowData]);

  // Real-time updates (every 10 seconds)
  useEffect(() => {
    const interval = setInterval(() => {
      fetchFlowData();
    }, 10000);

    return () => clearInterval(interval);
  }, [fetchFlowData]);

  const onNodesChange = useCallback((changes: NodeChange[]) => {
    setNodes((nds) =>
      nds.map((node) => {
        const change = changes.find((c) => "id" in c && c.id === node.id);
        if (change && change.type === "position" && "position" in change) {
          return { ...node, position: change.position };
        }
        return node;
      })
    );
  }, []);

  const onEdgesChange = useCallback((changes: EdgeChange[]) => {
    setEdges((eds) =>
      eds.map((edge) => {
        const change = changes.find((c) => "id" in c && c.id === edge.id);
        if (change && "id" in change) {
          return { ...edge, ...change };
        }
        return edge;
      })
    );
  }, []);

  return (
    <div className="min-h-screen bg-gradient-to-br from-[#F9FAFB] to-white dark:from-[#111827] dark:to-black">
      {/* Header */}
      <header className="relative z-10 border-b border-gray-200 dark:border-gray-800 bg-white/80 dark:bg-black/80 backdrop-blur-sm">
        <div className="max-w-7xl mx-auto px-6 py-4">
          <div className="flex items-center justify-between">
            <div className="flex items-center space-x-3">
              <div className="w-10 h-10 bg-gradient-to-br from-[#EC4899] to-[#F472B6] rounded-2xl flex items-center justify-center">
                <Command className="w-6 h-6 text-white" />
              </div>
              <div>
                <h1 className="text-xl font-semibold text-gray-900 dark:text-white">
                  DevLink CLI
                </h1>
                <p className="text-sm text-gray-500 dark:text-gray-400">
                  Real-time Development Workflow
                </p>
              </div>
            </div>
            <div className="flex items-center space-x-4">
              <div className="text-sm text-gray-600 dark:text-gray-400">
                <span
                  className={`inline-block w-2 h-2 rounded-full mr-2 ${
                    systemStatus === "online"
                      ? "bg-green-500"
                      : systemStatus === "syncing"
                      ? "bg-yellow-500"
                      : systemStatus === "error"
                      ? "bg-red-500"
                      : "bg-gray-500"
                  }`}
                ></span>
                {systemStatus === "online"
                  ? "Live Connection"
                  : systemStatus === "syncing"
                  ? "Syncing..."
                  : systemStatus === "error"
                  ? "Connection Error"
                  : "Offline"}
              </div>
              {lastUpdated && (
                <div className="text-xs text-gray-500 dark:text-gray-400">
                  Last updated: {new Date(lastUpdated).toLocaleTimeString()}
                </div>
              )}
              <a
                href="/"
                className="flex items-center space-x-2 px-4 py-2 text-gray-600 dark:text-gray-400 hover:text-[#EC4899] transition-colors"
              >
                <ArrowLeft className="w-4 h-4" />
                <span>Back to Home</span>
              </a>
            </div>
          </div>
        </div>
      </header>

      {/* Full-Screen React Flow Canvas */}
      <main className="relative z-10 h-[calc(100vh-80px)]">
        {isLoading ? (
          <div className="w-full h-full flex items-center justify-center">
            <div className="text-center">
              <div className="animate-spin rounded-full h-12 w-12 border-b-2 border-[#EC4899] mx-auto mb-4"></div>
              <p className="text-gray-600 dark:text-gray-400">
                Loading workflow...
              </p>
            </div>
          </div>
        ) : (
          <div className="w-full h-full">
            <ReactFlow
              nodes={nodes}
              edges={edges}
              onNodesChange={onNodesChange}
              onEdgesChange={onEdgesChange}
              onNodeClick={onNodeClick}
              fitView
              attributionPosition="bottom-left"
              defaultViewport={{ x: 0, y: 0, zoom: 0.8 }}
              minZoom={0.1}
              maxZoom={2}
              proOptions={{ hideAttribution: true }}
            >
              <Background
                color="#E5E7EB"
                gap={30}
                size={1}
                className="dark:bg-gray-800"
              />
              <Controls
                showZoom={true}
                showFitView={true}
                showInteractive={true}
                className="bg-white dark:bg-gray-800 border border-gray-200 dark:border-gray-700 rounded-lg shadow-lg"
              />
            </ReactFlow>
          </div>
        )}
      </main>

      {/* Metrics Modal */}
      {showMetricsModal && selectedNode && (
        <div className="fixed inset-0 z-50 flex items-center justify-center p-4">
          <div
            className="absolute inset-0 bg-black/50 backdrop-blur-sm"
            onClick={() => setShowMetricsModal(false)}
          />

          <div className="relative w-full max-w-2xl bg-white dark:bg-gray-900 border border-gray-200 dark:border-gray-700 rounded-2xl shadow-2xl overflow-hidden">
            {/* Header */}
            <div className="flex items-center justify-between p-6 border-b border-gray-200 dark:border-gray-700">
              <div className="flex items-center space-x-3">
                <div className="w-12 h-12 bg-gradient-to-br from-[#EC4899] to-[#F472B6] rounded-2xl flex items-center justify-center">
                  {selectedNode.type === "terminal" && (
                    <Terminal className="w-7 h-7 text-white" />
                  )}
                  {selectedNode.type === "database" && (
                    <Database className="w-7 h-7 text-white" />
                  )}
                  {selectedNode.type === "git" && (
                    <GitBranch className="w-7 h-7 text-white" />
                  )}
                  {selectedNode.type === "users" && (
                    <Users className="w-7 h-7 text-white" />
                  )}
                  {selectedNode.type === "package" && (
                    <Package className="w-7 h-7 text-white" />
                  )}
                  {selectedNode.type === "input" && (
                    <User className="w-7 h-7 text-white" />
                  )}
                  {selectedNode.type === "output" && (
                    <Server className="w-7 h-7 text-white" />
                  )}
                  {selectedNode.type === "default" && (
                    <Activity className="w-7 h-7 text-white" />
                  )}
                </div>
                <div>
                  <h2 className="text-2xl font-semibold text-gray-900 dark:text-white">
                    {selectedNode.label}
                  </h2>
                  <p className="text-sm text-gray-500 dark:text-gray-400 capitalize">
                    {selectedNode.type} • {selectedNode.status}
                  </p>
                </div>
              </div>
              <button
                onClick={() => setShowMetricsModal(false)}
                className="p-2 hover:bg-gray-100 dark:hover:bg-gray-800 rounded-full transition-colors"
              >
                <X size={20} className="text-gray-400" />
              </button>
            </div>

            {/* Content */}
            <div className="p-6">
              <div className="grid grid-cols-1 md:grid-cols-2 gap-6">
                {Object.entries(selectedNode.details).map(([key, value]) => (
                  <div key={key} className="space-y-2">
                    <label className="text-sm font-medium text-gray-500 dark:text-gray-400 uppercase tracking-wide">
                      {key}
                    </label>
                    <div className="text-sm text-gray-900 dark:text-white font-mono bg-gray-50 dark:bg-gray-800 px-3 py-2 rounded-lg">
                      {Array.isArray(value) ? value.join(", ") : String(value)}
                    </div>
                  </div>
                ))}
              </div>
            </div>
          </div>
        </div>
      )}
    </div>
  );
};

export default Flow;
