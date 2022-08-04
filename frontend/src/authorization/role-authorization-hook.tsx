import { useEffect, useState } from 'react';

import { useAuth0 } from '@auth0/auth0-react';

export const TRAININGS_SERVICE_ROLES = {
    TRAINER: 'Trainer',
    PARTICIPANT: 'Participant',
    USER_CLAIMS: 'trainings-service/roles'
};

export const useGetAuthorization = (roles: string[]): boolean => {
    const [userRoles, setUserRoles] = useState([""]);
    const { getIdTokenClaims } = useAuth0();
    useEffect(() => {
        const getUserRoles = async () => {
            try {
                const claims = await getIdTokenClaims();
                if (claims) setUserRoles(claims[TRAININGS_SERVICE_ROLES.USER_CLAIMS]);
            } catch (err) {
                console.log(err);
            }
        };
        getUserRoles();
    });
    return roles.some(r => userRoles.includes(r));
}
