import { createContext, useContext, useState, useEffect } from 'react';

const themes = {
  dark: {
    name: 'dark',
    // Backgrounds
    bg: '#1a1a2e',
    bgSecondary: '#2a2a3e',
    bgTertiary: '#252540',
    bgHover: '#3a3a5e',
    // Canvas
    canvasBg: '#16162a',
    canvasGrid: '#333',
    // Text
    text: '#fff',
    textSecondary: '#ccc',
    textMuted: '#888',
    textDim: '#666',
    // Borders
    border: '#444',
    borderLight: '#555',
    // Inputs
    inputBg: '#1a1a2e',
    inputBorder: '#444',
    // Overlay
    overlay: 'rgba(0,0,0,0.7)',
    // Shadows
    shadow: 'rgba(0,0,0,0.3)',
    // Status colors (same for both themes)
    success: '#4caf50',
    warning: '#ff9800',
    error: '#f44336',
    info: '#2196f3',
  },
  light: {
    name: 'light',
    // Backgrounds
    bg: '#f5f5f5',
    bgSecondary: '#ffffff',
    bgTertiary: '#e8e8e8',
    bgHover: '#e0e0e0',
    // Canvas
    canvasBg: '#fafafa',
    canvasGrid: '#ddd',
    // Text
    text: '#1a1a2e',
    textSecondary: '#333',
    textMuted: '#666',
    textDim: '#999',
    // Borders
    border: '#ccc',
    borderLight: '#ddd',
    // Inputs
    inputBg: '#fff',
    inputBorder: '#ccc',
    // Overlay
    overlay: 'rgba(0,0,0,0.5)',
    // Shadows
    shadow: 'rgba(0,0,0,0.15)',
    // Status colors (same for both themes)
    success: '#4caf50',
    warning: '#ff9800',
    error: '#f44336',
    info: '#2196f3',
  },
};

const ThemeContext = createContext();

export function ThemeProvider({ children }) {
  const [themeName, setThemeName] = useState(() => {
    // Check localStorage for saved preference
    const saved = localStorage.getItem('quest-editor-theme');
    return saved || 'light';
  });

  const theme = themes[themeName];

  const toggleTheme = () => {
    setThemeName(prev => prev === 'dark' ? 'light' : 'dark');
  };

  // Save preference to localStorage
  useEffect(() => {
    localStorage.setItem('quest-editor-theme', themeName);
    // Also set a data attribute on body for CSS
    document.body.setAttribute('data-theme', themeName);
  }, [themeName]);

  return (
    <ThemeContext.Provider value={{ theme, themeName, toggleTheme }}>
      {children}
    </ThemeContext.Provider>
  );
}

export function useTheme() {
  const context = useContext(ThemeContext);
  if (!context) {
    throw new Error('useTheme must be used within a ThemeProvider');
  }
  return context;
}

export { themes };
