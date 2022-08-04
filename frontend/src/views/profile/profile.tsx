import { Container, Image } from 'react-bootstrap';

import { useAuth0 } from '@auth0/auth0-react';

const Profile: React.FC = () => {
    const { user } = useAuth0();
    if (!user) {
        return null
    }
    return (
        <>
            <Container className='row algin-items-center profile-header'>
                <Container className='col-md-3 mb-3'>
                    <Image
                        src={user.picture}
                        alt='Profile'
                        className='rounded-circle img-fluid profile-picture mb-3 mb-md-0' />
                </Container>
                <Container className='col-md text-center text-md-left'>
                    <h2>{user.name}</h2>
                    <p className='lead text-muted'>{user.email}</p>
                </Container>
            </Container>
        </>
    );
};

export default Profile;
