import { BrowserRouter as Router } from 'react-router-dom';
import { HeaderProvider } from './context/HeaderContext';
import { AuthProvider } from './context/AuthContext';
import { Toaster } from "react-hot-toast";
import { AppRouter } from './route/AppRouter';
import { ThemeProvider } from './context/ThemeContext';

function App() {
  return (
    <AuthProvider>
      <Toaster position="bottom-right" />
      <Router>
        <ThemeProvider>
          <HeaderProvider>
            <AppRouter/>
          </HeaderProvider>
        </ThemeProvider>
      </Router>
    </AuthProvider>
  );
}

export default App;