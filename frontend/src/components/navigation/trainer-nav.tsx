import {Fragment, FC} from 'react';
import {NavLink} from "react-router-dom";
import {GetLinkClassName} from "./utils";

const TrainerNav: FC = () => {
    return(
        <Fragment>
            <NavLink to="/trainer/plan-group" className={({ isActive }) => GetLinkClassName(isActive)}>
                Plan Training Group
            </NavLink>
            <NavLink to="/trainer/list-groups" className={({ isActive }) => GetLinkClassName(isActive)}>
                List Groups
            </NavLink>
        </Fragment>
    )
}

export default TrainerNav;
