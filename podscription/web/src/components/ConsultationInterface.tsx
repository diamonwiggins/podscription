import React, { useState, useRef, useEffect } from 'react';
import { usePodscription } from '../context/PodscriptionContext';
import { MEDICAL_EMOJIS, CONFIG } from '../config/constants';
import { Button } from './ui/Button';
import PrescriptionCard from './PrescriptionCard';
import TreatmentHistoryPanel from './TreatmentHistoryPanel';

const ConsultationInterface: React.FC = () => {
  const { 
    currentSession, 
    sessions,
    isLoading,
    error,
    createSession, 
    selectSession, 
    sendMessage 
  } = usePodscription();
  
  const [inputMessage, setInputMessage] = useState('');
  const [showHistory, setShowHistory] = useState(false);
  const messagesEndRef = useRef<HTMLDivElement>(null);

  const scrollToBottom = () => {
    messagesEndRef.current?.scrollIntoView({ behavior: 'smooth' });
  };

  useEffect(() => {
    scrollToBottom();
  }, [currentSession?.messages]);

  const handleSubmit = (e: React.FormEvent) => {
    e.preventDefault();
    if (inputMessage.trim() && !isLoading) {
      if (!currentSession) {
        createSession('New Consultation');
      }
      sendMessage(inputMessage.trim());
      setInputMessage('');
    }
  };

  const handleNewConsultation = () => {
    createSession();
    setShowHistory(false);
  };

  return (
    <div className="max-w-6xl mx-auto">
      <div className="grid grid-cols-1 lg:grid-cols-4 gap-6">
        {/* Treatment History Sidebar */}
        <div className={`lg:col-span-1 ${showHistory ? 'block' : 'hidden lg:block'}`}>
          <TreatmentHistoryPanel 
            sessions={sessions}
            currentSession={currentSession}
            onSelectSession={selectSession}
            onNewSession={handleNewConsultation}
            onToggle={() => setShowHistory(!showHistory)}
          />
        </div>

        {/* Main Consultation Area */}
        <div className="lg:col-span-3">
          <div className="bg-white rounded-xl shadow-lg border border-prescription-200 overflow-hidden">
            {/* Chat Header */}
            <div className="bg-medical-50 px-6 py-4 border-b border-prescription-200">
              <div className="flex items-center justify-between">
                <div className="flex items-center space-x-3">
                  <div className="w-8 h-8 bg-medical-500 rounded-full flex items-center justify-center">
                    <span className="text-white text-sm">{MEDICAL_EMOJIS.DOCTOR}</span>
                  </div>
                  <div>
                    <h2 className="font-semibold text-prescription-800">
                      {currentSession?.name || 'New Consultation'}
                    </h2>
                    <p className="text-xs text-prescription-600">
                      Pod Doctor • Online
                    </p>
                  </div>
                </div>
                
                <Button
                  variant="ghost"
                  size="sm"
                  onClick={() => setShowHistory(!showHistory)}
                  className="lg:hidden"
                  aria-label="Toggle treatment history"
                >
                  {MEDICAL_EMOJIS.CLIPBOARD}
                </Button>
              </div>
            </div>

            {/* Messages Area */}
            <div style={{ height: `${CONFIG.CHAT_HEIGHT}px` }} className="overflow-y-auto p-6 space-y-4">
              {!currentSession?.messages.length ? (
                <div className="text-center py-16">
                  <div className="w-20 h-20 bg-medical-100 rounded-full flex items-center justify-center mx-auto mb-6">
                    <span className="text-3xl">{MEDICAL_EMOJIS.STETHOSCOPE}</span>
                  </div>
                  <h3 className="text-2xl font-bold text-prescription-800 mb-3">
                    Welcome to Podscription
                  </h3>
                  <p className="text-prescription-600 text-lg mb-2">
                    I'm your Kubernetes Pod Doctor
                  </p>
                  <p className="text-prescription-500 max-w-md mx-auto">
                    Describe your pod symptoms and I'll provide a diagnosis with treatment recommendations.
                  </p>
                </div>
              ) : (
                currentSession.messages.map((message) => (
                  <div key={message.id}>
                    {message.role === 'user' ? (
                      <div className="flex justify-end">
                        <div className="max-w-xs lg:max-w-md bg-medical-500 text-white rounded-lg px-4 py-2">
                          <p className="text-sm">{message.content}</p>
                          <span className="text-xs opacity-75">
                            {new Date(message.timestamp).toLocaleTimeString()}
                          </span>
                        </div>
                      </div>
                    ) : (
                      <PrescriptionCard message={message} />
                    )}
                  </div>
                ))
              )}
              
              {isLoading && (
                <div className="flex justify-start">
                  <div className="max-w-xs lg:max-w-md bg-prescription-100 rounded-lg px-4 py-3">
                    <div className="flex items-center space-x-2">
                      <div className="flex space-x-1">
                        <div className="w-2 h-2 bg-medical-500 rounded-full animate-bounce"></div>
                        <div className="w-2 h-2 bg-medical-500 rounded-full animate-bounce" style={{animationDelay: CONFIG.ANIMATION_DELAYS.DOT_2}}></div>
                        <div className="w-2 h-2 bg-medical-500 rounded-full animate-bounce" style={{animationDelay: CONFIG.ANIMATION_DELAYS.DOT_3}}></div>
                      </div>
                      <span className="text-sm text-prescription-600">
                        Pod Doctor is diagnosing...
                      </span>
                    </div>
                  </div>
                </div>
              )}
              <div ref={messagesEndRef} />
            </div>

            {/* Error Alert */}
            {error && (
              <div className="border-t border-red-200 bg-red-50 px-6 py-3">
                <div className="flex items-center space-x-2">
                  <span className="text-red-500">⚠️</span>
                  <p className="text-sm text-red-700">
                    <strong>Connection Error:</strong> {error}
                  </p>
                </div>
              </div>
            )}

            {/* Input Area */}
            <div className="border-t border-prescription-200 p-6">
              <form onSubmit={handleSubmit} className="flex space-x-4">
                <input
                  type="text"
                  value={inputMessage}
                  onChange={(e) => setInputMessage(e.target.value)}
                  placeholder="Describe your pod symptoms..."
                  className="flex-1 px-4 py-2 border border-prescription-300 rounded-lg focus:outline-none focus:ring-2 focus:ring-medical-500 focus:border-transparent"
                  disabled={isLoading}
                />
                <Button
                  type="submit"
                  disabled={isLoading || !inputMessage.trim()}
                  isLoading={isLoading}
                >
                  {isLoading ? MEDICAL_EMOJIS.PILL : '📝'}
                </Button>
              </form>
            </div>
          </div>
        </div>
      </div>
    </div>
  );
};

export default ConsultationInterface;