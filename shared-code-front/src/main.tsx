import * as React from 'react';
import * as ReactDOM from 'react-dom/client';
import App from './App';
import CustomThemeProvider from './theme/CustomThemeProvider';
import { BrowserRouter } from 'react-router-dom';
import { CssBaseline } from '@mui/material';
import { SocketProvider } from './context/SocketContext';

ReactDOM.createRoot(document.getElementById('root')!).render(
    <React.StrictMode>
        <SocketProvider>
            <CustomThemeProvider>
                <CssBaseline />
                <BrowserRouter>
                    <App />
                </BrowserRouter>
            </CustomThemeProvider>
        </SocketProvider>
    </React.StrictMode>,
);
