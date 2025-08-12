export interface Message {
  id: string;
  role: 'user' | 'assistant';
  content: string;
  timestamp: Date;
  intent?: PodIntent;
  prescription?: Prescription;
}

export interface Session {
  id: string;
  name: string;
  messages: Message[];
  createdAt: Date;
  updatedAt: Date;
}

export interface Prescription {
  diagnosis: string;
  treatment: string;
  commands?: string[];
  followUp?: string;
}

export interface PodIntent {
  category: 'networking' | 'storage' | 'pod-issues' | 'rbac' | 'performance' | 'general';
  confidence: number;
  symptoms: string[];
}

export interface PodscriptionContextType {
  currentSession: Session | null;
  sessions: Session[];
  isLoading: boolean;
  error: string | null;
  createSession: (name?: string) => Promise<void>;
  selectSession: (sessionId: string) => void;
  sendMessage: (content: string) => Promise<void>;
}