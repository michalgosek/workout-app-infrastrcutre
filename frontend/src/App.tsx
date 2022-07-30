import './styles/app.css';

import { Home, Profile, Trainings } from './views';
import { Route, Routes } from 'react-router-dom';

import { Loading, } from './components';
import NotFound from './views/not-found';
import ProtectedRoute from './authentication/protected-route';
import { useAuth0 } from '@auth0/auth0-react';

const App: React.FC = () => {
  const { isLoading } = useAuth0();
  if (isLoading) {
    return (
      <Loading />
    );
  }
  return (
    <Routes>
      <Route path="/" element={<Home />} />
      <Route path="/profile" element={<ProtectedRoute component={Profile} />} />
      <Route path="/trainings" element={<Trainings />} />
      <Route path="*" element={<NotFound />} />
    </Routes>
  )
};

export default App;
