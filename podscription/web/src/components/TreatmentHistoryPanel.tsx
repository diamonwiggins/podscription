import React from 'react';
import { Session } from '../types';

interface TreatmentHistoryPanelProps {
  sessions: Session[];
  currentSession: Session | null;
  onSelectSession: (sessionId: string) => void;
  onNewSession: () => void;
  onToggle: () => void;
}

const TreatmentHistoryPanel: React.FC<TreatmentHistoryPanelProps> = ({
  sessions,
  currentSession,
  onSelectSession,
  onNewSession,
  onToggle
}) => {
  const formatDate = (date: Date | string) => {
    const now = new Date();
    const dateObj = typeof date === 'string' ? new Date(date) : date;
    const diffInHours = (now.getTime() - dateObj.getTime()) / (1000 * 60 * 60);
    
    if (diffInHours < 1) {
      return 'Just now';
    } else if (diffInHours < 24) {
      return `${Math.floor(diffInHours)}h ago`;
    } else {
      return dateObj.toLocaleDateString();
    }
  };

  const getSessionIcon = (session: Session) => {
    if (session.messages.length === 0) return '📋';
    
    const lastMessage = session.messages[session.messages.length - 1];
    if (lastMessage.role === 'assistant' && lastMessage.intent) {
      switch (lastMessage.intent.category) {
        case 'pod-issues': return '🔴';
        case 'networking': return '🌐';
        case 'storage': return '💾';
        case 'performance': return '⚡';
        case 'rbac': return '🔐';
        default: return '📋';
      }
    }
    return '📋';
  };

  return (
    <div className="bg-white rounded-xl shadow-lg border border-prescription-200 overflow-hidden">
      {/* Header */}
      <div className="bg-prescription-50 px-4 py-3 border-b border-prescription-200">
        <div className="flex items-center justify-between">
          <h3 className="font-semibold text-prescription-800">Treatment History</h3>
          <button
            onClick={onToggle}
            className="lg:hidden text-prescription-600 hover:text-prescription-800"
          >
            ✕
          </button>
        </div>
      </div>

      {/* New Session Button */}
      <div className="p-4 border-b border-prescription-200">
        <button
          onClick={onNewSession}
          className="w-full flex items-center space-x-3 px-4 py-3 bg-medical-500 hover:bg-medical-600 text-white rounded-lg transition-colors"
        >
          <span>➕</span>
          <span className="font-medium">New Consultation</span>
        </button>
      </div>

      {/* Sessions List */}
      <div className="max-h-96 overflow-y-auto">
        {sessions.length === 0 ? (
          <div className="p-4 text-center text-prescription-600">
            <div className="w-12 h-12 bg-prescription-100 rounded-full flex items-center justify-center mx-auto mb-2">
              <span>📋</span>
            </div>
            <p className="text-sm">No consultations yet</p>
            <p className="text-xs text-prescription-500">Start your first session above</p>
          </div>
        ) : (
          <div className="divide-y divide-prescription-200">
            {sessions.map((session) => (
              <button
                key={session.id}
                onClick={() => onSelectSession(session.id)}
                className={`w-full text-left px-4 py-3 hover:bg-prescription-50 transition-colors ${
                  currentSession?.id === session.id ? 'bg-medical-50 border-r-2 border-medical-500' : ''
                }`}
              >
                <div className="flex items-start space-x-3">
                  <span className="text-lg mt-0.5">{getSessionIcon(session)}</span>
                  <div className="flex-1 min-w-0">
                    <p className="font-medium text-prescription-800 truncate">
                      {session.name}
                    </p>
                    <p className="text-xs text-prescription-600 mt-1">
                      {session.messages.length} messages • {formatDate(session.updatedAt)}
                    </p>
                    {session.messages.length > 0 && (
                      <p className="text-xs text-prescription-500 mt-1 truncate">
                        Last: {session.messages[session.messages.length - 1].content.substring(0, 30)}...
                      </p>
                    )}
                  </div>
                </div>
              </button>
            ))}
          </div>
        )}
      </div>

      {/* Footer */}
      <div className="px-4 py-3 border-t border-prescription-200 bg-prescription-50">
        <p className="text-xs text-prescription-600 text-center">
          💡 Medical records are stored locally
        </p>
      </div>
    </div>
  );
};

export default TreatmentHistoryPanel;