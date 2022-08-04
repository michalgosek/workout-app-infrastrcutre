import {FC} from 'react';
import {GetLinkClassName} from './utils';
import {NavLink} from 'react-router-dom';

const TrainerNav: FC = () => {
    return(
        <>
            <NavLink to='trainer/trainings/schedule' className={({ isActive }) => GetLinkClassName(isActive)}>
                Plan Training Group
            </NavLink>
            <NavLink to='trainer/trainings/list' className={({ isActive }) => GetLinkClassName(isActive)}>
                List Groups
            </NavLink>
        </>
    )
}

export default TrainerNav;
