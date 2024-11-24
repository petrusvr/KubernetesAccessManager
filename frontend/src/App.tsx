import { BrowserRouter as Router, Routes, Route } from 'react-router-dom';
import { MainPage } from "./pages/MainPage.tsx";
import Menu from "./components/Menu/Menu.tsx";
import EditorPage from "./pages/EditorPage.tsx";
import ResourcePage from "./pages/ResourcePage.tsx";
import CreateResourcePage from "./pages/CreateResourcePage.tsx";
import AuthCallbackPage from "./pages/AuthCallbackPage.tsx";
import { initializeAxiosInterceptors } from './config/axiosConfig.ts';
import { AuthProvider } from "./components/AuthProvider/AuthProvider.tsx";
import { LoginPage } from "./pages/LoginPage.tsx";
import PrivateRoute from "./components/PrivateRoute/PrivateRoute.tsx";

initializeAxiosInterceptors();

function App() {
    return (
        <Router>
            <AuthProvider>
                <Routes>
                    <Route path="/auth/callback" element={<AuthCallbackPage />} />
                    <Route path="/" element={<Menu />}>
                        <Route path="/login" element={<LoginPage />} />
                        <Route element={<PrivateRoute />}>
                            <Route index element={<MainPage />} />
                            <Route path=":resourceType" element={<ResourcePage />} />
                            <Route path="/editor" element={<EditorPage />} />
                            <Route path="/create" element={<CreateResourcePage />} />
                        </Route>
                    </Route>
                </Routes>
            </AuthProvider>
        </Router>
    );
}

export default App;
