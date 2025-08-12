import React, { useState, useCallback } from 'react';
import ReactMarkdown from 'react-markdown';
import remarkGfm from 'remark-gfm';
import { Message } from '../types';
import { MEDICAL_EMOJIS } from '../config/constants';
import { Button } from './ui/Button';

interface PrescriptionCardProps {
  message: Message;
}

const PrescriptionCard: React.FC<PrescriptionCardProps> = ({ message }) => {
  const [showCommands, setShowCommands] = useState(false);

  const copyCommand = useCallback(async (command: string) => {
    try {
      await navigator.clipboard.writeText(command);
      // Could add toast notification here
    } catch (error) {
      console.warn('Failed to copy to clipboard:', error);
      // Fallback for older browsers
      const textArea = document.createElement('textarea');
      textArea.value = command;
      document.body.appendChild(textArea);
      textArea.select();
      document.execCommand('copy');
      document.body.removeChild(textArea);
    }
  }, []);


  return (
    <div className="flex justify-start">
      <div className="max-w-2xl bg-gradient-to-r from-white to-prescription-50 border-l-4 border-medical-500 rounded-lg p-6 shadow-sm">
        {/* Prescription Header */}
        <div className="flex items-center space-x-3 mb-4">
          <div className="w-8 h-8 bg-medical-500 rounded-full flex items-center justify-center">
            <span className="text-white text-sm">{MEDICAL_EMOJIS.DOCTOR}</span>
          </div>
          <div className="flex-1">
            <div className="flex items-center space-x-2">
              <h4 className="font-semibold text-prescription-800">Pod Doctor</h4>
              {message.intent && (
                <span className={`px-2 py-1 text-xs rounded-full ${
                  message.intent.category === 'pod-issues' ? 'bg-red-100 text-red-700' :
                  message.intent.category === 'networking' ? 'bg-blue-100 text-blue-700' :
                  message.intent.category === 'storage' ? 'bg-yellow-100 text-yellow-700' :
                  'bg-gray-100 text-gray-700'
                }`}>
                  {message.intent.category.replace('-', ' ')}
                </span>
              )}
            </div>
            <p className="text-xs text-prescription-600">
              {new Date(message.timestamp).toLocaleTimeString()} • Confidence: {message.intent ? Math.round(message.intent.confidence * 100) : 70}%
            </p>
          </div>
        </div>

        {/* Prescription Content */}
        <div className="prose prose-sm max-w-none">
          <ReactMarkdown 
            remarkPlugins={[remarkGfm]}
            components={{
              h2: ({children}) => <h3 className="text-lg font-semibold text-prescription-800 mb-2">{children}</h3>,
              h3: ({children}) => <h4 className="text-md font-medium text-prescription-700 mb-2">{children}</h4>,
              strong: ({children}) => <strong className="font-semibold">{children}</strong>,
              code: ({children}) => <code className="bg-prescription-100 px-2 py-1 rounded text-sm font-mono text-prescription-800">{children}</code>,
            }}
          >
            {message.content}
          </ReactMarkdown>
        </div>

        {/* Commands Section */}
        {message.prescription?.commands && message.prescription.commands.length > 0 && (
          <div className="mt-4 pt-4 border-t border-prescription-200">
            <Button
              variant="ghost"
              size="sm"
              onClick={() => setShowCommands(!showCommands)}
              className="flex items-center space-x-2"
            >
              <span>{showCommands ? MEDICAL_EMOJIS.CLIPBOARD : MEDICAL_EMOJIS.PILL}</span>
              <span className="text-sm font-medium">
                {showCommands ? 'Hide' : 'Show'} Prescription Commands
              </span>
            </Button>
            
            {showCommands && (
              <div className="mt-3 space-y-2">
                {message.prescription.commands.map((command, index) => (
                  <div key={index} className="flex items-center space-x-2 bg-prescription-100 rounded-lg p-3">
                    <code className="flex-1 text-sm font-mono text-prescription-800">
                      {command}
                    </code>
                    <button
                      onClick={() => copyCommand(command)}
                      className="text-medical-600 hover:text-medical-700 transition-colors focus:outline-none focus:ring-2 focus:ring-medical-300 rounded p-1"
                      title="Copy command to clipboard"
                      aria-label="Copy command to clipboard"
                    >
                      {MEDICAL_EMOJIS.CLIPBOARD}
                    </button>
                  </div>
                ))}
              </div>
            )}
          </div>
        )}

        {/* Follow-up */}
        {message.prescription?.followUp && (
          <div className="mt-4 p-3 bg-medical-50 rounded-lg">
            <p className="text-sm text-prescription-700">
              <span className="font-medium">{MEDICAL_EMOJIS.MAGNIFYING_GLASS} Follow-up:</span> {message.prescription.followUp}
            </p>
          </div>
        )}
      </div>
    </div>
  );
};

export default PrescriptionCard;