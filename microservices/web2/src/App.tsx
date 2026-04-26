import React from 'react';
import { BrowserRouter as Router } from 'react-router-dom';
import { HeaderProvider } from './context/HeaderContext';
import { AuthProvider } from './context/AuthContext';
import { Toaster } from "react-hot-toast";
import { AppRouter } from './route/AppRouter';

function App() {
  return (
    <AuthProvider>
      <Toaster position="bottom-right" />
      <Router>
        <HeaderProvider>
          <AppRouter/>
        </HeaderProvider>
      </Router>
    </AuthProvider>
  );
}

export default App;