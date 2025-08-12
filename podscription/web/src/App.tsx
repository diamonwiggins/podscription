import React from 'react';
import { ErrorBoundary } from './components/ErrorBoundary';
import { PodscriptionProvider } from './context/PodscriptionContext';
import ConsultationInterface from './components/ConsultationInterface';
import Header from './components/Header';

function App() {
  return (
    <ErrorBoundary>
      <PodscriptionProvider>
        <div className="min-h-screen bg-gradient-to-br from-medical-50 to-prescription-100">
          <Header />
          <main className="container mx-auto px-4 py-8">
            <ConsultationInterface />
          </main>
        </div>
      </PodscriptionProvider>
    </ErrorBoundary>
  );
}

export default App;