import './styles/app.css';

import { Hero, Loading, } from './components';
import { NotAuthorized, NotFound } from './views/errors';
import { PageLayout, Profile, Trainings } from './views';
import { PlanTraining, Trainer } from './views/trainings/trainer';
import { ProtectedRoute, useGetAuthorization } from './authorization';
import { Route, Routes } from 'react-router-dom';

import { useAuth0 } from '@auth0/auth0-react';

const AppRouter = (): JSX.Element => {
  const isTrainer = useGetAuthorization(["Trainer"])

  return (
    <Routes>
      <Route element={<PageLayout />}>
        <Route path='' element={<Hero />} />
        <Route path='profile' element={<ProtectedRoute component={Profile} />} />
        <Route path='trainings' element={<Trainings />} />
        <Route path='trainer' element={isTrainer ? <Trainer /> : <NotAuthorized />}>
          <Route path='plan-group' element={<PlanTraining />} />
        </Route>
        <Route path='*' element={<NotFound />} />
      </Route>
    </Routes>
  );
}

const App: React.FC = () => {
  const { isLoading } = useAuth0();
  return isLoading ? <Loading /> : <AppRouter />
};

export default App;


