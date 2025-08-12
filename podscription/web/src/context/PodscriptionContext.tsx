import React, { createContext, useContext, useReducer, ReactNode, useEffect } from 'react';
import { Message, Session, PodscriptionContextType } from '../types';
import { apiService } from '../services/apiService';

interface State {
  currentSession: Session | null;
  sessions: Session[];
  isLoading: boolean;
  error: string | null;
}

type Action =
  | { type: 'SET_SESSIONS'; payload: { sessions: Session[] } }
  | { type: 'SET_CURRENT_SESSION'; payload: { session: Session } }
  | { type: 'SELECT_SESSION'; payload: { sessionId: string } }
  | { type: 'ADD_MESSAGE'; payload: { message: Message } }
  | { type: 'CREATE_SESSION'; payload: { name?: string } }
  | { type: 'SET_LOADING'; payload: { loading: boolean } }
  | { type: 'SET_ERROR'; payload: { error: string | null } };

const initialState: State = {
  currentSession: null,
  sessions: [],
  isLoading: false,
  error: null,
};

function podscriptionReducer(state: State, action: Action): State {
  switch (action.type) {
    case 'SET_SESSIONS': {
      return {
        ...state,
        sessions: action.payload.sessions,
      };
    }
    case 'SET_CURRENT_SESSION': {
      const session = action.payload.session;
      return {
        ...state,
        currentSession: session,
        sessions: state.sessions.some(s => s.id === session.id)
          ? state.sessions.map(s => s.id === session.id ? session : s)
          : [...state.sessions, session],
      };
    }
    case 'SELECT_SESSION': {
      const session = state.sessions.find(s => s.id === action.payload.sessionId);
      return {
        ...state,
        currentSession: session || null,
      };
    }
    case 'ADD_MESSAGE': {
      if (!state.currentSession) return state;
      
      const updatedSession: Session = {
        ...state.currentSession,
        messages: [...state.currentSession.messages, action.payload.message],
        updatedAt: new Date(),
      };
      
      return {
        ...state,
        currentSession: updatedSession,
        sessions: state.sessions.map(s => 
          s.id === updatedSession.id ? updatedSession : s
        ),
      };
    }
    case 'CREATE_SESSION': {
      const newSession: Session = {
        id: crypto.randomUUID(),
        name: action.payload.name || `Session ${new Date().toLocaleDateString()}`,
        messages: [],
        createdAt: new Date(),
        updatedAt: new Date(),
      };
      
      return {
        ...state,
        currentSession: newSession,
        sessions: [newSession, ...state.sessions],
      };
    }
    case 'SET_LOADING': {
      return {
        ...state,
        isLoading: action.payload.loading,
      };
    }
    case 'SET_ERROR': {
      return {
        ...state,
        error: action.payload.error,
      };
    }
    default:
      return state;
  }
}

const PodscriptionContext = createContext<PodscriptionContextType | undefined>(undefined);

export function PodscriptionProvider({ children }: { children: ReactNode }) {
  const [state, dispatch] = useReducer(podscriptionReducer, initialState);

  // Load sessions on mount
  useEffect(() => {
    const loadSessions = async () => {
      try {
        const sessions = await apiService.listSessions();
        dispatch({ type: 'SET_SESSIONS', payload: { sessions } });
      } catch (error) {
        console.warn('Failed to load sessions:', error);
        // Continue with empty sessions - not a critical error
      }
    };

    loadSessions();
  }, []);

  const createSession = async (name?: string) => {
    try {
      dispatch({ type: 'SET_LOADING', payload: { loading: true } });
      dispatch({ type: 'SET_ERROR', payload: { error: null } });
      
      const session = await apiService.createSession(name);
      dispatch({ type: 'SET_CURRENT_SESSION', payload: { session } });
    } catch (error) {
      const errorMessage = error instanceof Error ? error.message : 'Failed to create session';
      dispatch({ type: 'SET_ERROR', payload: { error: errorMessage } });
      // Fallback to local session creation if API fails
      dispatch({ type: 'CREATE_SESSION', payload: { name } });
    } finally {
      dispatch({ type: 'SET_LOADING', payload: { loading: false } });
    }
  };

  const selectSession = (sessionId: string) => {
    dispatch({ type: 'SELECT_SESSION', payload: { sessionId } });
  };

  const sendMessage = async (content: string) => {
    let sessionToUse = state.currentSession;
    
    // Create a new session if none exists
    if (!sessionToUse) {
      try {
        dispatch({ type: 'SET_LOADING', payload: { loading: true } });
        dispatch({ type: 'SET_ERROR', payload: { error: null } });
        
        const newSession = await apiService.createSession();
        dispatch({ type: 'SET_CURRENT_SESSION', payload: { session: newSession } });
        sessionToUse = newSession;
      } catch (error) {
        const errorMessage = error instanceof Error ? error.message : 'Failed to create session';
        dispatch({ type: 'SET_ERROR', payload: { error: errorMessage } });
        dispatch({ type: 'SET_LOADING', payload: { loading: false } });
        return;
      }
    }

    // Add user message
    const userMessage: Message = {
      id: crypto.randomUUID(),
      role: 'user',
      content,
      timestamp: new Date(),
    };
    
    dispatch({ type: 'ADD_MESSAGE', payload: { message: userMessage } });
    dispatch({ type: 'SET_LOADING', payload: { loading: true } });
    dispatch({ type: 'SET_ERROR', payload: { error: null } });

    try {
      const response = await apiService.sendMessage(content, sessionToUse.id);
      
      // Update the session with the response
      dispatch({ type: 'SET_CURRENT_SESSION', payload: { session: response.session } });
      
      // The API response should include the assistant message in the session
      // If not, add it separately
      if (!response.session.messages.find(m => m.id === response.message.id)) {
        dispatch({ type: 'ADD_MESSAGE', payload: { message: response.message } });
      }
    } catch (error) {
      const errorMessage = error instanceof Error ? error.message : 'Failed to send message';
      dispatch({ type: 'SET_ERROR', payload: { error: errorMessage } });
      
      // Add error message to chat
      const errorMessage_: Message = {
        id: crypto.randomUUID(),
        role: 'assistant',
        content: `I apologize, but I encountered an error: ${errorMessage}. Please try again.`,
        timestamp: new Date(),
      };
      dispatch({ type: 'ADD_MESSAGE', payload: { message: errorMessage_ } });
    } finally {
      dispatch({ type: 'SET_LOADING', payload: { loading: false } });
    }
  };

  const value: PodscriptionContextType = {
    currentSession: state.currentSession,
    sessions: state.sessions,
    isLoading: state.isLoading,
    error: state.error,
    createSession,
    selectSession,
    sendMessage,
  };

  return (
    <PodscriptionContext.Provider value={value}>
      {children}
    </PodscriptionContext.Provider>
  );
}

export function usePodscription() {
  const context = useContext(PodscriptionContext);
  if (context === undefined) {
    throw new Error('usePodscription must be used within a PodscriptionProvider');
  }
  return context;
}