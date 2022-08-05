import './assets/styles/app.css';

import { Hero, Loading, } from './components';
import { NotAuthorized, NotFound } from './views/errors';
import { PageLayout, Profile, Trainings } from './views';
import { PlanTraining, Trainer, TrainerGroup, TrainerGroups } from './views/trainings/trainer';
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

        <Route path='participants' element={isParticipant ? <Participant /> : <NotAuthorized />}>
          <Route path=':participantUUID'>
            <Route path='trainings' >
              <Route path='' element={<ParticipantTrainingGroups />} />
            </Route>
          </Route>
        </Route>

        <Route path='trainers' element={isTrainer ? <Trainer /> : <NotAuthorized />}>
          <Route path=':trainerUUID'>
            <Route path='trainings' >
              <Route path='' element={<TrainerGroups />} />
              <Route path=':trainingUUID' element={<TrainerGroup />} />
              <Route path='schedule' element={<PlanTraining />} />

            </Route>
          </Route>
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
