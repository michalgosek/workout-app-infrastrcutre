import { GetLinkClassName } from './utils';
import { NavLink } from 'react-router-dom';

const MainNav: React.FC = () => {
    return (
        <div className="navbar-nav mr-auto">
            <NavLink to="/" className={({ isActive }) => GetLinkClassName(isActive)}>
                Home
            </NavLink>
            <NavLink to="/profile" className={({ isActive }) => GetLinkClassName(isActive)}>
                Profile
            </NavLink>
            <NavLink to="/trainings" className={({ isActive }) => GetLinkClassName(isActive)}>
                Trainings
            </NavLink>
        </div >
    );
};

export default MainNav;
