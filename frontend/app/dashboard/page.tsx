'use client';

import { useState, useEffect, FormEvent } from 'react';
import api from '@/services/api';
import toast from 'react-hot-toast';
import {
  Loader2,
  Plus,
  Trash,
  Eye,
  EyeOff,
  Copy,
  ChevronRight,
  Lock,
} from 'lucide-react';

// --- Type Definitions ---
interface Secret {
  id: number;
  key: string;
  value: string; 
  project_id: number;
}

interface Project {
  ID: number;
  name: string;
  owner_id: number;
  secrets: Secret[]; 
}

// --- Main Dashboard Component ---
export default function DashboardPage() {
  // --- State ---
  const [projects, setProjects] = useState<Project[]>([]);
  const [secrets, setSecrets] = useState<Secret[]>([]);
  const [selectedProject, setSelectedProject] = useState<Project | null>(null);
  const [isLoadingProjects, setIsLoadingProjects] = useState(true);
  const [isLoadingSecrets, setIsLoadingSecrets] = useState(false);
  const [visibleSecrets, setVisibleSecrets] = useState<Record<number, boolean>>({});

  // --- Forms State ---
  const [newProjectName, setNewProjectName] = useState('');
  const [isProjectModalOpen, setIsProjectModalOpen] = useState(false);
  
  const [newSecretKey, setNewSecretKey] = useState('');
  const [newSecretValue, setNewSecretValue] = useState('');
  const [isSecretModalOpen, setIsSecretModalOpen] = useState(false);

  // --- Data Fetching ---
  const fetchProjects = async () => {
    setIsLoadingProjects(true);
    try {
      const response = await api.get('/api/projects');
      setProjects(response.data || []);
    } catch (error) {
      toast.error('Failed to load projects');
    } finally {
      setIsLoadingProjects(false);
    }
  };

  const fetchSecrets = async (projectId: number) => {
    setIsLoadingSecrets(true);
    setSecrets([]); // Clear old secrets
    try {
      const response = await api.get(`/api/projects/${projectId}/secrets`);
      setSecrets(response.data || []);
    } catch (error) {
      toast.error('Failed to load secrets');
    } finally {
      setIsLoadingSecrets(false);
    }
  };

  useEffect(() => {
    fetchProjects();
  }, []);

  // --- Event Handlers ---
  const handleProjectSelect = (project: Project) => {
    setSelectedProject(project);
    setVisibleSecrets({}); // Reset visibility
    fetchSecrets(project.ID);
  };

  const toggleSecretVisibility = (secretId: number) => {
    setVisibleSecrets((prev) => ({
      ...prev,
      [secretId]: !prev[secretId],
    }));
  };

  const copyToClipboard = (value: string) => {
    navigator.clipboard.writeText(value);
    toast.success('Copied to clipboard!');
  };

  // --- Form Submissions ---
  const handleCreateProject = async (e: FormEvent) => {
    e.preventDefault();
    try {
      const response = await api.post('/api/projects', { name: newProjectName });
      setProjects([response.data, ...projects]);
      toast.success('Project created!');
      setNewProjectName('');
      setIsProjectModalOpen(false);
    } catch (error) {
      toast.error('Failed to create project');
    }
  };

  const handleCreateSecret = async (e: FormEvent) => {
    e.preventDefault();
    if (!selectedProject) return;

    try {
      await api.post('/api/secrets', {
        project_id: selectedProject.ID,
        key: newSecretKey,
        value: newSecretValue,
      });
      toast.success('Secret created!');
      fetchSecrets(selectedProject.ID); // Refresh list
      setNewSecretKey('');
      setNewSecretValue('');
      setIsSecretModalOpen(false);
    } catch (error) {
      toast.error('Failed to create secret');
    }
  };

  const handleDeleteSecret = async (secretId: number) => {
    if (!window.confirm('Are you sure you want to delete this secret?')) {
      return;
    }
    try {
      await api.delete(`/api/secrets/${secretId}`);
      toast.success('Secret deleted');
      setSecrets(secrets.filter((s) => s.id !== secretId)); // Optimistic UI update
    } catch (error) {
      toast.error('Failed to delete secret');
    }
  };

  // --- Render ---
  return (
    <div className="grid grid-cols-1 md:grid-cols-3 gap-6 h-full">
      
      {/* --- Column 1: Projects List --- */}
      <div className="md:col-span-1 bg-gray-900 rounded-lg shadow-lg p-4">
        <div className="flex justify-between items-center mb-4">
          <h2 className="text-xl font-semibold">Projects</h2>
          <button
            onClick={() => setIsProjectModalOpen(true)}
            className="p-2 bg-blue-600 rounded-md hover:bg-blue-700"
          >
            <Plus className="h-4 w-4" />
          </button>
        </div>
        {isLoadingProjects ? (
          <div className="flex justify-center p-4">
            <Loader2 className="animate-spin" />
          </div>
        ) : (
          <ul className="space-y-2">
            {projects.map((project) => (
              <li key={project.ID}>
                <button
                  onClick={() => handleProjectSelect(project)}
                  className={`w-full text-left p-3 rounded-md flex justify-between items-center ${
                    selectedProject?.ID === project.ID
                      ? 'bg-blue-700 font-bold'
                      : 'bg-gray-800 hover:bg-gray-700'
                  }`}
                >
                  {project.name}
                  <ChevronRight className="h-4 w-4" />
                </button>
              </li>
            ))}
          </ul>
        )}
      </div>

      {/* --- Column 2: Secrets List --- */}
      <div className="md:col-span-2 bg-gray-900 rounded-lg shadow-lg p-4">
        {selectedProject ? (
          <>
            <div className="flex justify-between items-center mb-4">
              <h2 className="text-xl font-semibold">{selectedProject.name}</h2>
              <button
                onClick={() => setIsSecretModalOpen(true)}
                className="p-2 bg-blue-600 rounded-md hover:bg-blue-700 flex items-center gap-2 text-sm px-3"
              >
                <Plus className="h-4 w-4" /> Add Secret
              </button>
            </div>
            {isLoadingSecrets ? (
              <div className="flex justify-center p-4">
                <Loader2 className="animate-spin" />
              </div>
            ) : (
              <table className="w-full text-left">
                <thead>
                  <tr className="border-b border-gray-700">
                    <th className="p-2">Key</th>
                    <th className="p-2">Value</th>
                    <th className="p-2 text-right">Actions</th>
                  </tr>
                </thead>
                <tbody>
                  {secrets.map((secret) => (
                    <tr key={secret.id} className="border-b border-gray-800">
                      <td className="p-3 font-mono">{secret.key}</td>
                      <td className="p-3 font-mono">
                        {visibleSecrets[secret.id] ? secret.value : '••••••••••••'}
                      </td>
                      <td className="p-3 flex justify-end gap-2">
                        <button
                          onClick={() => toggleSecretVisibility(secret.id)}
                          className="p-2 hover:bg-gray-700 rounded-md"
                        >
                          {visibleSecrets[secret.id] ? (
                            <EyeOff className="h-4 w-4" />
                          ) : (
                            <Eye className="h-4 w-4" />
                          )}
                        </button>
                        <button
                          onClick={() => copyToClipboard(secret.value)}
                          className="p-2 hover:bg-gray-700 rounded-md"
                        >
                          <Copy className="h-4 w-4" />
                        </button>
                        <button
                          onClick={() => handleDeleteSecret(secret.id)}
                          className="p-2 text-red-500 hover:bg-gray-700 rounded-md"
                        >
                          <Trash className="h-4 w-4" />
                        </button>
                      </td>
                    </tr>
                  ))}
                </tbody>
              </table>
            )}
          </>
        ) : (
          <div className="flex flex-col items-center justify-center h-full text-gray-500">
            <Lock className="h-16 w-16 mb-4" />
            <p>Select a project to view its secrets.</p>
          </div>
        )}
      </div>

      {/* --- Modals --- */}
      
      {/* Project Modal */}
      {isProjectModalOpen && (
        <div className="fixed inset-0 bg-black/50 flex items-center justify-center">
          <form
            onSubmit={handleCreateProject}
            className="bg-gray-900 p-6 rounded-lg shadow-xl"
          >
            <h3 className="text-lg font-bold mb-4">Create New Project</h3>
            <input
              type="text"
              value={newProjectName}
              onChange={(e) => setNewProjectName(e.target.value)}
              placeholder="Project Name"
              className="w-full p-3 bg-gray-800 border border-gray-700 rounded-md mb-4"
              required
            />
            <div className="flex justify-end gap-2">
              <button
                type="button"
                onClick={() => setIsProjectModalOpen(false)}
                className="px-4 py-2 bg-gray-700 rounded-md"
              >
                Cancel
              </button>
              <button type="submit" className="px-4 py-2 bg-blue-600 rounded-md">
                Create
              </button>
            </div>
          </form>
        </div>
      )}

      {/* Secret Modal */}
      {isSecretModalOpen && (
        <div className="fixed inset-0 bg-black/50 flex items-center justify-center">
          <form
            onSubmit={handleCreateSecret}
            className="bg-gray-900 p-6 rounded-lg shadow-xl w-full max-w-md"
          >
            <h3 className="text-lg font-bold mb-4">Add Secret</h3>
            <div className="mb-4">
              <label className="block text-sm mb-1">Key</label>
              <input
                type="text"
                value={newSecretKey}
                onChange={(e) => setNewSecretKey(e.target.value)}
                placeholder="e.g., STRIPE_API_KEY"
                className="w-full p-3 bg-gray-800 border border-gray-700 rounded-md font-mono"
                required
              />
            </div>
            <div className="mb-4">
              <label className="block text-sm mb-1">Value</label>
              <input
                type="password"
                value={newSecretValue}
                onChange={(e) => setNewSecretValue(e.target.value)}
                placeholder="e.g., sk_live_..."
                className="w-full p-3 bg-gray-800 border border-gray-700 rounded-md font-mono"
                required
              />
            </div>
            <div className="flex justify-end gap-2">
              <button
                type="button"
                onClick={() => setIsSecretModalOpen(false)}
                className="px-4 py-2 bg-gray-700 rounded-md"
              >
                Cancel
              </button>
              <button type="submit" className="px-4 py-2 bg-blue-600 rounded-md">
                Save Secret
              </button>
            </div>
          </form>
        </div>
      )}
    </div>
  );
}