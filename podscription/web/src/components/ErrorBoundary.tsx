import React, { Component, ReactNode } from 'react';
import { MEDICAL_EMOJIS } from '../config/constants';

interface Props {
  children: ReactNode;
}

interface State {
  hasError: boolean;
  error?: Error;
}

export class ErrorBoundary extends Component<Props, State> {
  constructor(props: Props) {
    super(props);
    this.state = { hasError: false };
  }

  static getDerivedStateFromError(error: Error): State {
    return { hasError: true, error };
  }

  componentDidCatch(error: Error, errorInfo: any) {
    console.error('ErrorBoundary caught an error:', error, errorInfo);
  }

  render() {
    if (this.state.hasError) {
      return (
        <div className="min-h-screen flex items-center justify-center bg-prescription-50">
          <div className="text-center max-w-md mx-auto p-6">
            <div className="w-20 h-20 bg-red-100 rounded-full flex items-center justify-center mx-auto mb-6">
              <span className="text-3xl">{MEDICAL_EMOJIS.STETHOSCOPE}</span>
            </div>
            <h2 className="text-2xl font-bold text-prescription-800 mb-4">
              Pod Doctor Emergency!
            </h2>
            <p className="text-prescription-600 mb-6">
              Something unexpected happened. The Pod Doctor needs a moment to recover.
            </p>
            <button 
              onClick={() => this.setState({ hasError: false })}
              className="px-6 py-3 bg-medical-500 text-white rounded-lg hover:bg-medical-600 focus:outline-none focus:ring-2 focus:ring-medical-300 focus:ring-offset-2 transition-colors"
            >
              Restart Consultation
            </button>
            {process.env.NODE_ENV === 'development' && this.state.error && (
              <details className="mt-6 text-left">
                <summary className="cursor-pointer text-sm text-prescription-500">
                  Error Details
                </summary>
                <pre className="mt-2 text-xs bg-prescription-100 p-3 rounded-lg overflow-auto">
                  {this.state.error.stack}
                </pre>
              </details>
            )}
          </div>
        </div>
      );
    }

    return this.props.children;
  }
}