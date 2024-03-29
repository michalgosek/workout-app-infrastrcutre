import { Button } from "react-bootstrap";
import React from "react";
import { useAuth0 } from "@auth0/auth0-react";

const LogouButton: React.FC = () => {
    const { logout } = useAuth0();
    return (
        <Button
            className="btn btn-danger btn-block"
            onClick={() =>
                logout({
                    returnTo: window.location.origin,
                })
            }
        >
            Log Out
        </Button>
    );
};

export default LogouButton;
