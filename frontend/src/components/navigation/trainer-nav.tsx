import { FC } from 'react';
import { GetLinkClassName } from './utils';
import { NavLink } from 'react-router-dom';
import { useAuth0 } from '@auth0/auth0-react';

const TrainerNav: FC = () => {
    const { user } = useAuth0();
    if (!user) return null;

    const TRAINING_SCHEDULE_PATH = `trainers/${user.sub}/trainings/schedule`;
    const TRAINER_SCHEDULES_PATH = `trainers/${user.sub}/trainings`;

    return (
        <>
            <NavLink to={TRAINING_SCHEDULE_PATH} className={({ isActive }) => GetLinkClassName(isActive)}>
                Plan Training Group
            </NavLink>
            <NavLink to={TRAINER_SCHEDULES_PATH} className={({ isActive }) => GetLinkClassName(isActive)}>
                List Groups
            </NavLink>
        </>
    )
}

export default TrainerNav;
