import React, { createContext, useContext, useState, useCallback, useEffect } from 'react';
import type { ReactNode } from 'react';

// 👇 1. Интерфейс для данных хедера
export interface HeaderConfig {
  title: string;
  subtitle: string;
}

// 👇 2. Интерфейс для значения контекста (что доступно потребителям)
interface HeaderContextValue {
  headerConfig: HeaderConfig;
  setHeader: (title: string, subtitle?: string) => void;
  clearHeader: () => void;
}

// 👇 3. Интерфейс для пропсов провайдера (чтобы исправить ошибку children: any)
interface HeaderContextProps {
  children: ReactNode;
}

const HeaderContext = createContext<HeaderContextValue | undefined>(undefined);

export const useHeader = (): HeaderContextValue => {
  const context = useContext(HeaderContext);
  if (!context) {
    throw new Error('useHeader must be used within HeaderProvider');
  }
  return context;
};

// 👇 4. Типизированный хук для страниц
// Если заголовок не передан, он просто не обновится (или можно добавить дефолтное поведение)
export const usePageHeader = (title?: string, subtitle: string = ''): void => {
  const { setHeader, clearHeader } = useHeader();

  useEffect(() => {
    if (title) {
      setHeader(title, subtitle);
    }
    return () => clearHeader();
  }, [title, subtitle, setHeader, clearHeader]);
};

export const HeaderProvider: React.FC<HeaderContextProps> = ({ children }) => {
  const [headerConfig, setHeaderConfig] = useState<HeaderConfig>({
    title: '',
    subtitle: '',
  });

  const setHeader = useCallback((title: string, subtitle: string = '') => {
    setHeaderConfig({ title, subtitle });
  }, []);

  const clearHeader = useCallback(() => {
    setHeaderConfig({ title: '', subtitle: '' });
  }, []);

  return (
    <HeaderContext.Provider value={{ headerConfig, setHeader, clearHeader }}>
      {children}
    </HeaderContext.Provider>
  );
};