import { FC } from 'react';
import { GetLinkClassName } from './utils';
import { NavLink } from 'react-router-dom';

const ParticipantNav: FC = () => {
    return (
        <>
            <NavLink to='participants/trainings/list' className={({ isActive }) => GetLinkClassName(isActive)}>
                Planned Workouts
            </NavLink>
        </>
    )
}

export default ParticipantNav;
