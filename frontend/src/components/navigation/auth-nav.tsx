import AuthenticationButton from 'components/buttons/authentication-button';
import { Nav } from 'react-bootstrap';

const AuthNav: React.FC = () => {
    return (
        <Nav className='navbar-nav ml-auto'>
                <AuthenticationButton />
        </Nav>
    );
};

export default AuthNav;