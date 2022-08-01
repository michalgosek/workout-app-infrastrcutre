import { GetLinkClassName } from './utils';
import { NavLink } from 'react-router-dom';
import { useGetAuthorization } from '../../authorization/role-authorization-hook';

const MainNav: React.FC = () => {
    const isTrainer = useGetAuthorization(["Trainer"])
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
            {isTrainer ?
                <NavLink to="/trainer/plan-group" className={({ isActive }) => GetLinkClassName(isActive)}>
                    Plan Training Group
                </NavLink>

                : null}
        </div >
    );
};

export default MainNav;
