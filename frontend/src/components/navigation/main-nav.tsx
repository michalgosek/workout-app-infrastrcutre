import { FC } from 'react';
import { GetLinkClassName } from './utils';
import { Nav } from 'react-bootstrap';
import { NavLink } from 'react-router-dom';
import ParticipantNav from './participant-nav';
import { TRAININGS_SERVICE_ROLES } from 'authorization/role-authorization-hook';
import TrainerNav from './trainer-nav';
import { useAuth0 } from '@auth0/auth0-react';
import { useGetAuthorization } from 'authorization';

const Greeter: FC = () => {
    const { user } = useAuth0();
    const isTrainer = useGetAuthorization([TRAININGS_SERVICE_ROLES.TRAINER]);
    const isParticipant = useGetAuthorization([TRAININGS_SERVICE_ROLES.PARTICIPANT]);

    if (!user) return null;

    const name = user.name;
    const role = isTrainer ? TRAININGS_SERVICE_ROLES.TRAINER : isParticipant ? TRAININGS_SERVICE_ROLES.PARTICIPANT : "";
    return (
        <div className='text-left'>
            <p className='small'> Hello, <b>{name}</b> 👋 ! You're logged as <b>{role}</b>. Have fun! 💪</p>
        </div>
    );
}

const MainNav: FC = () => {
    const isTrainer = useGetAuthorization([TRAININGS_SERVICE_ROLES.TRAINER]);
    const isParticipant = useGetAuthorization([TRAININGS_SERVICE_ROLES.PARTICIPANT]);
    return (
        <>

            <Nav className='navbar-nav mr-auto'>
                <NavLink to='/' className={({ isActive }) => GetLinkClassName(isActive)}>
                    Home
                </NavLink>
                <NavLink to='/profile' className={({ isActive }) => GetLinkClassName(isActive)}>
                    Profile
                </NavLink>
                <NavLink to='/trainings' className={({ isActive }) => GetLinkClassName(isActive)}>
                    Trainings
                </NavLink>
                {isTrainer ? <TrainerNav /> : isParticipant ? <ParticipantNav /> : null}
            </Nav >
            <Greeter />
        </>




    );
};

export default MainNav;
