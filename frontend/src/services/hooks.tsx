import { useEffect, useState } from "react";

import { useAuth0 } from "@auth0/auth0-react";

export const useGetTrainingsServiceAccessToken = (): string => {
    const { getAccessTokenSilently } = useAuth0();
    const [token, setToken] = useState('');

    useEffect(() => {
        getAccessTokenSilently()
            .then(token => setToken(token))
            .catch(err => console.error(err))
    }, [getAccessTokenSilently]);

    return token;
}