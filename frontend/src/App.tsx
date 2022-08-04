import './assets/styles/app.css';

import { Hero, Loading, } from './components';
import { NotAuthorized, NotFound } from './views/errors';
import { PageLayout, Profile, Trainings } from './views';
import { PlanTraining, Trainer, TrainerTrainingGroups } from './views/trainings/trainer';
import { ProtectedRoute, useGetAuthorization } from './authorization';
import { Route, Routes } from 'react-router-dom';

import Participant from 'views/trainings/participant/participant';
import ParticipantTrainingGroups from 'views/trainings/participant/trainings';
import { TRAININGS_SERVICE_ROLES } from 'authorization/role-authorization-hook';
import { useAuth0 } from '@auth0/auth0-react';

const AppRouter = (): JSX.Element => {
  const isTrainer = useGetAuthorization([TRAININGS_SERVICE_ROLES.TRAINER]);
  const isParticipant = useGetAuthorization([TRAININGS_SERVICE_ROLES.PARTICIPANT]);

  return (
    <Routes>
      <Route element={<PageLayout />}>
        <Route path='' element={<Hero />} />
        <Route path='profile' element={<ProtectedRoute component={Profile} />} />
        <Route path='trainings' element={<Trainings />} />

        <Route path='participants/trainings' element={isParticipant ? <Participant /> : <NotAuthorized />}>
          <Route path='list' element={<ParticipantTrainingGroups />} />
        </Route>
        <Route path='trainer/trainings' element={isTrainer ? <Trainer /> : <NotAuthorized />}>
          <Route path='schedule' element={<PlanTraining />} />
          <Route path='list' element={<TrainerTrainingGroups />} />
        </Route>
        <Route path='*' element={<NotFound />} />
      </Route>
    </Routes >
  );
}

const App: React.FC = () => {
  const { isLoading } = useAuth0();
  return isLoading ? <Loading /> : <AppRouter />
};

export default App;
