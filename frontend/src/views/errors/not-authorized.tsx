import { Container } from 'react-bootstrap';
import React from 'react';

const NotAuthorized: React.FC = () => {
    return (
       <Container>
            <h2>You're not authorized 🤥 ⛔️ </h2>
        </Container>
    );
};

export default NotAuthorized;
