import { GetLinkClassName } from './utils';
import { Nav } from 'react-bootstrap';
import { NavLink } from 'react-router-dom';
import ParticipantNav from './participant-nav';
import { TRAININGS_SERVICE_ROLES } from 'authorization/role-authorization-hook';
import TrainerNav from './trainer-nav';
import { useGetAuthorization } from 'authorization';

const MainNav: React.FC = () => {
    const isTrainer = useGetAuthorization([TRAININGS_SERVICE_ROLES.TRAINER]);
    const isParticipant = useGetAuthorization([TRAININGS_SERVICE_ROLES.PARTICIPANT]);

    return (
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
    );
};

export default MainNav;
