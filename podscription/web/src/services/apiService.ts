import { Message, Session } from '../types';

interface ChatRequest {
  content: string;
  sessionId?: string;
}

interface ChatResponse {
  session: Session;
  message: Message;
}

interface CreateSessionRequest {
  name?: string;
}

interface SessionsResponse {
  sessions: Session[];
}

interface ApiError {
  error: string;
  code?: string;
  message: string;
}

class ApiService {
  private baseUrl: string;

  constructor() {
    // Use environment variable or default to localhost for development
    this.baseUrl = process.env.REACT_APP_API_URL || 'http://localhost:8080/api';
  }

  private async fetchWithErrorHandling<T>(
    endpoint: string,
    options?: RequestInit
  ): Promise<T> {
    const url = `${this.baseUrl}${endpoint}`;
    
    try {
      const response = await fetch(url, {
        headers: {
          'Content-Type': 'application/json',
          ...options?.headers,
        },
        ...options,
      });

      if (!response.ok) {
        const errorData: ApiError = await response.json().catch(() => ({
          error: 'NETWORK_ERROR',
          message: `HTTP ${response.status}: ${response.statusText}`,
        }));

        throw new Error(`API Error: ${errorData.message || errorData.error}`);
      }

      return await response.json();
    } catch (error) {
      if (error instanceof Error) {
        throw error;
      }
      throw new Error('Unknown API error occurred');
    }
  }

  async sendMessage(content: string, sessionId?: string): Promise<ChatResponse> {
    const request: ChatRequest = {
      content,
      ...(sessionId && { sessionId }),
    };

    return this.fetchWithErrorHandling<ChatResponse>('/chat', {
      method: 'POST',
      body: JSON.stringify(request),
    });
  }

  async createSession(name?: string): Promise<Session> {
    const request: CreateSessionRequest = {
      ...(name && { name }),
    };

    return this.fetchWithErrorHandling<Session>('/sessions', {
      method: 'POST',
      body: JSON.stringify(request),
    });
  }

  async getSession(sessionId: string): Promise<Session> {
    return this.fetchWithErrorHandling<Session>(`/sessions/${sessionId}`);
  }

  async listSessions(): Promise<Session[]> {
    const response = await this.fetchWithErrorHandling<SessionsResponse>('/sessions');
    return response.sessions;
  }

  async healthCheck(): Promise<{ status: string; service: string; version: string }> {
    const response = await fetch(`${this.baseUrl.replace('/api', '')}/health`);
    
    if (!response.ok) {
      throw new Error(`Health check failed: ${response.status}`);
    }

    return response.json();
  }
}

export const apiService = new ApiService();
export type { ChatResponse, ApiError };