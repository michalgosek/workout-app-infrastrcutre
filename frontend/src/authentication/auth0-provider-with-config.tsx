import React, { PropsWithChildren } from "react";

import { Auth0Provider } from "@auth0/auth0-react";
import StartupFailed from "views/startup-failed";

interface Auth0ProviderWithRedirectCallbackProps {
    children: React.ReactNode;
}

const Auth0ProviderWithConfig = ({ children }: PropsWithChildren<Auth0ProviderWithRedirectCallbackProps>): JSX.Element | null => {
    const AUTH0_CLIENT_DOMAIN: string = process.env.REACT_APP_AUTH0_DOMAIN || "";
    const AUTH0_CLIENT_ID: string = process.env.REACT_APP_AUTH0_CLIENT_ID || "";
    const AUTH0_CALLBACK_URL: string = process.env.REACT_APP_AUTH0_CALLBACK_URL || "";
    const AUTH0_AUDIENCE: string = process.env.REACT_APP_AUTH0_AUDIENCE || "";

    const evns: string[] = [
        AUTH0_CLIENT_DOMAIN,
        AUTH0_CLIENT_ID,
        AUTH0_CALLBACK_URL,
        AUTH0_AUDIENCE
    ]
    const isEnvsMissing: boolean = evns.some(el => el === "")
    if (isEnvsMissing) {
        return (
            <StartupFailed />
        );
    }
    return (
        <Auth0Provider
            domain={AUTH0_CLIENT_DOMAIN}
            clientId={AUTH0_CLIENT_ID}
            redirectUri={AUTH0_CALLBACK_URL}
            audience={AUTH0_AUDIENCE}
            scope={"read:current_user update:current_user_metadata"} 
        >
            {children}
        </Auth0Provider>
    );
};


export default Auth0ProviderWithConfig;