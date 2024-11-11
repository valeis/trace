import React from 'react';
import { BrowserRouter, Routes, Route } from 'react-router-dom';

import { ChakraProvider, Box } from '@chakra-ui/react';
import theme from './theme';
import { ColorModeSwitcher } from './ColorModeSwitcher';

import Landing from './Components/Landing';
import Register from './Components/Register';
import Login from './Components/Login';
import Chat from './Components/Chat/Chat';
import Protected from "./router/protected";
import UserManagement from "./Components/Admin/UserManagement";
function App() {
    return (
        <ChakraProvider theme={theme}>
            <Box textAlign="right">
                <ColorModeSwitcher justifySelf="flex-end" />
            </Box>
            <Box textAlign="center" fontSize="xl">
                <BrowserRouter>
                    <Routes>
                            <Route index element={<Login />} />
                        <Route element={<Protected />}>
                            <Route path="/chat" element={<Chat />} />
                            <Route path="/users" element={<UserManagement />} />
                        </Route>
                            <Route path="/register" element={<Register />} />
                            <Route path="/login" element={<Login />} />
                    </Routes>
                </BrowserRouter>
            </Box>
        </ChakraProvider>
    );
}

export default App;