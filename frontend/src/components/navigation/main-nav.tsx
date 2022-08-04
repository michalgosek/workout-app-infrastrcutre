import { GetLinkClassName } from './utils';
import { Nav } from 'react-bootstrap';
import { NavLink } from 'react-router-dom';
import TrainerNav from './trainer-nav';
import { useGetAuthorization } from 'authorization';

const MainNav: React.FC = () => {
    const isTrainer = useGetAuthorization(['Trainer'])
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
            {isTrainer ?
                    <TrainerNav />
                : null}
        </Nav >
    );
};

export default MainNav;
