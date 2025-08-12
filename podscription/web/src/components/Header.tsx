import React from 'react';
import { MEDICAL_EMOJIS } from '../config/constants';
import { Button } from './ui/Button';

const Header: React.FC = () => {
  return (
    <header className="bg-white shadow-sm border-b-2 border-medical-200">
      <div className="container mx-auto px-4 py-4">
        <div className="flex items-center justify-between">
          <div className="flex items-center space-x-3">
            <div className="w-10 h-10 bg-medical-500 rounded-lg flex items-center justify-center">
              <span className="text-white text-xl font-bold">{MEDICAL_EMOJIS.STETHOSCOPE}</span>
            </div>
            <div>
              <h1 className="text-2xl font-bold text-prescription-800">
                Podscription
              </h1>
              <p className="text-sm text-prescription-600">
                Your Kubernetes Pod Doctor
              </p>
            </div>
          </div>
          
          <div className="flex items-center space-x-4">
            <div className="hidden md:flex items-center space-x-2 text-sm text-prescription-600">
              <div className="w-2 h-2 bg-green-400 rounded-full"></div>
              <span>Pod Doctor Online</span>
            </div>
            
            <Button variant="secondary" size="sm" className="flex items-center space-x-2">
              <span>{MEDICAL_EMOJIS.SETTINGS}</span>
              <span className="hidden sm:inline">Settings</span>
            </Button>
          </div>
        </div>
      </div>
    </header>
  );
};

export default Header;