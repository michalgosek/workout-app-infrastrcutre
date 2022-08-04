import { Button } from 'react-bootstrap';
import { useAuth0 } from '@auth0/auth0-react';

const SignupButton: React.FC = () => {
    const { loginWithRedirect } = useAuth0();
    return (
        <Button
            className='btn btn-primary btn-block'
            onClick={() =>
                loginWithRedirect({
                    screen_hint: 'signup'
                })}
        >
            Sign Up
        </Button >
    );
};

export default SignupButton;