import { Container } from 'react-bootstrap';
import React from 'react';

const ServiceUnavailable: React.FC = () => {
    return (
        <Container className='text-center'>
            <div className='row'>
                <h2>Service Unavailable :( </h2>
            </div>
        </Container>
    );
};


export default ServiceUnavailable;
