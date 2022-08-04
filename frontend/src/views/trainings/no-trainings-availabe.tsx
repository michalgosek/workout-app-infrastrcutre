import { Container } from 'react-bootstrap';
import { FC } from 'react';

const NoTrainingsAvailable: FC = () => {
    return (
        <Container className='text-center'>
            <h1>There is no trainings available 😭</h1>
        </Container>
    )
};

export default NoTrainingsAvailable;