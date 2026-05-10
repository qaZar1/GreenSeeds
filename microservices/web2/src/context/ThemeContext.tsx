import React, {
  createContext,
  useContext,
  useEffect,
  useState,
} from 'react';

type Theme = 'light' | 'dark';

interface ThemeContextType {
  theme: Theme;
  toggleTheme: () => void;
}

const ThemeContext = createContext<ThemeContextType | null>(null);

const THEME_KEY = 'greenSeeds_theme';

export const ThemeProvider: React.FC<{
  children: React.ReactNode;
}> = ({ children }) => {
  const [theme, setTheme] = useState<Theme>(() => {
    const saved = localStorage.getItem(THEME_KEY);

    return saved === 'dark'
      ? 'dark'
      : 'light';
  });

  useEffect(() => {
    localStorage.setItem(THEME_KEY, theme);

    document.documentElement.classList.toggle(
      'dark',
      theme === 'dark'
    );
  }, [theme]);

  const toggleTheme = () => {
    setTheme(prev =>
      prev === 'light'
        ? 'dark'
        : 'light'
    );
  };

  return (
    <ThemeContext.Provider
      value={{
        theme,
        toggleTheme,
      }}
    >
      {children}
    </ThemeContext.Provider>
  );
};

export const useTheme = () => {
  const context = useContext(ThemeContext);

  if (!context) {
    throw new Error(
      'useTheme must be used within ThemeProvider'
    );
  }

  return context;
};