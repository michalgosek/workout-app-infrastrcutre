import { FC } from 'react';
import { GetLinkClassName } from './utils';
import { NavLink } from 'react-router-dom';
import { useAuth0 } from '@auth0/auth0-react';

const ParticipantNav: FC = () => {
    const { user } = useAuth0();
    if (!user) return null;

    const PARTICIPANTS_SCHEDULES_PATH = `participants/${user.sub}/trainings`;

    return (
        <NavLink to={PARTICIPANTS_SCHEDULES_PATH} className={({ isActive }) => GetLinkClassName(isActive)}>
            Planned Workouts
        </NavLink>

    );
}

export default ParticipantNav;
