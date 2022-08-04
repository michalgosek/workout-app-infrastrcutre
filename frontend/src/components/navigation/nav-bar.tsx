import { Container, Nav } from 'react-bootstrap';

import AuthNav from './auth-nav';
import MainNav from './main-nav';

const NavBar: React.FC = () => {
    return (
        <div className='nav-container mb-3'>
            <Nav className='navbar navbar-expand-md navbar-light bg-light'>
                <Container>
                    <div className='navbar-brand'/>
                    <MainNav />
                    <AuthNav />
                </Container>
            </Nav>
        </div >
    );
};

export default NavBar;